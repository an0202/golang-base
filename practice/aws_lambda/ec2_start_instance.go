/**
 * @Author: jie.an
 * @Description:
 * @File:  ec2_start_instance.go
 * @Version: 1.0.0
 * @Date: 2019/11/6 19:53
 */
package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func returnInstance() (instanceList []string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("cn-north-1")},
	)
	if err != nil {
		fmt.Println(err)
	}
	svc := ec2.New(sess)
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:AutoRestart"),
				Values: []*string{
					aws.String("True"),
				},
			},
		},
	}
	result, err := svc.DescribeInstances(input)
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
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instanceList = append(instanceList, *instance.InstanceId)
			fmt.Println(*instance.InstanceId)
		}
	}
	return instanceList
}

func startInstance(instanceList []string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("cn-north-1")},
	)
	if err != nil {
		fmt.Println(err)
	}
	svc := ec2.New(sess)
	if len(instanceList) != 0 {
		for _, instance := range instanceList {
			input := &ec2.StartInstancesInput{
				DryRun: aws.Bool(true),
				InstanceIds: []*string{
					aws.String(instance),
				},
			}
			result, err := svc.StartInstances(input)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(result)
		}
	} else {
		fmt.Println("none instance in list")
	}
}

func stopInstance(instanceList []string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("cn-north-1")},
	)
	if err != nil {
		fmt.Println(err)
	}
	svc := ec2.New(sess)
	if len(instanceList) != 0 {
		for _, instance := range instanceList {
			input := &ec2.StopInstancesInput{
				DryRun: aws.Bool(true),
				InstanceIds: []*string{
					aws.String(instance),
				},
			}
			result, err := svc.StopInstances(input)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(result)
		}
	} else {
		fmt.Println("none instance in list")
	}
}

func main() {
	a := returnInstance()
	//startInstance(a)
	stopInstance(a)
}
