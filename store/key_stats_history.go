package store

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-area-profiles-design-spike/models"
	log "github.com/daiLlew/funkylog"
	"github.com/pkg/errors"
	"time"
)

// Create tables SQl
var (
	// createKeyStatsHistoryTableSQL SQL statement to create the key statistics history table.
	createKeyStatsHistoryTableSQL = "CREATE TABLE IF NOT EXISTS key_stats_history (stat_id INT PRIMARY KEY NOT NULL, profile_id INT NOT NULL, version_id int NOT NULL, name VARCHAR(100) NOT NULL, value VARCHAR(100) NOT NULL, unit VARCHAR(25) NOT NULL, date_created TIMESTAMP NOT NULL, last_modified TIMESTAMP NOT NULL, UNIQUE (profile_id, version_id, name), CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES area_profiles (profile_id));"

	// createMetadataTableSQL SQL statement to create the key stat metadata table.
	createMetadataHistoryTableSQL = "CREATE TABLE IF NOT EXISTS stat_metadata_history (metadata_id INT PRIMARY KEY NOT NULL, stat_id INT NOT NULL, dataset_id VARCHAR(100) NOT NULL, dataset_name VARCHAR(100) NOT NULL, href VARCHAR(100), last_modified TIMESTAMP NOT NULL, CONSTRAINT fk_stat_id FOREIGN KEY (stat_id) REFERENCES key_stats_history (stat_id));"
)

// Create sequences SQL
var (
	//createKeyStatVersionIDSeqSQL is a SQL statement creating a sequence for generating key stat version ids.
	createKeyStatVersionIDSeqSQL = "CREATE SEQUENCE key_stat_version_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY key_stats.stat_id;"

	//createKeyStatsHistoryIDSeqSQL is a SQL statement creating a sequence for generating area profile ids.
	createKeyStatsHistoryIDSeqSQL = "CREATE SEQUENCE key_stat_history_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY key_stats_history.stat_id;"

	// createMetadataHistoryIDSeqSQL SQL statement to create the key stat metadata ID sequence.
	createMetadataHistoryIDSeqSQL = "CREATE SEQUENCE stat_metadata_history_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY stat_metadata_history.metadata_id;"
)

// Query SQL
var (
	// insertNewKeyStatHistorySQL is an SQL query to insert a new key stat version.
	insertNewKeyStatHistorySQL = "INSERT INTO key_stats_history (stat_id, profile_id, version_id, name, value, unit, date_created, last_modified) VALUES (nextval('key_stat_history_id'), $1, $2, $3, $4, $5, $6, $7) RETURNING stat_id;"

	// getKeyStatVersionIDSQL SQL query to get the next ID from the key stat version sequence.
	getKeyStatVersionIDSQL = "SELECT nextval('key_stat_version_id')"

	// getKeyStatVersionsSQL SQL query to get the key stat versions for the specified profile ID.
	getKeyStatVersionsSQL = "SELECT DISTINCT ON (date_created) date_created, stat_id, profile_id, version_id FROM key_stats_history WHERE profile_id = $1 ORDER by date_created DESC;"

	// getKeyStatVersionSQL SQL query returns key stat history data and related metadata belonging to the specified version;
	getKeyStatVersionSQL = "SELECT version_id, s.stat_id, s.profile_id, s.name, s.value, s.unit, s.date_created, s.last_modified, m.metadata_id, m.dataset_id, m.dataset_name, m.href FROM key_stats_history s INNER JOIN stat_metadata_history m ON m.stat_id = s.stat_id WHERE s.profile_id = $1 AND s.version_id = $2;"

	// insertMetadataSQL SQL to insert a new stat metadata entry
	insertMetadataHistorySQL = "INSERT INTO stat_metadata_history (metadata_id, stat_id, dataset_id, dataset_name, href, last_modified) VALUES (nextval('stat_metadata_history_id'), $1, $2, $3, $4, $5)"
)

