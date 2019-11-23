package main

import (
	"golang-base/aws"
	"golang-base/excel"
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
	//a := excel.ReadTest("C:\\Users\\jie.an\\Desktop\\tags2.xlsx", "EC2")
	a := excel.ReadTest("/Users/anjie/Desktop/tags2.xlsx", "EC2")
	for _, v := range a {
		b := aws.EC2InstanceMarshal(v)
		aws.EC2DeleteTags(sess, b)
	}
}
