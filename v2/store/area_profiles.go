package store

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

var (
	// createProfilesTableSQL is an SQL statement to create the area profiles table
	createProfilesTableSQL = "CREATE TABLE area_profiles (profile_id INT PRIMARY KEY NOT NULL, area_code VARCHAR(50) NOT NULL, name VARCHAR (100) NOT NULL, UNIQUE (area_code), CONSTRAINT fk_area_code FOREIGN KEY (area_code) REFERENCES areas (code));"

	//createAreaProfileIDSeqSQL is a SQL statement creating a sequence for generating area profile ids.
	createAreaProfileIDSeqSQL = "CREATE SEQUENCE area_profile_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY area_profiles.profile_id;"

	// getProfileByAreaCodeSQL SQL query returns the area profile for the specified area code.
	getProfileByAreaCodeSQL = "SELECT profile_id, name, area_code FROM area_profiles WHERE area_code = $1;"

	// getAreaProfilesSQL SQL query returning a list of all area profiles.
	getAreaProfilesSQL = "SELECT profile_id, name, area_code FROM area_profiles;"

	// insertProfileSQL is an SQL query to insert a new area profile, required area code and profile name.
	insertProfileSQL = "INSERT INTO area_profiles (profile_id, area_code, name) VALUES (nextval('area_profile_id'), $1, $2) RETURNING profile_id;"
)

// NewAreaProfile insert a new area profile returns the area profile ID.
func (s *AreaProfileStore) AddAreaProfile(areaCode, name string) (int, error) {
	var profileID int
	err := s.conn.QueryRow(context.Background(), insertProfileSQL, areaCode, name).Scan(&profileID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return profileID, nil
}

// GetAreaProfiles return a list of area profiles
func (s *AreaProfileStore) GetAreaProfiles() ([]AreaProfile, error) {
	rows, err := s.conn.Query(context.Background(), getAreaProfilesSQL)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	profiles, err := areaProfilesRowsMapper(rows)
	if err != nil {
		return nil, errors.Wrap(err, "error scanning get area profiles result rows")
	}

	return profiles, nil
}

// GetProfileIDByAreaCode return the area profile ID associated with the specified area code.
func (s *AreaProfileStore) GetProfileByAreaCode(areaCode string) (*AreaProfile, error) {
	var profileID int
	var name string
	var code string

	err := s.conn.QueryRow(context.Background(), getProfileByAreaCodeSQL, areaCode).Scan(&profileID, &name, &code)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &AreaProfile{
		ID:       profileID,
		Name:     name,
		AreaCode: code,
		Href:     fmt.Sprintf("http://localhost:8080/profiles/%s/stats", code),
	}, nil
}
