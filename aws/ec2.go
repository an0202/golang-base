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
	"golang-base/tools"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2InstanceDetail struct {
	InstanceID          string
	Region              string
	Tags                map[string]string
	BlockDeviceMappings []EBSDetail
}

type EBSDetail struct {
	VolumeId string
	Status   string
}

//Cretate AMI For EC2
func CreateImage(sess *session.Session, instanceid, suffix string) {
	// Create an EC2 service client.
	svc := ec2.New(sess)
	// Get instance tag name
	output, err := svc.DescribeTags(&ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("resource-id"),
				Values: []*string{aws.String(instanceid)},
			},
		},
	})
	if err != nil {
		tools.WarningLogger.Println(err)
		return
	}
	if len(output.Tags) == 0 {
		tools.WarningLogger.Printf("Instance : %s Does Not Exist Or Instance Does Not Have A Name Tag.\n", instanceid)
		return
	}
	var aminame string
	switch suffix {
	case "final":
		for _, tag := range output.Tags {
			if *tag.Key == "Name" {
				aminame = *tag.Value + "-FINALBACKUP"
			}
		}
	case "date":
		for _, tag := range output.Tags {
			if *tag.Key == "Name" {
				aminame = *tag.Value + time.Now().Format("-20060102")
			}
		}
	default:
		for _, tag := range output.Tags {
			if *tag.Key == "Name" {
				aminame = *tag.Value + "-" + suffix
			}
		}
	}
	//Create ami
	ami, err := svc.CreateImage(&ec2.CreateImageInput{
		DryRun:     aws.Bool(false),
		NoReboot:   aws.Bool(true),
		InstanceId: aws.String(instanceid),
		Name:       aws.String(aminame),
	})
	if err != nil {
		tools.WarningLogger.Println(err)
		return
	}
	tools.InfoLogger.Printf("Create AMI %s (%s) Successfully.", *ami.ImageId, aminame)
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

// CreateTags from excel with skip exist tag key
func EC2CreateTags(sess *session.Session, instance EC2InstanceDetail) {
	// instance reMarshal to aws ec2 type
	var resourceIDs = []string{instance.InstanceID}

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
	for _, v := range ebsmap.BlockDeviceMappings {
		var ebs = EBSDetail{
			VolumeId: *v.Ebs.VolumeId,
			Status:   *v.Ebs.Status,
		}
		resourceIDs = append(resourceIDs, ebs.VolumeId)
	}
	tools.InfoLogger.Println("Add Tags To :", resourceIDs)
	// Get Tag and modify
	curTags, err := svc.DescribeTags(&ec2.DescribeTagsInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("resource-id"),
				Values: []*string{aws.String(instance.InstanceID)},
			},
		},
	})
	if err != nil {
		tools.WarningLogger.Fatal(err)
	}
	// key list
	var curTagsKeyList []string
	if len(curTags.Tags) != 0 {
		for _, v := range curTags.Tags {
			curTagsKeyList = append(curTagsKeyList, *v.Key)
		}
		fmt.Println("Cur tag key  list:", curTagsKeyList)
	}
	for preTagKey, preTagValue := range instance.Tags {
		exist := tools.StringFind(curTagsKeyList, preTagKey)
		switch exist {
		case true:
			fmt.Println("exist abort")
		case false:
			curTags.Tags = append(curTags.Tags, &ec2.TagDescription{
				Key:   aws.String(preTagKey),
				Value: aws.String(preTagValue),
			})
		}
	}
	var tags []*ec2.Tag

	//// overide = true
	//	for tagKey, tagValue := range instance.Tags {
	//		tags = append(tags, &ec2.Tag{
	//			Key:   aws.String(tagKey),
	//			Value: aws.String(tagValue),
	//		})
	//	}
	for _, curtag := range curTags.Tags {
		tags = append(tags, &ec2.Tag{
			Key:   aws.String(*curtag.Key),
			Value: aws.String(*curtag.Value),
		})
	}
	fmt.Println(tags)
	//Create tag
	_, err = svc.CreateTags(&ec2.CreateTagsInput{
		DryRun:    aws.Bool(false),
		Resources: aws.StringSlice(resourceIDs),
		Tags:      tags,
	})
	if err != nil {
		tools.WarningLogger.Println(err)
	}
}

