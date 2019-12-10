///**
// * @Author: jie.an
// * @Description:
// * @File:  cmd-test.go
// * @Version: 1.0.0
// * @Date: 2019/12/09 16:00
// */
package cmd

import (
	"flag"
	"golang-base/aws"
	"golang-base/tools"
	"os"
)

func initami() {
	//	excelFile = flag.String("file", "tags.xlsx", "Source ExcelFile To Be Processed")
	//	sheetName = flag.String("sheet", "EC2", "Sheet In ExcelFile To Be Processed")
	srcFile = flag.String("f", "", "Read Instance ids From File Line By Line")
	region = flag.String("region", "cn-north-1", "AWS Region")
	suffix = flag.String("m", "date", "Add date/final/.. Suffix To AMI Name")
	instanceid = flag.String("i", "i-abc123", "AWS EC2 InstanceID")
	help = flag.Bool("h", false, "Print This Message")
}

func EC2CreateAMI() {
	initami()
	// Parse flag
	flag.Parse()
	if *help == true {
		flag.Usage()
		os.Exit(1)
	}
	if flag.NFlag() == 0 {
		flag.Usage()
	}
	// Read instance id from file
	if *srcFile != "" {
		instanceids := tools.GetRecords(*srcFile)
		sess := aws.InitSession(*region)
		for _, instanceid := range instanceids {
			aws.CreateImage(sess, instanceid, *suffix)
		}
	} else {
		sess := aws.InitSession(*region)
		aws.CreateImage(sess, *instanceid, *suffix)
	}
}
