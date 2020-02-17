package main

import (
    "fmt"
    "golang-base/aws"
)

func main() {
	////1. get policy
	//var headerline = []interface{}{"GroupName", "VpcId", "GroupId", "Protocol", "Source", "FromPort", "ToPort"}
	//sess := aws.InitSession("cn-north-1")
	//a := aws.GetSGPolicys(sess)
	//excel.CreateFile("/mnt/c/Users/jie.an/Desktop/output.xlsx", "NULL")
	//excel.SetHeaderLine("/mnt/c/Users/jie.an/Desktop/output.xlsx", "SecurityGroup", headerline)
	//excel.SetStructRows("/mnt/c/Users/jie.an/Desktop/output.xlsx", "SecurityGroup", a)
	////2. Get elasticach
	//elasticachlist := aws.ListElastiCache(sess)
	//excel.CreateFile("/mnt/c/Users/jie.an/Desktop/output.xlsx", "ECC")
	//excel.SetHeaderLine("/mnt/c/Users/jie.an/Desktop/output.xlsx", "ECC", headerline)
	//excel.SetListRows("/mnt/c/Users/jie.an/Desktop/output.xlsx", "ECC", elasticachlist)
	////3. Create Tag for ec2
	//cmd.EC2Tags()
	//fmt.Println(tagsmap[key1])
	//a := excel.ReadTest("C:\\Users\\jie.an\\Desktop\\tags2.xlsx", "EC2")
	//cmd.EBSaddtags()
	////4. CreateImage
	//aws.CreateImage(sess, "i-03177f7cffb8462be")
	//
	////5.Get snapshots
	//accountid := aws.GetAccountId(sess)
	//aws.ListSnapshots(sess, accountid)
	////6. Get EC2
	//var headerline = []interface{}{"AccountId", "Region", "Name", "InstanceId", "InstanceType", "Platform", "State",
	//"VPCId","Role","SubnetId","KeyPair","SecurityGroups","Tags"}
	//sess := aws.InitSession("cn-northwest-1")
	//ec2 := aws.ListInstances(sess)
	//excel.CreateFile("output.xlsx", "EC2")
	//excel.SetHeaderLine("output.xlsx", "EC2", headerline)
	//excel.SetListRows("output.xlsx", "EC2", ec2)
    //cmd.EC2()
    sess := aws.InitSession("cn-north-1")
    a := aws.ListAlarms(sess)
    for _, alarm := range(a) {
       fmt.Println(alarm)
       fmt.Println("------")
    }
}



////Functional Options
//
//var defaultStuffClient = stuffClient{
//   retries: 3,
//   timeout: 2,
//}
//type StuffClientOption func(*stuffClient)
//
//func WithRetries(r int) StuffClientOption {
//   return func(o *stuffClient) {
//       o.retries = r
//   }
//}
//func WithTimeout(t int) StuffClientOption {
//   return func(o *stuffClient) {
//       o.timeout = t
//   }
//}
//type StuffClient interface {
//   DoStuff() error
//}
//type stuffClient struct {
//   conn    Connection
//   timeout int
//   retries int
//}
//type Connection struct{}
//func NewStuffClient(conn Connection, opts ...StuffClientOption) StuffClient {
//   client := defaultStuffClient
//   for _, o := range opts {
//       o(&client)
//   }
//
//   client.conn = conn
//   return client
//}
//
//func (c stuffClient) DoStuff() error {
//   return nil
//}
//
//// The Test
//
//func main() {
//   x := NewStuffClient(Connection{})
//   fmt.Println(x) // prints &{{} 2 3}
//
//   x = NewStuffClient(Connection{}, WithRetries(1))
//   fmt.Println(x) // prints &{{} 2 1}
//}

