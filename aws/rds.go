/**
 * @Author: jie.an
 * @Description:
 * @File:  rds.go
 * @Version: 1.0.0
 * @Date: 2020/02/18 14:00
 */
package aws

import (
	"github.com/aws/aws-sdk-go/service/rds"
	"golang-base/tools"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

//List DBs
func ListDBs(se Session) (DBList [][]interface{}) {
	// Create an rds service client.
	svc := rds.New(se.Sess)
	// Get rds cluster
	output, err := svc.DescribeDBInstances(&rds.DescribeDBInstancesInput{
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
	for _, dbInstance := range output.DBInstances {
		var db []interface{}
		//handle securitygroups
		var sgs,pgs []string
		if len(dbInstance.VpcSecurityGroups) == 0 {
			sgs = append(sgs, "N/A")
		} else {
			for _, sg := range dbInstance.VpcSecurityGroups {
				sgs = append(sgs, *sg.VpcSecurityGroupId)
			}
		}
		// handle parametergroups
		if len(dbInstance.DBParameterGroups) == 0 {
			pgs = append(pgs, "N/A")
		} else {
			for _, pg := range dbInstance.DBParameterGroups {
				pgs = append(pgs, *pg.DBParameterGroupName +":"+*pg.ParameterApplyStatus+" ")
			}
		}
		if len(output.DBInstances) >= 100 {
			// todo cluster > 100
			tools.WarningLogger.Println("Number Of DB Clusters > 100 , Data May Loss.")
		}
		// handle account id
		accountId := GetARNDetail(*dbInstance.DBInstanceArn)["accountId"]
		db = append(db, accountId, se.UsedRegion,*dbInstance.DBInstanceIdentifier,*dbInstance.DBInstanceClass,
			*dbInstance.Endpoint.Address,*dbInstance.Engine, *dbInstance.EngineVersion,*dbInstance.Endpoint.Port,
			*dbInstance.DBSubnetGroup,*dbInstance.AvailabilityZone,*dbInstance.MultiAZ,
			*dbInstance.DBInstanceStatus, *dbInstance.StorageType, *dbInstance.PreferredBackupWindow,
			*dbInstance.BackupRetentionPeriod,*dbInstance.PreferredMaintenanceWindow,pgs, sgs)
		DBList = append(DBList, db)
	}
	return DBList
}
