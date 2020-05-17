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
	"golang-base/tools"
	"reflect"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

var rowList []interface{}

func CreateFile(path string) {
	f := excelize.NewFile()
	//https://xuri.me/excelize/en/workbook.html#SetDocProps
	f.SetDocProps(&excelize.DocProperties{
		Created:        time.Now().String(),
		Creator:        "Jie An",
		Description:    "This file created by Go Excelize",
		Identifier:     "xlsx",
		Keywords:       "Spreadsheet",
		LastModifiedBy: "Excelize",
		Modified:       time.Now().String(),
		Version:        "1.0.0",
	})
	err := f.SaveAs(path)
	if err != nil {
		fmt.Println(err)
	}
	tools.InfoLogger.Println("Create New File", path)
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
			//fmt.Println(i, v.Field(i))
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

// Set Rows From Golang Struct V2 （Start Cell Can Be Specified）
func SetStructRowsV2(path, sheetname, startColumn string, startRow int, rows []interface{}) {
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
			//fmt.Println(i, v.Field(i))
			rowList = append(rowList, v.Field(i).Interface())
		}
		// fmt.Println(rowList)
		err := f.SetSheetRow(sheetname, startColumn+strconv.Itoa(index+startRow), &rowList)
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

// Set Rows From Golang List V2 （Start Cell Can Be Specified）
func SetListRowsV2(path, sheetname, startColumn string, startRow int, rows [][]interface{}) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
	}
	if f.GetSheetIndex(sheetname) == 0 {
		f.NewSheet(sheetname)
	}
	for index, rowList := range rows {
		err := f.SetSheetRow(sheetname, startColumn+strconv.Itoa(index+startRow), &rowList)
		if err != nil {
			fmt.Println(err)
		}
	}
	err = f.Save()
	if err != nil {
		fmt.Println(err)
	}
}

// SetHeadLine From A List
func SetHeadLine(path, sheetname string, HeadLine []interface{}) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
	}
	if f.GetSheetIndex(sheetname) == 0 {
		f.NewSheet(sheetname)
	}
	err = f.SetSheetRow(sheetname, "A1", &HeadLine)
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
	err = f.SetCellStyle(sheetname, "A1", DescribeLastPosition(len(HeadLine)), style)
	//always set sheet 2 as active sheet , used to hidden "Sheet1" , "Sheet1" can not be delete for now.
	f.SetActiveSheet(2)
	err = f.SetSheetVisible("Sheet1", false)
	if err != nil {
		fmt.Println(err)
	}
	err = f.Save()
	if err != nil {
		fmt.Println(err)
	}
}
