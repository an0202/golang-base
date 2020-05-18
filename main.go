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
	//cmd.SamsungBillFilter2()
	var filePath = "project" + "-Galaxy-" + "-Export.xlsx"
	var projectHeadline = []interface{}{"事件ID", "公司名称", "状态", "影响范围", "紧急程度", "优先级", "环境",
		"类别", "一级分类", "二级分类", "三级分类", "标题", "描述", "处理组", "处理人", "提交人", "提交人所属部门",
		"创建事件", "响应时间", "响应是否超时", "响应超时原因", "预期解决时间", "实际解决时间", "处理是否超时", "根本原因",
		"挂起原因", "重新打开次数", "重新打开原因", "解决方案", "满意度", "评价内容", "忽略", "处理时效"}
	excel.CreateFile(filePath)
	excel.SetHeadLine(filePath, "Galaxy-"+"project", projectHeadline)
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
