package aws

import (
	"golang-base/tools"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Session struct {
	UsedRegion        string
	UsedAwsProfile 	  string
	Sess *session.Session
}

func (se *Session) InitSessionWithAWSProfile(region, awsProfile string) *session.Session {
	//Set AWS_SDK_LOAD_CONFIG="true"
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	os.Setenv("AWS_PROFILE", awsProfile)
	// Init a session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		tools.ErrorLogger.Fatalln(err)
	}
	tools.InfoLogger.Printf("Create A New Session By AWS_PROFILE: %s, Region: %s",awsProfile,region)
	se.UsedAwsProfile = awsProfile
	se.UsedRegion = region
	se.Sess = sess
	return sess
}

// compatible with old code

func InitSession(region string) *session.Session {
	//Set AWS_SDK_LOAD_CONFIG="true"
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	// Init a session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		tools.ErrorLogger.Fatalln(err)
	}
	return sess
}

func GetAccountId(sess *session.Session) string {
	svc := sts.New(sess)
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
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
	}
	tools.InfoLogger.Println("Get Caller Identity:", result)
	return *result.Account
}
