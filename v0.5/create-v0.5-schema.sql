-- 
-- SQL Script creates the v0.5 Area profiles schema design and populates it with some test data.
-- NOTE: This is not production ready SQL and is only intended to serve as an interactive illustration the schema design.
--

-- 
-- Clean up. Drop all sequences, data and tables.
-- 
DROP SEQUENCE IF EXISTS 
    area_profile_id,
    key_stats_id,
    key_stats_history_id,
    key_stat_version_id,
    key_stat_type_id,
    recipe_id,
    recipe_geography_id,
    geography_type_id;

DROP TABLE IF EXISTS 
    key_stats_history, 
    key_stats, 
    area_profiles, 
    areas, 
    key_stat_types, 
    key_stats_recipes, 
    recipe_geographies, 
    geography_types 
    CASCADE;

-- 
-- create the "areas" table and insert a test area.
--
CREATE TABLE IF NOT EXISTS areas (
    code VARCHAR (50) PRIMARY KEY NOT NULL, 
    name VARCHAR (100) NOT NULL
);

INSERT INTO areas (code, name) VALUES ('E05011362', 'Test area 1');

-- 
-- create the "area_profiles" table, sequence and insert some test data.
--
CREATE TABLE IF NOT EXISTS area_profiles (
    profile_id INT PRIMARY KEY NOT NULL, 
    area_code VARCHAR(50) NOT NULL, 
    name VARCHAR (100) NOT NULL, 
    UNIQUE (area_code), 
    CONSTRAINT fk_area_code 
        FOREIGN KEY (area_code) REFERENCES areas (code)
);

CREATE SEQUENCE area_profile_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY area_profiles.profile_id;
INSERT INTO area_profiles (profile_id, area_code, name) VALUES (nextval('area_profile_id'), 'E05011362', 'Test area 1');

-- 
-- create the "key_stats_type" table and populate it with some test data.
-- 
CREATE TABLE IF NOT EXISTS key_stat_types (
    type_id INT PRIMARY KEY NOT NULL, 
    name VARCHAR(100) NOT NULL,  
    UNIQUE (name)
);

CREATE SEQUENCE key_stat_type_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY key_stat_types.type_id;
INSERT INTO key_stat_types (type_id, name) VALUES (nextval('key_stat_type_id'), 'Resident population');
INSERT INTO key_stat_types (type_id, name) VALUES (nextval('key_stat_type_id'), 'Population density (Hectares)');
INSERT INTO key_stat_types (type_id, name) VALUES (nextval('key_stat_type_id'), 'Average (mean) age');
INSERT INTO key_stat_types (type_id, name) VALUES (nextval('key_stat_type_id'), 'People think their general health is good');
INSERT INTO key_stat_types (type_id, name) VALUES (nextval('key_stat_type_id'), 'Households where English is not the main language');
INSERT INTO key_stat_types (type_id, name) VALUES (nextval('key_stat_type_id'), 'Households owned with a mortgage, loan or shared ownership');


-- 
-- create the "key_stats" & "key_stats_history" tables and ID sequences.
-- 
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

CREATE SEQUENCE key_stat_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY key_stats.stat_id;

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

CREATE SEQUENCE key_stat_history_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY key_stats_history.stat_id;

-- Insert v1 key stats test data.

-- Resident population
INSERT INTO key_stats (stat_id, profile_id, stat_type, value, unit, date_created, dataset_id, dataset_name) VALUES (nextval('key_stat_id'), 1000, 1000, '1000', '', NOW(), 'abc123', 'Test dataset 1') ON CONFLICT ON CONSTRAINT key_stats_profile_id_stat_type_key DO UPDATE SET value = '1';
INSERT INTO key_stats_history (stat_id, profile_id, stat_type, value, unit, date_created, last_modified, dataset_id, dataset_name) VALUES (nextval('key_stat_history_id'), 1000, 1000, '1000', '', NOW(), NOW(), 'abc123', 'Test dataset 1');

-- Pop density
INSERT INTO key_stats (stat_id, profile_id, stat_type, value, unit, date_created, dataset_id, dataset_name) VALUES (nextval('key_stat_id'), 1000, 1100, '1100', '', NOW(), 'efg789', 'Test dataset 2') ON CONFLICT ON CONSTRAINT key_stats_profile_id_stat_type_key DO UPDATE SET value = '1100';
INSERT INTO key_stats_history (stat_id, profile_id, stat_type, value, unit, date_created, last_modified, dataset_id, dataset_name) VALUES (nextval('key_stat_history_id'), 1000, 1100, '1100', '', NOW(), NOW(), 'efg789', 'Test dataset 2');

-- Average age
INSERT INTO key_stats (stat_id, profile_id, stat_type, value, unit, date_created, dataset_id, dataset_name) VALUES (nextval('key_stat_id'), 1000, 1200, '1200', '', NOW(), 'abc123', 'Test dataset 1') ON CONFLICT ON CONSTRAINT key_stats_profile_id_stat_type_key DO UPDATE SET value = '1200';
INSERT INTO key_stats_history (stat_id, profile_id, stat_type, value, unit, date_created, last_modified, dataset_id, dataset_name) VALUES (nextval('key_stat_history_id'), 1000, 1200, '1200', '', NOW(), NOW(), 'abc123', 'Test dataset 1');

