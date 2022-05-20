package store

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

var (
	// createKeyStatsTableSQL SQL statement to create the area profiles key statistics table.
	createKeyStatsTableSQL = `
		CREATE TABLE IF NOT EXISTS key_stats (
			stat_id INT PRIMARY KEY NOT NULL, 
			profile_id INT NOT NULL,
			stat_type INT NOT NULL, 
			value VARCHAR(100) NOT NULL, 
			unit VARCHAR(25) NOT NULL, 
			date_created TIMESTAMP NOT NULL, 
			dataset_id VARCHAR(100) NOT NULL, 
			dataset_name VARCHAR(100) NOT NULL, 
			UNIQUE (profile_id, stat_type), 
			CONSTRAINT fk_profile_id 
				FOREIGN KEY (profile_id) REFERENCES area_profiles (profile_id),
			CONSTRAINT fk_stat_type 
				FOREIGN KEY (stat_type) REFERENCES key_stat_types (type_id) 
		);
	`

	//createAreaProfileIDSeqSQL is a SQL statement creating a sequence for generating area profile ids.
	createKeyStatsIDSeqSQL = `
		CREATE SEQUENCE 
			key_stat_id 
		START 
			1000 
		INCREMENT 
			100 
		MINVALUE 
			1000 
		OWNED BY 
			key_stats.stat_id;
	`

	// insertNewKeyStatSQL is an SQL query to insert a new key stat.
	insertNewKeyStatSQL = `
		INSERT INTO key_stats 
			(stat_id, profile_id, stat_type, value, unit, date_created, dataset_id, dataset_name) 
		VALUES 
			(nextval('key_stat_id'), $1, $2, $3, $4, $5, $6, $7) 
		ON CONFLICT ON CONSTRAINT 
			key_stats_profile_id_stat_type_key 
		DO UPDATE SET value = $3 RETURNING stat_id;
	`

	// getStatsByProfileIDSQL SQL query returns current version of the key statistics for the specified area profile.
	getStatsByProfileIDSQL = `
		SELECT 
			s.profile_id, s.stat_id, s.stat_type, t.name, s.value, s.unit, s.date_created, s.dataset_id, s.dataset_name 
		FROM 
			key_stats s
		INNER JOIN
			key_stat_types t
		ON
			t.type_id = s.stat_type
		WHERE 
			s.profile_id = $1;
	`
)

// NewKeyStat insert a key statistic for the specified area profile.
func (s *AreaProfileStore) InsertKeyStat(areaCode, name, value, unit, datasetID, datasetName string, dateCreated time.Time) (int, error) {
	profile, err := s.GetProfileByAreaCode(areaCode)
	if err != nil {
		return 0, err
	}

	statType, err := s.GetStatTypeByName(name)
	if err != nil {
		return 0, err
	}

	var keyStatID int

	err = s.conn.QueryRow(context.Background(), insertNewKeyStatSQL, profile.ID, statType, value, unit, dateCreated, datasetID, datasetName).Scan(&keyStatID)
	if err != nil {
		return 0, errors.Wrapf(err, "error inserting new key stat %q for profile_id=%d", name, profile.ID)
	}

	_, err = s.conn.Exec(context.Background(), insertNewKeyStatHistorySQL, profile.ID, statType, value, unit, dateCreated, dateCreated, datasetID, datasetName)
	if err != nil {
		return 0, errors.Wrapf(err, "error inserting key stat history %q for profile_id=%d", name, profile.ID)
	}

	return keyStatID, nil
}

// GetKeyStatsForProfile returns a list of the current Key stats associated with the specified area profile.
func (s *AreaProfileStore) GetKeyStatsForProfile(profile *AreaProfile) (KeyStatistics, error) {
	rows, err := s.conn.Query(context.Background(), getStatsByProfileIDSQL, profile.ID)
	if err != nil {
		return nil, err
	}

	stats, err := keyStatisticsRowsMapper(profile, rows)
	if err != nil {
		return nil, errors.Wrap(err, "error mapping result rows to keystatistics")
	}

	return stats, nil
}
