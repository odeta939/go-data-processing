package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/shopspring/decimal"
)

var EUCountries = []string{
	"Austria",
	"Belgium",
	"Bulgaria",
	"Croatia",
	"Cyprus",
	"Czechia",
	"Denmark",
	"Estonia",
	"Finland",
	"France",
	"Germany",
	"Greece",
	"Hungary",
	"Ireland",
	"Italy",
	"Latvia",
	"Lithuania",
	"Luxembourg",
	"Malta",
	"Netherlands",
	"Poland",
	"Portugal",
	"Romania",
	"Slovakia",
	"Slovenia",
	"Spain",
	"Sweden",
}

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

type Record struct {
	TimeRef     string
	Account     string
	Code        string
	CountryCode string
	ProductType string
	Value       decimal.Decimal
	Status      string
}

type Product struct {
	Name  string
	Value decimal.Decimal
}

type Report struct {
	County       string
	TradeBalance decimal.Decimal
	MostImported Product
	MostExported Product
}

func main() {

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
	reportNorway := ReportNorway(records)
	reportEU := ReportEU(records)

	reports := make([]Report, 0)

	for _, c := range EUCountries {
		r := ReportByCountry(records, c)
		reports = append(reports, r)
	}
	ReportToCSV(reports, "EU_reports.csv")
	reports = []Report{reportEU}
	ReportToCSV(reports, "EU_report_total.csv")
	reports = []Report{reportNorway}
	ReportToCSV(reports, "Norway_report.csv")

	fmt.Println(reports)
}

func ReportNorway(records []Record) Report {
	importTotal := decimal.NewFromInt(0)
	exportTotal := decimal.NewFromInt(0)
	importGoods := make(map[string]decimal.Decimal)
	exportGoods := make(map[string]decimal.Decimal)

	cCode, err := CountryCode("Norway")
	if err != nil {
		fmt.Println("failed getting country code", err)
	}
	for _, r := range records {
		if r.CountryCode != cCode {
			continue
		}
		if r.Account == "Imports" {
			importTotal = importTotal.Add(r.Value)
			importGoods[r.Code] = importGoods[r.Code].Add(r.Value)
		}
		if r.Account == "Exports" {
			exportTotal = exportTotal.Add(r.Value)
			exportGoods[r.Code] = exportGoods[r.Code].Add(r.Value)
		}
	}

	topImportVal := decimal.NewFromInt(0)
	topImportCode := ""
	for code, sum := range importGoods {
		if sum.GreaterThan(topImportVal) {
			topImportVal = sum
			topImportCode = code
		}
	}

	topExportVal := decimal.NewFromInt(0)
	topExportCode := ""
	for code, sum := range exportGoods {

		if sum.GreaterThan(topImportVal) {
			topExportVal.Add(sum)
			topExportCode = code
		}
	}

	tradeBalance := exportTotal.Sub(importTotal)
	name, err := ProductName(topImportCode)
	if err != nil {
		fmt.Println("failed getting product name", name, err)
		panic("failed getting product name, ")
	}
	importTop := Product{
		Name:  name,
		Value: topImportVal,
	}

	name, err = ProductName(topExportCode)
	if err != nil {
		fmt.Println("failed getting product name", err)
	}
	exportTop := Product{
		Name:  name,
		Value: topImportVal,
	}

	report := Report{
		County:       "Norway",
		TradeBalance: tradeBalance,
		MostImported: importTop,
		MostExported: exportTop,
	}
	return report
}

