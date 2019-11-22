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
	err = f.Save()
	if err != nil {
		fmt.Println(err)
	}
}
