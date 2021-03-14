/**
 * @Author: jie.an
 * @Description:
 * @File:  ecs.go
 * @Version: 1.0.0
 * @Date: 2021/3/13 21:05
 */
package cmd

import (
	"github.com/spf13/cobra"
	"golang-base/aliyun"
	"golang-base/tools"
)

// EcsCmd
var EcsCmd = &cobra.Command{
	Use:   "ecs",
	Short: "Aliyun ECS tool",
	Long: `Aliyun ECS tool set
`,
}

// GetUserDataCmd
var GetUserDataCmd = &cobra.Command{
	Use:   "get-userdata",
	Short: "Get user data from ecs instance",
	Long: `Get user data from ecs instance
`,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		id, _ := cmd.Flags().GetString("id")
		conf := aliyun.NewConfig(region)
		instance := new(aliyun.Instance)
		instance.Id = id
		instance.DescribeInstanceUserData(conf)
	},
}

// SetUserDataCmd
var SetUserDataCmd = &cobra.Command{
	Use:   "set-userdata",
	Short: "Set user data to ecs instance",
	Long: `Read user-data from file and set to ecs instance
`,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		id, _ := cmd.Flags().GetString("id")
		file, _ := cmd.Flags().GetString("file")
		conf := aliyun.NewConfig(region)
		instance := new(aliyun.Instance)
		instance.Id = id
		userData := tools.ReadFileToBase64String(file)
		instance.UserData = &userData
		instance.SetInstanceUserData(conf)
	},
}

// ChangeOSCmd
var ChangeOSCmd = &cobra.Command{
	Use:   "change-os",
	Short: "Change ecs instance os",
	Long: `Change ecs instance os and set new disk size and user-data
`,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		id, _ := cmd.Flags().GetString("id")
		file, _ := cmd.Flags().GetString("file")
		image, _ := cmd.Flags().GetString("image")
		size, _ := cmd.Flags().GetInt32("size")
		conf := aliyun.NewConfig(region)
		instance := new(aliyun.Instance)
		instance.Id = id
		instance.SystemDiskSize = size
		instance.ImageId = image
		userData := tools.ReadFileToBase64String(file)
		instance.UserData = &userData
		instance.SetInstanceUserData(conf)
		instance.ChangeOS(conf)
	},
}

func init() {
	EcsCmd.AddCommand(GetUserDataCmd)
	EcsCmd.AddCommand(SetUserDataCmd)
	EcsCmd.AddCommand(ChangeOSCmd)
	AliyunCmd.AddCommand(EcsCmd)
	ChangeOSCmd.Flags().String("image", "m-0xick0io750jvtlqy3ef", "image id")
	ChangeOSCmd.Flags().Int32("size", 40, "system image size")
	EcsCmd.PersistentFlags().StringP("file", "f", "user-data.txt", "path to file(user-data)")
	EcsCmd.PersistentFlags().StringP("id", "i", "i-xxxxx88888", `instance id`)
}
