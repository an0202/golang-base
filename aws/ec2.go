/**
 * @Author: jie.an
 * @Description:
 * @File:  ec2.go
 * @Version: 1.0.0
 * @Date: 2019/11/22 14:36
 */
package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"golang-base/tools"
)

type EC2InstanceDetail struct {
	InstanceID string
	Region     string
	Tags       map[string]string
	BlockDeviceMappings []EBSDetail
}

type EBSDetail struct{
	VolumeId string
	Status  string
}

// EC2Instance retype from excel
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
		tools.ErrorLogger.Fatalln("Missing InstanceID:", ec2instancedetail)
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
	// Get EBSid
	ebsmap, err := svc.DescribeInstanceAttribute(&ec2.DescribeInstanceAttributeInput{
		DryRun:     aws.Bool(false),
		InstanceId: aws.String(instance.InstanceID),
		Attribute:  aws.String("blockDeviceMapping"),
	})
	if err != nil {
		tools.WarningLogger.Fatal(err)
	}
	for _,v := range ebsmap.BlockDeviceMappings{
		var ebs = EBSDetail{
			VolumeId: *v.Ebs.VolumeId,
			Status:   *v.Ebs.Status,
		}
		resourceIDs = append(resourceIDs, ebs.VolumeId)
	}
	tools.InfoLogger.Println("Add Tags To :",resourceIDs)
	// Create tag
 	_, err = svc.CreateTags(&ec2.CreateTagsInput{
		DryRun:    aws.Bool(false),
		Resources: aws.StringSlice(resourceIDs),
		Tags:      tags,
	})
	if err != nil {
		tools.WarningLogger.Println(err)
	}
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
	// Get EBSid
	ebsmap, err := svc.DescribeInstanceAttribute(&ec2.DescribeInstanceAttributeInput{
		DryRun:     aws.Bool(false),
		InstanceId: aws.String(instance.InstanceID),
		Attribute:  aws.String("blockDeviceMapping"),
	})
	if err != nil {
		tools.WarningLogger.Fatal(err)
	}
	for _,v := range ebsmap.BlockDeviceMappings{
		var ebs = EBSDetail{
			VolumeId: *v.Ebs.VolumeId,
			Status:   *v.Ebs.Status,
		}
		resourceIDs = append(resourceIDs, ebs.VolumeId)
	}
	tools.InfoLogger.Println("Delete Tags From :",resourceIDs)
	// delete tag
	_, err = svc.DeleteTags(&ec2.DeleteTagsInput{
		DryRun:    aws.Bool(false),
		Resources: aws.StringSlice(resourceIDs),
		Tags:      tags,
	})
	if err != nil {
		tools.WarningLogger.Println(err)
	}
}
