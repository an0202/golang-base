package aws

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-xray-sdk-go/xray"
	"golang.org/x/net/context/ctxhttp"
	// Importing the plugins enables collection of AWS resource information at runtime.
	// Every plugin should be imported after "github.com/aws/aws-xray-sdk-go/xray" library.
	//	_ "github.com/aws/aws-xray-sdk-go/plugins/beanstalk"
	//	_ "github.com/aws/aws-xray-sdk-go/plugins/ec2"
	//	_ "github.com/aws/aws-xray-sdk-go/plugins/ecs"
)

func init() {
	xray.Configure(xray.Config{
		DaemonAddr:     "10.250.101.190:2000", // default
		ServiceVersion: "1.2.3",
	})
}

func listQueue() {
	
	xray.AWS(dy)
}
// func getExample(ctx context.Context) ([]byte, error) {
// 	resp, err := ctxhttp.Get(ctx, xray.Client(nil), "http://www.baidu.com/")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	fmt.Println(resp.Body)
// 	return ioutil.ReadAll(resp.Body)
// }

// func main() {
// 	// Start a segment
// 	ctx, seg := xray.BeginSegment(context.Background(), "service-name")
// 	// Start a subsegment
// 	// subCtx, subSeg := xray.BeginSubsegment(ctx, "subsegment-name")
// 	// ...
// 	// Add metadata or annotation here if necessary
// 	// ...
// 	// subSeg.Close(nil)
// 	// Close the segment

// 	// Http Rquest
// 	getExample(ctx)

// 	seg.Close(nil)
// }

