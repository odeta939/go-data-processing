package io

import (
	"encoding/csv"
	"fmt"
	"os"
)

func ReadCSV(path string) [][]string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("while opening file ", path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	r, err := reader.ReadAll()
	if err != nil {
		fmt.Println("while reading records", err)
	}
	return r
}
