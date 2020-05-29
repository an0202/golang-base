package aws

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"golang-base/tools"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Session struct {
	UsedRegion     string
	UsedAwsProfile string
	AccountId      string
	Sess           *session.Session
}

func InitSessionWithAKSK(ak, sk, region string) *session.Session {
	// Init a credentials
	cres := credentials.NewStaticCredentials(ak, sk, "")
	// Init a session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: cres,
	},
	)
	if err != nil {
		tools.ErrorLogger.Fatalln(err)
	}
	GetAccountId(sess)
	return sess
}

func (se *Session) InitSessionWithAWSProfile(region, awsProfile string) *Session {
	//Set AWS_SDK_LOAD_CONFIG="true"
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	if awsProfile == "" {
		tools.WarningLogger.Println("Config's AWS_PROFILE Is null , Use Previous AWS_PROFILE Or" +
			" Create Credential From OS Environment.")
	} else {
		tools.InfoLogger.Printf("Create Credential By : %s, Region: %s\n", awsProfile, region)
		os.Setenv("AWS_PROFILE", awsProfile)
	}
	// Init a session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		tools.ErrorLogger.Fatalln(err)
	}
	se.AccountId = GetAccountId(sess)
	se.UsedAwsProfile = awsProfile
	se.UsedRegion = region
	se.Sess = sess
	return se
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
		return "ErrorId"
	}
	tools.InfoLogger.Println("Get Caller Identity:", *result.Arn)
	return *result.Account
}
