/**
 * @Author: jie.an
 * @Description:
 * @File:  elb.go
 * @Version: 1.0.0
 * @Date: 2020/02/19 13:00
 */
package aws

import (
	"github.com/aws/aws-sdk-go/service/elb"
	"golang-base/tools"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

//List ElastiLBV1
func ListLBs(se Session) (LBList [][]interface{}) {
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
				tools.ErrorLogger.Fatalln(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			tools.ErrorLogger.Fatalln(err.Error())
		}
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
		LBList = append(LBList, loadBalancer)
	}
	return LBList
}
