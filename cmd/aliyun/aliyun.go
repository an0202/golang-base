/**
 * @Author: jie.an
 * @Description:
 * @File:  aliyun.go
 * @Version: 1.0.0
 * @Date: 2021/3/7 10:27 下午
 */

package cmd

import (
	"github.com/spf13/cobra"
	"golang-base/cmd"
)

// aliyunCmd
var AliyunCmd = &cobra.Command{
	Use:   "aliyun",
	Short: "Handle aliyun",
	Long: `Handle aliyun tool set,
`,
}

func init() {
	cmd.RootCmd.AddCommand(AliyunCmd)
	AliyunCmd.PersistentFlags().StringP("region", "r", "us-east-1", `region`)
}
