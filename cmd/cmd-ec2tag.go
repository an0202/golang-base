///**
// * @Author: jie.an
// * @Description:
// * @File:  cmd-test.go
// * @Version: 1.0.0
// * @Date: 2019/11/23 22:06
// */
package cmd

import (
	"flag"
	"golang-base/aws"
	"golang-base/excel"
	"golang-base/tools"
	"os"
)

//var (
//	excelFile *string
//	sheetName *string
//	region    *string
//	method    *string
//	help      *bool
//)

func inittag() {
	excelFile = flag.String("file", "tags.xlsx", "Source ExcelFile To Be Processed")
	sheetName = flag.String("sheet", "EC2", "Sheet In ExcelFile To Be Processed")
	region = flag.String("region", "cn-north-1", "AWS Region")
	method = flag.String("m", "add", "add/del Tags")
	help = flag.Bool("h", false, "Print This Message")
}

func EC2addtags() {
	inittag()
	// Parse flag
	flag.Parse()
	if *help == true {
		flag.Usage()
		os.Exit(1)
	}
	if flag.NFlag() == 0 {
		flag.Usage()
	}
	switch *method {
	case "add":
		sess := aws.InitSession(*region)
		a := excel.ReadTest(*excelFile, *sheetName)
		for _, v := range a {
			b := aws.EC2InstanceMarshal(v)
			aws.EC2CreateTags(sess, b)
		}
	case "del":
		sess := aws.InitSession(*region)
		a := excel.ReadTest(*excelFile, *sheetName)
		for _, v := range a {
			b := aws.EC2InstanceMarshal(v)
			aws.EC2DeleteTags(sess, b)
		}
	default:
		flag.Usage()
		tools.ErrorLogger.Fatalln("Illegal Method:", *method)
	}
}
