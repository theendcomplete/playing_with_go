package main

import (
	"fmt"
	"os"

	"github.com/tealeg/xlsx"
)

func main() {
	excelFileName := "1.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		os.Exit(1)
	}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				text := cell.String()
				fmt.Printf("%s\n", text)
			}
		}
	}
}
