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
	"golang-base/tools"
	"reflect"
	"strconv"
)

var rowList []interface{}

func CreateFile(path string) {
	f := excelize.NewFile()
	//https://xuri.me/excelize/en/workbook.html#SetDocProps
	f.SetDocProps(&excelize.DocProperties{
		Creator:     "Jie An",
		Description: "This file created by Go Excelize",
	})
	err := f.SaveAs(path)
	if err != nil {
		fmt.Println(err)
	}
	tools.InfoLogger.Println("Create New File", path)
}

// Set Rows From Map V2 （Start Cell Can Be Specified）
func SetMapRows(path, sheetName string, mapRow map[interface{}]interface{}) {
	SetMapRowsV2(path, sheetName, "A", 2, mapRow)
}

// Set Rows From Map V2 （Start Cell Can Be Specified）
func SetMapRowsV2(path, sheetName string, startColumn string, startRow int, mapRow map[interface{}]interface{}) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
	}
	if f.GetSheetIndex(sheetName) == -1 {
		f.NewSheet(sheetName)
	}
	val := reflect.TypeOf(mapRow)
	if val.Kind() == reflect.Map {
		index := 0
		for k, v := range mapRow {
			rowList = nil
			rowList = append(rowList, k, v)
			err := f.SetSheetRow(sheetName, startColumn+strconv.Itoa(index+startRow), &rowList)
			if err != nil {
				fmt.Println(err)
			}
			index++
		}
	} else {
		panic("input data is not a map")
	}
	err = f.Save()
	if err != nil {
		fmt.Println(err)
	}
}

// Set Rows From Golang Struct Type
func SetStructRows(path, sheetName string, rows []interface{}) {
	SetStructRowsV2(path, sheetName, "A", 2, rows)
}

// Set Rows From Golang Struct V2 （Start Cell Can Be Specified）
func SetStructRowsV2(path, sheetName, startColumn string, startRow int, rows []interface{}) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
	}
	if f.GetSheetIndex(sheetName) == -1 {
		f.NewSheet(sheetName)
	}
	for index, row := range rows {
		rowList = nil
		v := reflect.ValueOf(row)
		for i := 0; i < v.NumField(); i++ {
			//fmt.Println(i, v.Field(i))
			rowList = append(rowList, v.Field(i).Interface())
		}
		// fmt.Println(rowList)
		err := f.SetSheetRow(sheetName, startColumn+strconv.Itoa(index+startRow), &rowList)
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
func SetListRows(path, sheetName string, rows [][]interface{}) {
	SetListRowsV2(path, sheetName, "A", 2, rows)
}

// Set Rows From Golang List V2 （Start Cell Can Be Specified）
func SetListRowsV2(path, sheetName, startColumn string, startRow int, rows [][]interface{}) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
	}
	if f.GetSheetIndex(sheetName) == -1 {
		f.NewSheet(sheetName)
	}
	for index, rowList := range rows {
		err := f.SetSheetRow(sheetName, startColumn+strconv.Itoa(index+startRow), &rowList)
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
func SetHeadLine(path, sheetName string, HeadLine []interface{}) {
	SetHeadLineV2(path, sheetName, "A1", HeadLine)
}

// SetHeadLine From A List V2（Start Cell Can Be Specified）
func SetHeadLineV2(path, sheetName string, startCell string, HeadLine []interface{}) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
	}
	if f.GetSheetIndex(sheetName) == -1 {
		f.NewSheet(sheetName)
	}
	err = f.SetSheetRow(sheetName, startCell, &HeadLine)
	if err != nil {
		fmt.Println(err)
	}
	// Parameter calculation
	/// https://godoc.org/github.com/360EntSecGroup-Skylar/excelize#CellNameToCoordinates
	startX, startY, _ := excelize.CellNameToCoordinates(startCell)
	// Set Cell Style
	// https://xuri.me/excelize/zh-hans/cell.html#SetCellStyle
	// https://xuri.me/excelize/zh-hans/style.html#shading
	style, err := f.NewStyle(`{"font":{"bold":true,"family":"Microsoft YaHei Light","size":12,"color":"#000000"},
	"fill":{"type":"pattern","color":["#F9F900"],"pattern":1}}`)
	if err != nil {
		println(err.Error())
	}
	for i := 0; i < len(HeadLine); i++ {
		cellName, _ := excelize.ColumnNumberToName(startX)
		cellNum := startY
		cell := cellName + strconv.Itoa(cellNum)
		err = f.SetCellStyle(sheetName, cell, cell, style)
		startX += 1
	}
	//always set sheet 1 as active sheet , used to hidden "Sheet1" , "Sheet1" can not be delete for now.
	f.SetActiveSheet(1)
	err = f.SetSheetVisible("Sheet1", false)
	if err != nil {
		fmt.Println(err)
	}
	err = f.Save()
	if err != nil {
		fmt.Println(err)
	}
}
