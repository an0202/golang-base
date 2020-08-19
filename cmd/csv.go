/**
 * @Author: jie.an
 * @Description:
 * @File:  csv.go
 * @Version: 1.0.0
 * @Date: 2020/8/19 11:18
 */
package cmd

import (
	"encoding/csv"
	"github.com/spf13/cobra"
	csvtool "golang-base/csv"
	"golang-base/tools"
	"io"
	"os"
	"path"
	"strings"
)

// csvCmd
var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "Process CSV File",
	Long: `Extract specific data from csv file,
`,
	Run: func(cmd *cobra.Command, args []string) {
		inputFile, _ := cmd.Flags().GetString("file")
		column, _ := cmd.Flags().GetInt("col")
		include, _ := cmd.Flags().GetStringSlice("inc")
		title, _ := cmd.Flags().GetBool("title")
		if title == true {
			csvtool.PrintTitle(inputFile)
		} else {
			tools.InfoLogger.Println("Task Start")
			// PrintTitle(inputFile)
			baseRateCount := csvtool.RateOfProgress(inputFile)
			// Read file from csv.csv
			inputCSV, inputError := os.OpenFile(inputFile, os.O_RDONLY, 0666)

			if inputError != nil {
				tools.ErrorLogger.Fatalln(inputError)
			}
			defer inputCSV.Close()
			// init csv reader
			reader := csv.NewReader(inputCSV)
			// out put message to output.csv
			outputCSV, outputError := os.OpenFile(strings.TrimSuffix(path.Base(inputFile), ".csv")+"-output.csv", os.O_WRONLY|os.O_CREATE, 0666)
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
					if tools.StringFind(include, record[column]) {
						writer.Write(record)
					}
				}
				lineCount++
			}
			// flush to file
			writer.Flush()
			// End task
			tools.InfoLogger.Println("Task Done")
		}
	},
}

func init() {

	RootCmd.AddCommand(csvCmd)
	csvCmd.Flags().BoolP("title", "t", false, "whether to print title")
	csvCmd.Flags().StringP("file", "f", "origin.csv", "origin file to be process")
	csvCmd.Flags().IntP("col", "c", 0, `Condition field (Column To Be Processed)`)
	csvCmd.Flags().StringSliceP("inc", "i", []string{"405718244235", "0123456789"}, `Condition values`)
}
