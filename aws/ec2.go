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
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2InstanceDetail struct {
	InstanceId          string
	Region              string
	AWSProfile			string
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
	if _, ok := ec2instancedetail["InstanceId"]; ok {
		for k, v := range ec2instancedetail {
			switch k {
			case "InstanceId":
				instance.InstanceId = ec2instancedetail["InstanceId"]
			case "Region":
				instance.Region = ec2instancedetail["Region"]
			case "AWS_PROFILE":
				instance.AWSProfile = ec2instancedetail["AWS_PROFILE"]
			default:
				instance.Tags[k] = v
			}
		}
	} else {
		tools.ErrorLogger.Fatalln("Missing InstanceId:", ec2instancedetail)
	}
	return instance
}

// GetTags return ec2 tags information
// allTags : current tags for ec2
// queryResult : instanceId and the returned value of the queried tag key in list
func EC2GetTags(sess *session.Session, resource EC2InstanceDetail, queryKeys string) (allTags map[string]string,queryResult []interface{}){
	// Create an EC2 service client.
	svc := ec2.New(sess)

	// whether to override tag when tag exists
	// Get current ec2 Tags and modify
	curTags, err := svc.DescribeTags(&ec2.DescribeTagsInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("resource-id"),
				Values: []*string{aws.String(resource.InstanceId)},
			},
		},
	})
	if err != nil {
		tools.WarningLogger.Fatal(err)
	}
	// current ec2 tags key map
	if len(curTags.Tags) != 0 {
		allTags = make(map[string]string)
		for _, v := range curTags.Tags {
			allTags[*v.Key] = *v.Value
		}
		//fmt.Println("Cur ec2 tag key list:", allTags)
	}
	// queryResult with instanceid and value of the queried tag key
	queryResult = append(queryResult,resource.InstanceId)
	for _ , queryKey := range strings.Split(queryKeys, ",") {
		if _,ok := allTags[queryKey] ; ok {
			queryResult = append(queryResult, allTags[queryKey])
		} else {
			queryResult = append(queryResult, "N/A")
		}
	}
	tools.InfoLogger.Printf("Query Result For Resource %s : %s \n", resource.InstanceId,queryResult)
	return allTags, queryResult
}

// CreateTags from excel with skip exist tag key
func EC2CreateTags(sess *session.Session, instance EC2InstanceDetail, override bool) {
	// instance reMarshal to aws ec2 type
	var resourceIDs = []string{instance.InstanceId}

	// Create an EC2 service client.
	svc := ec2.New(sess)
	// Get EBSid
	ebsMap, err := svc.DescribeInstanceAttribute(&ec2.DescribeInstanceAttributeInput{
		DryRun:     aws.Bool(false),
		InstanceId: aws.String(instance.InstanceId),
		Attribute:  aws.String("blockDeviceMapping"),
	})
	if err != nil {
		tools.WarningLogger.Fatal(err)
	}
	for _, v := range ebsMap.BlockDeviceMappings {
		var ebs = EBSDetail{
			VolumeId: *v.Ebs.VolumeId,
			Status:   *v.Ebs.Status,
		}
		resourceIDs = append(resourceIDs, ebs.VolumeId)
	}
	tools.InfoLogger.Println("Starting Process Resource:", resourceIDs)
	var tags []*ec2.Tag
	// whether to override tag when tag exists
	if override {
		// override
		for tagKey, tagValue := range instance.Tags {
			tags = append(tags, &ec2.Tag{
				Key:   aws.String(tagKey),
				Value: aws.String(tagValue),
			})
		}
	} else {
		// Get current ec2 Tags and modify
		curTags, err := svc.DescribeTags(&ec2.DescribeTagsInput{
			DryRun: aws.Bool(false),
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("resource-id"),
					Values: []*string{aws.String(instance.InstanceId)},
				},
			},
		})
		if err != nil {
			tools.WarningLogger.Fatal(err)
		}
		// current ec2 Tags key list
		var curTagsKeyList []string
		if len(curTags.Tags) != 0 {
			for _, v := range curTags.Tags {
				curTagsKeyList = append(curTagsKeyList, *v.Key)
			}
			tools.InfoLogger.Println("Current EC2 Tag Key List Is:", curTagsKeyList)
		}
		//
		for preTagKey, preTagValue := range instance.Tags {
			exist := tools.StringFind(curTagsKeyList, preTagKey)
			switch exist {
			case true:
				tools.InfoLogger.Printf("Find %s In Current EC2 Tag List , Skip \n",preTagKey)
			case false:
				curTags.Tags = append(curTags.Tags, &ec2.TagDescription{
					Key:   aws.String(preTagKey),
					Value: aws.String(preTagValue),
				})
			}
		}
		for _, curTag := range curTags.Tags {
			tags = append(tags, &ec2.Tag{
				Key:   aws.String(*curTag.Key),
				Value: aws.String(*curTag.Value),
			})
		}
	}
	tools.InfoLogger.Printf("Create Tags For Resource %s Tag %s : \n",resourceIDs,tags)
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
	var resourceIDs = []string{instance.InstanceId}
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
		InstanceId: aws.String(instance.InstanceId),
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

