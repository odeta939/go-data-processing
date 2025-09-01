package internal

import (
	"fmt"

	"github.com/odeta939/go-data-processing/helpers"
	"github.com/odeta939/go-data-processing/model"
	"github.com/shopspring/decimal"
)

func ReportByCountry(records []model.Record, countryCode string, label string) model.Report {
	importTotal := decimal.NewFromInt(0)
	exportTotal := decimal.NewFromInt(0)
	importGoods := make(map[string]decimal.Decimal)
	exportGoods := make(map[string]decimal.Decimal)

	for _, r := range records {
		if r.CountryCode != countryCode {
			continue
		}
		switch r.Account {
		case "Imports":
			importTotal = importTotal.Add(r.Value)
			importGoods[r.Code] = importGoods[r.Code].Add(r.Value)
		case "Exports":
			exportTotal = exportTotal.Add(r.Value)
			exportGoods[r.Code] = exportGoods[r.Code].Add(r.Value)
		}
	}

	topImportVal := decimal.Zero
	topImportCode := ""
	for code, sum := range importGoods {
		if sum.GreaterThan(topImportVal) {
			topImportVal = sum
			topImportCode = code
		}
	}
	importName, err := helpers.ProductName(topImportCode)
	if err != nil {
		fmt.Println("failed while looking for product name", err)
	}

	topExportVal := decimal.Zero
	topExportCode := ""
	for code, sum := range exportGoods {
		if sum.GreaterThan(topExportVal) {
			topExportVal = sum
			topExportCode = code
		}
	}
	exportName, err := helpers.ProductName(topExportCode)
	if err != nil {
		fmt.Println("failed while looking for product name", err)
	}

	tradeBalance := exportTotal.Sub(importTotal)
	importTop := model.Product{
		Name:  importName,
		Value: topImportVal,
	}

	exportTop := model.Product{
		Name:  exportName,
		Value: topImportVal,
	}

	report := model.Report{
		Country:      label,
		TradeBalance: tradeBalance,
		MostImported: importTop,
		MostExported: exportTop,
	}

	return report
}

func ReportGroup(records []model.Record, countries []string) []model.Report {
	reports := []model.Report{}

	for _, cName := range countries {
		if cName == "EU" {
			euCodes := []string{}
			for _, euCountry := range model.EUCountries {
				code, err := helpers.CountryCode(euCountry)
				if err != nil {
					fmt.Println(err)
					continue
				}
				euCodes = append(euCodes, code)
			}

			report := ReportForMultipleCountries(records, euCodes, "EU")
			reports = append(reports, report)
			continue
		}

		cCode, err := helpers.CountryCode(cName)
		if err != nil {
			fmt.Println(err)
			continue
		}
		report := ReportByCountry(records, cCode, cName)
		reports = append(reports, report)
	}
	return reports
}

func ReportForMultipleCountries(records []model.Record, countryCodes []string, label string) model.Report {
	importTotal := decimal.Zero
	exportTotal := decimal.Zero
	importGoods := make(map[string]decimal.Decimal)
	exportGoods := make(map[string]decimal.Decimal)

	for _, r := range records {
		if !helpers.Contains(countryCodes, r.CountryCode) {
			continue
		}
		switch r.Account {
		case "Imports":
			importTotal = importTotal.Add(r.Value)
			importGoods[r.Code] = importGoods[r.Code].Add(r.Value)
		case "Exports":
			exportTotal = exportTotal.Add(r.Value)
			exportGoods[r.Code] = exportGoods[r.Code].Add(r.Value)
		}
	}

	topImportVal := decimal.Zero
	topImportCode := ""
	for code, sum := range importGoods {
		if sum.GreaterThan(topImportVal) {
			topImportVal = sum
			topImportCode = code
		}
	}
	importName, err := helpers.ProductName(topImportCode)
	if err != nil {
		fmt.Println("failed while looking for product name", err)
	}

	topExportVal := decimal.Zero
	topExportCode := ""
	for code, sum := range exportGoods {
		if sum.GreaterThan(topExportVal) {
			topExportVal = sum
			topExportCode = code
		}
	}
	exportName, err := helpers.ProductName(topExportCode)
	if err != nil {
		fmt.Println("failed while looking for product name", err)
	}
	tradeBalance := exportTotal.Sub(importTotal)
	importTop := model.Product{
		Name:  importName,
		Value: topImportVal,
	}

	exportTop := model.Product{
		Name:  exportName,
		Value: topImportVal,
	}

	report := model.Report{
		Country:      label,
		TradeBalance: tradeBalance,
		MostImported: importTop,
		MostExported: exportTop,
	}

	return report
}

//---------------------------//

// func ReportEU(records []model.Record) model.Report {
// 	importTotal := decimal.NewFromInt(0)
// 	exportTotal := decimal.NewFromInt(0)
// 	importGoods := make(map[string]decimal.Decimal)
// 	exportGoods := make(map[string]decimal.Decimal)

// 	cCodes, err := helpers.CountryCode(model.EUCountries)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	for _, r := range records {
// 		if !helpers.Contains(cCodes, r.CountryCode) {
// 			continue
// 		}
// 		if r.Account == "Imports" {
// 			importTotal = importTotal.Add(r.Value)
// 			importGoods[r.Code] = importGoods[r.Code].Add(r.Value)
// 		}
// 		if r.Account == "Exports" {
// 			exportTotal = exportTotal.Add(r.Value)
// 			exportGoods[r.Code] = exportGoods[r.Code].Add(r.Value)
// 		}
// 	}

// 	topImportVal := decimal.NewFromInt(0)
// 	topImportCode := ""
// 	for code, sum := range importGoods {
// 		if sum.GreaterThan(topImportVal) {
// 			topImportVal = sum
// 			topImportCode = code
// 		}
// 	}

// 	topExportVal := decimal.NewFromInt(0)
// 	topExportCode := ""
// 	for code, sum := range exportGoods {

// 		if sum.GreaterThan(topExportVal) {
// 			topExportVal.Add(sum)
// 			topExportCode = code
// 		}
// 	}

// 	tradeBalance := exportTotal.Sub(importTotal)
// 	name, err := helpers.ProductName(topImportCode)
// 	if err != nil {
// 		fmt.Println("failed getting product name", name, err)
// 		panic("failed getting product name, ")
// 	}
// 	importTop := model.Product{
// 		Name:  name,
// 		Value: topImportVal,
// 	}

// 	name, err = helpers.ProductName(topExportCode)
// 	if err != nil {
// 		fmt.Println("failed getting product name", err)
// 	}
// 	exportTop := model.Product{
// 		Name:  name,
// 		Value: topExportVal,
// 	}

// 	report := model.Report{
// 		Country:      "EU",
// 		TradeBalance: tradeBalance,
// 		MostImported: importTop,
// 		MostExported: exportTop,
// 	}
// 	return report
// }