func ReportByCountry(records []Record, country string) Report {
	importTotal := decimal.NewFromInt(0)
	exportTotal := decimal.NewFromInt(0)
	importGoods := make(map[string]decimal.Decimal)
	exportGoods := make(map[string]decimal.Decimal)

	cCode, err := CountryCode(country)
	if err != nil {
		fmt.Println("failed getting country code", err)
	}
	for _, r := range records {
		if r.CountryCode != cCode {
			continue
		}
		if r.Account == "Imports" {
			importTotal = importTotal.Add(r.Value)
			importGoods[r.Code] = importGoods[r.Code].Add(r.Value)
		}
		if r.Account == "Exports" {
			exportTotal = exportTotal.Add(r.Value)
			exportGoods[r.Code] = exportGoods[r.Code].Add(r.Value)
		}
	}

	topImportVal := decimal.NewFromInt(0)
	topImportCode := ""
	for code, sum := range importGoods {
		if sum.GreaterThan(topImportVal) {
			topImportVal = sum
			topImportCode = code
		}
	}

	topExportVal := decimal.NewFromInt(0)
	topExportCode := ""
	for code, sum := range exportGoods {

		if sum.GreaterThan(topImportVal) {
			topExportVal.Add(sum)
			topExportCode = code
		}
	}

	tradeBalance := exportTotal.Sub(importTotal)
	name, err := ProductName(topImportCode)
	if err != nil {
		fmt.Println("failed getting product name", name, err)
		panic("failed getting product name, ")
	}
	importTop := Product{
		Name:  name,
		Value: topImportVal,
	}

	name, err = ProductName(topExportCode)
	if err != nil {
		fmt.Println("failed getting product name", err)
	}
	exportTop := Product{
		Name:  name,
		Value: topImportVal,
	}

	report := Report{
		County:       country,
		TradeBalance: tradeBalance,
		MostImported: importTop,
		MostExported: exportTop,
	}
	return report
}

func ReportEU(records []Record) Report {
	importTotal := decimal.NewFromInt(0)
	exportTotal := decimal.NewFromInt(0)
	importGoods := make(map[string]decimal.Decimal)
	exportGoods := make(map[string]decimal.Decimal)

	cCodes := CountryCodes(EUCountries)

	for _, r := range records {
		if contains(cCodes, r.CountryCode) {
			continue
		}
		if r.Account == "Imports" {
			importTotal = importTotal.Add(r.Value)
			importGoods[r.Code] = importGoods[r.Code].Add(r.Value)
		}
		if r.Account == "Exports" {
			exportTotal = exportTotal.Add(r.Value)
			exportGoods[r.Code] = exportGoods[r.Code].Add(r.Value)
		}
	}

	topImportVal := decimal.NewFromInt(0)
	topImportCode := ""
	for code, sum := range importGoods {
		if sum.GreaterThan(topImportVal) {
			topImportVal = sum
			topImportCode = code
		}
	}

	topExportVal := decimal.NewFromInt(0)
	topExportCode := ""
	for code, sum := range exportGoods {

		if sum.GreaterThan(topImportVal) {
			topExportVal.Add(sum)
			topExportCode = code
		}
	}

	tradeBalance := exportTotal.Sub(importTotal)
	name, err := ProductName(topImportCode)
	if err != nil {
		fmt.Println("failed getting product name", name, err)
		panic("failed getting product name, ")
	}
	importTop := Product{
		Name:  name,
		Value: topImportVal,
	}

	name, err = ProductName(topExportCode)
	if err != nil {
		fmt.Println("failed getting product name", err)
	}
	exportTop := Product{
		Name:  name,
		Value: topImportVal,
	}

	report := Report{
		County:       "EU",
		TradeBalance: tradeBalance,
		MostImported: importTop,
		MostExported: exportTop,
	}
	return report
}

func ProductName(code string) (string, error) {
	file, err := os.Open(GoodsClassification)
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

func CountryCode(country string) (string, error) {
	file, err := os.Open(CountryClassification)
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(file)

	r, err := reader.ReadAll()

	for _, row := range r[1:] {
		if row[1] == country {
			return row[0], nil
		}
	}
	e := fmt.Errorf("country specified not found: %s", country)
	return "", e
}

func CountryCodes(countries []string) []string {
	file, err := os.Open(CountryClassification)
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(file)

	r, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	cCodes := make([]string, len(countries))
	for _, row := range r[1:] {
		for _, c := range countries {
			if row[1] == c {
				cCodes = append(cCodes, row[0])
			}
		}
	}
	return cCodes
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func ReportToCSV(reports []Report, filename string) error {
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
			r.County,
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
