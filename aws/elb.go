/**
 * @Author: jie.an
 * @Description:
 * @File:  elb.go
 * @Version: 1.0.0
 * @Date: 2020/02/19 13:00
 */
package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"golang-base/tools"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

type LoadBalancer struct {
	AccountId string
	Region    string
	VPCId     string
	Name 	  string
	DNSName   string
	ARN       string
	Type      string
	Scheme    string
	//FromPort  int64
	//ToPort    int64
	AvailabilityZones       []interface{}
	SecurityGroups			[]interface{}
	Listeners				[]interface{}
	TargetGroups            []interface{}
	Backends                []interface{}
}

//List ElastiLBV2
func Listv2LBv2s(se Session) (LBv2List []interface{}) {
	// Create an elb service client.
	svc := elbv2.New(se.Sess)
	// Get lb
	output, err := svc.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
		//MaxRecords: aws.Int64(100),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				tools.ErrorLogger.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			tools.ErrorLogger.Println(err.Error())
		}
		return
	}
	LB := new(LoadBalancer)
	for _, lb := range output.LoadBalancers {
		//handle securityGroups availabilityZones
		if len(lb.SecurityGroups) == 0 {
			LB.SecurityGroups = append(LB.SecurityGroups, "N/A")
		} else {
			for _, sg := range lb.SecurityGroups {
				LB.SecurityGroups = append(LB.SecurityGroups, *sg)
			}
		}
		if len(lb.AvailabilityZones) == 0 {
			LB.AvailabilityZones = append(LB.AvailabilityZones, "N/A")
		} else {
			for _, az := range lb.AvailabilityZones {
				LB.AvailabilityZones = append(LB.AvailabilityZones, *az)
			}
		}
		//if len(output.CacheClusters) >= 100 {
		//
		//	tools.WarningLogger.Println("Number Of Clusters > 100 , Data May Missing.")
		//}
		LB.AccountId = se.AccountId
		LB.Region = se.UsedRegion
		LB.VPCId = *lb.VpcId
		LB.Name = *lb.LoadBalancerName
		LB.DNSName = *lb.DNSName
		LB.ARN = *lb.LoadBalancerArn
		LB.Type = *lb.Type
		LB.Scheme = *lb.Scheme
		LB.Listeners = LB.ListListeners(se, LB.ARN)
		LB.TargetGroups,LB.Backends = LB.ListTargetGroups(se, LB.ARN)
		LBv2List = append(LBv2List, *LB)
	}
	return LBv2List
}

func (lb *LoadBalancer) ListListeners(se Session,LBv2ARN string) (ListenerList []interface{}) {
	// Create an elb service client.
	svc := elbv2.New(se.Sess)
	// Get lb listener
	output, err := svc.DescribeListeners(&elbv2.DescribeListenersInput{
		LoadBalancerArn: aws.String(LBv2ARN),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				tools.ErrorLogger.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			tools.ErrorLogger.Println(err.Error())
		}
		return
	}
	if len(output.Listeners) == 0 {
		ListenerList = append(ListenerList, "N/A")
	} else {
		for _ , listener := range output.Listeners {
			ListenerList = append(ListenerList, listener)
		}
	}
	return ListenerList
}

func (lb *LoadBalancer) ListTargetGroups(se Session,LBv2ARN string) (TargetGroupList,BackendList []interface{}) {
	// Create an elb service client.
	svc := elbv2.New(se.Sess)
	// Get lb targetGroups
	output, err := svc.DescribeTargetGroups(&elbv2.DescribeTargetGroupsInput{
		LoadBalancerArn: aws.String(LBv2ARN),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				tools.ErrorLogger.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			tools.ErrorLogger.Println(err.Error())
		}
		return
	}
	if len(output.TargetGroups) == 0 {
		TargetGroupList = append(TargetGroupList, "N/A")
		//handle target health
		BackendList = append(BackendList, "N/A")
	} else {
		for _ , targetGroup := range output.TargetGroups {
			TargetGroupList = append(TargetGroupList,targetGroup)
			//handle target health
			for _,backend := range lb.ListBackends(se, *targetGroup.TargetGroupArn) {
				BackendList = append(BackendList, backend)
			}
		}
	}
	return TargetGroupList,BackendList
}

func (lb *LoadBalancer) ListBackends(se Session,TargetARN string) (BackendList []interface{}) {
	// Create an elb service client.
	svc := elbv2.New(se.Sess)
	// Get lb targetGroups
	output, err := svc.DescribeTargetHealth(&elbv2.DescribeTargetHealthInput{
		TargetGroupArn: aws.String(TargetARN),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				tools.ErrorLogger.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			tools.ErrorLogger.Println(err.Error())
		}
		return
	}
	//
	backEnd := make(map[string]string)
	if len(output.TargetHealthDescriptions) == 0 {
		backEnd[TargetARN] = "N/A"
		BackendList = append(BackendList, backEnd)
	} else {
		for _ , targetHealth := range output.TargetHealthDescriptions {
			backEnd = make(map[string]string)
			backEnd[TargetARN] = *targetHealth.Target.Id
			BackendList = append(BackendList, backEnd)
		}
	}
	//fmt.Println(BackendList)
	return BackendList
}

//List ElastiLBV1
func ListCLBs(se Session) (CLBList [][]interface{}) {
	// Create an elb service client.
	svc := elb.New(se.Sess)
	// Get lb
	output, err := svc.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{
		//MaxRecords: aws.Int64(100),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				tools.ErrorLogger.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			tools.ErrorLogger.Println(err.Error())
		}
		return
	}
	for _, lb := range output.LoadBalancerDescriptions {
		var loadBalancer []interface{}
		//handle securityGroups listners instances availabilityZones
		var sgs ,listners ,instances ,azs []interface{}
		if len(lb.SecurityGroups) == 0 {
			sgs = append(sgs, "N/A")
		} else {
			for _, sg := range lb.SecurityGroups {
				sgs = append(sgs, *sg)
			}
		}
		if len(lb.ListenerDescriptions) == 0 {
			listners = append(listners, "N/A")
		} else {
			for _, listner := range lb.ListenerDescriptions {
				listners = append(listners, *listner.Listener)
			}
		}
		if len(lb.Instances) == 0 {
			instances = append(instances, "N/A")
		} else {
			for _, instance := range lb.Instances {
				instances = append(instances, *instance.InstanceId)
			}
		}
		if len(lb.AvailabilityZones) == 0 {
			azs = append(azs, "N/A")
		} else {
			for _, az := range lb.AvailabilityZones {
				azs = append(azs, *az)
			}
		}
		//if len(output.CacheClusters) >= 100 {
		//
		//	tools.WarningLogger.Println("Number Of Clusters > 100 , Data May Missing.")
		//}
		loadBalancer = append(loadBalancer, se.AccountId,se.UsedRegion,*lb.VPCId,*lb.LoadBalancerName, *lb.DNSName,
			*lb.Scheme,azs,sgs,listners,*lb.HealthCheck,instances)
		CLBList = append(CLBList, loadBalancer)
	}
	return CLBList
}

