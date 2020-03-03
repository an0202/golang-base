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

var DBInstanceHeadLine = []interface{}{"AccountId", "Region", "Name", "Type", "EndPoint", "Engine", "EngineVersion",
	"Port","SubnetGroup","AvailabilityZone","MultiAZ","Status","StorageType","BackupWindow","BackupPeriod","MaintenanceWindow",
	"ParameterGroups","SecurityGroups"}

var DBClusterHeadline = []interface{}{"AccountId", "Region", "Identifier", "Endpoint", "ReaderEndPoint", "EngineMode","Engine",
	"EngineVersion", "Port","MultiAZ","Status","MaintenanceWindow","BackupWindow","BackupPeriod", "ParameterGroups",
	"AvailabilityZone","SecurityGroups","Members"}

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

var NatGatewayHeadLine = []interface{}{"AccountId", "Region", "Name", "NatGatewayId", "VPCId", "SubNetId"}

var SGHeadLine = []interface{}{"GroupName", "VpcId", "GroupId", "Protocol", "Source", "FromPort", "ToPort"}

var AlarmHeadLine = []interface{}{"AccountId", "Region", "AlarmName", "NameSpace", "MetricName", "Actions", "Dimensions"}

var CLBHeadLine = []interface{}{"AccountId", "Region", "VPCId", "LoadBalancerName", "DNSName","Scheme", "AvailabilityZone",
	"SecurityGroups", "Listner","HealthCheck","Instances"}

var LBv2HeadLine = []interface{}{"AccountId", "Region", "VPCId", "LoadBalancerName", "DNSName","ARN","Type","Scheme",
	"AvailabilityZone", "SecurityGroups", "Listener","TargetGroups","Backends"}

var SNSHeadLine = []interface{}{"AccountId", "Region", "TopicName","Policy","ARN","Subscriptions"}

//var SQSHeadLine = []interface{}{"AccountId", "Region", "TopicName","Policy","ARN","Subscriptions"}

