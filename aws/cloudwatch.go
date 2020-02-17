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
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

//type EC2AlarmDetail struct {
//	AlarmID          string
//	Region              string
//	AWSProfile			string
//	Dimensions                map[string]string
//	BlockDeviceMappings []EBSDetail
//}
//
//type EBSDetail struct {
//	VolumeId string
//	Status   string
//}

//List Alarms
func ListAlarms(sess *session.Session) (AlarmList [][]interface{}) {
	// Create an cloudwatch service client.
	svc := cloudwatch.New(sess)
	// Set max records
	output, err := svc.DescribeAlarms(&cloudwatch.DescribeAlarmsInput{
		//todo next token for list > 100
		MaxRecords: aws.Int64(100),
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
	}
	for _, alarm := range output.MetricAlarms {
		var Alarm []interface{}
		//var platform, rolearn, alarmname ,keypair, publicip string
		//if alarm.Platform != nil {
		//	platform = *alarm.Platform
		//} else {
		//	platform = "linux"
		//}
		//if alarm.IamAlarmProfile == nil {
		//	rolearn = "N/A"
		//} else {
		//	rolearn = *alarm.IamAlarmProfile.Arn
		//}
		//if alarm.KeyName == nil {
		//	keypair = "N/A"
		//} else {
		//	keypair = *alarm.KeyName
		//}
		//if alarm.PublicIpAddress == nil{
		//	publicip = "N/A"
		//} else {
		//	publicip = *alarm.PublicIpAddress
		//}
		//handle securitygroups
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
		Alarm = append(Alarm,accountId,*sess.Config.Region,*alarm.AlarmArn, *alarm.AlarmName, *alarm.Namespace,
			*alarm.MetricName, actions, dimensions)
		//fmt.Println(Alarm)
		AlarmList = append(AlarmList, Alarm)
	}
	//for _, i := range AlarmList {
	//	fmt.Println(i)
	//}
	return AlarmList
}
