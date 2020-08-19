package csv

import (
	"encoding/csv"
	"fmt"
	"github.com/dimchansky/utfbom"
	"golang-base/tools"
	"io"
	"math"
	"os"
	"strings"
)

// PrintTitle return the first line of csv file
func PrintTitle(inputFile string) {
	// Read CSV File
	inputCSV, inputError := os.OpenFile(inputFile, os.O_RDONLY, 0666)
	if inputError != nil {
		tools.ErrorLogger.Fatalln("Error While Open CSV File")
	}
	defer inputCSV.Close()
	// Init CSV Reader
	reader := csv.NewReader(inputCSV)
	record, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			tools.ErrorLogger.Fatalln("No Title For This CSV File")
		}
		tools.ErrorLogger.Fatalln(err)
	}
	tools.InfoLogger.Println("Title:")
	for id, value := range record {
		fmt.Println(id, value)
	}
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

//
// Convert CSV Records to Dictionaries using Header Row as Keys
// https://gist.github.com/drernie/5684f9def5bee832ebc50cabb46c377a
//CSV:
//
//	A     B     C
//
//1  Name   Age   Sex
//
//2  Alice  18    Female
//
//3  Bob    22    Male
//
//ReturnData:
//[
//
//	{"Name":"Alice","Age":"18","Sex":"Female"},
//	{"Name":"Bob","Age":"22","Sex":"Male"}
//
//]
func ReadToMaps(inputFile string) (rowMaps []map[string]string) {
	// Read CSV File
	inputCSV, inputError := os.OpenFile(inputFile, os.O_RDONLY, 0666)
	if inputError != nil {
		tools.ErrorLogger.Fatalln("Error While Open CSV File")
	}
	defer inputCSV.Close()
	// Init CSV Reader and handle BOM
	bomReader := utfbom.SkipOnly(inputCSV)
	r := csv.NewReader(bomReader)
	var header []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			tools.WarningLogger.Panicln(err)
		}
		if header == nil {
			for _, v := range record {
				header = append(header, strings.TrimSpace(v))
			}
		} else {
			dict := map[string]string{}
			for i := range header {
				dict[header[i]] = record[i]
			}
			rowMaps = append(rowMaps, dict)
		}
	}
	return rowMaps
}
