/**
 * @Author: jie.an
 * @Description:
 * @File:  reader.go
 * @Version: 1.0.0
 * @Date: 2019/11/15 17:00
 */
package excel

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func ReadTest(path, sheetname string) (rowmaps []map[string]string) {
	// open file
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
	}
	// headline type
	headline := make(map[int]string)
	tmprows, tmperr := f.Rows(sheetname)
	if tmperr != nil {
		fmt.Println(tmperr)
	}
	headrow, headerr := tmprows.Columns()
	if headerr != nil {
		fmt.Println(headerr)
	}
	for k, v := range headrow {
		headline[k] = v
	}
	// iter all rows
	//var rowmaps []map[string]string
	rows, err := f.GetRows(sheetname)
	if err != nil {
		fmt.Println(err)
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
