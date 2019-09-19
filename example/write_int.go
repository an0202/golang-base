package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {

	outputFile, outPutError := os.OpenFile("text.log", os.O_WRONLY|os.O_CREATE, 0666)
	if outPutError != nil {
		fmt.Println("An error occurred with file opening or cration")
		return
	}
	defer outputFile.Close()

	outputWriter := bufio.NewWriter(outputFile)
	// outputString := "hello world! \n"

	for i := 0; i <= 400000000; i++ {
		// outputWriter.WriteString(outputString)

		outputWriter.WriteString(strconv.Itoa(i) + "\n")
	}
	outputWriter.Flush()
}
