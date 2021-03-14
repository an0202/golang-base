/**
 * @Author: jie.an
 * @Description:
 * @File:  huaweicloud.go
 * @Version: 1.0.0
 * @Date: 2021/3/14 19:17 下午
 */

package cmd

import (
	"github.com/spf13/cobra"
	"golang-base/cmd"
)

// huaweicloudCmd
var HuaweiCloudCmd = &cobra.Command{
	Use:   "huaweicloud",
	Short: "Handle huaweicloud",
	Long: `Handle huaweicloud tool set,
`,
}

func init() {
	cmd.RootCmd.AddCommand(HuaweiCloudCmd)
	HuaweiCloudCmd.PersistentFlags().StringP("region", "r", "ap-southeast-1", `region`)
}
