package main

import "golang-base/excel"

/*GCC
export GOOS="windows"
go build -o phoneix-windows.exe

export GOOS="linux"
go build -o phoneix-linux

export GOOS="darwin"
go build -o phoneix-darwin
*/
func main() {
	//echo "done!"
	////1. get policy
	//var headerline = []interface{}{"GroupName", "VpcId", "GroupId", "Protocol", "Source", "FromPort", "ToPort"}
	//sess := aws.InitSession("cn-north-1")
	////a := aws.GetSGPolicys(sess)
	//cmd.SamsungBillFilter2()
	var filePath = "monitorData.xlsx"
	var cpuAvgHeadLine = []interface{}{"ServerName", "CPU_AVG_PERCENT"}
	excel.CreateFile(filePath)
	excel.SetHeadLine(filePath, "CPU_AVG", cpuAvgHeadLine)
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
