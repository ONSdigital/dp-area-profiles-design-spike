package store

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-area-profiles-design-spike/models"
	"github.com/ONSdigital/dp-area-profiles-design-spike/testdata"
	log "github.com/daiLlew/funkylog"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"time"
)

var (
	// ErrNotFound is an error to represent the state where the requested record does not exist.
	ErrNotFound = errors.New("no rows exist matching your query parameters")
)

// Store represents the area profiles data store.
type Store interface {
	Init(areaCode, areaName, areaProfileName string) error
	Close() error
	GetProfiles() ([]*models.AreaProfileLink, error)
	GetProfileByAreaCode(areaCode string) (*models.AreaProfile, error)
	GetKeyStatsByProfileID(profileID int) ([]models.KeyStatistic, error)
	UpdateProfileKeyStats(profileID int, newStats []models.ImportRow) error
	GetKeyStatsVersions(areaCode string, profileID int) ([]models.KeyStatsVersion, error)
	GetKeyStatsVersion(profileID, versionID int) ([]models.KeyStatistic, error)
}

type areaProfileStore struct {
	conn *pgx.Conn
}

// New construct a new Area profile store.
func New(username, password, database string) (Store, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", username, password, database))
	if err != nil {
		return nil, errors.Wrap(err, "error opening postgres connection")
	}

	log.Info("successfully opened connection to database %q", database)
	return &areaProfileStore{conn: conn}, nil
}

// Init is an initialisation function. If dropSchema is true any existing tables, data and sequences will be dropped and recreated. If false no action is taken.
func (s *areaProfileStore) Init(areaCode, areaName, areaProfileName string) error {
	stmts := []string{
		dropSequencesSQL,
		dropTablesSQL,
	}

	log.Info("dropping database schema")
	if err := execStmts(context.Background(), s.conn, stmts...); err != nil {
		return err
	}

	stmts = []string{
		createAreasTableSQL,
		createProfilesTableSQL,
		createAreaProfileIDSeqSQL,
		createKeyStatsTableSQL,
		createKeyStatsIDSeqSQL,
		createKeyStatsHistoryTableSQL,
		createKeyStatsHistoryIDSeqSQL,
		createKeyStatVersionIDSeqSQL,
	}

	log.Info("recreating database schema")
	if err := execStmts(context.Background(), s.conn, stmts...); err != nil {
		return err
	}

	log.Info("adding area test data, name=%s, code=%s", areaName, areaCode)
	if _, err := s.AddArea(areaCode, areaName); err != nil {
		return err
	}

	log.Info("adding area profile test data, name=%s", areaProfileName)
	profileID, err := s.AddAreaProfile(areaCode, areaProfileName)
	if err != nil {
		return err
	}

	stats, err := testdata.ReadCVS("testdata/ex1.csv")
	if err != nil {
		return err
	}

	log.Info("adding key stat test data for area profile, id=%d", profileID)
	now := time.Now()
	for _, stat := range stats {
		if err := s.AddKeyStat(profileID, stat.Name, stat.Value, stat.Unit, now); err != nil {
			return err
		}
	}
	log.Info("database initialisation compeleted successfully :pizza:")
	return nil
}

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