//List KeysPairs
func ListKeyPairs(se Session) (KeyPairList[][]interface{}) {
	// Create an EC2 service client.
	svc := ec2.New(se.Sess)
	// Get instance tag name
	output, err := svc.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{
		DryRun: aws.Bool(false),
	})
	if err != nil {
		tools.WarningLogger.Println(err)
		return
	}
	for _, keypair := range output.KeyPairs {
		var keyPair []interface{}
		keyPair = append(keyPair,se.AccountId,se.UsedRegion,*keypair.KeyName,*keypair.KeyFingerprint)
		KeyPairList = append(KeyPairList, keyPair)
	}
	return KeyPairList
}

//ListSnapshots
func ListSnapshots(se Session) (SnapshotList [][]interface{}) {
	var maxResults = 300
	var token string
	var snapshots [][]interface{}
	var nextToken = "default"
	for nextToken != "" {
		//fmt.Println("Start Loop With Token:", token)
		snapshots, nextToken = listSnapshots(se, token, maxResults)
		for _, snapshot := range snapshots {
			SnapshotList = append(SnapshotList, snapshot)
		}
		if len(snapshots) == maxResults && nextToken != "" {
			snapshots, nextToken = listSnapshots(se, nextToken, maxResults)
			for _, snapshot := range snapshots {
				SnapshotList = append(SnapshotList, snapshot)
			}
			//fmt.Println("Generated New NextToken:      ",nextToken)
			token = nextToken
		} else {
			nextToken = ""
		}
	}
	return SnapshotList
}

