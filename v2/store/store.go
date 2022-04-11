package store

import (
	"context"
	"fmt"
	log "github.com/daiLlew/funkylog"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is an error to represent the state where the requested record does not exist.
	ErrNotFound = errors.New("no rows exist matching your query parameters")
)

// Drop sequences/tables.
var (
	// dropSequencesSQL is an SQL statement to drop the sequences created by this demo.
	dropSequencesSQL = "DROP SEQUENCE IF EXISTS area_profile_id, key_stats_id, key_stats_history_id, key_stat_version_id"

	// dropTablesSQL is an SQL statement to drop all tables created by this demo.
	dropTablesSQL = "DROP TABLE IF EXISTS key_stats_history, key_stats, area_profiles, areas CASCADE;"
)

// Store represents the area profiles data store.
type Store interface {
	Init(areaCode, areaName, areaProfileName string) error
	GetAreaProfiles() ([]AreaProfile, error)
	GetProfileByAreaCode(areaCode string) (*AreaProfile, error)
	Close() error
}

type AreaProfileStore struct {
	conn *pgx.Conn
}

// New construct a new Area profile store.
func New(username, password, database string) (*AreaProfileStore, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", username, password, database))
	if err != nil {
		return nil, errors.Wrap(err, "error opening postgres connection")
	}

	log.Info("successfully opened connection to database %q", database)
	return &AreaProfileStore{conn: conn}, nil
}

// Init is an initialisation function. If dropSchema is true any existing tables, data and sequences will be dropped and recreated. If false no action is taken.
func (s *AreaProfileStore) Init(areaCode, areaName, areaProfileName string) error {
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
	_, err := s.AddAreaProfile(areaCode, areaProfileName)
	if err != nil {
		return err
	}

	log.Info("database initialisation compeleted successfully :pizza:")
	return nil
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
func (s *AreaProfileStore) Close() error {
	return s.conn.Close(context.Background())
}
