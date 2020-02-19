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
	excelFile = flag.String("f", "query.xlsx", "Source ExcelFile To Be Processed")
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
		excel.CreateFile("output.xlsx")
		var headerLine = []interface{}{"AccountId", "Region", "Name", "InstanceId", "InstanceType", "Platform", "State",
			"VPCId","Role","SubnetId","KeyPair","SecurityGroups","PrivateIP","PublicIP","Tags"}
		// Total sheet , rowsNum is position for Total sheet written data
		excel.SetHeadLine("output.xlsx","Total", headerLine)
		rowsNum := 1
		// Get EC2 Instances From Excel (Excel Contains AWS_PROFILE And Region)
		if *excelFile != "" {
			se := new(aws.Session)
			a := excel.ReadToMaps(*excelFile, *sheetName)
			for _, v := range a {
				//fmt.Println("RowsNumis:      ",rowsNum)
				b := aws.EC2InstanceMarshal(v)
				// create a new session
				se.Sess = se.InitSessionWithAWSProfile(b.Region,b.AWSProfile)
				// fmt.Println(se.Sess.Config.Credentials)
				instances := aws.ListInstances(se.Sess)
				var outputSheetName = b.AWSProfile+"-"+b.Region
				if len(instances) != 0 {
					tools.InfoLogger.Printf("Found %d Instances In %s : %s \n", len(instances), b.AWSProfile,b.Region)
					excel.SetHeadLine("output.xlsx",outputSheetName, headerLine)
					excel.SetListRows("output.xlsx", outputSheetName, instances)
					// Write summary data to Total sheet
					excel.SetListRowsV2("output.xlsx","Total","A",rowsNum+1,instances)
					rowsNum += len(instances)
				} else {
					tools.InfoLogger.Printf("No EC2 Instnace In %s : %s \n", b.AWSProfile,b.Region)
				}
			}
		} else {
			// Does not read data from excel (used for export single account ec2 list)
			defaultSess := aws.InitSession(*region)
			instances := aws.ListInstances(defaultSess)
			if len(instances) != 0 {
				tools.InfoLogger.Printf("Found %d Instances In %s \n", len(instances),*region)
				excel.SetHeadLine("output.xlsx","Total", headerLine)
				excel.SetListRows("output.xlsx","Total", instances)
			} else {
				tools.InfoLogger.Printf("No EC2 Instnace In %s \n", *region)
			}
		}

	default:
		excel.CreateFile("output.xlsx")
		flag.Usage()
		tools.ErrorLogger.Fatalln("Illegal Method:", *method)
	}
}
