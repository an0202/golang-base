/**
 * @Author: jie.an
 * @Description:
 * @File:  get.go
 * @Version: 1.0.0
 * @Date: 2020/7/20 17:08
 */
package cmd

import (
	"github.com/spf13/cobra"
	"golang-base/aws"
	"golang-base/excel"
	"golang-base/tools"
	"strings"
)

// getCmd
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Assets",
	Long: `Get Assets using excel profile
`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		sheetName, _ := cmd.Flags().GetString("sheet")
		summary, _ := cmd.Flags().GetBool("s")
		// Read Config From File
		if summary == false {
			if file != "" {
				configs := excel.ReadToMaps(file, sheetName)
				var outputFile = "aws-get-output.xlsx"
				excel.CreateFile(outputFile)
				for _, config := range configs {
					c := aws.ExcelConfigMarshal(config)
					c.Do(outputFile)
				}
			} else {
				tools.ErrorLogger.Fatalln("Not Currently Supported, Please Use Excel Config File")
			}
		} else {
			if file != "" {
				configs := excel.ReadToMaps(file, sheetName)
				var operateList []string
				// check each config, if there is a different operate
				for _, config := range configs {
					c := aws.ExcelConfigMarshal(config)
					if c.Operate != "" {
						operateList = append(operateList, c.Operate)
					} else {
						continue
					}
				}
				op := tools.UniqueStringList(operateList)
				if len(op) != 1 {
					tools.ErrorLogger.Fatalln("OperateList Must Be Same , Current OperateList Is", op)
					return
				}
				// switchCase for List and Liv2
				switch strings.HasPrefix(op[0], "Liv2") {
				case true:
					rowsNum := 1
					var totalHeadLine []interface{}
					var outputFile = "aws-get-output.xlsx"
					excel.CreateFile(outputFile)
					for _, config := range configs {
						c := aws.ExcelConfigMarshal(config)
						results := c.ReturnResourcesV2()
						// use last result as totalHeadline
						if c.HeadLine != nil {
							totalHeadLine = c.HeadLine
						}
						if len(results) != 0 {
							tools.InfoLogger.Printf("Found %d Result In %s (%s) \n", len(results), c.AccountId, c.Region)
							excel.SetHeadLine(outputFile, c.OutputSheet, c.HeadLine)
							excel.SetStructRows(outputFile, c.OutputSheet, results)
							// Write summary data to Total sheet
							excel.SetStructRowsV2(outputFile, "Total", "A", rowsNum+1, results)
							rowsNum += len(results)
						} else {
							tools.InfoLogger.Printf("No Result In %s (%s) \n", c.AccountId, c.Region)
						}
					}
					excel.SetHeadLine(outputFile, "Total", totalHeadLine)
				case false:
					rowsNum := 1
					var totalHeadLine []interface{}
					var outputFile = "aws-get-output.xlsx"
					excel.CreateFile(outputFile)
					for _, config := range configs {
						c := aws.ExcelConfigMarshal(config)
						results := c.ReturnResources()
						// use last result as totalHeadline
						if c.HeadLine != nil {
							totalHeadLine = c.HeadLine
						}
						if len(results) != 0 {
							tools.InfoLogger.Printf("Found %d Result In %s (%s) \n", len(results), c.AccountId, c.Region)
							excel.SetHeadLine(outputFile, c.OutputSheet, c.HeadLine)
							excel.SetListRows(outputFile, c.OutputSheet, results)
							// Write summary data to Total sheet
							excel.SetListRowsV2(outputFile, "Total", "A", rowsNum+1, results)
							rowsNum += len(results)
						} else {
							tools.InfoLogger.Printf("No Result In %s (%s) \n", c.AccountId, c.Region)
						}
					}
					excel.SetHeadLine(outputFile, "Total", totalHeadLine)
				}
			} else {
				tools.ErrorLogger.Fatalln("Not Currently Supported, Please Use Excel Config File")
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("file", "f", "config.xlsx", "Read Config From Excel Line By Line")
	getCmd.Flags().StringP("sheet", "s", "default_config", `Sheet With Config To Be Process`)
	getCmd.Flags().Bool("s", false, `Summarize The Operation Results To Sheet "Total", Operate Must Be Same!`)
}
