/**
 * @Author: jie.an
 * @Description:
 * @File:  root.go
 * @Version: 0.0.2
 * @Date: 2021/3/7 14:34
 */
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const VERSION = "0.0.2"

var RootCmd = &cobra.Command{
	Use:   "base-tool",
	Short: "AWS/Aliyun resource processing tools",
	Long: fmt.Sprintf(`AWS/Aliyun resource processing tools.

Version: %s`,
		VERSION),
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