var S3HeadLine = []interface{}{"AccountId", "Region", "BucketName","ACL","Policy","CORS","LifeCycle","Versioning","WebSite"}

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
	if len(c.Operate) == 0 || len(c.Region) == 0 {
		tools.WarningLogger.Println("Missing Operation Or Region Instructions , Skip This Operate :", *c)
		return
	}
	c.init()
	c.outputFile = outputFile
	switch c.Operate {
	case "ListEC2":
		//c.OutputSheet = "EC2"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,EC2HeadLine)
		result := ListInstances(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListDBInstance":
		//c.OutputSheet = "RDS"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,DBInstanceHeadLine)
		result := ListDBInstances(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "Liv2DBCluster":
		excel.SetHeadLine(c.outputFile, c.OutputSheet,DBClusterHeadline)
		result := Listv2DBClusters(c.sess)
		excel.SetStructRows(c.outputFile, c.OutputSheet,result)
	case "ListECC":
		//c.OutputSheet = "ECC"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,ECCHeadLine)
		result := ListECCs(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListVolume":
		//c.OutputSheet = "EC2-Volume"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,VolumeHeadLine)
		result := ListVolumes(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListSnapshot":
		//c.OutputSheet = "EC2-Snapshot"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,SnapshotHeadLine)
		result := ListSnapshots(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListAMI":
		//c.OutputSheet = "EC2-AMI"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,AMIHeadLine)
		result := ListAMIs(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListEIP":
		//c.OutputSheet = "EC2-EIP"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,EIPHeadLine)
		result := ListEIPs(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListKeyPair":
		//c.OutputSheet = "EC2-KeyPairs"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,KeyPairHeadLine)
		result := ListKeyPairs(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListSubNet":
		//c.OutputSheet = "VPC-SubNet"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,SubnetHeadLine)
		result := ListSubNets(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListRouteTable":
		//c.OutputSheet = "VPC-RouteTable"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,RouteTableHeadLine)
		result := ListRouteTables(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListVPCPeering":
		excel.SetHeadLine(c.outputFile, c.OutputSheet,PeeringConnectionHeadLine)
		result := ListVPCPeering(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListNatGateway":
		excel.SetHeadLine(c.outputFile, c.OutputSheet,NatGatewayHeadLine)
		result := ListNatGateway(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "Liv2SG":
		//c.OutputSheet = "VPC-SG"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,SGHeadLine)
		result := Listv2SGs(c.sess)
		excel.SetStructRows(c.outputFile, c.OutputSheet,result)
	case "ListAlarm":
		//c.OutputSheet = "CloudWatch-Alarm"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,AlarmHeadLine)
		result := ListAlarms(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "ListCLB":
		//c.OutputSheet = "ELB"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,CLBHeadLine)
		result := ListCLBs(c.sess)
		excel.SetListRows(c.outputFile, c.OutputSheet,result)
	case "Liv2LBv2":
		//c.OutputSheet = "ELB"
		excel.SetHeadLine(c.outputFile, c.OutputSheet,LBv2HeadLine)
		result := Listv2LBv2s(c.sess)
		excel.SetStructRows(c.outputFile, c.OutputSheet,result)
	case "Liv2SNS":
		excel.SetHeadLine(c.outputFile, c.OutputSheet,SNSHeadLine)
		result := Listv2SNS(c.sess)
		excel.SetStructRows(c.outputFile, c.OutputSheet,result)
	//case "Liv2SQS":
	//	excel.SetHeadLine(c.outputFile, c.OutputSheet,SQSHeadLine)
	//	result := Listv2SQS(c.sess)
	//	excel.SetStructRows(c.outputFile, c.OutputSheet,result)
	case "Liv2S3":
		excel.SetHeadLine(c.outputFile, c.OutputSheet,S3HeadLine)
		result := Listv2S3(c.sess)
		excel.SetStructRows(c.outputFile, c.OutputSheet,result)
	default:
		tools.WarningLogger.Println("Unsupported Operate:", c.Operate)
	}
}

func (c *ExcelConfig) ReturnResources () (result [][]interface{}){
	if len(c.Operate) == 0 {
		tools.WarningLogger.Println("Missing Operation Instructions , Skip This Operate :", *c)
		return
	}
	c.init()
	switch c.Operate {
	case "ListEC2":
		c.HeadLine = EC2HeadLine
		result = ListInstances(c.sess)
	case "ListDBInstance":
		c.HeadLine = DBInstanceHeadLine
		result = ListDBInstances(c.sess)
	case "ListECC":
		c.HeadLine = ECCHeadLine
		result = ListECCs(c.sess)
	case "ListVolume":
		c.HeadLine = VolumeHeadLine
		result = ListVolumes(c.sess)
	case "ListSnapshot":
		c.HeadLine = SnapshotHeadLine
		result = ListSnapshots(c.sess)
	case "ListAMI":
		c.HeadLine = AMIHeadLine
		result = ListAMIs(c.sess)
	case "ListEIP":
		c.HeadLine = EIPHeadLine
		result = ListEIPs(c.sess)
	case "ListKeyPair":
		c.HeadLine = KeyPairHeadLine
		result = ListKeyPairs(c.sess)
	case "ListSubNet":
		c.HeadLine = SubnetHeadLine
		result = ListSubNets(c.sess)
	case "ListRouteTable":
		c.HeadLine = RouteTableHeadLine
		result = ListRouteTables(c.sess)
	case "ListVPCPeering":
		c.HeadLine = PeeringConnectionHeadLine
		result = ListVPCPeering(c.sess)
	case "ListNatGateway":
		c.HeadLine = NatGatewayHeadLine
		result = ListNatGateway(c.sess)
	case "ListAlarm":
		c.HeadLine = AlarmHeadLine
		result = ListAlarms(c.sess)
	case "ListCLB":
		c.HeadLine = CLBHeadLine
		result = ListCLBs(c.sess)
	default:
		tools.WarningLogger.Println("Unsupported Operate:", c.Operate)
	}
	return  result
}

func (c *ExcelConfig) ReturnResourcesV2 () (result []interface{}){
	if len(c.Operate) == 0 {
		tools.WarningLogger.Println("Missing Operation Instructions , Skip This Operate :", *c)
		return
	}
	c.init()
	switch c.Operate {
	case "Liv2LBv2":
		//c.OutputSheet = "ELB"
		c.HeadLine = LBv2HeadLine
		result = Listv2LBv2s(c.sess)
	case "Liv2SNS":
		c.HeadLine = SNSHeadLine
		result = Listv2SNS(c.sess)
	//case "Liv2SQS":
	//	excel.SetHeadLine(c.outputFile, c.OutputSheet,SQSHeadLine)
	//	result = Listv2SQS(c.sess)
	case "Liv2S3":
		c.HeadLine = S3HeadLine
		result = Listv2S3(c.sess)
	default:
		tools.WarningLogger.Println("Unsupported Operate:", c.Operate)
	}
	return result
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

