package xls

import (
	"log"
	"os"
	"time"

	"github.com/extrame/xls"
	"github.com/tealeg/xlsx"
)

var cell *xlsx.Cell

func Convert(file string) error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	targetPath := pwd + `\` + file + `x`
	xlsxFile := getXlsxFile(targetPath)
	xlsxSheet := xlsxFile.Sheets[0]
	xlsPath := pwd + `\` + file
	xlsFile, err := xls.Open(xlsPath, "")
	for err != nil {
		time.Sleep(1)
		xlsFile, err = xls.Open(xlsPath, "")
	}
	if err != nil {
		return err
	}
	sheet := xlsFile.GetSheet(0)
	for j := 0; j < int(sheet.MaxRow)+1; j++ {
		xlsRow := sheet.Row(j)
		rowColCount := xlsRow.LastCol()
		insertRowFromXls(xlsxSheet, xlsRow, rowColCount)
	}
	xlsxFile.Save(targetPath)
	return nil
}

func insertRowFromXls(sheet *xlsx.Sheet, rowDataPtr *xls.Row, rowColCount int) {
	row := sheet.AddRow()
	for i := 0; i < rowColCount; i++ {
		cell = row.AddCell()
		cell.Value = rowDataPtr.Col(i)
	}
}

func getXlsxFile(filePath string) *xlsx.File {
	file := xlsx.NewFile()
	_, err := file.AddSheet("Sheet1")
	if err != nil {
		log.Fatal(err)
	}
	return file
}
