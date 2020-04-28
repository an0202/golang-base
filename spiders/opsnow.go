/**
 * @Author: jie.an
 * @Description:
 * @File:  opsnow.go
 * @Version: 1.0.0
 * @Date: 2020/04/26 20:30
 */

package main

import "C"
import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/url"
)

//
var basicRequestHeader1 = map[string]string{
	"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"accept-encoding": "gzip, deflate, br",
	"accept-language": "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
	//"authorization":    "",
	"cookie":         "BSP_isRemember=false; BSP_isRememberEmail=jie.an%40bespinglobal.cn; BSP_LangCode=zh",
	"origin":         "https://sso.opsnow.cn",
	"referer":        "https://sso.opsnow.cn/loginForm.do",
	"sec-fetch-dest": "empty",
	"sec-fetch-site": "same-origin",
	"user-agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.122 Safari/537.36",
}

//
//var customRequestHeader1 = map[string]interface{}{}
//
//var serviceList1 = []string{"AUTH", "CPM", "FKP", "IMS", "KPI", "MAINPAGE", "PROSERVICE", "RM", "SAS", "SEARCH", "SECURITY",
//	"SMARTTHINGS", "SSO", "STS"}
//
//var respDataList1 []map[string]interface{}
//

func startJob() {
	// Create a Resty Client
	client := resty.New()
	client.SetDebug(true)
	for k, v := range basicRequestHeader1 {
		client.SetHeader(k, v)
	}
	// Using raw func into resty.SetRedirectPolicy
	client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
		// Implement your logic here
		if req.URL.RawQuery != "" {
			v, _ := url.ParseQuery(req.URL.RawQuery)
			if _, ok := v["access_token"]; ok {
				client.SetAuthToken(v["access_token"][0])
				client.SetCookie(&http.Cookie{
					Name:   "access_token",
					Value:  v["access_token"][0],
					Path:   "/",
					Domain: ".opsnow.cn",
				})
				client.SetCookie(&http.Cookie{
					Name:   "access_token_type",
					Value:  v["token_type"][0],
					Path:   "/",
					Domain: ".opsnow.cn",
				})
			}
		}
		//fmt.Println("Redirect Cookies:", client.GetClient().Jar)
		// return nil for continue redirect otherwise return error to stop/prevent redirect
		return nil
	}))
	//client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(15))
	client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		// Now you have access to Client and Request instance
		return nil // if its success otherwise return error
	})
	client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		// Now you have access to Client and Response instance
		// manipulate it as per your need
		return nil // if its success otherwise return error
	})
	loginResp, err := client.R().
		SetFormData(map[string]string{
			"username":      "email@email.com",
			"password":      "password@password",
			"client_id":     "wGBi_35lOOIxomIJRQp_cHxwBJka",
			"redirect_uri":  "https://www.opsnow.cn/dashboard",
			"scope":         "all",
			"response_type": "token",
		}).
		Post("https://sso.opsnow.cn/servicelogin")
	if err != nil {
		fmt.Println("post err", err)
	}
	fmt.Println(loginResp.StatusCode())
	//fmt.Println(client.Header)
	_, err = client.R().Get("https://project.opsnow.cn/")
	if err != nil {
		fmt.Println("post err", err)
	}
	_, err = client.R().Get("https://project.opsnow.cn/session/info")
	if err != nil {
		fmt.Println("post err", err)
	}
	_, err = client.R().SetHeader("customerid", "573a7278-aed8-4034-8e09-dff8ae40dd9b").
		SetHeader("accept", "application/json, text/javascript, */*; q=0.01").
		Get("https://project.opsnow.cn/prj/projectList?_=158803474931")
	if err != nil {
		fmt.Println("post err", err)
	}
}

func main() {
	startJob()
}
