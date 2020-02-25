/**
 * @Author: jie.an
 * @Description:
 * @File:  sns.go
 * @Version: 1.0.0
 * @Date: 2020/02/24 15:42
 */
package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sns"
	"golang-base/tools"
)

type Topic struct {
	AccountId string
	Region    string
	Name 	  string
	Policy    string
	ARN       string
	Subscription  map[string]string
}

//List SNS
func Listv2SNS(se Session) (SNSList []interface{}) {
	// Create an sns service client.
	svc := sns.New(se.Sess)
	// Get sns topics
	output, err := svc.ListTopics(&sns.ListTopicsInput{
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
		return
	}
	TP := new(Topic)
	for _, topic := range output.Topics {
		if len(output.Topics) >= 100 {
			tools.WarningLogger.Println("Number Of Topics > 100 , Data May Missing.")
		}
		TP.AccountId = se.AccountId
		TP.Region = se.UsedRegion
		TP.ARN = *topic.TopicArn
		atts := TP.GetTopicAttributes(se, TP.ARN)
		TP.Name = GetARNDetail(TP.ARN)["resource"]
		TP.Policy = *atts["Policy"]
		TP.Subscription = TP.ListSubscriptions(se, TP.ARN)
		SNSList = append(SNSList, *TP)
	}
	return SNSList
}

func (tp *Topic) GetTopicAttributes(se Session,TopicARN string) map[string]*string {
	// Create an sns service client.
	svc := sns.New(se.Sess)
	// Get topic attributes
	output, err := svc.GetTopicAttributes(&sns.GetTopicAttributesInput{
		TopicArn: aws.String(TopicARN),
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

func (tp *Topic) ListSubscriptions(se Session,TopicARN string) (Subscriptions map[string]string) {
	// Create an sns service client.
	svc := sns.New(se.Sess)
	// Get topic attributes
	output, err := svc.ListSubscriptionsByTopic(&sns.ListSubscriptionsByTopicInput{
		TopicArn: aws.String(TopicARN),
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
	if len(output.Subscriptions) == 100 {
		tools.InfoLogger.Printf("Subscription For Topic %s > 100 , Data May Missing.",TopicARN)
	}
	//{ endpoint : protocol }
	Subscriptions = make(map[string]string)
	if len(output.Subscriptions) == 0 {
		Subscriptions["N/A"] = "N/A"
	} else {
		for _, subscription := range output.Subscriptions {
			Subscriptions[*subscription.Endpoint] = *subscription.Protocol
		}
	}
	return Subscriptions
}