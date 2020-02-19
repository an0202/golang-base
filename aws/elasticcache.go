/**
 * @Author: jie.an
 * @Description:
 * @File:  elasticcache.go
 * @Version: 1.0.0
 * @Date: 2020/1/22 17:29
 */
package aws

import (
	"golang-base/tools"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

//List ElastiCaches
func ListECCs(se Session) (CacheList [][]interface{}) {
	// Create an elasticache service client.
	svc := elasticache.New(se.Sess)
	// Get elasticache cluster
	output, err := svc.DescribeCacheClusters(&elasticache.DescribeCacheClustersInput{
		// Todo: Unhandle error for numbers of cluster > 100
		MaxRecords: aws.Int64(100),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				tools.ErrorLogger.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			tools.ErrorLogger.Println(err.Error())
		}
		return
	}
	for _, cachecluster := range output.CacheClusters {
		var cache []interface{}
		//handle securitygroups
		var sgs []string
		if len(cachecluster.SecurityGroups) == 0 {
			sgs = append(sgs, "N/A")
		} else {
			for _, sg := range cachecluster.SecurityGroups {
				sgs = append(sgs, *sg.SecurityGroupId)
			}
		}
		if len(output.CacheClusters) >= 100 {
			// todo cluster > 100
			tools.WarningLogger.Println("Number Of Clusters > 100 , Data May Missing.")
		}
		cache = append(cache, se.AccountId, se.UsedRegion,*cachecluster.CacheClusterId, *cachecluster.NumCacheNodes, *cachecluster.CacheNodeType,
			*cachecluster.Engine, *cachecluster.EngineVersion, *cachecluster.CacheSubnetGroupName, *cachecluster.PreferredMaintenanceWindow,
			*cachecluster.SnapshotRetentionLimit, sgs)
		CacheList = append(CacheList, cache)
	}
	return CacheList
}
