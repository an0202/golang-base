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
	"github.com/aws/aws-sdk-go/aws/session"
)

//List ElastiCache
func ListElastiCache(sess *session.Session) (CacheList [][]interface{}) {
	// Create an elasticache service client.
	svc := elasticache.New(sess)
	// Get elasticache cluster
	output, err := svc.DescribeCacheClusters(&elasticache.DescribeCacheClustersInput{
		// Todo: Unhandle error for numbers of cluster > 100
		MaxRecords: aws.Int64(100),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				tools.ErrorLogger.Fatalln(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			tools.ErrorLogger.Fatalln(err.Error())
		}
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
		cache = append(cache, *cachecluster.CacheClusterId, *cachecluster.NumCacheNodes, *cachecluster.CacheNodeType,
			*cachecluster.Engine, *cachecluster.EngineVersion, *cachecluster.CacheSubnetGroupName, *cachecluster.PreferredMaintenanceWindow,
			*cachecluster.SnapshotRetentionLimit, sgs)
		CacheList = append(CacheList, cache)
	}
	return CacheList
}
