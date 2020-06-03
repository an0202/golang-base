/**
 * @Author: jie.an
 * @Description:
 * @File:  conf.go
 * @Version: 1.0.0
 * @Date: 2020/03/09 10:52
 */
package tools

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type YamlConf interface {
	getConf()
}

type HttpHeaderConf struct {
	//Host   string `yaml:"host"`
	Authorization string `yaml:"authorization"`
	Cookie        string `yaml:"cookie"`
}

func (c *HttpHeaderConf) GetConf(path string) *HttpHeaderConf {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}
