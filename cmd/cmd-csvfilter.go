/**
 * @Author: jie.an
 * @Description:
 * @File:  cmd-csvfilter
 * @Version: 1.0.0
 * @Date: 2020/04/17 12:30 下午
 */
package cmd

import (
	"encoding/csv"
	"flag"
	csvtool "golang-base/csv"
	"golang-base/tools"
	"io"
	"os"
	"strings"
)

func initCsvFilter() {
	inputFile = flag.String("file", "beijing.csv", "Source File To Be Processed")
	column = flag.Int("col", 0, "Column To Be Processed")
	include = flag.String("inc", "405718244235,0123456789", "Linked Account IDs Split With \",\"")
	help = flag.Bool("h", false, "Print This Message")
}

func CSVFilter() {
	initCsvFilter()
	// Parse flag and title line
	flag.Parse()
	if *help == true {
		flag.Usage()
		csvtool.PrintTitle(*inputFile)
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
	outputCSV, outputError := os.OpenFile("output.csv", os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		tools.ErrorLogger.Fatalln(outputError)
		return
	}
	defer outputCSV.Close()
	writer := csv.NewWriter(outputCSV)
	// count line from src file
	lineCount := 1
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
			includeList := strings.Split(*include, ",")
			if tools.StringFind(includeList, record[*column]) {
				writer.Write(record)
			}
		}
		lineCount++
	}
	// flush to file
	writer.Flush()
	// End task
	tools.InfoLogger.Println("Task Done")
}