//func EBSCreateTags(sess *session.Session, instance EC2InstanceDetail) {
//	// Create an EC2 service client.
//	svc := ec2.New(sess)
//	// Get EBSid
//	ebsmap, err := svc.DescribeInstanceAttribute(&ec2.DescribeInstanceAttributeInput{
//		DryRun:     aws.Bool(false),
//		InstanceId: aws.String(instance.InstanceID),
//		Attribute:  aws.String("blockDeviceMapping"),
//	})
//	if err != nil {
//		tools.WarningLogger.Fatal(err)
//	}
//	for _, v := range ebsmap.BlockDeviceMappings {
//		//var ebs = EBSDetail{
//		//	VolumeId: *v.Ebs.VolumeId,
//		//	Status:   *v.Ebs.Status,
//		//}
//		//ebsIDs = append(ebsIDs, ebs.VolumeId)
//		curTags, err := svc.DescribeTags(&ec2.DescribeTagsInput{
//			DryRun: aws.Bool(false),
//			Filters: []*ec2.Filter{
//				{
//					Name:   aws.String("resource-id"),
//					Values: []*string{aws.String(*v.Ebs.VolumeId)},
//				},
//			},
//		})
//		if err != nil {
//			tools.WarningLogger.Fatal(err)
//		}
//		// cur key list
//		var curTagsKeyList []string
//		for _, v := range curTags.Tags {
//			curTagsKeyList = append(curTagsKeyList, *v.Key)
//		}
//		// pretagkey
//		var tags []*ec2.Tag
//		// key list
//		//if len(curTags.Tags) != 0 {
//		//	for _, v := range curTags.Tags {
//		//		curTagsKeyList = append(curTagsKeyList, *v.Key)
//		//	}
//		//	fmt.Println("Cur tag key  list:", curTagsKeyList)
//		//}
//		for preTagKey, preTagValue := range instance.Tags {
//			exist := tools.StringFind(curTagsKeyList, preTagKey)
//			switch exist {
//			case true:
//				tools.InfoLogger.Println("TagKey Exist Abort To Add:", preTagKey, preTagValue)
//			case false:
//				curTags.Tags = append(curTags.Tags, &ec2.TagDescription{
//					Key:   aws.String(preTagKey),
//					Value: aws.String(preTagValue),
//				})
//			}
//		}
//		for _, curtag := range curTags.Tags {
//			tags = append(tags, &ec2.Tag{
//				Key:   aws.String(*curtag.Key),
//				Value: aws.String(*curtag.Value),
//			})
//		}
//		//Create tag
//		tools.InfoLogger.Println("Add Tags To :", *v.Ebs.VolumeId)
//		_, err = svc.CreateTags(&ec2.CreateTagsInput{
//			DryRun:    aws.Bool(false),
//			Resources: aws.StringSlice([]string{*v.Ebs.VolumeId}),
//			Tags:      tags,
//		})
//		if err != nil {
//			tools.WarningLogger.Println(err)
//		}
//	}
//}

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
	for _, v := range ebsmap.BlockDeviceMappings {
		var ebs = EBSDetail{
			VolumeId: *v.Ebs.VolumeId,
			Status:   *v.Ebs.Status,
		}
		resourceIDs = append(resourceIDs, ebs.VolumeId)
	}
	tools.InfoLogger.Println("Delete Tags From :", resourceIDs)
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

//List Snapshot
func ListSnapshots(sess *session.Session, accountid string) {
	// Create an EC2 service client.
	svc := ec2.New(sess)
	// Get instance tag name
	output, err := svc.DescribeSnapshots(&ec2.DescribeSnapshotsInput{
		OwnerIds: []*string{aws.String(accountid)},
	})
	if err != nil {
		tools.WarningLogger.Println(err)
		return
	}
	for _, snapshot := range output.Snapshots {
		fmt.Println(*snapshot.SnapshotId, *snapshot.Description)
	}
	//tools.InfoLogger.Println(&output)
}

//List AMI
func ListAMIs(sess *session.Session) {
	// Create an EC2 service client.
	svc := ec2.New(sess)
	// Get instance tag name
	output, err := svc.DescribeImages(&ec2.DescribeImagesInput{
		Owners: []*string{aws.String("self")},
	})
	if err != nil {
		tools.WarningLogger.Println(err)
		return
	}
	for _, image := range output.Images {
		fmt.Println(*image.ImageId, *image.Name)
	}
	//tools.InfoLogger.Println(&output)
}

//List Volumes
func ListVolumes(sess *session.Session) {
	// Create an EC2 service client.
	svc := ec2.New(sess)
	// Get instance tag name
	output, err := svc.DescribeVolumes(&ec2.DescribeVolumesInput{
		DryRun: aws.Bool(false),
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
	for _, volume := range output.Volumes {
		fmt.Println(*volume.VolumeId, *volume.State)
	}
	//tools.InfoLogger.Println(&output)
}

//List Instances
func ListInstances(sess *session.Session) (InstanceList [][]interface{}) {
	// Create an EC2 service client.
	svc := ec2.New(sess)
	// Get instance tag name
	output, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		DryRun: aws.Bool(false),
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
	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			var Instance []interface{}
			var platform, rolearn, instancename string
			if instance.Platform != nil {
				platform = *instance.Platform
			} else {
				platform = "linux"
			}
			if instance.IamInstanceProfile == nil {
				rolearn = "N/A"
			} else {
				rolearn = *instance.IamInstanceProfile.Arn
			}
			//handle securitygroups
			var sgs, tags []string
			if len(instance.SecurityGroups) == 0 {
				sgs = append(sgs, "N/A")
			} else {
				for _, sg := range instance.SecurityGroups {
					sgs = append(sgs, *sg.GroupId+*sg.GroupName)
				}
			}
			//handle tags
			if len(instance.Tags) == 0 {
				tags = append(tags, "N/A")
			} else {
				for _, tag := range instance.Tags {
					if *tag.Key == "Name" {
						instancename = *tag.Value
					}
					tags = append(tags, *tag.Key+":"+*tag.Value)
				}
				if len(instancename) == 0 {
					instancename = "N/A"
				}
			}
			Instance = append(Instance, instancename, *instance.InstanceId, *instance.InstanceType, platform, *instance.State.Name, *instance.VpcId,
				rolearn, *instance.SubnetId, *instance.KeyName, sgs, tags)
			InstanceList = append(InstanceList, Instance)
		}
	}
	//for _, i := range InstanceList {
	//	fmt.Println(i)
	//}
	return InstanceList
}
