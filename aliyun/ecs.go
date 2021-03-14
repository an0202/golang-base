/**
 * @Author: jie.an
 * @Description:
 * @File:  ecs.go
 * @Version: 1.0.0
 * @Date: 2021/3/13 20:45
 */

package aliyun

import (
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"golang-base/tools"
)

type Instance struct {
	Id             string
	UserData       *string
	SystemDiskSize int32
	ImageId        string
}

func createECSClient(c *Config) (_result *ecs20140526.Client, _err error) {
	// 访问的域名
	c.Conf.Endpoint = tea.String("ecs.aliyuncs.com")
	_result = &ecs20140526.Client{}
	_result, _err = ecs20140526.NewClient(c.Conf)
	return _result, _err
}

func (i *Instance) DescribeInstanceUserData(config *Config) (*Instance, error) {
	tools.InfoLogger.Println("Get User-data from instance:", i.Id)
	// Create an ecs client.
	client, _err := createECSClient(config)
	if _err != nil {
		tools.WarningLogger.Println(_err)
		return nil, _err
	}
	describeUserDataRequest := &ecs20140526.DescribeUserDataRequest{
		RegionId:   tea.String(config.UsedRegion),
		InstanceId: tea.String(i.Id),
	}
	_resp, _err := client.DescribeUserData(describeUserDataRequest)
	if _err != nil {
		tools.WarningLogger.Println(_err)
		return nil, _err
	}
	i.UserData = _resp.Body.UserData
	userData, _ := tools.DecodeBase64String(*i.UserData)
	tools.InfoLogger.Printf("User-data is: \n%s", userData)
	return i, nil
}

func (i *Instance) SetInstanceUserData(config *Config) error {
	// Create an ecs client.
	tools.InfoLogger.Println("Set User-data to instance:", i.Id)
	client, _err := createECSClient(config)
	if _err != nil {
		tools.WarningLogger.Println(_err)
		return _err
	}
	modifyInstanceAttributeRequest := &ecs20140526.ModifyInstanceAttributeRequest{
		InstanceId: tea.String(i.Id),
		UserData:   i.UserData,
	}
	_, _err = client.ModifyInstanceAttribute(modifyInstanceAttributeRequest)
	if _err != nil {
		tools.WarningLogger.Println(_err)
		return _err
	}
	// verify
	i.DescribeInstanceUserData(config)
	return nil
}

func (i *Instance) ChangeOS(config *Config) error {
	tools.InfoLogger.Println("ChangeOS with user-data on instance:", i.Id)
	// Create an ecs client.
	client, _err := createECSClient(config)
	if _err != nil {
		tools.WarningLogger.Println(_err)
		return _err
	}
	systemDisk := &ecs20140526.ReplaceSystemDiskRequestSystemDisk{
		Size: tea.Int32(i.SystemDiskSize),
	}
	replaceSystemDiskRequest := &ecs20140526.ReplaceSystemDiskRequest{
		InstanceId: tea.String(i.Id),
		ImageId:    tea.String(i.ImageId),
		SystemDisk: systemDisk,
	}
	_, _err = client.ReplaceSystemDisk(replaceSystemDiskRequest)
	if _err != nil {
		tools.WarningLogger.Println(_err)
		return _err
	}
	return nil
}
