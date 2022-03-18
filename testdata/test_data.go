package testdata

import (
	"encoding/csv"
	"github.com/ONSdigital/dp-area-profiles-design-spike/models"
	"github.com/pkg/errors"
	"io"
	"os"
)

// ReadCVS is a helper func to read an import csv file. Returns the data a list of ImportRow.
func ReadCVS(file string) ([]models.ImportRow, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	f, err = os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	_, err = r.Read() // header
	if err != nil {
		return nil, errors.Wrap(err, "error reading input CSV file")
	}

	stats := make([]models.ImportRow, 0)
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "error reading input CSV file")
		}

		stats = append(stats, models.ImportRow{
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
