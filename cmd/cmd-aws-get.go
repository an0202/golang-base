///**
// * @Author: jie.an
// * @Description:
// * @File:  cmd-getResources.go
// * @Version: 1.0.0
// * @Date: 2020/02/19 22:42
// */
package cmd

import (
	"flag"
	"golang-base/aws"
	"golang-base/excel"
	"golang-base/tools"
	"os"
	"strings"
)

func initResources() {
	configFile = flag.String("f", "config.xlsx", "Read Config From Excel Line By Line")
	sheetName = flag.String("sheet", "default_config", "Sheet With Config To Be Process")
	//region = flag.String("region", "cn-north-1", "AWS Region")
	summary = flag.Bool("s", false, "Summarize The Operation Results To Sheet \"Total\", Operate Must Be Same!")
	help = flag.Bool("h", false, "Print This Message")
}

func GetAWSResources() {
	initResources()
	// Parse flag
	flag.Parse()
	if *help == true {
		flag.Usage()
		os.Exit(1)
	}
	if flag.NFlag() == 0 {
		flag.Usage()
	}
	// Read Config From File
	if *summary == false {
		if *configFile != "" {
			configs := excel.ReadToMaps(*configFile, *sheetName)
			var outputFile = "output.xlsx"
			excel.CreateFile(outputFile)
			for _, config := range configs {
				c := aws.ExcelConfigMarshal(config)
				c.Do(outputFile)
			}
		} else {
			tools.ErrorLogger.Fatalln("Not Currently Supported, Please Use Excel Config File")
		}
	} else {
		if *configFile != "" {
			configs := excel.ReadToMaps(*configFile, *sheetName)
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
				var outputFile = "output.xlsx"
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
				var outputFile = "output.xlsx"
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
}
