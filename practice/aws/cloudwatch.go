/**
 * @Author: jie.an
 * @Description:
 * @File:  cloudwatch.go
 * @Version: 1.0.0
 * @Date: 2020/8/8 14:26
 */

package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type CloudWatchToFalcon struct {
	input         getMetricInput
	output        *cloudwatch.GetMetricStatisticsOutput
	awsDimensions []*cloudwatch.Dimension
	falconTags    string
	endPoint      string
}

func (ctf *CloudWatchToFalcon) Do(input getMetricInput) {
	ctf.input = input
	ctf.init()
	ctf.getMetric()
	ctf.falconPush()
}

//https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/cloudwatch_concepts.html
type getMetricInput struct {
	// Minimum, Maximum, Average, SampleCount, Sum
	Statistics string
	//https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/aws-services-cloudwatch-metrics.html
	Namespace  string
	MetricName string
	//e.g. "InstanceId": "i-01af376a15fa2f36d"
	Dimensions map[string]string
}

func (ctf *CloudWatchToFalcon) init() {
	// Set hostname as falcon endPoint
	hostName, err := os.Hostname()
	if err == nil {
		log.Warningln("Failed to get hostname, set endPoint to CloudWatch_To_Falcon")
		ctf.endPoint = "CloudWatch_To_Falcon"
	} else {
		ctf.endPoint = hostName
	}
	// Handle aws dimensions and falcon tags
	if len(ctf.input.Dimensions) == 0 {
		return
	}
	for k, v := range ctf.input.Dimensions {
		var awsDimension = &cloudwatch.Dimension{
			Name:  aws.String(k),
			Value: aws.String(v),
		}
		ctf.awsDimensions = append(ctf.awsDimensions, awsDimension)
		var tag = fmt.Sprintf("%s=%s", k, v)
		if ctf.falconTags == "" {
			ctf.falconTags = tag
		} else {
			ctf.falconTags = fmt.Sprintf("%s,%s", ctf.falconTags, tag)
		}
	}
}

// Get specific indicator data from cloudwatch one minute ago.
//
// CloudWatch should generated data each 60 seconds,otherwise no data can be obtained.
func (ctf *CloudWatchToFalcon) getMetric() {
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and configuration from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create new cloudwatch client.
	svc := cloudwatch.New(sess)
	result, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Dimensions: ctf.awsDimensions,
		EndTime:    aws.Time(time.Now()),
		MetricName: aws.String(ctf.input.MetricName),
		Namespace:  aws.String(ctf.input.Namespace),
		Period:     aws.Int64(60),
		StartTime:  aws.Time(time.Now().Add(-time.Minute * 1)),
		Statistics: aws.StringSlice([]string{ctf.input.Statistics}),
	})
	if err != nil {
		log.WithField("Error", err.Error()).Warnln("Error getting metrics")
	}
	ctf.output = result
}

//https://book.open-falcon.org/zh/usage/data-push.html
func (ctf *CloudWatchToFalcon) falconPush() {
	fmt.Println(ctf.output)
	type item struct {
		Endpoint    string  `json:"endpoint"`
		Metric      string  `json:"metric"`
		Timestamp   int64   `json:"timestamp"`
		Step        int     `json:"step"`
		Value       float64 `json:"value"`
		CounterType string  `json:"counterType"`
		Tags        string  `json:"tags"`
	}
	type message struct {
		Item []item `json:"item"`
	}
	// data body
	var post message
	if len(ctf.output.Datapoints) == 0 {
		log.Warningln("Found No Data From CloudWatch , Set Value To -1.")
		i := item{
			Endpoint:    ctf.endPoint,
			Metric:      *ctf.output.Label,
			Timestamp:   time.Now().Add(-time.Minute * 1).Unix(),
			Value:       -1,
			CounterType: "GAUGE",
			Tags:        ctf.falconTags,
			Step:        60,
		}
		post.Item = append(post.Item, i)
	} else {
		switch ctf.input.Statistics {
		case "Minimum":
			i := item{
				Endpoint:    ctf.endPoint,
				Metric:      *ctf.output.Label,
				Timestamp:   ctf.output.Datapoints[0].Timestamp.Unix(),
				Value:       *ctf.output.Datapoints[0].Minimum,
				CounterType: "GAUGE",
				Tags:        ctf.falconTags,
				Step:        60,
			}
			post.Item = append(post.Item, i)
		case "Maximum":
			i := item{
				Endpoint:    ctf.endPoint,
				Metric:      *ctf.output.Label,
				Timestamp:   ctf.output.Datapoints[0].Timestamp.Unix(),
				Value:       *ctf.output.Datapoints[0].Maximum,
				CounterType: "GAUGE",
				Tags:        ctf.falconTags,
				Step:        60,
			}
			post.Item = append(post.Item, i)
		case "Average":
			i := item{
				Endpoint:    ctf.endPoint,
				Metric:      *ctf.output.Label,
				Timestamp:   ctf.output.Datapoints[0].Timestamp.Unix(),
				Value:       *ctf.output.Datapoints[0].Average,
				CounterType: "GAUGE",
				Tags:        ctf.falconTags,
				Step:        60,
			}
			post.Item = append(post.Item, i)
		default:
			log.Warningln("Statistics not supported")
		}
	}
	fmt.Printf("%+v \n", post.Item)
	jsonStr, _ := json.Marshal(post.Item)
	// Create a Resty Client
	client := resty.New()
	//client.SetDebug(true)
	client.SetHeader("Content-Type", "application/json")
	// Post json data to falcon
	_, err := client.R().
		SetBody(jsonStr).
		Post("http://127.0.0.1:1988/v1/push")
	if err != nil {
		log.Errorln(err)
	}
}

func main() {
	//delay start,prevent data loss in the first 5 seconds
	time.Sleep(time.Duration(5) * time.Second)
	// init logger
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	// Broker 1 disk
	var ctf1 = new(CloudWatchToFalcon)
	ctf1.Do(getMetricInput{
		Statistics: "Minimum",
		Namespace:  "AWS/EC2",
		MetricName: "StatusCheckFailed_System",
		Dimensions: map[string]string{"InstanceId": "i-01af376a15fa2f36d"},
	})
	// Broker 2 disk
	var ctf2 = new(CloudWatchToFalcon)
	ctf2.Do(getMetricInput{
		Statistics: "Maximum",
		Namespace:  "AWS/EC2",
		MetricName: "StatusCheckFailed_System",
		Dimensions: map[string]string{"InstanceId": "i-01af376a15fa2f36d"},
	})
	// Broker 3 disk
	var ctf3 = new(CloudWatchToFalcon)
	ctf3.Do(getMetricInput{
		Statistics: "Average",
		Namespace:  "AWS/Kafka",
		MetricName: "KafkaDataLogsDiskUsed",
		Dimensions: map[string]string{"Cluster Name": "xxxxx", "Broker ID": "bbbb"},
	})
}
