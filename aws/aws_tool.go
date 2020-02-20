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

//Region Abbreviation
var RegionAbb = map[string]string{
	"cn-north-1": "CNN1",
	"cn-northwest-1": "CNW1",
}


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

var VolumeHeadLine = []interface{}{"AccountId", "Region", "Name","VolumeId","AttachedInstance","State", "Type", "Size",
	"AvailabilityZone"}

var AMIHeadLine = []interface{}{"AccountId", "Region", "ImageId", "Name", "State"}

var EIPHeadLine = []interface{}{"AccountId", "Region", "Name", "PublicIp"}

var SubnetHeadLine = []interface{}{"AccountId", "Region", "Name", "SubnetId", "VPCId", "CidrBlock", "AvailabilityZone",
	"DefaultForAz"}

var RouteTableHeadLine = []interface{}{"AccountId", "Region", "Name", "RouteTableId", "VPCId", "Assosications", "Routes"}

var PeeringConnectionHeadLine = []interface{}{"AccountId", "Region", "Name", "PeeringId", "RequesterInfo", "AccepterInfo"}

var SGHeadLine = []interface{}{"GroupName", "VpcId", "GroupId", "Protocol", "Source", "FromPort", "ToPort"}

var AlarmHeadLine = []interface{}{"AccountId", "Region", "AlarmName", "NameSpace", "MetricName", "Actions", "Dimensions"}

var LBHeadLine = []interface{}{"AccountId", "Region", "VPCId", "LoadBalancerName", "DNSName","Scheme", "AvailabilityZone",
	"SecurityGroups", "Listner","HealthCheck","Instances"}

type ExcelConfig struct {
	AWSProfile          string
	Region              string
	Operate   			string
	Comment             string
	outputFile          string
	OutputSheet		    string
	HeadLine            []interface{}
	sess                Session
	AccountId           string
}

func (c *ExcelConfig) init() {
	se := c.sess.InitSessionWithAWSProfile(c.Region,c.AWSProfile)
	c.Region = se.UsedRegion
	if regionAbb, ok := RegionAbb[se.UsedRegion]; ok{
		c.Region = regionAbb
	}else{
		c.Region = se.UsedRegion
	}
	c.AWSProfile = se.UsedAwsProfile
	c.AccountId = se.AccountId
	c.OutputSheet = c.Operate[4:]+"-"+se.AccountId+"-"+c.Region
	tools.InfoLogger.Printf("%s From Account: %s (%s) \n",c.Operate,c.AccountId,c.Region)
}

func (c *ExcelConfig) Do(outputFile string) {
	if len(c.Operate) == 0 {
		return
	}
	c.init()
	c.outputFile = outputFile
	switch c.Operate {
	case "ListInstances":
		//c.OutputSheet = "EC2"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,EC2HeadLine)
		result := ListInstances(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListDBs":
		//c.OutputSheet = "RDS"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,DBHeadLine)
		result := ListDBs(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListECCs":
		//c.OutputSheet = "ECC"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,ECCHeadLine)
		result := ListECCs(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListVolumes":
		//c.OutputSheet = "EC2-Volume"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,VolumeHeadLine)
		result := ListVolumes(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListSnapshots":
		//c.OutputSheet = "EC2-Snapshot"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,SnapshotHeadLine)
		result := ListSnapshots(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListAMIs":
		//c.OutputSheet = "EC2-AMI"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,AMIHeadLine)
		result := ListAMIs(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListEIPs":
		//c.OutputSheet = "EC2-EIP"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,EIPHeadLine)
		result := ListEIPs(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListKeyPairs":
		//c.OutputSheet = "EC2-KeyPairs"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,KeyPairHeadLine)
		result := ListKeyPairs(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListSubNets":
		//c.OutputSheet = "VPC-SubNet"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,SubnetHeadLine)
		result := ListSubNets(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListRouteTables":
		//c.OutputSheet = "VPC-RouteTable"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,RouteTableHeadLine)
		result := ListRouteTables(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListVPCPeering":
		excel.SetHeadLine(c.outputFile, c.OutputSheet,PeeringConnectionHeadLine)
		result := ListRouteTables(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListSGs":
		//c.OutputSheet = "VPC-SG"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,SGHeadLine)
		result := ListSGs(c.sess)
		excel.SetStructRows(c.outputFile, c.OutputSheet,result)
	case "ListAlarms":
		//c.OutputSheet = "CloudWatch-Alarm"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,AlarmHeadLine)
		result := ListAlarms(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListLBs":
		//c.OutputSheet = "ELB"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,LBHeadLine)
		result := ListLBs(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	default:
		tools.WarningLogger.Println("Unsupported Operate:", c.Operate)
	}
}

func (c *ExcelConfig) ReturnResources () (result [][]interface{}){
	if len(c.Operate) == 0 {
		return
	}
	c.init()
	switch c.Operate {
	case "ListInstances":
		c.HeadLine = EC2HeadLine
		result = ListInstances(c.sess)
	case "ListDBs":
		c.HeadLine = DBHeadLine
		result = ListDBs(c.sess)
	case "ListECCs":
		c.HeadLine = ECCHeadLine
		result = ListECCs(c.sess)
	case "ListVolumes":
		c.HeadLine = VolumeHeadLine
		result = ListVolumes(c.sess)
	case "ListSnapshots":
		c.HeadLine = SnapshotHeadLine
		result = ListSnapshots(c.sess)
	case "ListAMIs":
		c.HeadLine = AMIHeadLine
		result = ListAMIs(c.sess)
	case "ListEIPs":
		c.HeadLine = EIPHeadLine
		result = ListEIPs(c.sess)
	case "ListKeyPairs":
		c.HeadLine = KeyPairHeadLine
		result = ListKeyPairs(c.sess)
	case "ListSubNets":
		c.HeadLine = SubnetHeadLine
		result = ListSubNets(c.sess)
	case "ListRouteTables":
		c.HeadLine = RouteTableHeadLine
		result = ListRouteTables(c.sess)
	case "ListVPCPeering":
		c.HeadLine = PeeringConnectionHeadLine
		result = ListRouteTables(c.sess)
	case "ListAlarms":
		c.HeadLine = AlarmHeadLine
		result = ListAlarms(c.sess)
	case "ListLBs":
		c.HeadLine = LBHeadLine
		result = ListLBs(c.sess)
	default:
		tools.WarningLogger.Println("Unsupported Operate:", c.Operate)
	}
	return  result
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

