/**
 * @Author: jie.an
 * @Description:
 * @File:  loadbalancer.go
 * @Version: 1.0.0
 * @Date: 2019/12/6 17:44
 */
package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/elb"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func GetLoadBalancers(sess *session.Session) {
	svc := elb.New(sess)
	input := &elb.DescribeLoadBalancersInput{
		PageSize: aws.Int64(100),
	}

	result, err := svc.DescribeLoadBalancers(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elb.ErrCodeAccessPointNotFoundException:
				fmt.Println(elb.ErrCodeAccessPointNotFoundException, aerr.Error())
			case elb.ErrCodeDependencyThrottleException:
				fmt.Println(elb.ErrCodeDependencyThrottleException, aerr.Error())
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

	fmt.Println(result)
}

func GetLoadBalancersv2(sess *session.Session) {
	svc := elbv2.New(sess)
	input := &elbv2.DescribeLoadBalancersInput{
		PageSize: aws.Int64(100),
	}

	result, err := svc.DescribeLoadBalancers(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elbv2.ErrCodeLoadBalancerNotFoundException:
				fmt.Println(elbv2.ErrCodeLoadBalancerNotFoundException, aerr.Error())
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

	fmt.Println(result)
}
