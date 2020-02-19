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
)

func initResources() {
	configFile = flag.String("f", "config.xlsx", "Read Config From Excel Line By Line")
	sheetName = flag.String("sheet", "default_config", "Sheet With Config To Be Process")
	region = flag.String("region", "cn-north-1", "AWS Region")
	suffix = flag.String("m", "date", "Add date/final/.. Suffix To AMI Name")
	instanceid = flag.String("i", "i-abc123", "AWS EC2 InstanceID")
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
}
