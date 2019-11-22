/**
 * @Author: jie.an
 * @Description:
 * @File:  ec2.go
 * @Version: 1.0.0
 * @Date: 2019/11/22 14:36
 */
package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2InstanceDetail struct {
	InstanceID string
	Region     string
	Tags       map[string]string
}

func EC2InstanceMarshal(ec2instancedetail map[string]string) (instance EC2InstanceDetail) {
	instance.Tags = make(map[string]string)
	if _, ok := ec2instancedetail["InstanceID"]; ok {
		for k, v := range ec2instancedetail {
			switch k {
			case "InstanceID":
				instance.InstanceID = ec2instancedetail["InstanceID"]
			case "Region":
				instance.Region = ec2instancedetail["Region"]
			default:
				instance.Tags[k] = v
			}
		}
	} else {
		fmt.Println("Error While Processing:", ec2instancedetail)
		os.Exit(2)
	}
	return instance
}

func EC2CreateTags(sess *session.Session, instance EC2InstanceDetail) {
	// instance reMarshal to aws ec2 type
	var resourceIDs = []string{instance.InstanceID}
	var tags []*ec2.Tag
	for k, v := range instance.Tags {
		tags = append(tags, &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	// Create an EC2 service client.
	svc := ec2.New(sess)

	result, err := svc.CreateTags(&ec2.CreateTagsInput{
		DryRun:    aws.Bool(false),
		Resources: aws.StringSlice(resourceIDs),
		Tags:      tags,
	})
	fmt.Println(err)
	if err != nil {
	}
	fmt.Println(result)
}

func EC2DeleteTags(sess *session.Session, instance EC2InstanceDetail) {
	// instance reMarshal to aws ec2 type
	var resourceIDs = []string{instance.InstanceID}
	var tags []*ec2.Tag
	for k, v := range instance.Tags {
		tags = append(tags, &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	// Create an EC2 service client.
	svc := ec2.New(sess)

	result, err := svc.DeleteTags(&ec2.DeleteTagsInput{
		DryRun:    aws.Bool(false),
		Resources: aws.StringSlice(resourceIDs),
		Tags:      tags,
	})
	fmt.Println(err)
	if err != nil {
	}
	fmt.Println(result)
}

func Testdescribeinstance(sess *session.Session, instanceid string) {
	// Create an EC2 service client.
	svc := ec2.New(sess)
	result, err := svc.DescribeInstanceAttribute(&ec2.DescribeInstanceAttributeInput{
		DryRun:     aws.Bool(false),
		InstanceId: aws.String(instanceid),
		Attribute:  aws.String("blockDeviceMapping"),
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}
