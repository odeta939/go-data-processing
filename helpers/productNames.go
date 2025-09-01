package helpers

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/odeta939/go-data-processing/model"
)

func ProductName(code string) (string, error) {
	file, err := os.Open(model.GoodsClassification)
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(file)

	r, err := reader.ReadAll()
	for _, row := range r[1:] {
		if code == row[0] {
			return row[2], nil
		}
	}
	e := fmt.Errorf("code specified not found: %s", code)
	return "", e
}
