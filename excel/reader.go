/**
 * @Author: jie.an
 * @Description:
 * @File:  reader.go
 * @Version: 1.0.0
 * @Date: 2019/11/15 17:00
 */
package excel

import (
	"golang-base/tools"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

//
//ReadToMaps read excel file line by line and return data with map in a list
//
//Excel:
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

func ReadToMaps(excelFile, sheetName string) (rowMaps []map[string]string) {
	// open file
	tools.InfoLogger.Printf("Start Processing File: %s ,Sheet: %s", excelFile, sheetName)
	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		tools.ErrorLogger.Fatalln(err)
	}
	// headline type
	headline := make(map[int]string)
	tmpRows, tmpErr := f.Rows(sheetName)
	if tmpErr != nil {
		tools.ErrorLogger.Fatalln(tmpErr)
	}
	tmpRows.Next()
	headRow, headErr := tmpRows.Columns()
	if headErr != nil {
		tools.ErrorLogger.Fatalln(headErr)
	}
	for k, v := range headRow {
		headline[k] = v
	}
	tools.InfoLogger.Println("HeadLine:", headline)
	// iter all rows
	//var rowMaps []map[string]string
	rows, err := f.GetRows(sheetName)
	if err != nil {
		tools.ErrorLogger.Fatalln(err)
	}
	for _, row := range rows {
		rowMap := make(map[string]string)
		for i, v := range row {
			//skip head line (title)
			if row[i] == headline[i] {
				continue
			} else {
				rowMap[headline[i]] = v
			}
		}
		if len(rowMap) != 0 {
			rowMaps = append(rowMaps, rowMap)
		}
	}
	// return rowdata in map with out headeline (title)
	return rowMaps
}
