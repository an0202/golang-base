package aws

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-xray-sdk-go/xray"
	"golang.org/x/net/context/ctxhttp"
)

// func init() {
// 	xray.Configure(xray.Config{
// 		DaemonAddr:     "10.250.101.190:2000", // default
// 		ServiceVersion: "1.2.3",
// 	})
// }

func getExample(ctx context.Context) ([]byte, error) {
	resp, err := ctxhttp.Get(ctx, xray.Client(nil), "http://baiduaabb.com/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println(resp.Body)
	return ioutil.ReadAll(resp.Body)
}

func xrayDemo() {
	// Init a session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)
	if err != nil {
		fmt.Println(err)
	}

	// Create a SQS service client.
	svc := sqs.New(sess)

	xray.AWS(svc.Client)

	// Start a segment
	ctx, seg := xray.BeginSegment(context.Background(), "xray-sqs")
	defer seg.Close(nil)
	result, err := svc.ListQueuesWithContext(ctx, nil)
	// List queue
	// result, err := svc.ListQueues(nil)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	for i, urls := range result.QueueUrls {
		// Avoid dereferencing a nil pointer.
		if urls == nil {
			continue
		}
		fmt.Printf("%d: %s\n", i, *urls)
	}

	getExample(ctx)

	// 	// Http Rquest
	// 	getExample(ctx)

	// 	seg.Close(nil)

	// 	// Get Message
	// 	qURL := "https://sqs.us-east-2.amazonaws.com/281525879386/xray-demo-sqs"
	// 	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
	// 		AttributeNames: []*string{
	// 			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
	// 		},
	// 		MessageAttributeNames: []*string{
	// 			aws.String(sqs.QueueAttributeNameAll),
	// 		},
	// 		QueueUrl:            &qURL,
	// 		MaxNumberOfMessages: aws.Int64(10),
	// 		VisibilityTimeout:   aws.Int64(60), // 60 seconds
	// 		WaitTimeSeconds:     aws.Int64(0),
	// 	})
	// 	if err != nil {
	// 		fmt.Println("Error", err)
	// 		return
	// 	}
	// 	if len(result.Messages) == 0 {
	// 		fmt.Println("Received no messages")
	// 		return
	// 	}

	// 	fmt.Printf("Success: %+v\n", result.Messages)
}
