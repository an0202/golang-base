/**
 * @Author: jie.an
 * @Description:
 * @File:  monitorDatav2.go
 * @Version: 1.0.0
 * @Date: 2020/04/28 18:18
 */

package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"golang-base/tools"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

//
var basicRequestHeader2 = map[string]string{
	"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"accept-encoding": "gzip, deflate, br",
	"accept-language": "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
	"authorization":   "Bearer 4446851e-cfcf-48d4-9be7-22f8f0cafb47",
	//"cookie":         "BSP_isRemember=false; BSP_isRememberEmail=jie.an%40bespinglobal.cn; BSP_LangCode=zh",
	"referer":        "https://project.opsnow.cn/",
	"sec-fetch-dest": "empty",
	"sec-fetch-site": "same-origin",
	"customerid":     "573a7278-aed8-4034-8e09-dff8ae40dd9b",
	//"user-agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.122 Safari/537.36",
}

func startJob2() {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.122 Safari/537.36"),
	)
	customHeader2 := make(map[string]string)
	//
	jar, _ := cookiejar.New(nil)
	var cookies, cookies02 []*http.Cookie
	cookie := &http.Cookie{
		Name:   "access_token_type",
		Value:  "Bearer",
		Path:   "/",
		Domain: ".opsnow.cn",
	}
	cookies = append(cookies, cookie)
	cookie = &http.Cookie{
		Name:   "access_token",
		Value:  "4446851e-cfcf-48d4-9be7-22f8f0cafb47",
		Path:   "/",
		Domain: ".opsnow.cn",
	}
	cookies = append(cookies, cookie)
	url, _ := url.Parse("http://opsnow.cn")
	jar.SetCookies(url, cookies)
	fmt.Println("jar001:", jar)
	//jar 002
	cookie02 := &http.Cookie{
		Name:   "JSESSIONID",
		Value:  "0BA7F097A9CA674762093D5CDD867863",
		Path:   "/",
		Domain: "project.opsnow.cn",
	}
	cookies02 = append(cookies02, cookie02)
	url, _ = url.Parse("http://project.opsnow.cn")
	jar.SetCookies(url, cookies02)
	c.SetCookieJar(jar)
	fmt.Println("jar002:", jar)
	c.SetRedirectHandler(func(req *http.Request, via []*http.Request) error {
		fmt.Println("redirect to header:", req.Header)
		//if req.URL.RawQuery != "" {
		//	//fmt.Println("print cookies formarw")
		//	//c.OnRequest(func(r *colly.Request) {
		//	//	v, _ := url.ParseQuery(req.URL.RawQuery)
		//	//	if _, ok := v["access_token"]; ok {
		//	//		fmt.Println("get coolies", r.Headers.Get("cookie"))
		//	//	}
		//	//})
		//}
		return nil
	})
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*opsnow.*",
		Delay:       3 * time.Second,
		RandomDelay: 2 * time.Second,
	})
	c.SetRequestTimeout(30 * time.Second)
	c.OnRequest(func(r *colly.Request) {
		for k, v := range basicRequestHeader2 {
			r.Headers.Set(k, v)
		}
		for k, v := range customHeader2 {
			r.Headers.Set(k, v)
		}
		fmt.Println("request header", r.Headers)
	})
	c.OnError(func(r *colly.Response, err error) {
		tools.InfoLogger.Println("Request URL:", r.Request.URL, "failed with responses:", r, "\nError", err)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
		fmt.Println("Cookies:", c.Cookies(r.Request.URL.String()))
		fmt.Println(string(r.Body))
	})
	//c.Post("https://sso.opsnow.cn/servicelogin", map[string]string{
	//	"username":      "email@email",
	//	"password":      "password",
	//	"client_id":     "wGBi_35lOOIxomIJRQp_cHxwBJka",
	//	"redirect_uri":  "https://www.opsnow.cn/dashboard",
	//	"scope":         "all",
	//	"response_type": "token",
	//})
	//c.Visit("https://project.opsnow.cn/")
	//c.Visit("https://project.opsnow.cn/session/info")
	c.Visit("https://project.opsnow.cn/prj/projectList?_=158803474931")
}

func main() {
	//parse conf
	startJob2()
}
