package io

import (
	"encoding/csv"
	"os"

	"github.com/odeta939/go-data-processing/model"
)

func ReportToCSV(reports []model.Report, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{
		"country",
		"trade_balance",
		"most_imported_product",
		"import_value",
		"most_exported_product",
		"exported_value",
	})

	for _, r := range reports {
		row := []string{
			r.Country,
			r.TradeBalance.StringFixed(2),
			r.MostImported.Name,
			r.MostImported.Value.StringFixed(2),
			r.MostExported.Name,
			r.MostExported.Value.StringFixed(2),
		}
		writer.Write(row)
	}

	return nil
}
