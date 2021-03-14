/**
 * @Author: jie.an
 * @Description:
 * @File:  slb.go
 * @Version: 1.0.0
 * @Date: 2021/03/12 18:50
 */
package aliyun

import (
	"encoding/json"
	slb20140515 "github.com/alibabacloud-go/slb-20140515/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"golang-base/tools"
)

type LoadBalancer struct {
	Id             string
	BackendServers []*slb20140515.DescribeLoadBalancerAttributeResponseBodyBackendServersBackendServer
}

func createSLBClient(c *Config) (_result *slb20140515.Client, _err error) {
	// 访问的域名
	c.Conf.Endpoint = tea.String("slb.aliyuncs.com")
	_result = &slb20140515.Client{}
	_result, _err = slb20140515.NewClient(c.Conf)
	return _result, _err
}

func DescribeSLB(config *Config, slbId string) (*LoadBalancer, error) {
	lb := new(LoadBalancer)
	// Create an slb client.
	client, _err := createSLBClient(config)
	if _err != nil {
		tools.WarningLogger.Println(_err)
		return nil, _err
	}
	describeLoadBalancerAttributeRequest := &slb20140515.DescribeLoadBalancerAttributeRequest{
		RegionId:       client.RegionId,
		LoadBalancerId: tea.String(slbId),
	}
	_resp, _err := client.DescribeLoadBalancerAttribute(describeLoadBalancerAttributeRequest)
	if _err != nil {
		tools.WarningLogger.Println(_err)
		return nil, _err
	}
	lb.Id = *_resp.Body.RegionId
	lb.BackendServers = _resp.Body.BackendServers.BackendServer
	return lb, _err
}

func AddBackEndServer(config *Config, balancer *LoadBalancer) error {
	// Create an slb client.
	client, _err := createSLBClient(config)
	if _err != nil {
		tools.WarningLogger.Println(_err)
		return _err
	}
	// Add backend server
	var backendServerList []slb20140515.DescribeLoadBalancerAttributeResponseBodyBackendServersBackendServer
	for i, v := range balancer.BackendServers {
		backendServerList = append(backendServerList, *v)
		if len(backendServerList) == 20 {
			serverList, err := json.Marshal(backendServerList)
			if err != nil {
				tools.ErrorLogger.Fatalln(err)
			}
			addBackendServersRequest := &slb20140515.AddBackendServersRequest{
				RegionId:       client.RegionId,
				LoadBalancerId: tea.String(balancer.Id),
				BackendServers: tea.String(string(serverList)),
			}
			backendServerList = []slb20140515.DescribeLoadBalancerAttributeResponseBodyBackendServersBackendServer{}
			_, err = client.AddBackendServers(addBackendServersRequest)
			if err != nil {
				return err
			}
		}
		if i == len(balancer.BackendServers)-1 {
			serverList, err := json.Marshal(backendServerList)
			if err != nil {
				tools.ErrorLogger.Fatalln(err)
			}
			addBackendServersRequest := &slb20140515.AddBackendServersRequest{
				RegionId:       client.RegionId,
				LoadBalancerId: tea.String(balancer.Id),
				BackendServers: tea.String(string(serverList)),
			}
			_, err = client.AddBackendServers(addBackendServersRequest)
			if err != nil {
				return err
			}
		}
	}
	return _err
}
