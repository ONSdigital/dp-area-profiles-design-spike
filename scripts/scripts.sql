
CREATE TABLE stat_metadata (metadat_id INT PRIMARY KEY NOT NULL, dataset_name VARCHAR(100) NOT NULL, href VARCHAR(100) NOT NULL);





CREATE TABLE IF NOT EXISTS key_stats (stat_id INT PRIMARY KEY NOT NULL, profile_id INT NOT NULL, metadat_id INT NOT NULL, name VARCHAR(100) NOT NULL, value VARCHAR(100) NOT NULL, unit VARCHAR(25) NOT NULL, date_created TIMESTAMP NOT NULL, UNIQUE (profile_id, name), CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES area_profiles (profile_id)), CONSTRAINT fk_metadata_id FOREIGN KEY (metadat_id) REFERENCES stat_metadata (metadat_id));