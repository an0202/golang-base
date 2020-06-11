/**
 * @Author: jie.an
 * @Description:
 * @File:  elb.go
 * @Version: 1.0.1
 * @Date: 2020/02/25 13:42
 */
package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"golang-base/tools"
	"os"
	"path/filepath"
	"strings"
)

type Bucket struct {
	AccountId  string
	Region     string
	Name       string
	Env        string
	ACL        []interface{}
	Policy     string
	CORS       []interface{}
	LifeCycle  []interface{}
	Versioning string
	Website    interface{}
	Tags       map[string]string
}

type Object struct {
	AccountId    string
	Region       string
	Bucket       string
	Key          string
	Size         int64
	StorageClass string
	LastModify   string
}

//List S3Bucket
func Listv2S3(se Session) (S3List []interface{}) {
	// Create an s3 service client.
	svc := s3.New(se.Sess)
	// Get sns topics
	output, err := svc.ListBuckets(&s3.ListBucketsInput{})
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
	for _, bucket := range output.Buckets {
		BK := new(Bucket)
		BK.AccountId = se.AccountId
		BK.Name = *bucket.Name
		//handle tags
		tags := ListBucketTags(se, BK.Name)
		if v, ok := tags["Env"]; ok {
			BK.Env = v
		} else {
			BK.Env = "N/A"
		}
		BK.Tags = tags
		BK.Region = BK.GetBucketLocation(se, BK.Name)
		//permission
		BK.ACL = BK.GetBucketACLs(se, BK.Name)
		BK.Policy = BK.GetBucketPolicy(se, BK.Name)
		BK.CORS = BK.GetCORSRules(se, BK.Name)
		//property
		BK.Versioning = BK.GetBucketVersioning(se, BK.Name)
		BK.Website = BK.GetBucketWebSite(se, BK.Name)
		//management
		BK.LifeCycle = BK.GetLifeCycleRules(se, BK.Name)
		S3List = append(S3List, *BK)
	}
	return S3List
}

