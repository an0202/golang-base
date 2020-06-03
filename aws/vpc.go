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

//List SecurityGroups With Policy
func Listv2SGs(se Session) (PolicyList []interface{}) {
	// Create an EC2 service client.
	svc := ec2.New(se.Sess)

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
func ListEIPs(se Session) (EIPList [][]interface{}) {
	// Create an EC2 service client.
	svc := ec2.New(se.Sess)
	// get eips
	output, err := svc.DescribeAddresses(&ec2.DescribeAddressesInput{
		DryRun: aws.Bool(false),
	})
	if err != nil {
		tools.WarningLogger.Println(err)
		return
	}
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
		EIP = append(EIP, se.AccountId, se.UsedRegion, name, *eip.PublicIp)
		EIPList = append(EIPList, EIP)
	}
	return EIPList
}

//List SubNets
func ListSubNets(se Session) (SubNetList [][]interface{}) {
	// Create an EC2 service client.
	svc := ec2.New(se.Sess)
	// get eips
	output, err := svc.DescribeSubnets(&ec2.DescribeSubnetsInput{
		MaxResults: aws.Int64(100),
		DryRun:     aws.Bool(false),
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
		SubNet = append(SubNet, *subnet.OwnerId, se.UsedRegion, name, *subnet.SubnetId, *subnet.VpcId, *subnet.CidrBlock,
			*subnet.AvailabilityZone, *subnet.DefaultForAz)
		SubNetList = append(SubNetList, SubNet)
	}
	return SubNetList
}

//List RouteTables
func ListRouteTables(se Session) (RouteTableList [][]interface{}) {
	// Create an EC2 service client.
	svc := ec2.New(se.Sess)
	// get routetables
	output, err := svc.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		MaxResults: aws.Int64(100),
		DryRun:     aws.Bool(false),
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
		var routes, assosications []interface{}
		if len(routeTable.Routes) == 0 {
			routes = append(routes, "N/A")
		} else {
			for _, route := range routeTable.Routes {
				routes = append(routes, *route, " ")
			}
		}
		if len(routeTable.Associations) == 0 {
			assosications = append(assosications, "N/A")
		} else {
			for _, assosication := range routeTable.Associations {
				assosications = append(assosications, *assosication, " ")
			}
		}
		if len(output.RouteTables) == 100 {
			tools.WarningLogger.Println("Subnet Number > 100 , Data May Loss")
		}
		RouteTable = append(RouteTable, *routeTable.OwnerId, se.UsedRegion, name, *routeTable.RouteTableId, *routeTable.VpcId,
			assosications, routes)
		RouteTableList = append(RouteTableList, RouteTable)
	}
	return RouteTableList
}

//List VpcPeering
func ListVPCPeering(se Session) (PeeringConnectionList [][]interface{}) {
	// Create an EC2 service client.
	svc := ec2.New(se.Sess)
	// get VpcPeering
	output, err := svc.DescribeVpcPeeringConnections(&ec2.DescribeVpcPeeringConnectionsInput{
		MaxResults: aws.Int64(100),
		DryRun:     aws.Bool(false),
	})
	if err != nil {
		tools.WarningLogger.Println(err)
		return
	}
	for _, vpcPerring := range output.VpcPeeringConnections {
		var PeeringConnection []interface{}
		var name string
		//get name tag
		for _, tag := range vpcPerring.Tags {
			if *tag.Key == "Name" {
				name = *tag.Value
			}
		}
		if len(name) == 0 {
			name = "N/A "
		}
		if len(output.VpcPeeringConnections) == 100 {
			tools.WarningLogger.Println("VpcPeeringConnections Number > 100 , Data May Loss")
		}
		PeeringConnection = append(PeeringConnection, se.AccountId, se.UsedRegion, name, *vpcPerring.VpcPeeringConnectionId,
			*vpcPerring.RequesterVpcInfo, *vpcPerring.AccepterVpcInfo)
		PeeringConnectionList = append(PeeringConnectionList, PeeringConnection)
	}
	return PeeringConnectionList
}

//List NatGateway
func ListNatGateway(se Session) (NatGatewayList [][]interface{}) {
	// Create an EC2 service client.
	svc := ec2.New(se.Sess)
	// get nat gateway
	output, err := svc.DescribeNatGateways(&ec2.DescribeNatGatewaysInput{
		MaxResults: aws.Int64(100),
	})
	if err != nil {
		tools.WarningLogger.Println(err)
		return
	}
	for _, natGateway := range output.NatGateways {
		var NatGateway []interface{}
		var name string
		//get name tag
		for _, tag := range natGateway.Tags {
			if *tag.Key == "Name" {
				name = *tag.Value
			}
		}
		if len(name) == 0 {
			name = "N/A "
		}
		if len(output.NatGateways) == 100 {
			tools.WarningLogger.Println("NatGateways Number > 100 , Data May Loss")
		}
		NatGateway = append(NatGateway, se.AccountId, se.UsedRegion, name, *natGateway.NatGatewayId, *natGateway.VpcId,
			*natGateway.SubnetId)
		NatGatewayList = append(NatGatewayList, NatGateway)
	}
	return NatGatewayList
}
