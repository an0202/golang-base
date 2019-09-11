package tools

import (
	"bufio"
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
	count := 0
	// Count
	for scanner.Scan() {
		count++
	}
	// fmt.Println("[Info]:", count)
	InfoLogger.Println("Number Of File Lines:", count)
	return count
}
