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
	"strings"
)

//var (
//	excelFile *string
//	sheetName *string
//	region    *string
//	method    *string
//	help      *bool
//)

func initTag() {
	excelFile = flag.String("file", "tags.xlsx", "Source ExcelFile To Be Processed")
	sheetName = flag.String("sheet", "EC2", "Sheet In ExcelFile To Be Processed")
	region = flag.String("region", "cn-north-1", "Used For Init A AWS Default Session")
	method = flag.String("m", "get", "add/del/get Tags")
	tags = flag.String("tags", "Name,Env,Project", "Require:[ method = get],Get Specific Tags From Resource And Write To Excel")
	overide = flag.Bool("o", true, "Overide Exist Tags")
	help = flag.Bool("h", false, "Print This Message")
}

func EC2Tags() {
	initTag()
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
		defaultSess := aws.InitSession(*region)
		se := new(aws.Session)
		a := excel.ReadToMaps(*excelFile, *sheetName)
		for _, v := range a {
			b := aws.EC2InstanceMarshal(v)
			if b.AWSProfile != "" && b.Region != "" {
				// Use the old session when the current resource and the previous resource belong to the same account and region
				if b.AWSProfile == se.UsedAwsProfile && b.Region == se.UsedRegion {
					//fmt.Println(se.Sess.Config.Credentials)
					aws.EC2CreateTags(se.Sess, b, *overide)
				} else {
					// create a new session
					se.InitSessionWithAWSProfile(b.Region,b.AWSProfile)
					//fmt.Println(se.Sess.Config.Credentials)
					aws.EC2CreateTags(se.Sess, b, *overide)
				}
				//se.Sess = se.InitSessionWithAWSProfile(b.Region,b.AWSProfile)
				//fmt.Println(se.Sess.Config.Credentials)
				//aws.EC2CreateTags(se.Sess, b, *overide)
			} else {
				tools.InfoLogger.Println("Use Default Session")
				aws.EC2CreateTags(defaultSess, b, *overide)
			}
		}
	case "del":
		sess := aws.InitSession(*region)
		a := excel.ReadToMaps(*excelFile, *sheetName)
		for _, v := range a {
			b := aws.EC2InstanceMarshal(v)
			aws.EC2DeleteTags(sess, b)
		}
	case "get":
		defaultSess := aws.InitSession(*region)
		se := new(aws.Session)
		a := excel.ReadToMaps(*excelFile, *sheetName)
		var results [][]interface{}
		for _, v := range a {
			b := aws.EC2InstanceMarshal(v)
			if b.AWSProfile != "" && b.Region != "" {
				// Use the old session when the current resource and the previous resource belong to the same account and region
				if b.AWSProfile == se.UsedAwsProfile && b.Region == se.UsedRegion {
					//fmt.Println(se.Sess.Config.Credentials)
					_, result := aws.EC2GetTags(se.Sess, b, *tags)
					results = append(results, result)
				} else {
					// create a new session
					se.InitSessionWithAWSProfile(b.Region,b.AWSProfile)
					// fmt.Println(se.Sess.Config.Credentials)
					_, result := aws.EC2GetTags(se.Sess, b, *tags)
					results = append(results, result)
				}
			} else {
				tools.InfoLogger.Println("Use Default Session")
				_, result := aws.EC2GetTags(defaultSess, b, *tags)
				results = append(results, result)
			}
		}
		var headerline = []interface{}{"ResourceId"}
		for _, v := range strings.Split(*tags,",") {
			headerline = append(headerline, v)
		}
		excel.CreateFile("output-"+*excelFile)
		excel.SetHeadLine("output-"+*excelFile,"result", headerline)
		excel.SetListRows("output-"+*excelFile,"result", results)
	default:
		flag.Usage()
		tools.ErrorLogger.Fatalln("Illegal Method:", *method)
	}
}
