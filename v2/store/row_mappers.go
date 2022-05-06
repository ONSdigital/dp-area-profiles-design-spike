package store

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"time"
)

// areaProfilesRowsMapper maps a postgres results rows to a list of AreaProfile structs
func areaProfilesRowsMapper(rows pgx.Rows) ([]AreaProfile, error) {
	profiles := make([]AreaProfile, 0)

	for rows.Next() {
		p, err := mapRowsToAreaProfile(rows)
		if err != nil {
			return nil, err
		}

		profiles = append(profiles, p)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return profiles, nil
}

func mapRowsToAreaProfile(rows pgx.Rows) (AreaProfile, error) {
	profile := AreaProfile{}

	if err := rows.Scan(&profile.ID, &profile.Name, &profile.AreaCode); err != nil {
		return profile, err
	}

	profile.Href = fmt.Sprintf("http://localhost:8080/profiles/%s", profile.AreaCode)

	return profile, nil
}

// keyStatisticsRowsMapper maps postgres result rows to a list of KeyStatistics
func keyStatisticsRowsMapper(p *AreaProfile, rows pgx.Rows) (KeyStatistics, error) {
	stats := make(KeyStatistics, 0)

	for rows.Next() {
		s, err := mapRowsToKeyStats(p, rows)
		if err != nil {
			return nil, err
		}

		stats = append(stats, s)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return stats, nil
}

func mapRowsToKeyStats(p *AreaProfile, rows pgx.Rows) (KeyStatistic, error) {
	s := KeyStatistic{AreaCode: p.AreaCode}

	if err := rows.Scan(&s.ProfileID, &s.StatID, &s.StatType, &s.Name, &s.Value, &s.Unit, &s.DateCreated, &s.Metadata.DatasetID, &s.Metadata.DatasetName); err != nil {
		return s, err
	}

	s.Metadata.Href = fmt.Sprintf("http://localhost:8080/datasets/%s", s.Metadata.DatasetID)

	return s, nil
}

func mapRowToKeyStat(row pgx.Row) (KeyStatistic, error) {
	s := KeyStatistic{}

	if err := row.Scan(&s.ProfileID, &s.StatID, &s.StatType, &s.Name, &s.Value, &s.Unit, &s.DateCreated, &s.Metadata.DatasetID, &s.Metadata.DatasetName); err != nil {
		return s, err
	}

	s.Metadata.Href = fmt.Sprintf("http://localhost:8080/datasets/%s", s.Metadata.DatasetID)

	return s, nil
}

// versionsRowsMapper maps pgx.Rows results to an array of time.Time
func versionsRowsMapper(rows pgx.Rows) ([]time.Time, error) {
	versions := make([]time.Time, 0)

	for rows.Next() {
		var dateCreated time.Time

		if err := rows.Scan(&dateCreated); err != nil {
			return nil, errors.Wrap(err, "error scanning date created into Time struct")
		}

		versions = append(versions, dateCreated)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return versions, nil
}
