package xls
""" Модуль для преобразования файлов из xls в xlsx """

import (
	"log"
	"os"

	"github.com/extrame/xls"
	"github.com/tealeg/xlsx"
)

var cell *xlsx.Cell // объявление переменной типа excel ячейки

""" функция Convert для преобразования xls файла в xlsx """
func Convert(file string) error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	targetPath := pwd + `\` + file + `x` // путь к файлу xlsx
	xlsxFile := getXlsxFile(targetPath)
	xlsxSheet := xlsxFile.Sheets[0]       // запись в первый лист xlsx файла
	xlsPath := pwd + `\` + file           // путь к xls файлу
	xlsFile, err := xls.Open(xlsPath, "") // открытие xls файла
	if err != nil {
		return err
	}
	sheet := xlsFile.GetSheet(0) // чтение из первого листа xls файла
	for j := 0; j < int(sheet.MaxRow)+1; j++ {
		xlsRow := sheet.Row(j)
		rowColCount := xlsRow.LastCol()
		insertRowFromXls(xlsxSheet, xlsRow, rowColCount) // создание строк в xlsx
	}
	xlsxFile.Save(targetPath) // сохранение xlsx файла
	return nil
}

// добавление строк в xlsx файл
func insertRowFromXls(sheet *xlsx.Sheet, rowDataPtr *xls.Row, rowColCount int) {
	row := sheet.AddRow()
	for i := 0; i < rowColCount; i++ {
		cell = row.AddCell()
		cell.Value = rowDataPtr.Col(i)
	}
}

// создание xlsx файла
func getXlsxFile(filePath string) *xlsx.File {
	file := xlsx.NewFile()
	_, err := file.AddSheet("Sheet1")
	if err != nil {
		log.Fatal(err)
	}
	return file
}
