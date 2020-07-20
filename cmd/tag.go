/**
 * @Author: jie.an
 * @Description:
 * @File:  tag
 * @Version: 1.0.0
 * @Date: 2020/7/20 17:37
 */
package cmd

import (
	"flag"
	"github.com/spf13/cobra"
	"golang-base/aws"
	"golang-base/excel"
	"golang-base/tools"
	"strings"
)

// tagCmd
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Tag Assets",
	Long: `Get/Add/Del Tag For Assets,using local AWS_PROFILE
`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		sheetName, _ := cmd.Flags().GetString("sheet")
		region, _ := cmd.Flags().GetString("region")
		method, _ := cmd.Flags().GetString("method")
		tags, _ := cmd.Flags().GetString("tags")
		override, _ := cmd.Flags().GetBool("override")
		switch method {
		case "add":
			defaultSess := aws.InitSession(region)
			se := new(aws.Session)
			a := excel.ReadToMaps(file, sheetName)
			for _, v := range a {
				b := aws.EC2InstanceMarshal(v)
				if b.AWSProfile != "" && b.Region != "" {
					// Use the old session when the current resource and the previous resource belong to the same account and region
					if b.AWSProfile == se.UsedAwsProfile && b.Region == se.UsedRegion {
						//fmt.Println(se.Sess.Config.Credentials)
						aws.EC2CreateTags(se.Sess, b, override)
					} else {
						// create a new session
						se.InitSessionWithAWSProfile(b.Region, b.AWSProfile)
						//fmt.Println(se.Sess.Config.Credentials)
						aws.EC2CreateTags(se.Sess, b, override)
					}
					//se.Sess = se.InitSessionWithAWSProfile(b.Region,b.AWSProfile)
					//fmt.Println(se.Sess.Config.Credentials)
					//aws.EC2CreateTags(se.Sess, b, *overide)
				} else {
					tools.InfoLogger.Println("Use Default Session")
					aws.EC2CreateTags(defaultSess, b, override)
				}
			}
		case "del":
			sess := aws.InitSession(region)
			a := excel.ReadToMaps(file, sheetName)
			for _, v := range a {
				b := aws.EC2InstanceMarshal(v)
				aws.EC2DeleteTags(sess, b)
			}
		case "get":
			defaultSess := aws.InitSession(region)
			se := new(aws.Session)
			a := excel.ReadToMaps(file, sheetName)
			var results [][]interface{}
			for _, v := range a {
				b := aws.EC2InstanceMarshal(v)
				if b.AWSProfile != "" && b.Region != "" {
					// Use the old session when the current resource and the previous resource belong to the same account and region
					if b.AWSProfile == se.UsedAwsProfile && b.Region == se.UsedRegion {
						//fmt.Println(se.Sess.Config.Credentials)
						_, result := aws.EC2GetTags(se.Sess, b, tags)
						results = append(results, result)
					} else {
						// create a new session
						se.InitSessionWithAWSProfile(b.Region, b.AWSProfile)
						// fmt.Println(se.Sess.Config.Credentials)
						_, result := aws.EC2GetTags(se.Sess, b, tags)
						results = append(results, result)
					}
				} else {
					tools.InfoLogger.Println("Use Default Session")
					_, result := aws.EC2GetTags(defaultSess, b, tags)
					results = append(results, result)
				}
			}
			var headerline = []interface{}{"ResourceId"}
			for _, v := range strings.Split(tags, ",") {
				headerline = append(headerline, v)
			}
			excel.CreateFile("output-" + file)
			excel.SetHeadLine("output-"+file, "result", headerline)
			excel.SetListRows("output-"+file, "result", results)
		default:
			flag.Usage()
			tools.ErrorLogger.Fatalln("Illegal Method:", method)
		}
	},
}

func init() {
	RootCmd.AddCommand(tagCmd)
	tagCmd.Flags().StringP("file", "f", "tags.xlsx", "Read Config From Excel Line By Line")
	tagCmd.Flags().StringP("sheet", "s", "get", `Sheet With Config To Be Process`)
	tagCmd.Flags().StringP("region", "r", "cn-north-1", `Used For Init A AWS Default Session`)
	tagCmd.Flags().StringP("method", "m", "get", "add/del/get Tags")
	tagCmd.Flags().String("tags", "Name,Env,Project", `Require:[ method = get],Get Specific Tags From Resource And Write To Excel`)
	tagCmd.Flags().Bool("o", false, `Override Exist Tags`)
}
