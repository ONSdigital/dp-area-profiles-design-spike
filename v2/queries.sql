-- Create Area table
CREATE TABLE IF NOT EXISTS areas (
    code VARCHAR (50) PRIMARY KEY NOT NULL, 
    name VARCHAR (100) NOT NULL
);


-- Create Area profile table
CREATE TABLE IF NOT EXISTS area_profiles (
    profile_id INT PRIMARY KEY NOT NULL, 
    area_code VARCHAR(50) NOT NULL, 
    name VARCHAR (100) NOT NULL, 
    UNIQUE (area_code), 
    CONSTRAINT fk_area_code 
        FOREIGN KEY (area_code) REFERENCES areas (code)
);


-- Area profile.ID seq.
CREATE SEQUENCE area_profile_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY area_profiles.profile_id;


-- Create key statistics table
CREATE TABLE IF NOT EXISTS key_stats (
    stat_id INT PRIMARY KEY NOT NULL, 
    profile_id INT NOT NULL, 
    name VARCHAR(100) NOT NULL, 
    value VARCHAR(100) NOT NULL, 
    unit VARCHAR(25) NOT NULL, 
    date_created TIMESTAMP NOT NULL, 
    dataset_id VARCHAR(100) NOT NULL, 
    dataset_name VARCHAR(100) NOT NULL, 
    UNIQUE (profile_id, name), 
    CONSTRAINT fk_profile_id 
        FOREIGN KEY (profile_id) REFERENCES area_profiles (profile_id)
);


-- Create key stats ID sequence.
CREATE SEQUENCE key_stat_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY key_stats.stat_id;

-- Insert key stat. If a key stat with the specified name already exists for this profile ID then update the value.
INSERT INTO key_stats 
    (stat_id, profile_id, name, value, unit, date_created, dataset_id, dataset_name) 
VALUES 
    (nextval('key_stat_id'), $1, $2, $3, $4, $5, $6, $7) 
ON CONFLICT ON CONSTRAINT 
    key_stats_profile_id_name_key 
DO UPDATE SET value = $3 
RETURNING stat_id;


-- Create key stats history table.
CREATE TABLE IF NOT EXISTS key_stats_history (
    stat_id INT PRIMARY KEY NOT NULL, 
    profile_id INT NOT NULL, 
    name VARCHAR(100) NOT NULL, 
    value VARCHAR(100) NOT NULL, 
    unit VARCHAR(25) NOT NULL, 
    date_created TIMESTAMP NOT NULL, 
    last_modified TIMESTAMP NOT NULL, 
    dataset_id VARCHAR(100) NOT NULL, 
    dataset_name VARCHAR(100) NOT NULL, 
    UNIQUE (profile_id, last_modified, name), 
    CONSTRAINT fk_profile_id 
        FOREIGN KEY (profile_id) REFERENCES area_profiles (profile_id)
);


-- Create key_stats_history_id sequence.
CREATE SEQUENCE key_stat_history_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY key_stats_history.stat_id;


-- Insert key stat into history table.
INSERT INTO key_stats_history 
    (stat_id, profile_id, name, value, unit, date_created, last_modified, dataset_id, dataset_name) 
VALUES 
    (nextval('key_stat_history_id'), $1, $2, $3, $4, $5, $6, $7, $8) 
RETURNING stat_id;


-- Get a list of key stat versions for the specified area profile
SELECT DISTINCT 
    s.date_created 
FROM 
    key_stats_history s 
WHERE 
    s.profile_id = $1 
ORDER BY 
    s.date_created 
DESC;


-- Get all key stats that belong to the specified version (versions == date_created timestamp)
SELECT DISTINCT ON 
    (s.name) s.profile_id, s.stat_id, s.name, s.value, s.unit, s.date_created, s.dataset_id, s.dataset_name 
FROM 
    key_stats_history s 
WHERE 
    s.profile_id = $1 AND s.date_created <= $2 
ORDER BY 
    s.name, s.date_created 
DESC;


-- Example
-- Returns all key_stat_history rows for profile_id 1000 that were created on or before 2022-04-07 16:18:05.6306.
-- DISTINCT on name and ordering by date_created will return the lastest value for each unique stat name created on/before thet specified date.
SELECT DISTINCT ON (s.name)
    s.stat_id, 
    s.profile_id,
    s.name,
    s.value, 
    s.unit, 
    s.date_created, 
    s.last_modified, 
    s.dataset_id, 
    s.dataset_name 
FROM 
    key_stats_history s
WHERE
    s.profile_id = 1000 AND s.date_created <= '2022-04-07 16:18:05.6306'
ORDER BY 
    s.name, s.date_created DESC;
