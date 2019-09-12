package main

import (
	"encoding/csv"
	"io"
	"math"
	"os"
	"regexp"
	"tools"
)

// RateOfProgress is a counter for progeress, such as a progress bar
func RateOfProgress(inputFile string) int {
	lineCount := tools.CountRecord(inputFile)
	if lineCount > 100000 {
		oneTenthCount := float64(lineCount / 10)
		return int(math.Ceil(oneTenthCount))
	}
	return 0
}

func main() {
	inputFile := "aws_bill0.csv"
	accountID := "869869223565"

	tools.InfoLogger.Println("Task Start")
	// PrintTitle(inputFile)
	baseRateCount := RateOfProgress(inputFile)
	// Read file from csv.csv
	inputCSV, inputError := os.OpenFile(inputFile, os.O_RDONLY, 0666)

	if inputError != nil {
		tools.ErrorLogger.Fatalln(inputError)
	}
	defer inputCSV.Close()
	// init csv reader
	reader := csv.NewReader(inputCSV)
	// out put message to output.csv
	outputCSV, outputError := os.OpenFile("csv_fliter_output.csv", os.O_WRONLY|os.O_CREATE, 0666)
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
			match, _ := regexp.MatchString(accountID, record[2])
			if match == true {
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
