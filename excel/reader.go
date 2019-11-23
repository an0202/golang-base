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

func ReadTest(excelfile, sheetname string) (rowmaps []map[string]string) {
	// open file
	tools.InfoLogger.Printf("Start Processing File: %s ,Sheet: %s", excelfile, sheetname)
	f, err := excelize.OpenFile(excelfile)
	if err != nil {
		tools.ErrorLogger.Fatalln(err)
	}
	// headline type
	headline := make(map[int]string)
	tmprows, tmperr := f.Rows(sheetname)
	if tmperr != nil {
		tools.ErrorLogger.Fatalln(tmperr)
	}
	headrow, headerr := tmprows.Columns()
	if headerr != nil {
		tools.ErrorLogger.Fatalln(headerr)
	}
	for k, v := range headrow {
		headline[k] = v
	}
	tools.InfoLogger.Println("HeadLine:",headline)
	// iter all rows
	//var rowmaps []map[string]string
	rows, err := f.GetRows(sheetname)
	if err != nil {
		tools.ErrorLogger.Fatalln(err)
	}
	for _, row := range rows {
		rowmap := make(map[string]string)
		for k, v := range row {
			if row[k] == headline[k] {
				continue
			} else {
				rowmap[headline[k]] = v
			}
		}
		if len(rowmap) != 0 {
			rowmaps = append(rowmaps, rowmap)
		}
	}
	return rowmaps
}
