/**
 * @Author: jie.an
 * @Description:
 * @File:  cmd-csv2
 * @Version: 1.0.0
 * @Date: 2020/04/26 10:06
 */
package cmd

import (
	"encoding/csv"
	"flag"
	"fmt"
	csvtool "golang-base/csv"
	"golang-base/tools"
	"io"
	"os"
	"regexp"
	"strings"
)

func initSamsungBill2() {
	inputFile = flag.String("file", "beijing.csv", "Source File To Be Processed")
	accountIDs = flag.String("id", "405718244235,0123456789", "Linked Account IDs Split With \",\"")
	help = flag.Bool("h", false, "Print This Message")
}

func SamsungBillFilter2() {
	initSamsungBill2()
	// Parse flag
	flag.Parse()
	if *help == true {
		flag.Usage()
		os.Exit(1)
	}
	if flag.NFlag() == 0 {
		flag.Usage()
	}
	tools.InfoLogger.Println("Task Start")
	// PrintTitle(inputFile)
	baseRateCount := csvtool.RateOfProgress(*inputFile)
	// Read file from csv.csv
	inputCSV, inputError := os.OpenFile(*inputFile, os.O_RDONLY, 0666)

	if inputError != nil {
		tools.ErrorLogger.Fatalln(inputError)
	}
	defer inputCSV.Close()
	// init csv reader
	reader := csv.NewReader(inputCSV)
	// out put message to output.csv
	outputCSV, outputError := os.OpenFile(*inputFile+"-Export.csv", os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		tools.ErrorLogger.Fatalln(outputError)
		return
	}
	defer outputCSV.Close()
	writer := csv.NewWriter(outputCSV)
	// count line from src file
	lineCount := 1
	// mspBill filter
	ms := csvtool.NewMSPBillings()
	for {
		// if lineCount == 200000 {
		// 	break
		// }
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			tools.ErrorLogger.Fatalln(err)
		}
		// out put fist line "title" to output.csv
		if lineCount == 1 {
			writer.Write(record)
		} else {
			// show progressBar if necessary , out put specific to output.csv
			if (baseRateCount != 0) && (lineCount%baseRateCount == 0) {
				tools.InfoLogger.Println("Processing , Processed Rows :", lineCount)
			}
			accountList := strings.Split(*accountIDs, ",")
			if tools.StringFind(accountList, record[2]) {
				operateMatch, _ := regexp.MatchString(`.*Run.*`, record[10])
				if operateMatch == true {
					resourceMatch, _ := regexp.MatchString(`^i-.*`, record[21])
					if resourceMatch == true {
						if record[22] != "No" {
							writer.Write(record)
							// mspBill filter
							ms.ProcessMSPBillings(record)
						}
					}
				}
			}
		}
		lineCount++
	}
	// flush to file
	writer.Flush()
	// mspBill filter
	for _, v := range ms.MSPBillings {
		fmt.Println(v.LinkedAccountId, v.ResourceId, v.UsageType, v.UserTag1, v.UserTag2, v.RunningDays)
	}
	// End task
	tools.InfoLogger.Println("Task Done")
}
