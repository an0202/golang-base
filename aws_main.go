package main

import "golang-base/aws"

func main() {
	sess := aws.InitSession()
	aws.S3Download(sess, "cloudtrail-697333127309", "AWSLogs/281525879386/CloudTrail/ap-northeast-1/2019/10/28/281525879386_CloudTrail_ap-northeast-1_20191028T0325Z_yRuQd8yp96N1efGV.json.gz")
}
