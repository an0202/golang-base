package csv

import (
	"encoding/csv"
	"fmt"
	"golang-base/tools"
	"io"
	"os"
)

// PrintTitle return the first line of csv file
func PrintTitle(inputFile string) {
	// Read CSV File
	inputCSV, inputError := os.OpenFile(inputFile, os.O_RDONLY, 0666)
	if inputError != nil {
		tools.ErrorLogger.Fatalln("Error While Open CSV File")
	}
	defer inputCSV.Close()
	// Init CSV Reader
	reader := csv.NewReader(inputCSV)
	record, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			tools.ErrorLogger.Fatalln("No Title For This CSV File")
		}
		tools.ErrorLogger.Fatalln(err)
	}
	for id, value := range record {
		tools.InfoLogger.Println("Title:")
		fmt.Println(id, value)
	}
}

//
//ReadToMaps read csv file line by line and return data with map in a list
//
//CSV:
//
//	A     B     C
//
//1  Name   Age   Sex
//
//2  Alice  18    Female
//
//3  Bob    22    Male
//
//ReturnData:
//[
//
//	{"Name":"Alice","Age":"18","Sex":"Female"},
//	{"Name":"Bob","Age":"22","Sex":"Male"}
//
//]
//
//todo: https://gist.github.com/drernie/5684f9def5bee832ebc50cabb46c377a
//func ReadToMaps(csvfile) (rowmaps []map[string]string) {
//	// Read CSV File
//	inputCSV, inputError := os.OpenFile(inputFile, os.O_RDONLY, 0666)
//	if inputError != nil {
//		tools.ErrorLogger.Fatalln("Error While Open CSV File")
//	}
//	defer inputCSV.Close()
//	// Init CSV Reader
//	reader := csv.NewReader(inputCSV)
//	f, err := excelize.OpenFile(excelfile)
//	if err != nil {
//		tools.ErrorLogger.Fatalln(err)
//	}
//	// headline type
//	headline := make(map[int]string)
//	tmprows, tmperr := f.Rows(sheetname)
//	if tmperr != nil {
//		tools.ErrorLogger.Fatalln(tmperr)
//	}
//	headrow, headerr := tmprows.Columns()
//	if headerr != nil {
//		tools.ErrorLogger.Fatalln(headerr)
//	}
//	for k, v := range headrow {
//		headline[k] = v
//	}
//	tools.InfoLogger.Println("HeadLine:", headline)
//	// iter all rows
//	//var rowmaps []map[string]string
//	rows, err := f.GetRows(sheetname)
//	if err != nil {
//		tools.ErrorLogger.Fatalln(err)
//	}
//	for _, row := range rows {
//		rowmap := make(map[string]string)
//		for k, v := range row {
//			//skip head line (title)
//			if row[k] == headline[k] {
//				continue
//			} else {
//				rowmap[headline[k]] = v
//			}
//		}
//		if len(rowmap) != 0 {
//			rowmaps = append(rowmaps, rowmap)
//		}
//	}
//	// return rowdata in map with out headeline (title)
//	return rowmaps
//}
