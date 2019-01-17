package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/tealeg/xlsx"
)

func TestParsing(t *testing.T) {
	i := 0
	excelFileName := "1.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		os.Exit(1)
	}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			// for _, cell := range row.Cells {
			for _, cell := range row.Cells {
				text := cell.String()
				// fmt.Printf("%s\n", text)
				if len(text) > 0 {
					// fmt.Println(i)
					i++
				}

			}
		}
	}
	if i == 0 {
		t.Errorf("Expected i>0 %v", i)

	}
	fmt.Println(i)
}
