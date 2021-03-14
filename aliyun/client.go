/**
 * @Author: jie.an
 * @Description:
 * @File:  client.go
 * @Version: 1.0.0
 * @Date: 2021/3/12 20:20
 */

package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
	"os"
)

type Config struct {
	UsedRegion string
	Conf       *openapi.Config
}

func NewConfig(region string, opts ...ConfigOption) *Config {
	aliConfig := openapi.Config{
		RegionId: tea.String(region),
	}
	c := &Config{
		UsedRegion: region,
		Conf:       &aliConfig,
	}
	if len(opts) == 0 {
		// use env Config as default Config
		WithEnv()(c)
	} else {
		// Loop through each option
		for _, opt := range opts {
			// Call the option giving the instantiated
			// *client as the argument
			opt(c)
		}
	}

	// return the modified client instance
	return c
}

type ConfigOption func(*Config)

func WithEnv() ConfigOption {
	return func(c *Config) {
		c.Conf.RegionId = tea.String(c.UsedRegion)
		c.Conf.AccessKeyId = tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"))
		c.Conf.AccessKeySecret = tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"))
	}
}
