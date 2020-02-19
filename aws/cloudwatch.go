/**
 * @Author: jie.an
 * @Description:
 * @File:  cloudwatch.go
 * @Version: 1.0.0
 * @Date: 2020/02/17 14:57
 */
package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

//List Alarms
func ListAlarms(sess *session.Session) (AlarmList [][]interface{}) {
	var maxResults = 100
	var token string
	var alarms [][]interface{}
	var nextToken = "default"
	for nextToken != "" {
		//fmt.Println("Start Loop With Token:", token)
		alarms, nextToken = listAlarms(sess, token, maxResults)
		for _, snapshot := range alarms {
			AlarmList = append(AlarmList, snapshot)
		}
		if len(alarms) == maxResults && nextToken != "" {
			alarms, nextToken = listAlarms(sess, nextToken, maxResults)
			for _, alarm := range alarms {
				AlarmList = append(AlarmList, alarm)
			}
			//fmt.Println("Generated New NextToken:      ",nextToken)
			token = nextToken
		} else {
			nextToken = ""
		}
	}
	return AlarmList
}
//List Alarm Internal
func listAlarms(sess *session.Session, token string,maxResults int) (AlarmList [][]interface{},nextToken string) {
	// Create an cloudwatch service client.
	svc := cloudwatch.New(sess)
	// Get alarms
	output, err := svc.DescribeAlarms(&cloudwatch.DescribeAlarmsInput{
		MaxRecords: aws.Int64(int64(maxResults)),
		NextToken: aws.String(token),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	for _, alarm := range output.MetricAlarms {
		var Alarm []interface{}
		var actions,dimensions []string
		if len(alarm.AlarmActions) == 0 {
			actions = append(actions, "N/A ")
		} else {
			for _, action := range alarm.AlarmActions {
				actions = append(actions, *action)
			}
		}
		//handle dimensions
		if len(alarm.Dimensions) == 0 {
			dimensions = append(dimensions, "N/A ")
		} else {
			for _, dimension := range alarm.Dimensions {
				dimensions = append(dimensions, *dimension.Name+":"+*dimension.Value+" ")
			}
		}
		// handel accountid
		arnMap := GetARNDetail(*alarm.AlarmArn)
		accountId := arnMap["accountId"]
		Alarm = append(Alarm,accountId,*sess.Config.Region, *alarm.AlarmName, *alarm.Namespace,
			*alarm.MetricName, actions, dimensions)
		//fmt.Println(Alarm)
		AlarmList = append(AlarmList, Alarm)
	}
	if output.NextToken == nil {
		nextToken = ""
	} else {
		nextToken = *output.NextToken
	}
	return AlarmList, nextToken
}
