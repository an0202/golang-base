/**
 * @Author: jie.an
 * @Description:
 * @File:  ess.go
 * @Version: 1.0.0
 * @Date: 2021/10/18 18:13
 */
package cmd

import (
	"github.com/spf13/cobra"
	"golang-base/aliyun"
	"golang-base/tools"
)

// EssCmd
var EssCmd = &cobra.Command{
	Use:   "ess",
	Short: "Aliyun ESS tool",
	Long: `Aliyun ESS tool set
`,
}

// DisableScaleDownCmd
var DisableScaleDownCmd = &cobra.Command{
	Use:   "disable-scaledown",
	Short: "Disable Scaling Down for auto scaling group",
	Long: `Disable Scaling Down for auto scaling group
`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		region, _ := cmd.Flags().GetString("region")
		id, _ := cmd.Flags().GetString("id")
		conf := aliyun.NewConfig(region)
		if file != "" {
			asgIds := tools.GetRecords(file)
			asg := new(aliyun.AutoScaling)
			for _, asgId := range asgIds {
				asg.Id = asgId
				asg.DisableScaleDown(conf)
			}
		} else {
			asg := new(aliyun.AutoScaling)
			asg.Id = id
			asg.DisableScaleDown(conf)
		}
	},
}

// EnableScaleDownCmd
var EnableScaleDownCmd = &cobra.Command{
	Use:   "enable-scaledown",
	Short: "Enable Scaling Down for auto scaling group",
	Long: `Enable Scaling Down for auto scaling group
`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		region, _ := cmd.Flags().GetString("region")
		id, _ := cmd.Flags().GetString("id")
		conf := aliyun.NewConfig(region)
		if file != "" {
			asgIds := tools.GetRecords(file)
			asg := new(aliyun.AutoScaling)
			for _, asgId := range asgIds {
				asg.Id = asgId
				asg.EnableScaleDown(conf)
			}
		} else {
			asg := new(aliyun.AutoScaling)
			asg.Id = id
			asg.EnableScaleDown(conf)
		}
	},
}

// ExecuteScaleUpCmd
var ExecuteScaleUpCmd = &cobra.Command{
	Use:   "execute-scaleup",
	Short: "Execute Scale Up for auto scaling group",
	Long: `Execute Scale Up for auto scaling group
`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		region, _ := cmd.Flags().GetString("region")
		id, _ := cmd.Flags().GetString("id")
		conf := aliyun.NewConfig(region)
		if file != "" {
			asgIds := tools.GetRecords(file)
			asg := new(aliyun.AutoScaling)
			for _, asgId := range asgIds {
				asg.Id = asgId
				asg.ExecuteScaleUp(conf)
			}
		} else {
			asg := new(aliyun.AutoScaling)
			asg.Id = id
			asg.ExecuteScaleUp(conf)
		}
	},
}

func init() {
	EssCmd.AddCommand(DisableScaleDownCmd)
	EssCmd.AddCommand(EnableScaleDownCmd)
	EssCmd.AddCommand(ExecuteScaleUpCmd)
	AliyunCmd.AddCommand(EssCmd)
	EssCmd.PersistentFlags().StringP("file", "f", "", "path to file(asg id)")
	EssCmd.PersistentFlags().StringP("id", "i", "asg-0xi5phe7f8y943bv0tjn", `auto scaling group id`)
}
