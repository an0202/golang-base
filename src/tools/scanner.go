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
	buf := make([]byte, 4*1024)
	scanner.Buffer(buf, 10*1024)
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
	buf := make([]byte, 4*1024)
	scanner.Buffer(buf, 10*1024)
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
