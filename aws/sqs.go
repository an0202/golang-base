/**
 * @Author: jie.an
 * @Description:
 * @File:  sqs.go
 * @Version: 1.0.0
 * @Date: 2020/02/24 19:32
 */
package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sqs"
	"golang-base/tools"
)

type Queue struct {
	AccountId string
	Region    string
	Name      string
	Policy    string
	URL       string
}

//List SQS , not test yet , some atts(DelaySeconds,MaximumMessageSize...) need  be  add
func Listv2SQS(se Session) (SQSList []interface{}) {
	// Create an sqs service client.
	svc := sqs.New(se.Sess)
	// Get sns topics
	output, err := svc.ListQueues(&sqs.ListQueuesInput{})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				tools.ErrorLogger.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			tools.ErrorLogger.Println(err.Error())
		}
		return
	}
	for _, queue := range output.QueueUrls {
		QE := new(Queue)
		QE.AccountId = se.AccountId
		QE.Region = se.UsedRegion
		QE.URL = *queue
		atts := QE.GetQueueAttributes(se, QE.URL)
		QE.Policy = *atts["Policy"]
		QE.Name = GetARNDetail(*atts["QueueArn"])["resource"]
		SQSList = append(SQSList, *QE)
	}
	return SQSList
}

func (qe *Queue) GetQueueAttributes(se Session, QueueURL string) map[string]*string {
	// Create an sqs service client.
	svc := sqs.New(se.Sess)
	// Get queue attributes
	output, err := svc.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(QueueURL),
		AttributeNames: []*string{aws.String("ALL")},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				tools.ErrorLogger.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			tools.ErrorLogger.Println(err.Error())
		}
		return nil
	}
	return output.Attributes
}
