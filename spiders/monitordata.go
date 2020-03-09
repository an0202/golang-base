/**
 * @Author: jie.an
 * @Description:
 * @File:  monitorData.go
 * @Version: 1.0.1
 * @Date: 2020/03/08 21:15
 */

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"golang-base/excel"
	"golang-base/tools"
	"reflect"
	"sync"
	"time"
)

var basicRequestHeader = map[string]string{
	"accept": "application/json, text/javascript, */*; q=0.01",
	"accept-encoding": "gzip, deflate, br",
	"accept-language": "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
	"authorization": "",
	"cookie": "",
	"referer": "https://md.opsnow.cn/",
	"sec-fetch-dest": "empty",
	"sec-fetch-mode": "cors",
	"sec-fetch-site": "same-origin",
	"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.122 Safari/537.36",
	"x-requested-with": "XMLHttpRequest",
}

var customRequestHeader = map[string]interface{}{}

var serviceList = []string{"AUTH","CPM","FKP","IMS","KPI","MAINPAGE","PROSERVICE","RM","SAS","SEARCH","SECURITY",
	"SMARTTHINGS","SSO","STS"}

var respDataList []map[string]interface{}

func getServerList()(serverList []interface{}){
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.122 Safari/537.36"),
	)
	c.Limit(&colly.LimitRule{
		DomainGlob: "*opsnow.*",
		Delay:        3 * time.Second,
		RandomDelay:  2 * time.Second,
	})
	c.SetRequestTimeout(30 * time.Second)
	c.OnRequest(func(r *colly.Request) {
		tools.InfoLogger.Println("Request:",r.URL)
		for k,v := range basicRequestHeader {
			r.Headers.Set(k,v)
		}
		if len(customRequestHeader) != 0 {
			for ck,cv := range customRequestHeader {
				r.Headers.Set(ck,cv.(string))
			}
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		tools.InfoLogger.Println("Request URL:", r.Request.URL,"failed with responses:",r,"\nError",err)
	})
	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			tools.ErrorLogger.Println("error while response")
		}
		error := json.Unmarshal(r.Body, &respDataList)
		if error != nil {
			fmt.Println(error)
		}
		tools.InfoLogger.Printf("Found %d Server\n",len(respDataList))
		for _,v := range respDataList {
			if v["serverName"] != nil {
				serverList = append(serverList, v["serverName"])
			}
		}
	})
	for _,service := range serviceList {
		c.Visit("https://md.opsnow.cn/pl/serviceServer/SEC-SRCN_"+service+"/status/customer/573a7278-aed8-4034-8e09-dff8ae40dd9b")
		c.Wait()
	}
	return serverList
}
/*
MetricName:QOS_CPU_USAGE,QOS_MEMORY_PHYSICAL_PERC
Method:max,avg
 */
func getMericData(serverList []interface{},metricName string,method string)(serverMetricList [][]interface{}) {
	//make a queue
	var ch = make(chan string, 500)
	if len(serverList) == 0 {
		return
	}
	//add in channel
	for _, i := range serverList {
		ch <- i.(string)
	}
	// read channel
	wg := sync.WaitGroup{}
	wg.Add(1)
	for k:= 0; k<len(serverList) ; k++ {
		var serverNames string
		if len(ch) == 0 {
			break
		}
		//each loop will read 20 items from channel
		for i := 0; i < 20; i++ {
			if len(ch) == 0 {
				break
			}
			if i == 0 {
				serverNames = serverNames + <-ch
			} else {
				serverNames = serverNames + "," + <-ch
			}
		}
		//handle metric
		c := colly.NewCollector(
			colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.122 Safari/537.36"),
		)
		c.Limit(&colly.LimitRule{
			DomainGlob: "*opsnow.*",
			Delay:        3 * time.Second,
			RandomDelay:  2 * time.Second,
		})
		c.SetRequestTimeout(30 * time.Second)
		c.OnRequest(func(r *colly.Request) {
			tools.InfoLogger.Println("Request:",r.URL)
			for k,v := range basicRequestHeader {
				r.Headers.Set(k,v)
			}
			if len(customRequestHeader) != 0 {
				for ck,cv := range customRequestHeader {
					r.Headers.Set(ck,cv.(string))
				}
			}
		})
		c.OnError(func(r *colly.Response, err error) {
			tools.InfoLogger.Println("Request URL:", r.Request.URL,"failed with responses:",r,"\nError",err)
		})
		c.OnResponse(func(r *colly.Response) {
			if r.StatusCode != 200 {
				tools.ErrorLogger.Println("error while response")
			}
			error := json.Unmarshal(r.Body, &respDataList)
			if error != nil {
				fmt.Println(error)
			}
			tools.InfoLogger.Printf("Found %d Metric\n",len(respDataList))
			for _,v := range respDataList {
				var serverMetric []interface{}
				if v["serverName"] != "" {
					serverMetric = append(serverMetric, v["serverName"], v["sampleValue"])
					serverMetricList = append(serverMetricList,serverMetric)
				}
			}
		})
		c.Visit("https://md.opsnow.cn/pl/serviceServerTop5/qos/hbase?start=720h-ago&aggregator=" +method+"&qosList=" +
			metricName + "&serverList=" + serverNames + "&downsample=650h-"+ method)
		c.Wait()
	}
	if len(ch) == 0 {
		wg.Done()
	}
	wg.Wait()
	return serverMetricList
}

func main() {
	//parse conf
	var c tools.HttpHeaderConf
	conf := c.GetConf("conf.yaml")
	//https://my.oschina.net/solate/blog/715681
	t := reflect.TypeOf(*conf)
	v := reflect.ValueOf(*conf)
	for k := 0; k < t.NumField(); k++ {
		if v.Field(k).CanInterface() {
			customRequestHeader[t.Field(k).Name] = v.Field(k).Interface()
		}
	}
	var filePath = "monitorData.xlsx"
	var cpuAvgHeadLine = []interface{}{"ServerName","CPU_AVG_PERCENT"}
	var cpuMaxHeadLine = []interface{}{"ServerName","CPU_MAX_PERCENT"}
	var memAvgHeadLine = []interface{}{"ServerName","MEM_AVG_PERCENT"}
	var memMaxHeadLine = []interface{}{"ServerName","MEM_MAX_PERCENT"}
	//get data
	serverList := getServerList()
	cpuAvg := getMericData(serverList,"QOS_CPU_USAGE","avg")
	cpuMax := getMericData(serverList,"QOS_CPU_USAGE","max")
	memAvg := getMericData(serverList,"QOS_MEMORY_PHYSICAL_PERC","avg")
	memMax := getMericData(serverList,"QOS_MEMORY_PHYSICAL_PERC","max")
	//write to excel
	excel.CreateFile(filePath)
	excel.SetHeadLine(filePath,"CPU_AVG",cpuAvgHeadLine)
	excel.SetListRows(filePath,"CPU_AVG",cpuAvg)
	excel.SetHeadLine(filePath,"CPU_MAX",cpuMaxHeadLine)
	excel.SetListRows(filePath,"CPU_MAX",cpuMax)
	excel.SetHeadLine(filePath,"MEM_AVG",memAvgHeadLine)
	excel.SetListRows(filePath,"MEM_AVG",memAvg)
	excel.SetHeadLine(filePath,"MEM_MAX",memMaxHeadLine)
	excel.SetListRows(filePath,"MEM_MAX",memMax)
}