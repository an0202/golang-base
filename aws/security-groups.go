/**
 * @Author: jie.an
 * @Description:
 * @File:  security-sgs.go
 * @Version: 1.0.0
 * @Date: 2019/11/20 13:34
 */
package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type PolicyDetail struct {
	GroupName string
	VpcId     string
	GroupId   string
	Protocol  string
	Source    string
	FromPort  int64
	ToPort    int64
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func GetSGPolicys(sess *session.Session) (PolicyList []interface{}) {
	// Create an EC2 service client.
	svc := ec2.New(sess)

	// Retrieve the security sg descriptions
	result, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		DryRun: aws.Bool(false),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "InvalidGroupId.Malformed":
				fallthrough
			case "InvalidGroup.NotFound":
				exitErrorf("%s.", aerr.Message())
			}
		}
		exitErrorf("Unable to get descriptions for security sgs, %v", err)
	}
	//
	policy := new(PolicyDetail)
	for _, sg := range result.SecurityGroups {
		for _, ippermission := range sg.IpPermissions {
			if *ippermission.IpProtocol == "-1" {
				if len(ippermission.IpRanges) != 0 {
					for _, permission := range ippermission.IpRanges {
						// fmt.Println("  ALL IpRanges information：", *sg.GroupName, *sg.VpcId, *sg.GroupId, "ALL PROTOCOL", "from port ALL", "end port ALL", *permission.CidrIp)
						policy.GroupName = *sg.GroupName
						policy.VpcId = *sg.VpcId
						policy.GroupId = *sg.GroupId
						policy.Source = *permission.CidrIp
						policy.Protocol = "ALL"
						policy.FromPort = 0
						policy.ToPort = 65535
						PolicyList = append(PolicyList, *policy)
					}
				}
				if len(ippermission.UserIdGroupPairs) != 0 {
					for _, permission := range ippermission.UserIdGroupPairs {
						// fmt.Println("  ALL GroupPairs information：", *sg.GroupName, *sg.VpcId, *sg.GroupId, "ALLPROTOCOL", "fromportALL", "endportALL", *permission.GroupId)
						policy.GroupName = *sg.GroupName
						policy.VpcId = *sg.VpcId
						policy.GroupId = *sg.GroupId
						policy.Source = *permission.GroupId
						policy.Protocol = "ALL"
						policy.FromPort = 0
						policy.ToPort = 65535
						PolicyList = append(PolicyList, *policy)
					}
				}
				if len(ippermission.PrefixListIds) != 0 {
					for _, permission := range ippermission.PrefixListIds {
						//fmt.Println("  ===all===Prefix information:", permission)
						policy.GroupName = *sg.GroupName
						policy.VpcId = *sg.VpcId
						policy.GroupId = *sg.GroupId
						policy.Source = *permission.PrefixListId
						policy.Protocol = "Unknow"
						policy.FromPort = 0
						policy.ToPort = 65535
						PolicyList = append(PolicyList, *policy)
					}
				}
			} else {
				if len(ippermission.IpRanges) != 0 {
					for _, permission := range ippermission.IpRanges {
						//fmt.Println("  IpRanges information：", *sg.GroupName, *sg.VpcId, *sg.GroupId, *ippermission.IpProtocol, *ippermission.FromPort, *ippermission.ToPort, *permission.CidrIp)
						policy.GroupName = *sg.GroupName
						policy.VpcId = *sg.VpcId
						policy.GroupId = *sg.GroupId
						policy.Source = *permission.CidrIp
						policy.Protocol = *ippermission.IpProtocol
						policy.FromPort = *ippermission.FromPort
						policy.ToPort = *ippermission.ToPort
						PolicyList = append(PolicyList, *policy)
					}
				}
				if len(ippermission.UserIdGroupPairs) != 0 {
					for _, permission := range ippermission.UserIdGroupPairs {
						//fmt.Println("  GroupPairs information：", *sg.GroupName, *sg.VpcId, *sg.GroupId, *ippermission.IpProtocol, *ippermission.FromPort, *ippermission.ToPort, *permission.GroupId )
						policy.GroupName = *sg.GroupName
						policy.VpcId = *sg.VpcId
						policy.GroupId = *sg.GroupId
						policy.Source = *permission.GroupId
						policy.Protocol = *ippermission.IpProtocol
						policy.FromPort = *ippermission.FromPort
						policy.ToPort = *ippermission.ToPort
						PolicyList = append(PolicyList, *policy)
					}
				}
				if len(ippermission.PrefixListIds) != 0 {
					for _, permission := range ippermission.PrefixListIds {
						//fmt.Println("  ======Prefix information:", permission)
						policy.GroupName = *sg.GroupName
						policy.VpcId = *sg.VpcId
						policy.GroupId = *sg.GroupId
						policy.Source = *permission.PrefixListId
						policy.Protocol = "Unknow"
						policy.FromPort = -1
						policy.ToPort = -1
						PolicyList = append(PolicyList, *policy)
					}
				}
			}
		}
	}
	return PolicyList
}
