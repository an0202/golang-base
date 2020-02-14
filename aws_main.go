package main

import "golang-base/cmd"

func main() {
	// get policy
	//var headerline = []interface{}{"GroupName", "VpcId", "GroupId", "Protocol", "Source", "FromPort", "ToPort"}
	//sess := aws.InitSession("cn-north-1")

	//a := aws.GetSGPolicys(sess)
	//excel.CreateFile("/mnt/c/Users/jie.an/Desktop/output.xlsx", "NULL")
	//excel.SetHeaderLine("/mnt/c/Users/jie.an/Desktop/output.xlsx", "SecurityGroup", headerline)
	//excel.SetStructRows("/mnt/c/Users/jie.an/Desktop/output.xlsx", "SecurityGroup", a)
	// Get Instances
	//elasticachlist := aws.ListElastiCache(sess)
	//excel.CreateFile("/mnt/c/Users/jie.an/Desktop/output.xlsx", "EC2")
	//excel.SetHeaderLine("/mnt/c/Users/jie.an/Desktop/output.xlsx", "EC2", headerline)
	//excel.SetListRows("/mnt/c/Users/jie.an/Desktop/output.xlsx", "EC2", elasticachlist)
	//// 5. Create Tag for ec2
	cmd.EC2addTags()
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
