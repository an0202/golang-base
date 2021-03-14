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
	"golang-base/huaweicloud"
	"golang-base/tools"
)

// EcsCmd
var EcsCmd = &cobra.Command{
	Use:   "ecs",
	Short: "HuaweiCloud ECS tool",
	Long: `HuaweiCloud ECS tool set
`,
}

// ChangeOSCmd
var ChangeOSCmd = &cobra.Command{
	Use:   "change-os",
	Short: "Change ecs instance os",
	Long: `Change ecs instance os and set new user-data
`,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		id, _ := cmd.Flags().GetString("id")
		file, _ := cmd.Flags().GetString("file")
		image, _ := cmd.Flags().GetString("image")
		auth := huaweicloud.NewAuth(region)
		instance := new(huaweicloud.Instance)
		instance.Id = id
		instance.ImageId = image
		userData := tools.ReadFileToBase64String(file)
		instance.UserData = &userData
		_ = instance.ChangeOSWithUserData(auth)
	},
}

func init() {
	EcsCmd.AddCommand(ChangeOSCmd)
	HuaweiCloudCmd.AddCommand(EcsCmd)
	ChangeOSCmd.Flags().String("image", "m-0xick0io750jvtlqy3ef", "image id")
	EcsCmd.PersistentFlags().StringP("file", "f", "user-data.txt", "path to file(user-data)")
	EcsCmd.PersistentFlags().StringP("id", "i", "49b54826-axxxx-426c-ae6c-35dc45c498xx", `instance id`)
}
