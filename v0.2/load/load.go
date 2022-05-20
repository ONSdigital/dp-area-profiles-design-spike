package load

import (
	"encoding/csv"
	"github.com/ONSdigital/dp-area-profiles-design-spike/v2/store"
	"github.com/pkg/errors"
	"io"
	"os"
	"time"
)

// Store represents the area profiles data store.
type Store interface {
	Init(areaCode, areaName, areaProfileName string) error
	GetProfileByAreaCode(areaCode string) (*store.AreaProfile, error)
	InsertKeyStat(areaCode, name, value, unit string, datasetID, datasetName string, dateCreated time.Time) (int, error)
	Close() error
}

// ImportRow is a Go representation of an area profiles key statistic in an import cvs row.
type RowData struct {
	AreaCode    string
	Title       string
	Name        string
	Value       string
	Unit        string
	DatasetID   string
	DatasetName string
}

// DataFromFile load test data into the postgres database from the specified file.
func DataFromFile(filename string, store Store) error {
	rows, err := readFile(filename)
	if err != nil {
		return err
	}

	created := time.Now()
	for _, r := range rows {
		if _, err := store.InsertKeyStat(r.AreaCode, r.Name, r.Value, r.Unit, r.DatasetID, r.DatasetName, created); err != nil {
			return err
		}
	}

	return nil
}

func readFile(filename string) ([]RowData, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	_, err = r.Read() // read/discard the header row
	if err != nil {
		return nil, errors.Wrap(err, "error reading input CSV file")
	}

	stats := make([]RowData, 0)
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "error reading input CSV file")
		}

		stats = append(stats, RowData{
			AreaCode:    row[0],
			Title:       row[1],
			Name:        row[2],
			Value:       row[3],
			Unit:        row[4],
			DatasetID:   row[5],
			DatasetName: row[6],
		})
	}

	return stats, nil
}
