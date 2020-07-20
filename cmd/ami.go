/**
 * @Author: jie.an
 * @Description:
 * @File:  ami.go
 * @Version: 1.0.0
 * @Date: 2020/7/20 14:56
 */
package cmd

import (
	"github.com/spf13/cobra"
	"golang-base/aws"
	"golang-base/tools"
)

// amiCmd
var amiCmd = &cobra.Command{
	Use:   "ami",
	Short: "Create AMI",
	Long: `Create EC2 AMI using local AWS_PROFILE, 
use config file to pass in multiple InstanceIDs at once(One InstanceId per line),
`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		region, _ := cmd.Flags().GetString("region")
		instanceId, _ := cmd.Flags().GetString("id")
		suffix, _ := cmd.Flags().GetString("suffix")
		if file != "" {
			instanceIds := tools.GetRecords(file)
			sess := aws.InitSession(region)
			for _, instanceId := range instanceIds {
				aws.CreateImage(sess, instanceId, suffix)
			}
		} else {
			sess := aws.InitSession(region)
			aws.CreateImage(sess, instanceId, suffix)
		}
	},
}

func init() {
	RootCmd.AddCommand(amiCmd)
	amiCmd.Flags().StringP("file", "f", "", "file with instanceId line by line")
	amiCmd.Flags().StringP("region", "r", "cn-north-1", `region`)
	amiCmd.Flags().StringP("id", "i", "i-abc123", `instance id`)
	amiCmd.Flags().String("suffix", "date", "Add date/final/.. Suffix To AMI Name")
}
