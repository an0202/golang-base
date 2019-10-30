package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func sqsDemo() {
	// Init a session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)
	if err != nil {
		fmt.Println(err)
	}

	// Create a SQS service client.
	// svc := sqs.New(sess)
	// // List the queues available in a given region.
	// result, err := svc.ListQueues(nil)
	// if err != nil {
	// 	fmt.Println("Error", err)
	// 	return
	// }

	// fmt.Println("Success")
	// // As these are pointers, printing them out directly would not be useful.
	// for i, urls := range result.QueueUrls {
	// 	// Avoid dereferencing a nil pointer.
	// 	if urls == nil {
	// 		continue
	// 	}
	// 	fmt.Printf("%d: %s\n", i, *urls)
	// }
	// URL to our queue
	svc := sqs.New(sess)
	qURL := "https://sqs.us-east-2.amazonaws.com/281525879386/xray-demo-sqs"

	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &qURL,
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   aws.Int64(60), // 60 seconds
		WaitTimeSeconds:     aws.Int64(0),
	})
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	if len(result.Messages) == 0 {
		fmt.Println("Received no messages")
		return
	}

	fmt.Printf("Success: %+v\n", result.Messages)
}
