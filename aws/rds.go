/**
 * @Author: jie.an
 * @Description:
 * @File:  rds.go
 * @Version: 1.0.1
 * @Date: 2020/03/02 11:00
 */
package aws

import (
	"github.com/aws/aws-sdk-go/service/rds"
	"golang-base/tools"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

type DBCluster struct {
	AccountId         string
	Region            string
	Identifier        string
	Endpoint          string
	ReaderEndpoint    string
	EngineMode        string
	Engine            string
	EngineVersion     string
	Port              int64
	MultiAZ           bool
	Status            string
	MaintenanceWindow string
	BackupWindow      string
	BackupRetention   int64
	ParameterGroup    string
	AvailabilityZones []interface{}
	SecurityGroups    []interface{}
	Members           []interface{}
}

//List DBClusters
func Listv2DBClusters(se Session) (DBClusterList []interface{}) {
	// Create an rds service client.
	svc := rds.New(se.Sess)
	// Get rds cluster
	output, err := svc.DescribeDBClusters(&rds.DescribeDBClustersInput{
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
	for _, dbCluster := range output.DBClusters {
		DBC := new(DBCluster)
		if len(output.DBClusters) >= 100 {
			// todo cluster > 100
			tools.WarningLogger.Println("Number Of DB Clusters > 100 , Data May Loss.")
		}
		DBC.AccountId = se.AccountId
		DBC.Region = se.UsedRegion
		DBC.Identifier = *dbCluster.DBClusterIdentifier
		DBC.Endpoint = *dbCluster.Endpoint
		if dbCluster.ReaderEndpoint != nil {
			DBC.ReaderEndpoint = *dbCluster.ReaderEndpoint
		} else {
			DBC.ReaderEndpoint = "N/A"
		}
		DBC.EngineMode = *dbCluster.EngineMode
		DBC.Engine = *dbCluster.Engine
		DBC.EngineVersion = *dbCluster.EngineVersion
		DBC.Port = *dbCluster.Port
		DBC.MultiAZ = *dbCluster.MultiAZ
		DBC.Status = *dbCluster.Status
		DBC.MaintenanceWindow = *dbCluster.PreferredMaintenanceWindow
		DBC.BackupWindow = *dbCluster.PreferredBackupWindow
		DBC.BackupRetention = *dbCluster.BackupRetentionPeriod
		DBC.ParameterGroup = *dbCluster.DBClusterParameterGroup
		//handle availabilityZones securityGroups Members
		if len(dbCluster.AvailabilityZones) == 0 {
			DBC.AvailabilityZones = append(DBC.AvailabilityZones, "N/A")
		} else {
			for _, az := range dbCluster.AvailabilityZones {
				DBC.AvailabilityZones = append(DBC.AvailabilityZones, *az)
			}
		}
		if len(dbCluster.VpcSecurityGroups) == 0 {
			DBC.SecurityGroups = append(DBC.SecurityGroups, "N/A")
		} else {
			for _, sg := range dbCluster.VpcSecurityGroups {
				DBC.SecurityGroups = append(DBC.SecurityGroups, *sg.VpcSecurityGroupId)
			}
		}
		if len(dbCluster.DBClusterMembers) == 0 {
			DBC.Members = append(DBC.Members, "N/A")
		} else {
			for _, mber := range dbCluster.DBClusterMembers {
				DBC.Members = append(DBC.Members, *mber)
			}
		}
		DBClusterList = append(DBClusterList, *DBC)
	}
	return DBClusterList
}

//List DBInstances
func ListDBInstances(se Session) (DBInstanceList [][]interface{}) {
	// Create an rds service client.
	svc := rds.New(se.Sess)
	// Get rds instances
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
	if len(output.DBInstances) >= 100 {
		// todo cluster > 100
		tools.WarningLogger.Println("Number Of DB Instances > 100 , Data May Loss.")
	}
	for _, dbInstance := range output.DBInstances {
		var db []interface{}
		//handle securitygroups
		var sgs, pgs []string
		var env string
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
				pgs = append(pgs, *pg.DBParameterGroupName+":"+*pg.ParameterApplyStatus+" ")
			}
		}
		// handle tags
		tags := ListDBTags(se, *dbInstance.DBInstanceArn)
		if v, ok := tags["Env"]; ok {
			env = v
		} else {
			env = "N/A"
		}
		accountId := se.AccountId
		// accountId := GetARNDetail(*dbInstance.DBInstanceArn)["accountId"]
		db = append(db, accountId, se.UsedRegion, *dbInstance.DBInstanceIdentifier, env, *dbInstance.DBInstanceClass,
			*dbInstance.Endpoint.Address, *dbInstance.Engine, *dbInstance.EngineVersion, *dbInstance.Endpoint.Port,
			*dbInstance.DBSubnetGroup, *dbInstance.AvailabilityZone, *dbInstance.MultiAZ,
			*dbInstance.DBInstanceStatus, *dbInstance.StorageType, *dbInstance.PreferredBackupWindow,
			*dbInstance.BackupRetentionPeriod, *dbInstance.PreferredMaintenanceWindow, pgs, sgs, tags)
		DBInstanceList = append(DBInstanceList, db)
	}
	return DBInstanceList
}

//List DBInstances Tags
func ListDBTags(se Session, dbName string) (DBTags map[string]string) {
	// Create an rds service client.
	svc := rds.New(se.Sess)
	// Get rds instances
	output, err := svc.ListTagsForResource(&rds.ListTagsForResourceInput{
		ResourceName: aws.String(dbName),
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
	DBTags = make(map[string]string)
	for _, tag := range output.TagList {
		DBTags[*tag.Key] = *tag.Value
	}
	return DBTags
}
