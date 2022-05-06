package store

import (
	"context"
	"github.com/pkg/errors"
)

var (
	// createKeyStatTypeSQL SQL statement to create the key stat type table.
	createKeyStatTypeSQL = `
		CREATE TABLE IF NOT EXISTS key_stat_types (
			type_id INT PRIMARY KEY NOT NULL,
			name VARCHAR(100) NOT NULL,
			UNIQUE (name)
		);
	`

	//createKeyStatTypeSeqSQL SQL statement to create the key stat id sequence.
	createKeyStatTypeSeqSQL = `
		CREATE SEQUENCE 
			key_stat_type_id 
		START 
			1000 
		INCREMENT 
			100 
		MINVALUE 
			1000 
		OWNED BY 
			key_stat_types.type_id;
	`

	// insertKeyStatTypeSQL SQL statement to insert a new key stat type entry.
	insertKeyStatTypeSQL = `
		INSERT INTO key_stat_types 
			(type_id, name) 
		VALUES 
			(nextval('key_stat_type_id'), $1);
	`

	//getStatTypeByName SQL query returns key stat type id for the type with the specified name.
	getStatTypeByName = `
		SELECT 
			t.type_id 
		FROM 
			key_stat_types t 
		WHERE 
			name = $1;
	`
)

// InsertKeyStatTypes create a new key stat type for each of the name values provided.
func (s *AreaProfileStore) InsertKeyStatTypes(names ...string) error {
	for _, name := range names {
		_, err := s.conn.Exec(context.Background(), insertKeyStatTypeSQL, name)
		if err != nil {
			return errors.Wrapf(err, "error inserting key_stat_type: %q", name)
		}
	}
	return nil
}

// GetStatTypeByName return the stat type if for the name with the specified name value.
func (s *AreaProfileStore) GetStatTypeByName(name string) (int, error) {
	var typeID int
	err := s.conn.QueryRow(context.Background(), getStatTypeByName, name).Scan(&typeID)
	if err != nil {
		return 0, errors.Wrapf(err, "error getting stat type for name %q", name)
	}
	return typeID, nil
}
