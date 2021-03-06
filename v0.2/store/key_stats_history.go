package store

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

// Key stats history queries/statments.
var (
	// createKeyStatsHistoryTableSQL SQL statement to create the key stats version table.
	createKeyStatsHistoryTableSQL = `
		CREATE TABLE IF NOT EXISTS key_stats_history (
			stat_id INT PRIMARY KEY NOT NULL, 
			profile_id INT NOT NULL, 
			stat_type INT NOT NULL,  
			value VARCHAR(100) NOT NULL, 
			unit VARCHAR(25) NOT NULL, 
			date_created TIMESTAMP NOT NULL, 
			last_modified TIMESTAMP NOT NULL, 
			dataset_id VARCHAR(100) NOT NULL, 
			dataset_name VARCHAR(100) NOT NULL, 
			UNIQUE (profile_id, last_modified, stat_type), 
			CONSTRAINT fk_profile_id 
				FOREIGN KEY (profile_id) REFERENCES area_profiles (profile_id),
			CONSTRAINT fk_stat_type 
				FOREIGN KEY (stat_type) REFERENCES key_stat_types (type_id)
		);
	`
	//createKeyStatsHistoryIDSeqSQL is a SQL statement creating a sequence for generating area profile ids.
	createKeyStatsHistoryIDSeqSQL = `
		CREATE SEQUENCE 
			key_stat_history_id 
		START 
			1000 
		INCREMENT 
			100 
		MINVALUE 
			1000 
		OWNED BY
			key_stats_history.stat_id;
	`

	// insertNewKeyStatHistorySQL is an SQL query to insert a new key stat version.
	insertNewKeyStatHistorySQL = `
		INSERT INTO key_stats_history 
			(stat_id, profile_id, stat_type, value, unit, date_created, last_modified, dataset_id, dataset_name) 
		VALUES 
			(nextval('key_stat_history_id'), $1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING stat_id;`

	// listVersionsSQL SQL query returns a list of key stats versions for an area profile.
	listVersionsSQL = `
		SELECT DISTINCT 
			s.date_created 
		FROM 
			key_stats_history s 
		WHERE 
			s.profile_id = $1 
		ORDER BY 
			s.date_created DESC
	`

	// getKeyStatsVersionSQL SQL query returning key stats for the specified area profile ID and version.
	getKeyStatsVersionSQL = `
		SELECT DISTINCT ON 
			(s.stat_type) s.profile_id, s.stat_id, s.stat_type, t.name, s.value, s.unit, s.date_created, s.dataset_id, s.dataset_name
		FROM 
			key_stats_history s 
		INNER JOIN
			key_stat_types t
		ON
			t.type_id = s.stat_type
		WHERE 
			s.profile_id = $1 AND s.date_created <= $2
		ORDER BY 
			s.stat_type, s.date_created 
		DESC;
	`
)

// GetKeyStatsVersionsForProfile list all versions of the key stats for this area profile
func (s *AreaProfileStore) GetKeyStatsVersionsForProfile(profile *AreaProfile) ([]time.Time, error) {
	rows, err := s.conn.Query(context.Background(), listVersionsSQL, profile.ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	versions, err := versionsRowsMapper(rows)
	if err != nil {
		return nil, errors.Wrap(err, "error mapping rows to key stats versions list")
	}

	return versions, nil
}

// GetKeyStatsVersion returns a list of key stats belonging to the specified version of the area profile.
func (s *AreaProfileStore) GetKeyStatsVersion(profile *AreaProfile, date string) (KeyStatistics, error) {
	rows, err := s.conn.Query(context.Background(), getKeyStatsVersionSQL, profile.ID, date)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	stats, err := keyStatisticsRowsMapper(profile, rows)
	if err != nil {
		return nil, errors.Wrap(err, "error mapping stats version result rows")
	}

	return stats, nil
}
