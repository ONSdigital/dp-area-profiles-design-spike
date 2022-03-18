package store

import (
	"context"
	"github.com/ONSdigital/dp-area-profiles-design-spike/models"
	log "github.com/daiLlew/funkylog"
	"github.com/pkg/errors"
	"time"
)

// Create tables SQL
var (
	// createKeyStatsTableSQL SQL statement to create the area profiles key statistics table.
	createKeyStatsTableSQL = "CREATE TABLE IF NOT EXISTS key_stats (stat_id INT PRIMARY KEY NOT NULL, profile_id INT NOT NULL, name VARCHAR(100) NOT NULL, value VARCHAR(100) NOT NULL, unit VARCHAR(25) NOT NULL, date_created TIMESTAMP NOT NULL, UNIQUE (profile_id, name), CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES area_profiles (profile_id));"

	// createMetadataTableSQL SQL statement to create the key stat metadata table.
	createMetadataTableSQL = "CREATE TABLE IF NOT EXISTS stat_metadata (metadata_id INT PRIMARY KEY NOT NULL, stat_id INT NOT NULL, dataset_id VARCHAR(100) NOT NULL, dataset_name VARCHAR(100) NOT NULL, href VARCHAR(100), CONSTRAINT fk_stat_id FOREIGN KEY (stat_id) REFERENCES key_stats (stat_id));"
)

// Create sequences SQL
var (
	//createAreaProfileIDSeqSQL is a SQL statement creating a sequence for generating area profile ids.
	createKeyStatsIDSeqSQL = "CREATE SEQUENCE key_stat_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY key_stats.stat_id;"

	// createMetadataIDSeqSQL SQL statement to create the key stat metadata ID sequence.
	createMetadataIDSeqSQL = "CREATE SEQUENCE stat_metadata_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY stat_metadata.metadata_id;"
)

// Query SQL
var (
	// insertNewKeyStatSQL is an SQL query to insert a new key stat.
	insertNewKeyStatSQL = "INSERT INTO key_stats (stat_id, profile_id, name, value, unit, date_created) VALUES (nextval('key_stat_id'), $1, $2, $3, $4, $5) RETURNING stat_id;"

	// insertMetadataSQL SQL to insert a new stat metadata entry
	insertMetadataSQL = "INSERT INTO stat_metadata (metadata_id, stat_id, dataset_id, dataset_name, href) VALUES (nextval('stat_metadata_id'), $1, $2, $3, $4)"

	// updateKeyStatSQL SQL statement to update key stats value, unit and date created fields.
	updateKeyStatSQL = "UPDATE key_stats SET value = $1, unit = $2, date_created = $3 WHERE profile_id = $4 AND name = $5 RETURNING stat_id;"

	// getKeyStatsForProfileIDSQL SQL statement that returning a list of key stats with the related metadata for the specified area profile
	getKeyStatsForProfileIDSQL = "SELECT s.stat_id, s.profile_id, s.name, s.value, s.unit, s.date_created, m.metadata_id, m.dataset_id, m.dataset_name, m.href FROM key_stats s INNER JOIN stat_metadata m ON m.stat_id = s.stat_id WHERE s.profile_id = $1;"

	// updateMetadataSQL SQL statement updating stat metadata fields withe the specified values.
	updateMetadataSQL = "UPDATE stat_metadata SET dataset_id = $1, dataset_name = $2, href = $3 WHERE stat_id = $4;"
)

// GetKeyStatsForProfileID return a list of key stats for the specified area profile ID.
func (s *areaProfileStore) GetKeyStatsByProfileID(profileID int) ([]models.KeyStatistic, error) {
	stats := make([]models.KeyStatistic, 0)

	rows, err := s.conn.Query(context.Background(), getKeyStatsForProfileIDSQL, profileID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			statID      int
			profileID   int
			name        string
			value       string
			unit        string
			dateCreated time.Time
			metadataID  int
			datasetID   string
			datasetName string
			href        string
		)

		if err := rows.Scan(&statID, &profileID, &name, &value, &unit, &dateCreated, &metadataID, &datasetID, &datasetName, &href); err != nil {
			return nil, err
		}

		stats = append(stats, models.KeyStatistic{
			StatID:      statID,
			ProfileID:   profileID,
			Name:        name,
			Value:       value,
			Unit:        unit,
			DateCreated: dateCreated,
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

// NewKeyStat insert a key statistic for the specified area profile.
func (s *areaProfileStore) AddKeyStat(profileID int, name, value, unit string, dateCreated time.Time, datasetID, datasetName, href string) (int, error) {
	statID, err := s.insertKeyStatistic(profileID, name, value, unit, dateCreated)
	if err != nil {
		return 0, err
	}

	if err = s.insertKeyStatisticMetadata(statID, datasetID, datasetName, href); err != nil {
		return 0, nil
	}

	return statID, nil
}

func (s *areaProfileStore) insertKeyStatistic(profileID int, name, value, unit string, dateCreated time.Time) (int, error) {
	var keyStatID int
	err := s.conn.QueryRow(context.Background(), insertNewKeyStatSQL, profileID, name, value, unit, dateCreated).Scan(&keyStatID)
	if err != nil {
		return 0, errors.Wrapf(err, "error inserting new key stat %q for profile_id=%d", name, profileID)
	}

	return keyStatID, nil
}

func (s *areaProfileStore) insertKeyStatisticMetadata(statID int, datasetID, name, href string) error {
	_, err := s.conn.Exec(context.Background(), insertMetadataSQL, statID, datasetID, name, href)
	if err != nil {
		return errors.Wrapf(err, "error inserting metadata for key stat %d", statID)
	}

	return nil
}

// UpdateProfileKeyStats creates a new version of key statistics for the specified area profile.
// The current key stats are superceded and added to the key stats history table before updating
// the new value for each of the key stats in the new current version.
func (s *areaProfileStore) UpdateProfileKeyStats(profileID int, newStats []models.ImportRow) error {
	currentKeyStats, err := s.GetKeyStatsByProfileID(profileID)
	if err != nil {
		return err
	}

	if err := s.versionCurrentProfileKeyStats(profileID, currentKeyStats); err != nil {
		return err
	}

	if err := s.addNewProfileKeyStats(profileID, newStats); err != nil {
		return err
	}

	return nil
}

func (s *areaProfileStore) addNewProfileKeyStats(profileID int, newStats []models.ImportRow) error {
	log.Info("inserting latest key stats for for profile ID: %s", profileID)
	now := time.Now()

	for _, stat := range newStats {
		var statID int

		// Update the key stats figures with the new values.
		err := s.conn.QueryRow(context.Background(), updateKeyStatSQL, stat.Value, stat.Unit, now, profileID, stat.Name).Scan(&statID)
		if err != nil {
			return errors.Wrapf(err, "error updating key stat %q for profile_id=%d", stat.Name, profileID)
		}

		// Update ket stats metadata with the latest values.
		_, err = s.conn.Exec(context.Background(), updateMetadataSQL, stat.DatasetID, stat.DatasetName, stat.GetDatasetHref(), statID)
		if err != nil {
			return errors.Wrapf(err, "error updating metadata for key stat %d", statID)
		}
	}

	return nil
}
