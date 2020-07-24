/**
 * @Author: jie.an
 * @Description:
 * @File:  ec2.go
 * @Version: 1.0.0
 * @Date: 2020/7/24 14:34
 */
package cmd

import (
	"github.com/spf13/cobra"
	"golang-base/aws"
)

// ec2Cmd
var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "RunInstance",
	Long: `Launching an instance using parameters from an existing instance
`,
	Run: func(cmd *cobra.Command, args []string) {
		amiId, _ := cmd.Flags().GetString("ami")
		region, _ := cmd.Flags().GetString("region")
		instanceId, _ := cmd.Flags().GetString("id")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		sess := aws.InitSession(region)
		aws.EC2LaunchMoreLikeThis(sess, instanceId, amiId, dryRun)
	},
}

func init() {
	RootCmd.AddCommand(ec2Cmd)
	ec2Cmd.Flags().StringP("ami", "a", "ami-abc123", `ami to use`)
	ec2Cmd.Flags().StringP("id", "i", "i-abc123", `origin instance id`)
	ec2Cmd.Flags().StringP("region", "r", "cn-north-1", `region to use`)
	ec2Cmd.Flags().Bool("dry-run", true, "dry run test")
}
