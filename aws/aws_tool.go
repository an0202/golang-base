/**
 * @Author: jie.an
 * @Description:
 * @File:  aws_tool.go
 * @Version: 1.0.0
 * @Date: 2020/02/17 16:41
 */
package aws

import (
	"golang-base/excel"
	"golang-base/tools"
	"regexp"
)

// Excel HeadLine For AWS Resource
var EC2HeadLine = []interface{}{"AccountId", "Region", "Name", "InstanceId", "InstanceType", "Platform", "State",
	"VPCId","Role","SubnetId","PrivateIp","PublicIp","KeyPair","SecurityGroups","Tags"}

var DBHeadLine = []interface{}{"AccountId", "Region", "Name", "Type", "EndPoint", "Engine", "EngineVersion",
	"Port","SubnetGroup","AvailabilityZone","MultiAZ","Status","StorageType","BackupWindow","BackupPeriod","MaintenanceWindow",
	"ParameterGroups","SecurityGroups"}

var ECCHeadLine = []interface{}{"AccountId", "Region", "CacheClusterId", "CacheNodesNumber", "CacheNodeType", "Engine",
	"EngineVersion", "CacheSubnetGroup","MaintenanceWindow","SnapshotRetention","SecurityGroups"}

var KeyPairHeadLine = []interface{}{"AccountId", "Region", "KeyName", "Fingerprint"}

var SnapshotHeadLine = []interface{}{"AccountId", "Region", "SnapshotId", "VolumeId", "Description", "State"}

var VolumeHeadLine = []interface{}{"AccountId", "Region", "VolumeId", "State", "Type", "Size", "AvailabilityZone"}

var AMIHeadLine = []interface{}{"AccountId", "Region", "ImageId", "Name", "State"}

var EIPHeadLine = []interface{}{"AccountId", "Region", "Name", "PublicIp"}

var SubnetHeadLine = []interface{}{"AccountId", "Region", "Name", "SubnetId", "VPCId", "CidrBlock", "AvailabilityZone",
	"DefaultForAz"}

var RouteTableHeadLine = []interface{}{"AccountId", "Region", "Name", "RouteTableId", "VPCId", "Assosications", "Routes"}

var SGHeadLine = []interface{}{"GroupName", "VpcId", "GroupId", "Protocol", "Source", "FromPort", "ToPort"}

var AlarmHeadLine = []interface{}{"AccountId", "Region", "AlarmName", "NameSpace", "MetricName", "Actions", "Dimensions"}

var LBHeadLine = []interface{}{"AccountId", "Region", "VPCId", "LoadBalancerName", "Scheme", "AvailabilityZone", "SecurityGroups",
	"Listner","HealthCheck","Instances"}

type ExcelConfig struct {
	AWSProfile          string
	Region              string
	Operate   			string
	Comment             string
	outputFile          string
	outputSheet		    string
}

func (c *ExcelConfig) init() {
	// todo : config check
	return
}