// NewArea insert a new area, returns the area code.
func (s *areaProfileStore) AddArea(code, name string) (string, error) {
	var areaCode string
	err := s.conn.QueryRow(context.Background(), insertAreaSQL, code, name).Scan(&areaCode)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return areaCode, nil
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

// NewKeyStat insert a key statistic for the specified area profile.
func (s *areaProfileStore) AddKeyStat(profileID int, name, value, unit string, dateCreated time.Time) error {
	_, err := s.conn.Exec(context.Background(), insertNewKeyStatSQL, profileID, name, value, unit, dateCreated)
	if err != nil {
		return errors.Wrapf(err, "error inserting new key stat %q for profile_id=%d", name, profileID)
	}
	return nil
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

// GetKeyStatsForProfileID return a list of key stats for the specified area profile ID.
func (s *areaProfileStore) GetKeyStatsByProfileID(profileID int) ([]models.KeyStatistic, error) {
	stats := make([]models.KeyStatistic, 0)

	rows, err := s.conn.Query(context.Background(), getKeyStatsForProfileIDSQL, profileID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			statID      int
			profileID   int
			name        string
			value       string
			unit        string
			dateCreated time.Time
		)

		if err := rows.Scan(&statID, &profileID, &name, &value, &unit, &dateCreated); err != nil {
			return nil, err
		}

		stats = append(stats, models.KeyStatistic{
			StatID:      statID,
			ProfileID:   profileID,
			Name:        name,
			Value:       value,
			Unit:        unit,
			DateCreated: dateCreated,
		})
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return stats, nil
}

// GetKeyStatsVersions returns a list of KeyStatsVersion available for the specified area profile ID.
func (s *areaProfileStore) GetKeyStatsVersions(areaCode string, profileID int) ([]models.KeyStatsVersion, error) {
	rows, err := s.conn.Query(context.Background(), getKeyStatVersionsSQL, profileID)
	if err != nil {
		return nil, errors.Wrap(err, "error querying for versions")
	}

	defer rows.Close()

	versions := make([]models.KeyStatsVersion, 0)
	for rows.Next() {
		var (
			dateCreated time.Time
			statID      int
			profileID   int
			versionID   int
		)

		if err := rows.Scan(&dateCreated, &statID, &profileID, &versionID); err != nil {
			return nil, err
		}

		versions = append(versions, models.KeyStatsVersion{
			StatID:      statID,
			ProfileID:   profileID,
			VersionID:   versionID,
			DateCreated: dateCreated,
			Href:        fmt.Sprintf("http://localhost:8080/profile/%s/versions/%d", areaCode, versionID),
		})
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return versions, nil
}

// GetKeyStatsVersion return a list of KeyStatistics associated with the specified version of the area profile ID.
func (s *areaProfileStore) GetKeyStatsVersion(profileID, versionID int) ([]models.KeyStatistic, error) {
	rows, err := s.conn.Query(context.Background(), getKeyStatVersionSQL, profileID, versionID)
	if err != nil {
		return nil, errors.Wrap(err, "error querying for versions")
	}

	defer rows.Close()

	stats := make([]models.KeyStatistic, 0)
	for rows.Next() {
		var (
			versionID    int
			statID       int
			profileID    int
			name         string
			value        string
			unit         string
			dateCreated  time.Time
			lastModified time.Time
		)

		if err := rows.Scan(&versionID, &statID, &profileID, &name, &value, &unit, &dateCreated, &lastModified); err != nil {
			return nil, err
		}

		stats = append(stats, models.KeyStatistic{
			VersionID:    versionID,
			StatID:       statID,
			ProfileID:    profileID,
			Name:         name,
			Value:        value,
			Unit:         unit,
			DateCreated:  dateCreated,
			LastModified: lastModified,
		})
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return stats, nil
}

// UpdateProfileKeyStats creates a new version of key statistics for the specified area profile.
// The current key stats are superceded and added to the key stats version history table before updating
// setting the new value for each of the key stats in the new current version.
func (s *areaProfileStore) UpdateProfileKeyStats(profileID int, newStats []models.ImportRow) error {
	currentKeyStats, err := s.GetKeyStatsByProfileID(profileID)
	if err != nil {
		return err
	}

	if err := s.versionProfileKeyStats(profileID, currentKeyStats); err != nil {
		return nil
	}

	if err := s.addNewProfileKeyStats(profileID, newStats); err != nil {
		return nil
	}

	return nil
}

// versionProfileKeyStats moves the current area profile key stats into the version history table
func (s *areaProfileStore) versionProfileKeyStats(profileID int, currentKeyStats []models.KeyStatistic) error {
	log.Info("adding current key stats to version history table, profile ID: %s", profileID)
	versionID, err := s.getNextVersionID()
	if err != nil {
		return err
	}

	lastUpdated := time.Now()
	for _, stat := range currentKeyStats {
		_, err := s.conn.Exec(context.Background(), insertNewKeyStatHistorySQL, profileID, versionID, stat.Name, stat.Value, stat.Unit, stat.DateCreated, lastUpdated)
		if err != nil {
			return errors.Wrapf(err, "error inserting new key stat version %q for profile_id=%d", stat.Name, profileID)
		}
	}

	return nil
}

func (s *areaProfileStore) addNewProfileKeyStats(profileID int, newStats []models.ImportRow) error {
	log.Info("inserting latest key stats for for profile ID: %s", profileID)
	now := time.Now()

	for _, stat := range newStats {
		_, err := s.conn.Exec(context.Background(), updateKeyStatSQL, stat.Value, stat.Unit, now, profileID, stat.Name)
		if err != nil {
			return errors.Wrapf(err, "error updating key stat %q for profile_id=%d", stat.Name, profileID)
		}
	}

	return nil
}

func (s *areaProfileStore) getNextVersionID() (int, error) {
	var id int
	if err := s.conn.QueryRow(context.Background(), getKeyStatVersionIDSQL).Scan(&id); err != nil {
		return 0, errors.Wrap(err, "error getting next key stat version id")
	}

	return id, nil
}

func execStmts(ctx context.Context, conn *pgx.Conn, statements ...string) error {
	for _, stmt := range statements {

		_, err := conn.Exec(ctx, stmt)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error executing sql statement: %q", stmt))
		}
	}
	return nil
}

// Close closes the underlying postgres connection
func (s *areaProfileStore) Close() error {
	log.Warn("closing area profile store connection")
	return s.conn.Close(context.Background())
}
