/**
 * @Author: jie.an
 * @Description:
 * @File:  ecs.go
 * @Version: 1.0.0
 * @Date: 2021/3/14 17:00
 */

package huaweicloud

import (
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/region"
	"golang-base/tools"
)

type Instance struct {
	Id       string
	UserData *string
	ImageId  string
}

func createECSClient(a *Auth) (_result *ecs.EcsClient, _err error) {
	_result = ecs.NewEcsClient(
		ecs.EcsClientBuilder().
			WithRegion(region.ValueOf(a.UsedRegion)).
			WithCredential(a.Cred).
			Build())
	return _result, _err
}

func (i *Instance) ChangeOSWithUserData(auth *Auth) error {
	tools.InfoLogger.Println("ChangeOS with user-data on instance:", i.Id)
	// Create an ecs client.
	client, _err := createECSClient(auth)
	if _err != nil {
		tools.WarningLogger.Println(_err)
		return _err
	}
	request := &model.ChangeServerOsWithCloudInitRequest{}
	request.ServerId = i.Id
	metadataOsChange := &model.ChangeSeversOsMetadata{
		UserData: i.UserData,
	}
	osChangeBody := &model.ChangeServerOsWithCloudInitOption{
		Imageid:  i.ImageId,
		Metadata: metadataOsChange,
	}
	request.Body = &model.ChangeServerOsWithCloudInitRequestBody{
		OsChange: osChangeBody,
	}
	_, _err = client.ChangeServerOsWithCloudInit(request)
	if _err != nil {
		tools.WarningLogger.Println(_err)
		return _err
	}
	return nil
}
