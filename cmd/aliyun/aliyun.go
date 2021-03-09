/**
 * @Author: jie.an
 * @Description:
 * @File:  aliyun.go
 * @Version: 1.0.0
 * @Date: 2021/3/7 10:27 下午
 */

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang-base/cmd"
)

// aliyunCmd
var aliyunCmd = &cobra.Command{
	Use:   "aliyun",
	Short: "Handle aliyun",
	Long: `Handle aliyun tool set,
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("aliyun")
	},
}

func init() {
	cmd.RootCmd.AddCommand(aliyunCmd)
}
