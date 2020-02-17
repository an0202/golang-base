/**
 * @Author: jie.an
 * @Description:
 * @File:  aws_tool.go
 * @Version: 1.0.0
 * @Date: 2020/02/17 16:41
 */
package aws

import (
	"golang-base/tools"
	"regexp"
)

//GetARNDetail return a map contains ARN base information
//map[accountId:123456789012 region: resource:user/Development/product_1234/* service:iam]
//https://blog.csdn.net/butterfly5211314/article/details/82532970
//https://golang.org/pkg/regexp/syntax/
func GetARNDetail(arn string) (arnMap map[string]string) {
	re := regexp.MustCompile(`^arn:(?:aws|aws-cn|\s):(?P<service>[a-z\s]+):(?P<region>[a-z0-9\-\s]+):(?P<accountId>[0-9\s]+):(?P<resource>.*)`)
	matched := re.MatchString(arn)
	arnMap = make(map[string]string)
	if matched == false {
		tools.ErrorLogger.Fatalln(arn," Is Not A AWS ARN.")
	} else {
		groupNames := re.SubexpNames()
		match := re.FindStringSubmatch(arn)
		//for k, v := range match {
		//	fmt.Println(k,v)
		//}
		for i, name := range groupNames {
			if i != 0 && name != "" {
				arnMap[name] = match[i]
			}
		}
	}
	return arnMap
}
