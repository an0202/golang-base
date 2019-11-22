package main

import (
	"golang-base/aws"
)

func main() {
	// get policy
	//var headerline = []interface{}{"GroupName", "VpcId", "GroupId", "Protocol", "Source", "FromPort", "ToPort"}
	sess := aws.InitSession("cn-north-1")
	//a := aws.GetSGPolicys(sess)
	//excel.CreateFile("C:\\Users\\jie.an\\Desktop\\test.xlsx", "test")
	//excel.SetHeaderLine("C:\\Users\\jie.an\\Desktop\\test.xlsx", "test2", headerline)
	//excel.SetStructRows("C:\\Users\\jie.an\\Desktop\\test.xlsx", "test2", a)

	//// create tag for ec2
	//var InstanceIds = []string{"i-03177f7cffb8462be"}
	//var tagsmap = map[string]string{
	//	"InstanceID": "i-03177f7cffb8462be",
	//	"key10":      "compressnet",
	//	"key19":      "test19",
	//}
	//fmt.Println(tagsmap[key1])
	//b := aws.EC2InstanceMarshal(tagsmap)
	//fmt.Println(b.InstanceID)
	//aws.EC2CreateTags(sess, b)
	//a := excel.ReadTest("C:\\Users\\jie.an\\Desktop\\tags2.xlsx", "EC2")
	//for k, v := range a {
	//	fmt.Println(k, v)
	//	b := aws.EC2InstanceMarshal(v)
	//	aws.EC2DeleteTags(sess, b)
	//}
	aws.Testdescribeinstance(sess, "i-03177f7cffb8462be")
}
