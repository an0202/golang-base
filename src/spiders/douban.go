package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
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

// printMoviesList by regrex
func printMoviesListv2(contents []byte) []string {
	// regexp
	var movieURLs []string
	re := regexp.MustCompile(`<a href="(https://movie.douban.com/subject/[0-9]+/)".*\n.*<span class="title">(.*)</span>\n(.*\n){1,2}.*</a>`)
	matches := re.FindAllSubmatch(contents, -1)

	for _, m := range matches {
		// fmt.Printf("Title: %s,URL: %s\n", m[2], m[1])
		movieURLs = append(movieURLs, string(m[1]))
	}
	return movieURLs
}

//printMoviesList by goquery
func printMoviesList(contents io.Reader) []string {
	// regexp
	var movieURLs []string
	dom, err := goquery.NewDocumentFromReader(contents)
	if err != nil {
		log.Fatalln(err)
	}

	dom.Find("div.hd").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Find("a").Attr("href"))
		fmt.Println(selection.Find("a>span.title").First().Text())
	})

	dom.Find("div")

	return movieURLs
}

// describeURL return json data in URL
func describeURL(URL string) {
	// http get
	resp, err := http.Get(URL)
	if err != nil {
		tools.ErrorLogger.Fatalln(err)
	}
	contents, err1 := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err1 != nil {
		tools.ErrorLogger.Fatalln(err)
	}
	// <script type="application/ld\+json">\n(.*\n)+?</script>
	re := regexp.MustCompile(`{\n\s+"@context":(.*\n)+?}`)
	matches := re.Find(contents)
	md := MovieDetail{}
	err2 := json.Unmarshal(matches, &md)
	if err2 != nil {
		fmt.Println(err2, "handel error ...")
		// handel error for unkonw charactor \n
		result := bytes.Replace(matches, []byte("\n"), []byte(""), -1)
		err3 := json.Unmarshal(result, &md)
		if err3 != nil {
			fmt.Println("Unkonwn Error:", err3)
		}
	}
	movieURL := "https://movie.douban.com" + md.URL
	fmt.Printf("Name: %s URL: %s Rating: %s (RatingCount: %s)\n", md.Name, movieURL, md.AggregateRating["ratingValue"], md.AggregateRating["ratingCount"])
	// time.Sleep(2 * time.Second)
}

func main() {
	for i := 0; i <= 0; {
		URL := "https://movie.douban.com/top250/?start=" + strconv.Itoa(i)
		resp, err := http.Get(URL)
		if err != nil {
			tools.ErrorLogger.Println(err)
		}
		// use for printMoviesListv2
		// record, err := ioutil.ReadAll(resp.Body)
		// defer resp.Body.Close()
		// if err != nil {
		// 	tools.ErrorLogger.Println(err)
		// }
		// moviesurllist := printMoviesList(record)
		// for _, movieurl := range moviesurllist {
		// 	// fmt.Println("start to process ", movieurl)
		// 	describeURL(movieurl)
		// }
		// i += 25

		// use for printMoviesList
		printMoviesList(resp.Body)
		i++
	}
	// resp, err := http.Get("https://movie.douban.com/subject/1292052/")
	// if err != nil {
	// 	panic(err)
	// }
	// record, err := ioutil.ReadAll(resp.Body)
	// defer resp.Body.Close()
	// if err != nil {
	// 	panic(err)
	// }
	// describeURL("https://movie.douban.com/subject/1292720/")
}
