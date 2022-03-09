package store

// Drop sequences/tables.
var (
	// dropSequencesSQL is an SQL statement to drop the sequences created by this demo.
	dropSequencesSQL = "DROP SEQUENCE IF EXISTS area_profile_id, key_stats_id, key_stats_history_id, key_stat_version_id"

	// dropTablesSQL is an SQL statement to drop all tables created by this demo.
	dropTablesSQL = "DROP TABLE IF EXISTS key_stats_history, key_stats, area_profiles, areas CASCADE;"
)

// Create database tables
var (
	// createAreasTableSQL is an SQL statement to create the areas table
	createAreasTableSQL = "CREATE TABLE IF NOT EXISTS areas (code VARCHAR (50) PRIMARY KEY NOT NULL, name VARCHAR (100) NOT NULL);"

	// createProfilesTableSQL is an SQL statement to create the area profiles table
	createProfilesTableSQL = "CREATE TABLE area_profiles (profile_id INT PRIMARY KEY NOT NULL, area_code VARCHAR(50) NOT NULL, name VARCHAR (100) NOT NULL, UNIQUE (area_code), CONSTRAINT fk_area_code FOREIGN KEY (area_code) REFERENCES areas (code));"

	// createKeyStatsTableSQL SQL statement to create the area profiles key statistics table.
	createKeyStatsTableSQL = "CREATE TABLE IF NOT EXISTS key_stats (stat_id INT PRIMARY KEY NOT NULL, profile_id INT NOT NULL, name VARCHAR(100) NOT NULL, value VARCHAR(100) NOT NULL, unit VARCHAR(25) NOT NULL, date_created TIMESTAMP NOT NULL, UNIQUE (profile_id, name), CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES area_profiles (profile_id));"

	// createKeyStatsHistoryTableSQL SQL statement to create the key statistics history table.
	createKeyStatsHistoryTableSQL = "CREATE TABLE IF NOT EXISTS key_stats_history (stat_id INT PRIMARY KEY NOT NULL, profile_id INT NOT NULL, version_id int NOT NULL, name VARCHAR(100) NOT NULL, value VARCHAR(100) NOT NULL, unit VARCHAR(25) NOT NULL, date_created TIMESTAMP NOT NULL, last_modified TIMESTAMP NOT NULL, UNIQUE (profile_id, version_id, name), CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES area_profiles (profile_id));"
)

// Create sequences from primary keys.
var (
	//createAreaProfileIDSeqSQL is a SQL statement creating a sequence for generating area profile ids.
	createAreaProfileIDSeqSQL = "CREATE SEQUENCE area_profile_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY area_profiles.profile_id;"

	//createAreaProfileIDSeqSQL is a SQL statement creating a sequence for generating area profile ids.
	createKeyStatsIDSeqSQL = "CREATE SEQUENCE key_stat_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY key_stats.stat_id;"

	//createKeyStatVersionIDSeqSQL is a SQL statement creating a sequence for generating key stat version ids.
	createKeyStatVersionIDSeqSQL = "CREATE SEQUENCE key_stat_version_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY key_stats.stat_id;"

	//createKeyStatsHistoryIDSeqSQL is a SQL statement creating a sequence for generating area profile ids.
	createKeyStatsHistoryIDSeqSQL = "CREATE SEQUENCE key_stat_history_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY key_stats_history.stat_id;"
)

// Inserts statements.
var (
	// insertAreaSQL is an SQL query to insert a new area - requires area code and name.
	insertAreaSQL = "INSERT INTO areas (code, name) VALUES ($1, $2) RETURNING code;"

	// insertProfileSQL is an SQL query to insert a new area profile, required area code and profile name.
	insertProfileSQL = "INSERT INTO area_profiles (profile_id, area_code, name) VALUES (nextval('area_profile_id'), $1, $2) RETURNING profile_id;"

	// insertNewKeyStatSQL is an SQL query to insert a new key stat.
	insertNewKeyStatSQL = "INSERT INTO key_stats (stat_id, profile_id, name, value, unit, date_created) VALUES (nextval('key_stat_id'), $1, $2, $3, $4, $5);"

	// updateKeyStatSQL SQL statement to update key stats value, unit and date created fields.
	updateKeyStatSQL = "UPDATE key_stats SET value = $1, unit = $2, date_created = $3 WHERE profile_id = $4 AND name = $5"

	// insertNewKeyStatHistorySQL is an SQL query to insert a new key stat version.
	insertNewKeyStatHistorySQL = "INSERT INTO key_stats_history (stat_id, profile_id, version_id, name, value, unit, date_created, last_modified) VALUES (nextval('key_stat_history_id'), $1, $2, $3, $4, $5, $6, $7);"
)

// Queries
var (
	// getProfileByAreaCodeSQL SQL query returns the area profile for the specified area code.
	getProfileByAreaCodeSQL = "SELECT profile_id, name, area_code FROM area_profiles WHERE area_code = $1;"

	// getKeyStatsforProfileIDSQL SQL query returning the key stat details for the specified area profile ID.
	getKeyStatsforProfileIDSQL = "SELECT stat_id, profile_id, name, value, unit, date_created FROM key_stats WHERE profile_id = $1"

	// getKeyStatsForProfileIDSQL SQL query returning key statistics field for the specified area profile.
	getKeyStatsForProfileIDSQL = "SELECT stat_id, profile_id, name, value, unit, date_created FROM key_stats WHERE profile_id = $1;"

	// getKeyStatVersionIDSQL SQL query to get the next ID from the key stat version sequence.
	getKeyStatVersionIDSQL = "SELECT nextval('key_stat_version_id')"

	// getKeyStatVersionsSQL SQL query to get the key stat versions for the specified profile ID.
	getKeyStatVersionsSQL = "SELECT DISTINCT ON (date_created) date_created, stat_id, profile_id, version_id FROM key_stats_history WHERE profile_id = $1 ORDER by date_created DESC;"

	// getKeyStatVersionSQL SQL query to get the key stat for the specified profile ID /version ID.
	getKeyStatVersionSQL = "SELECT version_id, stat_id, profile_id, name, value, unit, date_created, last_modified FROM key_stats_history WHERE profile_id = $1 AND version_id = $2;"
)
