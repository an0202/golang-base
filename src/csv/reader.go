package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"tools"
)

// PrintTitle return the first line of csv file
func PrintTitle(inputFile string) {
	// Read CSV File
	inputCSV, inputError := os.OpenFile(inputFile, os.O_RDONLY, 0666)
	if inputError != nil {
		fmt.Println("[Error] Error While Open CSV File")
		os.Exit(2)
	}
	defer inputCSV.Close()
	// Init CSV Reader
	reader := csv.NewReader(inputCSV)
	record, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			fmt.Println("[Error] No Title For This CSV File")
			os.Exit(1)
		}
		fmt.Println(err)
	}
	for id, value := range record {
		fmt.Println(id, value)
	}
}

// CountRecord return total records number of csv file
// func CountRecord(inputFile string) int {
// 	// Read CSV File
// 	inputCSV, inputError := os.OpenFile(inputFile, os.O_RDONLY, 0666)
// 	if inputError != nil {
// 		fmt.Println("[Error] Error While Open CSV File")
// 		os.Exit(2)
// 	}
// 	defer inputCSV.Close()
// 	// Init CSV Reader
// 	reader := csv.NewReader(inputCSV)
// 	// Count Line
// 	linecount := 1
// 	for {
// 		_, err := reader.Read()
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			fmt.Println(err)
// 		}
// 		linecount++
// 	}
// 	fmt.Println("Total Record(Exclude Title):", linecount-1)
// 	return linecount - 1
// }

// PrintNRecord return the Nth row of data
func PrintNRecord(inputFile string, N int) {
	// Handle Error
	recordNum := tools.CountRecord(inputFile)
	if N > recordNum {
		fmt.Println("[Error] Out Of Range")
		os.Exit(2)
	}
	// Read CSV File
	inputCSV, inputError := os.OpenFile(inputFile, os.O_RDONLY, 0666)
	if inputError != nil {
		fmt.Println("[Error] Error While Open CSV File")
		os.Exit(2)
	}
	defer inputCSV.Close()
	// Init CSV Reader
	reader := csv.NewReader(inputCSV)
	//
	linecount := 1
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
		}
		if linecount == N {
			fmt.Println("Record :", record)
			break
		}
		linecount++
	}
}

// func main() {
// 	inputFile := "aws_bill.csv"
// 	// PrintTitle(inputFile)
// 	// PrintNRecord(inputFile, 3000)
// 	tools.CountRecord(inputFile)
// }
