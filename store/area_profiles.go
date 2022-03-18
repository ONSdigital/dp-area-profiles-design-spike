package store

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-area-profiles-design-spike/models"
	"github.com/jackc/pgx/v4"
)

var (
	// createProfilesTableSQL is an SQL statement to create the area profiles table
	createProfilesTableSQL = "CREATE TABLE area_profiles (profile_id INT PRIMARY KEY NOT NULL, area_code VARCHAR(50) NOT NULL, name VARCHAR (100) NOT NULL, UNIQUE (area_code), CONSTRAINT fk_area_code FOREIGN KEY (area_code) REFERENCES areas (code));"

	//createAreaProfileIDSeqSQL is a SQL statement creating a sequence for generating area profile ids.
	createAreaProfileIDSeqSQL = "CREATE SEQUENCE area_profile_id START 1000 INCREMENT 100 MINVALUE 1000 OWNED BY area_profiles.profile_id;"

	// insertProfileSQL is an SQL query to insert a new area profile, required area code and profile name.
	insertProfileSQL = "INSERT INTO area_profiles (profile_id, area_code, name) VALUES (nextval('area_profile_id'), $1, $2) RETURNING profile_id;"

	// getProfilesSQL SQL query returning all area profiles.
	getProfilesSQL = "SELECT profile_id, area_code, name FROM area_profiles"

	// getProfileByAreaCodeSQL SQL query returns the area profile for the specified area code.
	getProfileByAreaCodeSQL = "SELECT profile_id, name, area_code FROM area_profiles WHERE area_code = $1;"
)

// GetProfiles return an array of existing area profiles.
func (s *areaProfileStore) GetProfiles() ([]*models.AreaProfileLink, error) {
	profiles := make([]*models.AreaProfileLink, 0)

	rows, err := s.conn.Query(context.Background(), getProfilesSQL)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			profileID int
			areaCode  string
			name      string
		)

		if err := rows.Scan(&profileID, &areaCode, &name); err != nil {
			return nil, err
		}

		profiles = append(profiles, &models.AreaProfileLink{
			ProfileID: profileID,
			AreaCode:  areaCode,
			Name:      name,
			Href:      fmt.Sprintf("http://localhost:8080/profiles/%s", areaCode),
		})
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return profiles, nil
}

// GetProfileIDByAreaCode return the area profile ID associated with the specified area code.
func (s *areaProfileStore) GetProfileByAreaCode(areaCode string) (*models.AreaProfile, error) {
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
	return &models.AreaProfile{
		ID:       profileID,
		Name:     name,
		AreaCode: code,
	}, nil
}

// NewAreaProfile insert a new area profile returns the area profile ID.
func (s *areaProfileStore) AddAreaProfile(areaCode, name string) (int, error) {
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
