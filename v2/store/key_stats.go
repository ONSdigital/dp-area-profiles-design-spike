package store

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

// Key stats queries/statments.
var (
	// createKeyStatsTableSQL SQL statement to create the area profiles key statistics table.
	createKeyStatsTableSQL = "CREATE TABLE IF NOT EXISTS key_stats (stat_id INT PRIMARY KEY NOT NULL, profile_id INT NOT NULL, name VARCHAR(100) NOT NULL, value VARCHAR(100) NOT NULL, unit VARCHAR(25) NOT NULL, date_created TIMESTAMP NOT NULL, dataset_id VARCHAR(100) NOT NULL, dataset_name VARCHAR(100) NOT NULL, UNIQUE (profile_id, name), CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES area_profiles (profile_id));"

	//createAreaProfileIDSeqSQL is a SQL statement creating a sequence for generating area profile ids.
	createKeyStatsIDSeqSQL = "CREATE SEQUENCE key_stat_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY key_stats.stat_id;"

	// insertNewKeyStatSQL is an SQL query to insert a new key stat.
	insertNewKeyStatSQL = "INSERT INTO key_stats (stat_id, profile_id, name, value, unit, date_created, dataset_id, dataset_name) VALUES (nextval('key_stat_id'), $1, $2, $3, $4, $5, $6, $7) ON CONFLICT ON CONSTRAINT key_stats_profile_id_name_key DO UPDATE SET value = $3 RETURNING stat_id;"

	// getStatsByProfileIDSQL SQL query returns current version of the key statistics for the specified area profile.
	getStatsByProfileIDSQL = "SELECT s.profile_id, s.stat_id, s.name, s.value, s.unit, s.date_created, s.dataset_id, s.dataset_name FROM key_stats s WHERE s.profile_id = $1;"
)

// NewKeyStat insert a key statistic for the specified area profile.
func (s *AreaProfileStore) InsertKeyStat(areaCode, name, value, unit, datasetID, datasetName string, dateCreated time.Time) (int, error) {
	profile, err := s.GetProfileByAreaCode(areaCode)
	if err != nil {
		return 0, err
	}

	var keyStatID int

	err = s.conn.QueryRow(context.Background(), insertNewKeyStatSQL, profile.ID, name, value, unit, dateCreated, datasetID, datasetName).Scan(&keyStatID)
	if err != nil {
		return 0, errors.Wrapf(err, "error inserting new key stat %q for profile_id=%d", name, profile.ID)
	}

	_, err = s.conn.Exec(context.Background(), insertNewKeyStatHistorySQL, profile.ID, name, value, unit, dateCreated, dateCreated, datasetID, datasetName)
	if err != nil {
		return 0, errors.Wrapf(err, "error inserting key stat history %q for profile_id=%d", name, profile.ID)
	}

	return keyStatID, nil
}

// GetKeyStatsForProfile returns a list of the current Key stats associated with the specified area profile.
func (s *AreaProfileStore) GetKeyStatsForProfile(profileID int) (KeyStatistics, error) {
	rows, err := s.conn.Query(context.Background(), getStatsByProfileIDSQL, profileID)
	if err != nil {
		return nil, err
	}

	stats, err := keyStatisticsRowsMapper(rows)
	if err != nil {
		return nil, errors.Wrap(err, "error mapping result rows to keystatistics")
	}

	return stats, nil
}
