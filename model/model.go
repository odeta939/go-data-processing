package model

import "github.com/shopspring/decimal"

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
	Country      string
	TradeBalance decimal.Decimal
	MostImported Product
	MostExported Product
}