func (c *ExcelConfig) Do(outputFile string) {
	//c.init()
	c.outputFile = outputFile
	sess := InitSession(c.Region)
	switch c.Operate {
	case "ListInstances":
		c.outputSheet = "EC2"
		excel.SetHeadLine(c.outputFile, c.outputSheet,EC2HeadLine)
		result := ListInstances(sess)
		excel.SetListRows(c.outputFile, c.outputSheet,result)
	case "ListDBs":
		c.outputSheet = "RDS"
		excel.SetHeadLine(c.outputFile, c.outputSheet,DBHeadLine)
		result := ListDBs(sess)
		excel.SetListRows(c.outputFile, c.outputSheet,result)
	case "ListECCs":
		c.outputSheet = "ECC"
		excel.SetHeadLine(c.outputFile, c.outputSheet,ECCHeadLine)
		result := ListECCs(sess)
		excel.SetListRows(c.outputFile, c.outputSheet,result)
	case "ListVolumes":
		c.outputSheet = "EC2-Volume"
		excel.SetHeadLine(c.outputFile, c.outputSheet,VolumeHeadLine)
		result := ListVolumes(sess)
		excel.SetListRows(c.outputFile, c.outputSheet,result)
	case "ListSnapshots":
		c.outputSheet = "EC2-Snapshot"
		excel.SetHeadLine(c.outputFile, c.outputSheet,SnapshotHeadLine)
		result := ListSnapshots(sess)
		excel.SetListRows(c.outputFile, c.outputSheet,result)
	case "ListAMIs":
		c.outputSheet = "EC2-AMI"
		excel.SetHeadLine(c.outputFile, c.outputSheet,AMIHeadLine)
		result := ListAMIs(sess)
		excel.SetListRows(c.outputFile, c.outputSheet,result)
	case "ListEIPs":
		c.outputSheet = "EC2-EIP"
		excel.SetHeadLine(c.outputFile, c.outputSheet,EIPHeadLine)
		result := ListEIPs(sess)
		excel.SetListRows(c.outputFile, c.outputSheet,result)
	case "ListKeyPairs":
		c.outputSheet = "EC2-KeyPairs"
		excel.SetHeadLine(c.outputFile, c.outputSheet,KeyPairHeadLine)
		result := ListKeyPairs(sess)
		excel.SetListRows(c.outputFile, c.outputSheet,result)
	case "ListSubNets":
		c.outputSheet = "VPC-SubNet"
		excel.SetHeadLine(c.outputFile, c.outputSheet,SubnetHeadLine)
		result := ListSubNets(sess)
		excel.SetListRows(c.outputFile, c.outputSheet,result)
	case "ListRouteTables":
		c.outputSheet = "VPC-RouteTable"
		excel.SetHeadLine(c.outputFile, c.outputSheet,RouteTableHeadLine)
		result := ListRouteTables(sess)
		excel.SetListRows(c.outputFile, c.outputSheet,result)
	case "ListSGs":
		c.outputSheet = "VPC-SG"
		excel.SetHeadLine(c.outputFile, c.outputSheet,SGHeadLine)
		result := ListSGs(sess)
		excel.SetStructRows(c.outputFile, c.outputSheet,result)
	case "ListAlarms":
		c.outputSheet = "CloudWatch-Alarm"
		excel.SetHeadLine(c.outputFile, c.outputSheet,AlarmHeadLine)
		result := ListAlarms(sess)
		excel.SetListRows(c.outputFile, c.outputSheet,result)
	case "ListLBs":
		c.outputSheet = "ELB"
		excel.SetHeadLine(c.outputFile, c.outputSheet,LBHeadLine)
		result := ListLBs(sess)
		excel.SetListRows(c.outputFile, c.outputSheet,result)
	default:
		tools.WarningLogger.Println("Unsupported Operate:", c.Operate)
	}
}

//GetARNDetail return a map contains ARN base information
//map[accountId:123456789012 region: resource:user/Development/product_1234/* service:iam]
//https://blog.csdn.net/butterfly5211314/article/details/82532970
//https://golang.org/pkg/regexp/syntax/
func GetARNDetail(arn string) (arnMap map[string]string) {
	re := regexp.MustCompile(`^arn:(?:aws|aws-cn|\s):(?P<service>[a-z\s]+):(?P<region>[a-z0-9\-\s]+):(?P<accountId>[0-9\s]+):(?P<resource>.*)`)
	matched := re.MatchString(arn)
	arnMap = make(map[string]string)
	if matched == false {
		tools.ErrorLogger.Fatalln(arn," Is Not A AWS ARN.")
	} else {
		groupNames := re.SubexpNames()
		match := re.FindStringSubmatch(arn)
		//for k, v := range match {
		//	fmt.Println(k,v)
		//}
		for i, name := range groupNames {
			if i != 0 && name != "" {
				arnMap[name] = match[i]
			}
		}
	}
	return arnMap
}

// retype from excel
func ExcelConfigMarshal(excelConfig map[string]string) (config ExcelConfig) {
	for k, v := range excelConfig {
		switch k {
		case "AWS_PROFILE":
			config.AWSProfile = excelConfig["AWS_PROFILE"]
		case "Region":
			config.Region = excelConfig["Region"]
		case "Operate":
			config.Operate = excelConfig["Operate"]
		case "Comment":
			config.Comment = excelConfig["Comment"]
		default:
			tools.WarningLogger.Printf("Unsupported Configuration: %s - %s",k ,v)
		}
	}
	return config
}