// GetKeyStatsVersions returns a list of KeyStatsVersion available for the specified area profile ID.
func (s *areaProfileStore) GetKeyStatsVersions(areaCode string, profileID int) ([]models.KeyStatsVersion, error) {
	rows, err := s.conn.Query(context.Background(), getKeyStatVersionsSQL, profileID)
	if err != nil {
		return nil, errors.Wrap(err, "error querying for versions")
	}

	defer rows.Close()

	versions := make([]models.KeyStatsVersion, 0)
	for rows.Next() {
		var (
			dateCreated time.Time
			statID      int
			profileID   int
			versionID   int
		)

		if err := rows.Scan(&dateCreated, &statID, &profileID, &versionID); err != nil {
			return nil, err
		}

		versions = append(versions, models.KeyStatsVersion{
			StatID:      statID,
			ProfileID:   profileID,
			VersionID:   versionID,
			DateCreated: dateCreated,
			Href:        fmt.Sprintf("http://localhost:8080/profiles/%s/versions/%d", areaCode, versionID),
		})
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return versions, nil
}

// GetKeyStatsVersion return a list of KeyStatistics associated with the specified version of the area profile ID.
func (s *areaProfileStore) GetKeyStatsVersion(profileID, versionID int) ([]models.KeyStatistic, error) {
	rows, err := s.conn.Query(context.Background(), getKeyStatVersionSQL, profileID, versionID)
	if err != nil {
		return nil, errors.Wrap(err, "error querying for versions")
	}

	defer rows.Close()

	stats := make([]models.KeyStatistic, 0)
	for rows.Next() {
		var (
			versionID    int
			statID       int
			profileID    int
			name         string
			value        string
			unit         string
			dateCreated  time.Time
			lastModified time.Time
			metadataID   int
			datasetID    string
			datasetName  string
			href         string
		)

		if err := rows.Scan(&versionID, &statID, &profileID, &name, &value, &unit, &dateCreated, &lastModified, &metadataID, &datasetID, &datasetName, &href); err != nil {
			return nil, err
		}

		stats = append(stats, models.KeyStatistic{
			VersionID:    versionID,
			StatID:       statID,
			ProfileID:    profileID,
			Name:         name,
			Value:        value,
			Unit:         unit,
			DateCreated:  dateCreated,
			LastModified: lastModified,
			Metadata: models.KeyStatisticMetadata{
				MetadataID:  metadataID,
				DatasetID:   datasetID,
				DatasetName: datasetName,
				Link:        href,
			},
		})
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return stats, nil
}

// versionCurrentProfileKeyStats moves the current area profile key stats and related metadata into the history tables.
func (s *areaProfileStore) versionCurrentProfileKeyStats(profileID int, currentKeyStats []models.KeyStatistic) error {
	log.Info("adding current key stats to version history table, profile ID: %s", profileID)

	// stat version ID groups together all key statistics that belong to a particular version.
	statVersionID, err := s.getNextKeyStatHistoryVersionID()
	if err != nil {
		return err
	}

	lastUpdated := time.Now()
	for _, stat := range currentKeyStats {
		var statHistoryID int
		err := s.conn.QueryRow(context.Background(), insertNewKeyStatHistorySQL, profileID, statVersionID, stat.Name, stat.Value, stat.Unit, stat.DateCreated, lastUpdated).Scan(&statHistoryID)
		if err != nil {
			return errors.Wrapf(err, "error inserting new key stat version %q for profile_id=%d", stat.Name, profileID)
		}

		_, err = s.conn.Exec(context.Background(), insertMetadataHistorySQL, statHistoryID, stat.Metadata.DatasetID, stat.Metadata.DatasetName, stat.Metadata.Link, lastUpdated)
		if err != nil {
			return errors.Wrapf(err, "error inserting stat metadata history dataset_id: %s, stat id: %d", stat.Metadata.DatasetID, stat.StatID)
		}
	}

	return nil
}

func (s *areaProfileStore) getNextKeyStatHistoryVersionID() (int, error) {
	var id int
	if err := s.conn.QueryRow(context.Background(), getKeyStatVersionIDSQL).Scan(&id); err != nil {
		return 0, errors.Wrap(err, "error getting next key stat version id")
	}

	return id, nil
}
