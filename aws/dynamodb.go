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

type Table struct {
	AccountId     string
	Region        string
	Name          string
	Status        string
	ARN           string
	SizeBytes     int64
	ReadCapacity  int64
	WriteCapacity int64
	KeyScheme     []interface{}
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
		BK := new(Table)
		BK.AccountId = se.AccountId
		BK.describeTable(se, *tableName)
		DynamoDBList = append(DynamoDBList, *BK)
	}
	return DynamoDBList
}

func (tb *Table) describeTable(se Session, TableName string) Table {
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
		return *tb
	}
	tb.Name = *output.Table.TableName
	tb.Region = se.UsedRegion
	tb.Status = *output.Table.TableStatus
	tb.ARN = *output.Table.TableArn
	tb.SizeBytes = *output.Table.TableSizeBytes
	tb.ReadCapacity = *output.Table.ProvisionedThroughput.ReadCapacityUnits
	tb.WriteCapacity = *output.Table.ProvisionedThroughput.WriteCapacityUnits
	if len(output.Table.KeySchema) != 0 {
		for _, key := range output.Table.KeySchema {
			tb.KeyScheme = append(tb.KeyScheme, *key)
		}
	}
	return *tb
}
