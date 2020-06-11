/**
 * @Author: jie.an
 * @Description:
 * @File:  dynamodb.go
 * @Version: 1.0.0
 * @Date: 2020/6/10 14:18
 */
package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"golang-base/tools"
)

type DynamoDBTable struct {
	AccountId     string
	Region        string
	Name          string
	Env           string
	Status        string
	ARN           string
	SizeBytes     int64
	ReadCapacity  int64
	WriteCapacity int64
	KeyScheme     []interface{}
	Tags          map[string]string
}

//List DynamoDBs
func Listv2DynamoDB(se Session) (DynamoDBList []interface{}) {
	// Create an dynamoDB service client.
	svc := dynamodb.New(se.Sess)
	// Get sns topics
	output, err := svc.ListTables(&dynamodb.ListTablesInput{})
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
	for _, tableName := range output.TableNames {
		dt := new(DynamoDBTable)
		dt.AccountId = se.AccountId
		dt.describeTable(se, *tableName)
		DynamoDBList = append(DynamoDBList, *dt)
	}
	return DynamoDBList
}

func (dt *DynamoDBTable) describeTable(se Session, TableName string) DynamoDBTable {
	// Create an dynamoDB service client.
	svc := dynamodb.New(se.Sess)
	// describe table
	output, err := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(TableName),
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
		return *dt
	}
	dt.Name = *output.Table.TableName
	dt.Region = se.UsedRegion
	dt.Status = *output.Table.TableStatus
	dt.ARN = *output.Table.TableArn
	//handle tags
	tags := ListDynamoDBTags(se, dt.ARN)
	if v, ok := tags["Env"]; ok {
		dt.Env = v
	} else {
		dt.Env = "N/A"
	}
	dt.Tags = tags
	dt.SizeBytes = *output.Table.TableSizeBytes
	dt.ReadCapacity = *output.Table.ProvisionedThroughput.ReadCapacityUnits
	dt.WriteCapacity = *output.Table.ProvisionedThroughput.WriteCapacityUnits
	if len(output.Table.KeySchema) != 0 {
		for _, key := range output.Table.KeySchema {
			dt.KeyScheme = append(dt.KeyScheme, *key)
		}
	}
	return *dt
}

//List DynamoDBTags Tags
func ListDynamoDBTags(se Session, DynamoDBName string) (DynamoDBTags map[string]string) {
	// Create an DynamoDB service client.
	svc := dynamodb.New(se.Sess)
	// Get DynamoDB tags
	output, err := svc.ListTagsOfResource(&dynamodb.ListTagsOfResourceInput{
		ResourceArn: aws.String(DynamoDBName),
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
	DynamoDBTags = make(map[string]string)
	for _, tag := range output.Tags {
		DynamoDBTags[*tag.Key] = *tag.Value
	}
	return DynamoDBTags
}