//List S3Bucket Tags
func ListBucketTags(se Session, BucketName string) (BucketTags map[string]string) {
	// Create an s3 service client.
	svc := s3.New(se.Sess)
	// Get Bucket tags
	output, err := svc.GetBucketTagging(&s3.GetBucketTaggingInput{
		Bucket: aws.String(BucketName),
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
	BucketTags = make(map[string]string)
	for _, tag := range output.TagSet {
		BucketTags[*tag.Key] = *tag.Value
	}
	return BucketTags
}

func (bk *Bucket) GetBucketLocation(se Session, BucketName string) string {
	// Create an s3 service client.
	svc := s3.New(se.Sess)
	// Get bucket location
	output, err := svc.GetBucketLocation(&s3.GetBucketLocationInput{
		Bucket: aws.String(BucketName),
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
		return ""
	}
	return *output.LocationConstraint
}

func (bk *Bucket) GetBucketVersioning(se Session, BucketName string) string {
	// Create an s3 service client.
	svc := s3.New(se.Sess)
	// Get bucket location
	output, err := svc.GetBucketVersioning(&s3.GetBucketVersioningInput{
		Bucket: aws.String(BucketName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				tools.ErrorLogger.Println(BucketName, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			tools.ErrorLogger.Println(err.Error())
		}
		return "N/A"
	}
	if output.Status == nil {
		return "N/A"
	}
	return *output.Status
}

func (bk *Bucket) GetBucketWebSite(se Session, BucketName string) interface{} {
	// Create an s3 service client.
	svc := s3.New(se.Sess)
	// Get bucket location
	output, err := svc.GetBucketWebsite(&s3.GetBucketWebsiteInput{
		Bucket: aws.String(BucketName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				tools.WarningLogger.Println(BucketName, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			tools.ErrorLogger.Println(err.Error())
		}
		return "N/A"
	}
	return *output
}

func (bk *Bucket) GetBucketPolicy(se Session, BucketName string) string {
	// Create an s3 service client.
	svc := s3.New(se.Sess)
	// Get bucket location
	output, err := svc.GetBucketPolicy(&s3.GetBucketPolicyInput{
		Bucket: aws.String(BucketName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				tools.WarningLogger.Println(BucketName, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			tools.ErrorLogger.Println(err.Error())
		}
		return "N/A"
	}
	return *output.Policy
}

func (bk *Bucket) GetBucketACLs(se Session, BucketName string) (GrantList []interface{}) {
	// Create an s3 service client.
	svc := s3.New(se.Sess)
	// Get bucket location
	output, err := svc.GetBucketAcl(&s3.GetBucketAclInput{
		Bucket: aws.String(BucketName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				tools.ErrorLogger.Println(BucketName, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			tools.ErrorLogger.Println(err.Error())
		}
		GrantList = append(GrantList, "N/A")
		return GrantList
	}
	if len(output.Grants) == 0 {
		GrantList = append(GrantList, "N/A")
	} else {
		for _, grant := range output.Grants {
			GrantList = append(GrantList, *grant)
		}
	}
	return GrantList
}

func (bk *Bucket) GetCORSRules(se Session, BucketName string) (RuleList []interface{}) {
	// Create an s3 service client.
	svc := s3.New(se.Sess)
	// Get bucket location
	output, err := svc.GetBucketCors(&s3.GetBucketCorsInput{
		Bucket: aws.String(BucketName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				tools.WarningLogger.Println(BucketName, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			tools.ErrorLogger.Println(err.Error())
		}
		RuleList = append(RuleList, "N/A")
		return RuleList
	}
	if len(output.CORSRules) == 0 {
		RuleList = append(RuleList, "N/A")
	} else {
		for _, corsRule := range output.CORSRules {
			RuleList = append(RuleList, *corsRule)
		}
	}
	return RuleList
}

func (bk *Bucket) GetLifeCycleRules(se Session, BucketName string) (RuleList []interface{}) {
	// Create an s3 service client.
	svc := s3.New(se.Sess)
	// Get bucket location
	output, err := svc.GetBucketLifecycleConfiguration(&s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(BucketName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				tools.WarningLogger.Println(BucketName, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			tools.ErrorLogger.Println(err.Error())
		}
		RuleList = append(RuleList, "N/A")
		return RuleList
	}
	if len(output.Rules) == 0 {
		RuleList = append(RuleList, "N/A")
	} else {
		for _, lifecyleRule := range output.Rules {
			RuleList = append(RuleList, *lifecyleRule)
		}
	}
	return RuleList
}

// Return objectList witch contains all object name
// Return nil while no content
// prefix should not start with "/"
func S3ListObjects(sess *session.Session, bucket string, prefix string) (ObjectsList []Object) {
	tools.InfoLogger.Println("ListBucket File:", bucket+prefix)
	svc := s3.New(sess)
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int64(1000),
	}
	output, err := svc.ListObjectsV2(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				tools.WarningLogger.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				tools.WarningLogger.Println(aerr.Error())
			}
		} else {
			tools.WarningLogger.Println(err.Error())
		}
		return
	}
	accountId := GetAccountId(sess)
	region := *sess.Config.Region
	if len(output.Contents) == 0 {
		return nil
	} else {
		for _, content := range output.Contents {
			if strings.HasSuffix(*content.Key, "/") {
				continue
			}
			obj := new(Object)
			obj.AccountId = accountId
			obj.Region = region
			obj.Bucket = bucket
			obj.Key = *content.Key
			obj.StorageClass = *content.StorageClass
			obj.Size = *content.Size
			obj.LastModify = content.LastModified.String()
			ObjectsList = append(ObjectsList, *obj)
		}
	}
	return ObjectsList
}

// Head S3 Object
func headS3Object(sess *session.Session, bucket string, item string) {
	svc := s3.New(sess)
	output, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				tools.WarningLogger.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			tools.ErrorLogger.Println(err.Error())
		}
	}
	tools.InfoLogger.Println("Object Size:", *output.ContentLength, "bytes.")
}

// S3Download Used Download Content From S3
func S3Download(sess *session.Session, bucket string, item string, dest string) {
	headS3Object(sess, bucket, item)
	tools.InfoLogger.Println("Downloading File:", bucket+item)
	// Store filename/path for returning and using later on
	fPath := filepath.Join(dest, item)

	if !strings.HasPrefix(fPath, filepath.Clean(dest)+string(os.PathSeparator)) {
		tools.ErrorLogger.Fatalf("%s: illegal file path", fPath)
	}
	// create directory
	if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
		tools.ErrorLogger.Fatalf(err.Error())
	}
	//Create local file
	file, err := os.Create(fPath)
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
	tools.InfoLogger.Println("Downloaded", file.Name(), numBytes, "bytes.")
}
