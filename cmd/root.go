/**
 * @Author: jie.an
 * @Description:
 * @File:  root.go
 * @Version: 1.0.0
 * @Date: 2020/7/20 14:34
 */
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const VERSION = "0.0.1"

var RootCmd = &cobra.Command{
	Use:   "aws-base",
	Short: "AWS resource processing tools",
	Long: fmt.Sprintf(`AWS resource processing tools: get assets,  create ami, create tag.

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
