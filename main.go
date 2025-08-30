package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

const (
	//-------------- DATA FILES -------------//
	CountryClassification = "/home/odeta/Documents/projects/go-data-processing/data/country_classification.csv"
	GoodsClassification   = "/home/odeta/Documents/projects/go-data-processing/data/goods_classification.csv"
	OutputCSVFull         = "/home/odeta/Documents/projects/go-data-processing/data/output_csv_full.csv"
	Revised               = "/home/odeta/Documents/projects/go-data-processing/data/revised.csv"
	ServiceClassification = "/home/odeta/Documents/projects/go-data-processing/data/service_classification.csv"

	//--------------OTHER CONSTANTS-----------------//
	YearToReport = "2024"
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
	type Record struct {
		TimeRef     string
		Account     string
		Code        string
		CountryCode string
		ProductType string
		Value       float64
		Status      string
	}
	file, err := os.Open(OutputCSVFull)
	if err != nil {
		fmt.Println("while opening file ", OutputCSVFull, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	r, err := reader.ReadAll()
	if err != nil {
		fmt.Println("while reading records", err)
	}

	records := make([]Record, 0)
	for _, row := range r[1:] {

		isYearValid := row[time_ref][:4] == YearToReport
		codeValid := len(row[code]) == 4
		isTypeGoods := row[product_type] == "goods"
		valueValid := row[value] != ""

		if !isYearValid || !codeValid || !isTypeGoods || !valueValid {
			continue
		}

		valString := row[value]
		val, err := strconv.ParseFloat(valString, 64)
		if err != nil {
			fmt.Println("parsing value to float", err)
		}
		record := Record{
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
}
