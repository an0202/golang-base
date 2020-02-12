package main

import (
	"golang-base/aws"
	"golang-base/excel"
)

func main() {
	// get policy
	var headerline = []interface{}{"GroupName", "VpcId", "GroupId", "Protocol", "Source", "FromPort", "ToPort"}
	sess := aws.InitSession("cn-north-1")

	//a := aws.GetSGPolicys(sess)
	//excel.CreateFile("/mnt/c/Users/jie.an/Desktop/output.xlsx", "NULL")
	//excel.SetHeaderLine("/mnt/c/Users/jie.an/Desktop/output.xlsx", "SecurityGroup", headerline)
	//excel.SetStructRows("/mnt/c/Users/jie.an/Desktop/output.xlsx", "SecurityGroup", a)
	// Get Instances
	elasticachlist := aws.ListElastiCache(sess)
	excel.CreateFile("/mnt/c/Users/jie.an/Desktop/output.xlsx", "EC2")
	excel.SetHeaderLine("/mnt/c/Users/jie.an/Desktop/output.xlsx", "EC2", headerline)
	excel.SetListRows("/mnt/c/Users/jie.an/Desktop/output.xlsx", "EC2", elasticachlist)
	//// create tag for ec2
	//var InstanceIds = []string{"i-03177f7cffb8462be"}
	//var tagsmap = map[string]string{
	//	"InstanceID": "i-03177f7cffb8462be",
	//	"key10":      "compressnet",
	//	"key19":      "test19",
	//}
	//fmt.Println(tagsmap[key1])
	//a := excel.ReadTest("C:\\Users\\jie.an\\Desktop\\tags2.xlsx", "EC2")
	//cmd.EBSaddtags()

	//CreateImage
	//aws.CreateImage(sess, "i-03177f7cffb8462be")
	//
	// Get snapshots
	//accountid := aws.GetAccountId(sess)
	//aws.ListSnapshots(sess, accountid)
}
