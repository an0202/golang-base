/**
 * @Author: jie.an
 * @Description:
 * @File:  slb.go
 * @Version: 1.0.0
 * @Date: 2021/3/7 18:40
 */
package cmd

import (
	"github.com/spf13/cobra"
	"golang-base/aliyun"
	"golang-base/tools"
)

// SlbCmd
var SlbCmd = &cobra.Command{
	Use:   "slb",
	Short: "Aliyun SLB migrate tool",
	Long: `Migrate the backend server from src slb to dest slb
`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		region, _ := cmd.Flags().GetString("region")
		destId, _ := cmd.Flags().GetString("dest")
		srcId, _ := cmd.Flags().GetString("src")
		if file != "" {
			tools.ErrorLogger.Fatalln("Not support yet")
			slbIds := tools.GetRecords(file)
			pro := aliyun.NewProvider(region)
			for _, slbId := range slbIds {
				aliyun.DescribeSLB(*pro, slbId)
			}
		} else {
			tools.InfoLogger.Printf("Migrate the backend server from %s to %s .", srcId, destId)
			pro := aliyun.NewProvider(region)
			lb, err := aliyun.DescribeSLB(*pro, srcId)
			if err != nil {
				tools.ErrorLogger.Fatalln(err)
			}
			// set lb id
			lb.Id = destId
			err = aliyun.AddBackEndServer(*pro, *lb)
			if err != nil {
				tools.WarningLogger.Println(err)
			} else {
				tools.InfoLogger.Println("Migrate done.")
			}
		}
	},
}

func init() {
	aliyunCmd.AddCommand(SlbCmd)
	SlbCmd.Flags().StringP("file", "f", "", "file with slbId line by line")
	SlbCmd.Flags().StringP("region", "r", "us-east-1", `region`)
	SlbCmd.Flags().StringP("src", "s", "lb-0xia3woxxxxxepdy76t", `src slb id`)
	SlbCmd.Flags().StringP("dest", "d", "lb-0xil9xoxxxxxd8ty9hyc", `dest slb id`)
}
