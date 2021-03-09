/**
 * @Author: jie.an
 * @Description:
 * @File:  slb.go
 * @Version: 1.0.0
 * @Date: 2021/03/07 18:50
 */
package aliyun

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"golang-base/tools"
)

type LoadBalancer struct {
	Id            string
	BackendServer []slb.BackendServerInDescribeLoadBalancerAttribute
}

func DescribeSLB(pro Provider, slbId string) (*LoadBalancer, error) {
	lb := new(LoadBalancer)
	// Create an slb service client.
	client, err := slb.NewClientWithProvider(pro.UsedRegion, pro.Pro)
	if err != nil {
		tools.WarningLogger.Println(err)
		return nil, err
	}
	// Get lb attribute
	request := slb.CreateDescribeLoadBalancerAttributeRequest()
	request.Scheme = "https"
	request.LoadBalancerId = slbId
	response, err := client.DescribeLoadBalancerAttribute(request)
	if err != nil {
		return nil, err
	}
	for _, v := range response.BackendServers.BackendServer {
		lb.BackendServer = append(lb.BackendServer, v)
	}
	return lb, nil
}

func AddBackEndServer(pro Provider, balancer LoadBalancer) error {
	// Create an slb service client.
	client, err := slb.NewClientWithProvider(pro.UsedRegion, pro.Pro)
	if err != nil {
		tools.WarningLogger.Println(err)
		return err
	}
	// Add backend server
	request := slb.CreateAddBackendServersRequest()
	request.Scheme = "https"
	request.LoadBalancerId = balancer.Id
	var backendServerList []slb.BackendServerInDescribeLoadBalancerAttribute
	for i, v := range balancer.BackendServer {
		backendServerList = append(backendServerList, v)
		if len(backendServerList) == 20 {
			serverList, err := json.Marshal(backendServerList)
			if err != nil {
				tools.ErrorLogger.Fatalln(err)
			}
			request.BackendServers = string(serverList)
			backendServerList = []slb.BackendServerInDescribeLoadBalancerAttribute{}
			_, err = client.AddBackendServers(request)
			if err != nil {
				return err
			}
		}
		if i == len(balancer.BackendServer)-1 {
			serverList, err := json.Marshal(backendServerList)
			if err != nil {
				tools.ErrorLogger.Fatalln(err)
			}
			request.BackendServers = string(serverList)
			_, err = client.AddBackendServers(request)
			if err != nil {
				return err
			}
		}
	}
	return err
}
