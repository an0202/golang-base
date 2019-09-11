package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"
)

func main() {
	// Read file from csv.csv
	inputCSV, inputError := os.OpenFile("origin.csv", os.O_RDONLY, 0666)

	if inputError != nil {
		fmt.Println("Error while open csv file")
	}
	defer inputCSV.Close()
	// init csv reader
	reader := csv.NewReader(inputCSV)
	// out put message to output.csv
	outputCSV, outputError := os.OpenFile("csv_output.csv", os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Println("An error occurred with file opening or cration")
		return
	}
	defer outputCSV.Close()
	writer := csv.NewWriter(outputCSV)
	// count line from src file
	linecount := 1
	for {
		// if linecount == 200000 {
		// 	break
		// }

		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
		}
		// out put fist line "title" to output.csv
		if linecount == 1 {
			writer.Write(record)
			// fmt.Println(record[2])
		} else {
			// out put specific line to output.csv
			// for _, value := range record {
			// 	match, _ := regexp.MatchString("^869869223565$", value)
			// 	if match == true {
			// 		// fmt.Println(record, linecount)
			// 		writer.Write(record)
			// 	}
			// }
			if linecount%120000 == 0 {
				fmt.Println("[INFO] Be Patient, Processed Rows :", linecount, time.Now())
			}
			match, _ := regexp.MatchString("^869869223565$", record[2])
			if match == true {
				// fmt.Println(record, linecount)
				writer.Write(record)
			}
		}
		linecount++
		// fmt.Println("Count Line:", linecount, "Value:", record)
		// fmt.Println(reflect.TypeOf(record))
	}
	// flush to file
	writer.Flush()
	// Number of processed rows
	fmt.Println("[INFO] Task Done ,Number Of Processed Rows:", linecount)
}
