/**
 * @Author: jie.an
 * @Description:
 * @File:  lambda.go
 * @Version: 1.0.0
 * @Date: 2020/6/11 12:15
 */
package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
	"golang-base/tools"
)

type LambdaFunction struct {
	AccountId  string
	Region     string
	Name       string
	Env        string
	ARN        string
	MemorySize int64
	CodeSize   int64
	Handler    string
	Runtime    string
	Timeout    int64
	Tags       map[string]string
}

//List Lambda
func Listv2Lambda(se Session) (LambdaList []interface{}) {
	// Create an lambda service client.
	svc := lambda.New(se.Sess)
	// Get lambda functions
	output, err := svc.ListFunctions(&lambda.ListFunctionsInput{})
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
	for _, function := range output.Functions {
		lf := new(LambdaFunction)
		lf.AccountId = se.AccountId
		lf.Region = se.UsedRegion
		lf.Name = *function.FunctionName
		lf.ARN = *function.FunctionArn
		lf.CodeSize = *function.CodeSize
		lf.MemorySize = *function.MemorySize
		lf.Handler = *function.Handler
		lf.Runtime = *function.Runtime
		lf.Timeout = *function.Timeout
		//handle tags
		tags := ListLambdaTags(se, lf.ARN)
		if v, ok := tags["Env"]; ok {
			lf.Env = *v
		} else {
			lf.Env = "N/A"
		}
		lf.Tags = make(map[string]string)
		for a, b := range tags {
			lf.Tags[a] = *b
		}
		LambdaList = append(LambdaList, *lf)
	}
	return LambdaList
}

//List Lambda Tags
func ListLambdaTags(se Session, LambdaName string) (LambdaTags map[string]*string) {
	// Create an DynamoDB service client.
	svc := lambda.New(se.Sess)
	// Get DynamoDB tags
	output, err := svc.ListTags(&lambda.ListTagsInput{
		Resource: aws.String(LambdaName),
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
	LambdaTags = make(map[string]*string)
	LambdaTags = output.Tags
	return LambdaTags
}
