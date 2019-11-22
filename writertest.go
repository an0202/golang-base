/**
 * @Author: jie.an
 * @Description:
 * @File:  reader.go
 * @Version: 1.0.0
 * @Date: 2019/11/15 17:00
 */
package main

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func main() {
	f, err := excelize.OpenFile("C:\\Users\\jie.an\\OneDrive - BESPIN GLOBAL CHINA\\test.xlsx")
	if err != nil {
		fmt.Println(err)
	}
	f.NewSheet("aaaaa")
	err = f.SetSheetRow("aaaaa", "B6", &[]interface{}{"1", nil, 2})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(err)
	err = f.Save()
	if err != nil {
		fmt.Println(err)
	}
}
