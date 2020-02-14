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
	"reflect"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

var rowList []interface{}

func CreateFile(path, sheetname string) {
	f := excelize.NewFile()
	f.NewSheet(sheetname)
	err := f.SaveAs(path)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Create New File", path, sheetname)
}

// Set Rows From Golang Struct Type
func SetStructRows(path, sheetname string, rows []interface{}) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
	}
	if f.GetSheetIndex(sheetname) == 0 {
		f.NewSheet(sheetname)
	}
	for index, row := range rows {
		rowList = nil
		v := reflect.ValueOf(row)
		for i := 0; i < v.NumField(); i++ {
			// fmt.Println(i, v.Field(i))
			rowList = append(rowList, v.Field(i).Interface())
		}
		// fmt.Println(rowList)
		err := f.SetSheetRow(sheetname, "A"+strconv.Itoa(index+2), &rowList)
		if err != nil {
			fmt.Println(err)
		}
	}
	err = f.Save()
	if err != nil {
		fmt.Println(err)
	}
}

// Set Rows From Golang List
func SetListRows(path, sheetname string, rows [][]interface{}) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
	}
	if f.GetSheetIndex(sheetname) == 0 {
		f.NewSheet(sheetname)
	}
	for index, rowList := range rows {
		err := f.SetSheetRow(sheetname, "A"+strconv.Itoa(index+2), &rowList)
		if err != nil {
			fmt.Println(err)
		}
	}
	err = f.Save()
	if err != nil {
		fmt.Println(err)
	}
}

// SetHeaderLine From A List
func SetHeaderLine(path, sheetname string, HeaderLine []interface{}) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
	}
	if f.GetSheetIndex(sheetname) == 0 {
		f.NewSheet(sheetname)
	}
	err = f.SetSheetRow(sheetname, "A1", &HeaderLine)
	if err != nil {
		fmt.Println(err)
	}
	// Set Cell Style
	// https://xuri.me/excelize/zh-hans/cell.html#SetCellStyle
	// https://xuri.me/excelize/zh-hans/style.html#shading
	style, err := f.NewStyle(`{"font":{"bold":true,"family":"Microsoft YaHei Light","size":12,"color":"#000000"},
	"fill":{"type":"pattern","color":["#F9F900"],"pattern":1}}`)
	if err != nil {
		println(err.Error())
	}
	err = f.SetCellStyle(sheetname, "A1", "J1", style)
	err = f.Save()
	if err != nil {
		fmt.Println(err)
	}
}
