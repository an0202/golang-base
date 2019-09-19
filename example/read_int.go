package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {

	inputFile, inputError := os.OpenFile("text.log", os.O_RDONLY, 0666)
	if inputError != nil {
		fmt.Println("An error occurred with file opening")
		return
	}
	defer inputFile.Close()
	reader := bufio.NewScanner(inputFile)
	buf := make([]byte, 4*1024)
	reader.Buffer(buf, 10*1024)
	//
	for reader.Scan() {
		// fmt.Println(reader.Text())
		convToint, _ := strconv.Atoi(reader.Text())
		if convToint%86132 == 0 {
			println(convToint)
			// fmt.Println(reflect.TypeOf(convToint))
		}
	}
	// for reader.Scan() {
	// 	fmt.Println(reader.Text())
	// }
}