//List snapshot Internal
func listSnapshots(se Session, token string,maxResults int) (SnapshotList [][]interface{},nextToken string) {
	// Create an EC2 service client.
	svc := ec2.New(se.Sess)
	// Get snapshots
	output, err := svc.DescribeSnapshots(&ec2.DescribeSnapshotsInput{
		DryRun: aws.Bool(false),
		MaxResults: aws.Int64(int64(maxResults)),
		NextToken: aws.String(token),
		OwnerIds: []*string{aws.String(se.AccountId)},
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
	for _, snapshot := range output.Snapshots {
		var Snapshot []interface{}
		Snapshot = append(Snapshot,*snapshot.OwnerId,se.UsedRegion,*snapshot.SnapshotId,*snapshot.VolumeId,
			*snapshot.Description,*snapshot.State)
		SnapshotList = append(SnapshotList, Snapshot)
	}
	if output.NextToken == nil {
		nextToken = ""
	} else {
		nextToken = *output.NextToken
	}
	return SnapshotList, nextToken
}

//List AMI
func ListAMIs(se Session) (AMIList[][]interface{}) {
	// Create an EC2 service client.
	svc := ec2.New(se.Sess)
	// Get instance tag name
	output, err := svc.DescribeImages(&ec2.DescribeImagesInput{
		DryRun: aws.Bool(false),
		Owners: []*string{aws.String("self")},
	})
	if err != nil {
		tools.WarningLogger.Println(err)
		return
	}
	for _, image := range output.Images {
		var Image []interface{}
		Image = append(Image,*image.OwnerId,se.UsedRegion,*image.ImageId,*image.Name,*image.State)
		AMIList = append(AMIList, Image)
	}
	return AMIList
}

//ListVolumes
func ListVolumes(se Session) (VolumeList [][]interface{}) {
	var maxResults = 100
	var token string
	var vols [][]interface{}
	var nextToken = "default"
	for nextToken != "" {
		//fmt.Println("use nextToken:",nextToken)
		vols, nextToken = listVolumes(se, token, maxResults)
		for _, vol := range vols {
			VolumeList = append(VolumeList, vol)
		}
		if len(vols) == maxResults && nextToken != "" {
			//tools.WarningLogger.Println("Get More Volumes ......")
			vols, nextToken = listVolumes(se, nextToken, maxResults)
			for _, vol := range vols {
				VolumeList = append(VolumeList, vol)
			}
			//fmt.Println("nextTokenGenerate:      ",nextToken)
			token = nextToken
		} else {
			nextToken = ""
		}
	}
	return VolumeList
}

//List Volumes Internal
func listVolumes(se Session, token string,maxResults int) (VolumeList [][]interface{},nextToken string) {
	// Create an EC2 service client.
	svc := ec2.New(se.Sess)
	// Get volumes
	output, err := svc.DescribeVolumes(&ec2.DescribeVolumesInput{
		DryRun: aws.Bool(false),
		MaxResults: aws.Int64(int64(maxResults)),
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
	for _, volume := range output.Volumes {
		var Volume []interface{}
		var name , attachedInstance string
		var tags []string
		//handle tags
		if len(volume.Tags) == 0 {
			tags = append(tags, "N/A ")
		} else {
			for _, tag := range volume.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
				}
				tags = append(tags, *tag.Key+":"+*tag.Value+" ")
			}
			if len(name) == 0 {
				name = "N/A "
			}
		}
		//handle attached instance
		if len(volume.Attachments) == 0 {
			attachedInstance = "N/A"
		} else {
			for _ ,attach := range volume.Attachments {
				attachedInstance = *attach.InstanceId
			}
		}
		Volume = append(Volume,se.AccountId,se.UsedRegion,name ,*volume.VolumeId,attachedInstance,*volume.State, *volume.VolumeType,
			*volume.Size,*volume.AvailabilityZone)
		VolumeList = append(VolumeList, Volume)
	}
	if output.NextToken == nil {
		nextToken = ""
	} else {
		nextToken = *output.NextToken
	}
	return VolumeList, nextToken
}

//List Instances
func ListInstances(se Session) (InstanceList [][]interface{}) {
	// Create an EC2 service client.
	svc := ec2.New(se.Sess)
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
			var platform, rolearn, instancename ,keypair, publicip string
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
			if instance.KeyName == nil {
				keypair = "N/A"
			} else {
				keypair = *instance.KeyName
			}
			if instance.PublicIpAddress == nil{
				publicip = "N/A"
			} else {
				publicip = *instance.PublicIpAddress
			}
			//handle securitygroups
			var sgs, tags []string
			if len(instance.SecurityGroups) == 0 {
				sgs = append(sgs, "N/A ")
			} else {
				for _, sg := range instance.SecurityGroups {
					sgs = append(sgs, *sg.GroupId+"("+*sg.GroupName+")"+" ")
				}
			}
			//handle tags
			if len(instance.Tags) == 0 {
				tags = append(tags, "N/A ")
			} else {
				for _, tag := range instance.Tags {
					if *tag.Key == "Name" {
						instancename = *tag.Value
					}
					tags = append(tags, *tag.Key+":"+*tag.Value+" ")
				}
				if len(instancename) == 0 {
					instancename = "N/A "
				}
			}
			Instance = append(Instance, se.AccountId,se.UsedRegion,instancename, *instance.InstanceId,
				*instance.InstanceType, platform, *instance.State.Name, *instance.VpcId,
				rolearn, *instance.SubnetId, *instance.PrivateIpAddress, publicip,keypair, sgs, tags)
			//fmt.Println(Instance)
			InstanceList = append(InstanceList, Instance)
		}
	}
	//for _, i := range InstanceList {
	//	fmt.Println(i)
	//}
	return InstanceList
}
