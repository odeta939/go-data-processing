package helpers

import (
	"encoding/csv"
	"os"

	"github.com/odeta939/go-data-processing/model"
)

func CountryCode(country string) (string, error) {
	file, err := os.Open(model.CountryClassification)
	if err != nil {
		return "", err
	}

	reader := csv.NewReader(file)

	r, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	cCode := ""
	for _, row := range r[1:] {
		if row[1] == country {
			cCode = row[0]
		}
	}
	return cCode, nil
}
