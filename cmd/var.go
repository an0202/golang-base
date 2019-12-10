/**
 * @Author: jie.an
 * @Description:
 * @File:  var.go
 * @Version: 1.0.0
 * @Date: 2019/12/9 16:31
 */
package cmd

// General
var (
	region *string
	help   *bool
)

// ec2tag
var (
	excelFile *string
	sheetName *string
	method    *string
)

//ec2 ami
var (
	suffix     *string
	instanceid *string
	srcFile    *string
)
