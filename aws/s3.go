package aws

import (
	"fmt"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Download Used Download Content From S3
func S3Download(sess *session.Session, bucket string, item string) {
	//Create local file
	file, err := os.Create(path.Base(item))
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	downloader := s3manager.NewDownloader(sess)

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}
