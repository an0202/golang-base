/**
 * @Author: jie.an
 * @Description:
 * @File:  vpc.go
 * @Version: 1.0.1
 * @Date: 2019/11/20 13:34
 */
package aws

import (
	"fmt"
	"golang-base/tools"
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

//List EIPs
func ListEIPs(sess *session.Session) (EIPList[][]interface{}) {
	// Create an EC2 service client.
	svc := ec2.New(sess)
	// get eips
	output, err := svc.DescribeAddresses(&ec2.DescribeAddressesInput{
		DryRun: aws.Bool(false),
	})
	if err != nil {
		tools.WarningLogger.Println(err)
		return
	}
	//handel accountId
	accountId := GetAccountId(sess)
	//
	for _, eip := range output.Addresses {
		var EIP []interface{}
		var name string
		//get name tag
		for _, tag := range eip.Tags {
			if *tag.Key == "Name" {
				name = *tag.Value
			}
		}
		if len(name) == 0 {
			name = "N/A "
		}
		EIP = append(EIP,accountId,*sess.Config.Region,name,*eip.PublicIp)
		EIPList = append(EIPList, EIP)
	}
	return EIPList
}

//List SubNets
func ListSubNets(sess *session.Session) (SubNetList[][]interface{}) {
	// Create an EC2 service client.
	svc := ec2.New(sess)
	// get eips
	output, err := svc.DescribeSubnets(&ec2.DescribeSubnetsInput{
		MaxResults: aws.Int64(100),
		DryRun: aws.Bool(false),
	})
	if err != nil {
		tools.WarningLogger.Println(err)
		return
	}
	for _, subnet := range output.Subnets {
		var SubNet []interface{}
		var name string
		//get name tag
		for _, tag := range subnet.Tags {
			if *tag.Key == "Name" {
				name = *tag.Value
			}
		}
		if len(name) == 0 {
			name = "N/A "
		}
		//
		if len(output.Subnets) == 100 {
			tools.WarningLogger.Println("Subnet Number > 100 , Data May Loss")
		}
		SubNet = append(SubNet,*subnet.OwnerId,*sess.Config.Region,name,*subnet.SubnetId,*subnet.VpcId,*subnet.CidrBlock,
			*subnet.AvailabilityZone, *subnet.DefaultForAz)
		SubNetList = append(SubNetList, SubNet)
	}
	return SubNetList
}

//List RouteTables
func ListRouteTables(sess *session.Session) (RouteTableList[][]interface{}) {
	// Create an EC2 service client.
	svc := ec2.New(sess)
	// get routetables
	output, err := svc.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		MaxResults: aws.Int64(100),
		DryRun: aws.Bool(false),
	})
	if err != nil {
		tools.WarningLogger.Println(err)
		return
	}
	for _, routeTable := range output.RouteTables {
		var RouteTable []interface{}
		var name string
		//get name tag
		for _, tag := range routeTable.Tags {
			if *tag.Key == "Name" {
				name = *tag.Value
			}
		}
		if len(name) == 0 {
			name = "N/A "
		}
		//handle route and assosication
		var routes,assosications []interface{}
		if len(routeTable.Routes) == 0 {
			routes = append(routes,"N/A")
		} else {
			for _, route := range routeTable.Routes {
				routes = append(routes, *route," ")
			}
		}
		if len(routeTable.Associations) == 0 {
			assosications = append(assosications,"N/A")
		} else {
			for _, assosication := range routeTable.Associations {
				assosications = append(assosications, *assosication," ")
			}
		}
		if len(output.RouteTables) == 100 {
			tools.WarningLogger.Println("Subnet Number > 100 , Data May Loss")
		}
		RouteTable = append(RouteTable,*routeTable.OwnerId,*sess.Config.Region,name,*routeTable.RouteTableId,*routeTable.VpcId,
		assosications,routes)
		RouteTableList = append(RouteTableList, RouteTable)
	}
	return RouteTableList
}

