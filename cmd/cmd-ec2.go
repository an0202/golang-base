///**
// * @Author: jie.an
// * @Description:
// * @File:  cmd-test.go
// * @Version: 1.0.0
// * @Date: 2020/02/15 21:17
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

func initEC2() {
	excelFile = flag.String("f", "", "Source ExcelFile To Be Processed")
	sheetName = flag.String("sheet", "EC2", "Sheet In ExcelFile To Be Processed")
	region = flag.String("region", "cn-north-1", "Used For Init A AWS Default Session")
	method = flag.String("m", "get", "Get EC2 Instance")
	help = flag.Bool("h", false, "Print This Message")
}

func EC2() {
	initEC2()
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
	case "get":
		excel.CreateFile("output.xlsx", "Default")
		var headerLine = []interface{}{"AccountId", "Region", "Name", "InstanceId", "InstanceType", "Platform", "State",
			"VPCId","Role","SubnetId","KeyPair","SecurityGroups","Tags"}
		// Get EC2 Instances From Excel (Excel Contains AWS_PROFILE And Region)
		if *excelFile != "" {
			se := new(aws.Session)
			a := excel.ReadToMaps(*excelFile, *sheetName)
			for _, v := range a {
				b := aws.EC2InstanceMarshal(v)
				// create a new session
				se.Sess = se.InitSessionWithAWSProfile(b.Region,b.AWSProfile)
				// fmt.Println(se.Sess.Config.Credentials)
				instances := aws.ListInstances(se.Sess)
				var outputSheetName = b.AWSProfile+"-"+b.Region
				if len(instances) != 0 {
					tools.InfoLogger.Printf("Found %d Instances In %s : %s \n", len(instances), b.AWSProfile,b.Region)
					excel.SetHeaderLine("output.xlsx",outputSheetName, headerLine)
					excel.SetListRows("output.xlsx", outputSheetName, instances)
				} else {
					tools.InfoLogger.Printf("No EC2 Instnace In %s : %s \n", b.AWSProfile,b.Region)
				}
			}
		} else {
			defaultSess := aws.InitSession(*region)
			instances := aws.ListInstances(defaultSess)
			if len(instances) != 0 {
				tools.InfoLogger.Printf("Found %d Instances In %s \n", len(instances),*region)
				excel.SetHeaderLine("output.xlsx","Default", headerLine)
				excel.SetListRows("output.xlsx","Default", instances)
			} else {
				tools.InfoLogger.Printf("No EC2 Instnace In %s \n", *region)
			}
		}

	default:
		flag.Usage()
		tools.ErrorLogger.Fatalln("Illegal Method:", *method)
	}
}
