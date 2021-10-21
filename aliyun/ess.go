/**
 * @Author: jie.an
 * @Description:
 * @File:  ess.go
 * @Version: 1.0.0
 * @Date: 2021/10/18 18:05
 */

package aliyun

import (
	ess20140828 "github.com/alibabacloud-go/ess-20140828/client"
	"github.com/alibabacloud-go/tea/tea"
	"golang-base/tools"
	"strings"
)

type ScaleRule struct {
	Ari         string
	Id          string
	AlarmTaskId string
}

type AutoScaling struct {
	Id            string
	ScaleUpRule   ScaleRule
	ScaleDownRule ScaleRule
	Config        *Config
	Client        *ess20140828.Client
}

func (a *AutoScaling) Init(c *Config) *AutoScaling {
	// 访问的域名
	a.Config = c
	a.Config.Conf.Endpoint = tea.String("ess.aliyuncs.com")
	a.Client = &ess20140828.Client{}
	client, err := ess20140828.NewClient(a.Config.Conf)
	if err != nil {
		tools.WarningLogger.Println(err)
	}
	a.Client = client
	return a
}

func (a *AutoScaling) describeScalingRules() (*AutoScaling, error) {
	describeScalingRulesRequest := &ess20140828.DescribeScalingRulesRequest{
		RegionId:       tea.String(a.Config.UsedRegion),
		ScalingGroupId: tea.String(a.Id),
		ShowAlarmRules: tea.Bool(true),
	}
	_resp, _err := a.Client.DescribeScalingRules(describeScalingRulesRequest)
	if _err != nil {
		tools.WarningLogger.Println(_err)
		return nil, _err
	}
	if _resp.Body.ScalingRules != nil {
		for _, rule := range _resp.Body.ScalingRules.ScalingRule {
			if strings.Contains(*rule.ScalingRuleName, "remove") {
				a.ScaleDownRule.Id = *rule.ScalingRuleId
				a.ScaleDownRule.Ari = *rule.ScalingRuleAri
				if rule.Alarms != nil {
					a.ScaleDownRule.AlarmTaskId = *rule.Alarms.Alarm[0].AlarmTaskId
				}
			} else if strings.Contains(*rule.ScalingRuleName, "add") {
				a.ScaleUpRule.Id = *rule.ScalingRuleId
				a.ScaleUpRule.Ari = *rule.ScalingRuleAri
				if rule.Alarms != nil {
					a.ScaleUpRule.AlarmTaskId = *rule.Alarms.Alarm[0].AlarmTaskId
				}
			} else {
				tools.WarningLogger.Println("Illegal Rule Name:", *rule.ScalingRuleName)
			}
		}
	}
	return a, nil
}

func (a *AutoScaling) DisableScaleDown(config *Config) (*AutoScaling, error) {
	tools.InfoLogger.Println("Disable scale-down for asg:", a.Id)
	// Create an ess client.
	a.Init(config)
	a.describeScalingRules()
	disableAlarmRequest := &ess20140828.DisableAlarmRequest{
		RegionId:    tea.String(config.UsedRegion),
		AlarmTaskId: tea.String(a.ScaleDownRule.AlarmTaskId),
	}
	_, _err := a.Client.DisableAlarm(disableAlarmRequest)
	if _err != nil {
		tools.WarningLogger.Println(_err)
		return nil, _err
	}
	tools.InfoLogger.Printf("Disable scale-down for asg %s successfully. \n", a.Id)
	return a, nil
}

func (a *AutoScaling) EnableScaleDown(config *Config) (*AutoScaling, error) {
	tools.InfoLogger.Println("Enable scale-down for asg:", a.Id)
	// Create an ess client.
	a.Init(config)
	a.describeScalingRules()
	enableAlarmRequest := &ess20140828.EnableAlarmRequest{
		RegionId:    tea.String(config.UsedRegion),
		AlarmTaskId: tea.String(a.ScaleDownRule.AlarmTaskId),
	}
	_, _err := a.Client.EnableAlarm(enableAlarmRequest)
	if _err != nil {
		tools.WarningLogger.Println(_err)
		return nil, _err
	}
	tools.InfoLogger.Printf("Enable scale-down for asg %s successfully. \n", a.Id)
	return a, nil
}

func (a *AutoScaling) ExecuteScaleUp(config *Config) (*AutoScaling, error) {
	tools.InfoLogger.Println("Execute scale-up for asg:", a.Id)
	// Create an ess client.
	a.Init(config)
	a.describeScalingRules()
	executeScalingRuleRequest := &ess20140828.ExecuteScalingRuleRequest{
		ScalingRuleAri: tea.String(a.ScaleUpRule.Ari),
	}
	_, _err := a.Client.ExecuteScalingRule(executeScalingRuleRequest)
	if _err != nil {
		tools.WarningLogger.Println(_err)
		return nil, _err
	}
	tools.InfoLogger.Printf("Execute scale-up for asg %s successfully. \n", a.Id)
	return a, nil
}
