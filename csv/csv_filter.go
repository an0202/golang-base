package main

import (
	"encoding/csv"
	"flag"
	"golang-base/tools"
	"io"
	"math"
	"os"
	"regexp"
	"strings"
)

var (
	inputFile *string
	accountIDs *string
	help      *bool
)

func init() {
	inputFile = flag.String("file", "beijing.csv", "Source File To Be Processed")
	accountIDs = flag.String("id", "405718244235,0123456789", "Linked Account IDs Split With \",\"")
	/*
	567306684220,405509095605,405625928234,405718244235,405757921744,406117769501,406173781768,406419040736,535285806536,405351385095,406144704895,406752347000,559484794032,140084596652,292521336001,405660919495,302234245394,243643571971,240824288423,241300967417,241488999070,241716109803,242004344475,242265030755,242352519975,242485912731,242671812748,242771374863,243462577727,140682541603,776981220315,519723696092
	*/
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
			accountList := strings.Split(*accountIDs,",")
			if tools.StringFind(accountList, record[2]) {
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
