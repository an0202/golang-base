package main

import "golang-base/aws"

func main() {
	// get policy
	//var headerline = []interface{}{"GroupName", "VpcId", "GroupId", "Protocol", "Source", "FromPort", "ToPort"}
	sess := aws.InitSession("cn-north-1")
	//a := aws.GetSGPolicys(sess)
	//excel.CreateFile("C:\\Users\\jie.an\\Desktop\\test.xlsx", "test")
	//excel.SetHeaderLine("C:\\Users\\jie.an\\Desktop\\test.xlsx", "test", headerline)
	//excel.SetStructRows("C:\\Users\\jie.an\\Desktop\\test.xlsx", "test", a)

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
	aws.CreateImage(sess, "i-03177f7cffb8462be")
}
