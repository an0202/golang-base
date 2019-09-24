package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"tools"

	"github.com/PuerkitoBio/goquery"
)

type DishesDetail struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	Rating        string `json:"rating"`
	Cooked        string `json:"cooked"`
	DatePublished string `json:"datepublished"`
	Collected     string `json:"collected"`
}

//printDishesList by goquery
func printDishesList(URL string) []string {
	// http get
	fmt.Println("start")
	resp, err := http.Get(URL)
	if err != nil {
		tools.ErrorLogger.Fatalln(err)
	}
	switch resp.StatusCode {
	case 404:
		tools.ErrorLogger.Fatalln(resp.StatusCode)
	}
	dom, err1 := goquery.NewDocumentFromReader(resp.Body)
	defer resp.Body.Close()
	if err1 != nil {
		tools.ErrorLogger.Fatalln(err)
	}
	// regexp
	var DishesURLs []string
	dom.Find("p.name>a[href]").Each(func(i int, selection *goquery.Selection) {
		re := regexp.MustCompile(`\S+(\s\S+)*`)
		title := strings.ReplaceAll(re.FindString(selection.Text()), "\n", " ")
		url, _ := selection.Attr("href")
		fmt.Printf("Title: %s, URL: %s\n", title, "http://www.xiachufang.com"+url)
		DishesURLs = append(DishesURLs, "http://www.xiachufang.com"+url)
	})
	return DishesURLs
}

// describeDishes return json data in URL
func describeDishes(URL string) {
	// http get with custom header
	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		tools.ErrorLogger.Fatalln(err)
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36")
	resp, err := http.DefaultClient.Do(request)
	dom, err1 := goquery.NewDocumentFromReader(resp.Body)
	defer resp.Body.Close()
	if err1 != nil {
		tools.ErrorLogger.Fatalln(err)
	}
	//
	Dishes := new(DishesDetail)
	re := regexp.MustCompile(`\S+(\s\S+)*`)
	Dishes.Name = re.FindString(dom.Find("h1.page-title[itemprop=name]").Text())
	Dishes.URL = URL
	Dishes.Rating = dom.Find("span.number[itemprop=ratingValue]").Text()
	Dishes.Cooked = dom.Find("div.cooked.float-left>.number").Text()
	Dishes.DatePublished = dom.Find("span[itemprop=datePublished]").Text()
	Dishes.Collected = strings.Split(dom.Find("div.pv").Text(), " ")[0]
	//
	b, err := json.Marshal(Dishes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)
	// time.Sleep(3 * time.Second)
}

func main() {
	// for i := 0; i <= 0; {
	// 	URL := "https://maoyan.com/board/4?offset=" + strconv.Itoa(i)
	// 	resp, err := http.Get(URL)
	// 	if err != nil {
	// 		tools.ErrorLogger.Println(err)
	// 	}
	// 	defer resp.Body.Close()
	// 	// use for printMoviesListv2
	// 	// record, err := ioutil.ReadAll(resp.Body)
	// 	// defer resp.Body.Close()
	// 	// if err != nil {
	// 	// 	tools.ErrorLogger.Println(err)
	// 	// }
	// moviesurllist := printMoviesList(record)
	// for _, movieurl := range moviesurllist {
	// 	// fmt.Println("start to process ", movieurl)
	// 	describeURL(movieurl)
	// }
	// i += 25

	// 	// use for printMoviesList
	// 	printMoviesList(resp.Body)
	// 	i += 10
	// }
	// resp, err := http.Get("https://movie.douban.com/subject/1292052/")
	// if err != nil {
	// 	panic(err)
	// }
	// record, err := ioutil.ReadAll(resp.Body)
	// defer resp.Body.Close()
	// if err != nil {
	// 	panic(err)
	// }
	for i := 1; i <= 2; {
		URL := "http://www.xiachufang.com/explore/monthhonor/201812/?page=" + strconv.Itoa(i)
		DishesList := printDishesList(URL)
		wg := sync.WaitGroup{}
		for _, dishesurl := range DishesList {
			wg.Add(1)
			go func(dishesurl string) {
				defer wg.Done()
				describeDishes(dishesurl)
			}(dishesurl)
		}
		wg.Wait()
		i++
	}
	// describeURL("http://www.xiachufang.com/recipe/102180897/")
}