-- Good health
INSERT INTO key_stats (stat_id, profile_id, stat_type, value, unit, date_created, dataset_id, dataset_name) VALUES (nextval('key_stat_id'), 1000, 1300, '1300', '%', NOW(), 'efg789', 'Test dataset 2') ON CONFLICT ON CONSTRAINT key_stats_profile_id_stat_type_key DO UPDATE SET value = '1300';
INSERT INTO key_stats_history (stat_id, profile_id, stat_type, value, unit, date_created, last_modified, dataset_id, dataset_name) VALUES (nextval('key_stat_history_id'), 1000, 1300, '1300', '', NOW(), NOW(), 'efg789', 'Test dataset 2');

-- English is main language
INSERT INTO key_stats (stat_id, profile_id, stat_type, value, unit, date_created, dataset_id, dataset_name) VALUES (nextval('key_stat_id'), 1000, 1400, '1400', '%', NOW(), 'abc123', 'Test dataset 1') ON CONFLICT ON CONSTRAINT key_stats_profile_id_stat_type_key DO UPDATE SET value = '1400';
INSERT INTO key_stats_history (stat_id, profile_id, stat_type, value, unit, date_created, last_modified, dataset_id, dataset_name) VALUES (nextval('key_stat_history_id'), 1000, 1400, '1400', '', NOW(), NOW(), 'abc123', 'Test dataset 1');

-- Houses with a mortgage
INSERT INTO key_stats (stat_id, profile_id, stat_type, value, unit, date_created, dataset_id, dataset_name) VALUES (nextval('key_stat_id'), 1000, 1500, '1500', '%', NOW(), 'xxx666', 'Test dataset 3') ON CONFLICT ON CONSTRAINT key_stats_profile_id_stat_type_key DO UPDATE SET value = '1500';
INSERT INTO key_stats_history (stat_id, profile_id, stat_type, value, unit, date_created, last_modified, dataset_id, dataset_name) VALUES (nextval('key_stat_history_id'), 1000, 1500, '1500', '', NOW(), NOW(), 'xxx666', 'Test dataset 3');

-- 
-- Create the "geography_types" table, ID sequence and add some test data.
-- 
CREATE TABLE IF NOT EXISTS geography_types (
    id INT PRIMARY KEY NULL NULL,
    name VARCHAR(100) NOT NULL,
    UNIQUE (name)
);

CREATE SEQUENCE geography_type_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY geography_types.id;
INSERT INTO geography_types (id, name) VALUES (nextval('geography_type_id'), 'output area');
INSERT INTO geography_types (id, name) VALUES (nextval('geography_type_id'), 'higher level output area');
INSERT INTO geography_types (id, name) VALUES (nextval('geography_type_id'), 'lowever output area');

--
-- create the "key_stats_recipes" table and ID sequence.
-- 
CREATE TABLE IF NOT EXISTS key_stats_recipes (
    recipe_id INT PRIMARY KEY NOT NULL,
    dataset_id VARCHAR(100) NOT NULL,
    dataset_edition VARCHAR(100) NOT NULL,
    cantabular_query VARCHAR(100) NOT NULL,
    stat_type INT NOT NULL,
    CONSTRAINT fk_stat_type 
        FOREIGN KEY (stat_type) REFERENCES key_stat_types (type_id)

);

CREATE SEQUENCE recipe_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY key_stats_recipes.recipe_id;

--
-- Create the "recipe_geographies" junction table and ID sequence.
--
CREATE TABLE IF NOT EXISTS recipe_geographies (
    id INT PRIMARY KEY NOT NULL,
    recipe_id INT NOT NULL,
    geography_type_id INT NOT NULL,
    CONSTRAINT fk_recipe_id 
        FOREIGN KEY (recipe_id) REFERENCES key_stats_recipes (recipe_id),
    CONSTRAINT fk_geography_type_id
        FOREIGN KEY (geography_type_id) REFERENCES geography_types (id)
);

CREATE SEQUENCE recipe_geography_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY recipe_geographies.id;

--
-- Insert test recipes along with their geograhy types
--
INSERT INTO key_stats_recipes (recipe_id, dataset_id, dataset_edition, cantabular_query, stat_type) VALUES (nextval('recipe_id'),'Test dataset 1','abc123','<catabular query 1 template goes here>',(SELECT t.type_id FROM key_stat_types t WHERE t.name = 'Resident population'));
INSERT INTO recipe_geographies (id, recipe_id, geography_type_id) VALUES (nextval('recipe_geography_id'), 1000, 1000);
INSERT INTO recipe_geographies (id, recipe_id, geography_type_id) VALUES (nextval('recipe_geography_id'), 1000, 1100);
INSERT INTO recipe_geographies (id, recipe_id, geography_type_id) VALUES (nextval('recipe_geography_id'), 1000, 1200);

INSERT INTO key_stats_recipes (recipe_id, dataset_id, dataset_edition, cantabular_query, stat_type) VALUES (nextval('recipe_id'), 'Test dataset 1', 'abc123', '<catabular query 2 template goes here>', (SELECT t.type_id FROM key_stat_types t WHERE t.name = 'Population density (Hectares)'));
INSERT INTO recipe_geographies (id, recipe_id, geography_type_id) VALUES (nextval('recipe_geography_id'), 1100, 1000);
INSERT INTO recipe_geographies (id, recipe_id, geography_type_id) VALUES (nextval('recipe_geography_id'), 1100, 1100);

INSERT INTO key_stats_recipes (recipe_id, dataset_id, dataset_edition, cantabular_query, stat_type) VALUES (nextval('recipe_id'), 'Test dataset 2', 'efg789', '<catabular query 3 template goes here>', (SELECT t.type_id FROM key_stat_types t WHERE t.name = 'Average (mean) age'));
INSERT INTO recipe_geographies (id, recipe_id, geography_type_id) VALUES (nextval('recipe_geography_id'), 1200, 1200);