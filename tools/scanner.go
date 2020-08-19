package tools

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// CountRecord return number of lines
func CountRecord(srcFile string) int {
	// Print Time Duration
	defer TimeTrack(time.Now(), "CountRecord")
	// Open File
	inputFile, inputError := os.OpenFile(srcFile, os.O_RDONLY, 0666)
	if inputError != nil {
		ErrorLogger.Fatalln("Error While Open File :", inputError)
	}
	defer inputFile.Close()
	// Init Scanner
	scanner := bufio.NewScanner(inputFile)
	buf := make([]byte, 512*1024)
	scanner.Buffer(buf, 1024*1024)
	lineCount := 0
	// Count
	for scanner.Scan() {
		lineCount++
	}
	// fmt.Println("[Info]:", count)
	InfoLogger.Println("Number Of File Lines:", lineCount)
	return lineCount
}

// PrintNRecord return the Nth row of data
func PrintNRecord(srcFile string, N int) {
	// Handle Error
	recordNum := CountRecord(srcFile)
	if N > recordNum {
		ErrorLogger.Fatalln("Out Of Range")
	}
	// Print Time Duration
	defer TimeTrack(time.Now(), "PrintNRecord")
	// Open File
	inputFile, inputError := os.OpenFile(srcFile, os.O_RDONLY, 0666)
	if inputError != nil {
		ErrorLogger.Fatalln("Error While Open File :", inputError)
	}
	defer inputFile.Close()
	// Init Scanner
	scanner := bufio.NewScanner(inputFile)
	buf := make([]byte, 512*1024)
	scanner.Buffer(buf, 1024*1024)
	lineCount := 1
	// Count
	for scanner.Scan() {
		if lineCount == N {
			InfoLogger.Println("The Nth Row Of Data:")
			fmt.Println("Record :", scanner.Text())
			break
		}
		lineCount++
	}
}

// GetRecords return a list which value read line by line from a file.
func GetRecords(srcFile string) []string {
	// Print Time Duration
	defer TimeTrack(time.Now(), "GetRecords")
	// Open File
	inputFile, inputError := os.OpenFile(srcFile, os.O_RDONLY, 0666)
	if inputError != nil {
		ErrorLogger.Fatalln("Error While Open File :", inputError)
	}
	defer inputFile.Close()
	// Init Scanner
	scanner := bufio.NewScanner(inputFile)
	buf := make([]byte, 512*1024)
	scanner.Buffer(buf, 1024*1024)
	var values []string
	// Count
	for scanner.Scan() {
		values = append(values, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		ErrorLogger.Fatalln(err)
	}
	InfoLogger.Println("Value List:", values)
	return values
}
