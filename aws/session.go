package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func InitSession(region string) *session.Session {
	//Set AWS_SDK_LOAD_CONFIG="true"
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	// Init a session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		fmt.Println(err)
	}
	return sess
}
