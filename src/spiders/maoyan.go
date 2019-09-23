package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"tools"

	"github.com/PuerkitoBio/goquery"
)

type MovieDetail struct {
	Context         string                 `json:"@context"`
	Name            string                 `json:"name"`
	URL             string                 `json:"url"`
	Image           string                 `json:"image"`
	Director        []map[string]string    `json:"director"`
	Author          []map[string]string    `json:"author"`
	Actor           []map[string]string    `json:"actor"`
	DatePublished   string                 `json:"datePublished"`
	Genre           []string               `json:"genre"`
	Duration        string                 `json:"duration"`
	Description     string                 `json:"description"`
	Type            string                 `json:"@type"`
	AggregateRating map[string]interface{} `json:"aggregateRating"`
}

//printMoviesList by goquery
func printMoviesList(contents io.Reader) []string {
	// regexp
	var movieURLs []string
	dom, err := goquery.NewDocumentFromReader(contents)
	if err != nil {
		log.Fatalln(err)
	}

	dom.Find("p.name>a[href]").Each(func(i int, selection *goquery.Selection) {
		title, _ := selection.Attr("title")
		url, _ := selection.Attr("href")
		fmt.Printf("Title: %s, URL: %s\n", title, "https//maoyan.com"+url)
		movieURLs = append(movieURLs, "https//maoyan.com"+url)
	})
	return movieURLs
}

// describeURL return json data in URL
func describeURL(URL string) {
	// http get
	resp, err := http.Get(URL)
	if err != nil {
		tools.ErrorLogger.Fatalln(err)
	}
	dom, err1 := goquery.NewDocumentFromReader(resp.Body)
	defer resp.Body.Close()
	if err1 != nil {
		tools.ErrorLogger.Fatalln(err)
	}
	// <script type="application/ld\+json">\n(.*\n)+?</script>
	dom.Find("div.banner").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})
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
	// 	// moviesurllist := printMoviesList(record)
	// 	// for _, movieurl := range moviesurllist {
	// 	// 	// fmt.Println("start to process ", movieurl)
	// 	// 	describeURL(movieurl)
	// 	// }
	// 	// i += 25

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
	describeURL("https://maoyan.com/films/1383")
}
