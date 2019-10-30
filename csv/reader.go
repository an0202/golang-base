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
	for id, value := range record {
		tools.InfoLogger.Println("Title:")
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
// func PrintNRecord(inputFile string, N int) {
// 	// Handle Error
// 	recordNum := tools.CountRecord(inputFile)
// 	if N > recordNum {
// 		tools.ErrorLogger.Fatalln("Out Of Range")
// 	}
// 	// Read CSV File
// 	inputCSV, inputError := os.OpenFile(inputFile, os.O_RDONLY, 0666)
// 	if inputError != nil {
// 		tools.ErrorLogger.Fatalln("Error While Open CSV File")
// 	}
// 	defer inputCSV.Close()
// 	// Init CSV Reader
// 	reader := csv.NewReader(inputCSV)
// 	//
// 	linecount := 1
// 	for {
// 		record, err := reader.Read()
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			fmt.Println(err)
// 		}
// 		if linecount == N {
// 			tools.InfoLogger.Println("The Nth Row Of Data:")
// 			fmt.Println("Record :", record)
// 			break
// 		}
// 		linecount++
// 	}
// }

// func main() {
// 	inputFile := "azure_bill.csv"
// 	// PrintTitle(inputFile)
// 	tools.PrintNRecord(inputFile, 11332)
// 	// tools.CountRecord(inputFile)
// }
