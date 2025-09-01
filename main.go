package main

import (
	"fmt"

	"github.com/odeta939/go-data-processing/internal"
	"github.com/odeta939/go-data-processing/io"
	"github.com/odeta939/go-data-processing/model"
	"github.com/shopspring/decimal"
)

const (
	time_ref = iota
	account
	code
	country_code
	product_type
	value
	status
)

func main() {

	r := io.ReadCSV(model.OutputCSVFull)

	records := make([]model.Record, 0)
	for _, row := range r[1:] {

		isYearValid := row[time_ref][:4] == model.YearToReport
		valueValid := row[value] != ""
		codeValid := len(row[code]) == 4
		isTypeGoods := row[product_type] == "Goods"

		if !isYearValid || !valueValid || !codeValid || !isTypeGoods {
			continue
		}

		valString := row[value]

		val, err := decimal.NewFromString(valString)
		if err != nil {
			fmt.Println("parsing value to float", err)
		}
		record := model.Record{
			TimeRef:     row[time_ref],
			Account:     row[account],
			Code:        row[code],
			CountryCode: row[country_code],
			ProductType: row[product_type],
			Value:       val,
			Status:      row[status],
		}

		records = append(records, record)
	}

	reports := []model.Report{}
	reports = internal.ReportGroup(records, []string{"Norway"})
	io.ReportToCSV(reports, "output/Norway_report.csv")

	reports = internal.ReportGroup(records, model.EUCountries)
	io.ReportToCSV(reports, "output/EU_reports.csv")

	reports = internal.ReportGroup(records, []string{"EU"})
	io.ReportToCSV(reports, "output/EU_report_total.csv")

}
