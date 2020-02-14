package main

import (
	"encoding/csv"
	"flag"
	"golang-base/tools"
	"io"
	"math"
	"os"
	"regexp"
)

var (
	inputFile *string
	accountID *string
	help      *bool
)

func init() {
	inputFile = flag.String("file", "beijing.csv", "Source File To Be Processed")
	accountID = flag.String("id", "216745712527", "Linked Account ID")
	help = flag.Bool("h", false, "Print This Message")
}

// RateOfProgress is a counter for progress, such as a progress bar
func RateOfProgress(inputFile string) int {
	lineCount := tools.CountRecord(inputFile)
	if lineCount > 100000 {
		oneTenthCount := float64(lineCount / 10)
		return int(math.Ceil(oneTenthCount))
	}
	return 0
}

func main() {
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
	baseRateCount := RateOfProgress(*inputFile)
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
			match, _ := regexp.MatchString(*accountID, record[1])
			if match == true {
				operateMatch, _:= regexp.MatchString(`.*Run.*`, record[10])
				if operateMatch == true {
					resourceMatch, _ := regexp.MatchString(`^i-.*`, record[21])
					if resourceMatch == true {
						writer.Write(record)
					}
				}
			}
		}
		lineCount++
	}
	// flush to file
	writer.Flush()
	// End task
	tools.InfoLogger.Println("Task Done")
}
