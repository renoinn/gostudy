package main

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	Crawl("https://cookpad.com/recipe/7091250", 1)
}

// リクエスト
type Request struct {
	url   string
	depth int
}

// 結果
type Result struct {
	err error
	url string
}

// チャンネル
type Channels struct {
	req  chan Request
	res  chan Result
	quit chan int
}

func Crawl(url string, depth int) {
	body, urls, err := Fetch(url)
	if err != nil {
		fmt.Printf("fetch error %s", url)
	}

	fmt.Printf("body: \n%s", body)

	for _, u := range urls {
		fmt.Printf("%s\n", u)
	}
}

func Fetch(u string) (string, []string, error) {
	body := ""
	urls := make([]string, 0)

	baseUrl, err := url.Parse(u)
	if err != nil {
		return body, urls, err
	}

	resp, err := http.Get(baseUrl.String())
	if err != nil {
		return body, urls, err
	}
	defer resp.Body.Close()

	docs, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return body, urls, err
	}
	body = docs.Find("body").Text()

	docs.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, exist := s.Attr("href")
		if exist {
			isSameHost := check_regexp(`^\/`, href)
			if isSameHost && !Contains(urls, href) {
				urls = append(urls, href)
			}
		}
	})

	return body, urls, err
}

func check_regexp(reg, str string) bool {
	return regexp.MustCompile(reg).Match([]byte(str))
}

func Contains(list interface{}, elem interface{}) bool {
	listV := reflect.ValueOf(list)

	if listV.Kind() == reflect.Slice {
		for i := 0; i < listV.Len(); i++ {
			item := listV.Index(i).Interface()
			// 型変換可能か確認する
			if !reflect.TypeOf(elem).ConvertibleTo(reflect.TypeOf(item)) {
				continue
			}
			// 型変換する
			target := reflect.ValueOf(elem).Convert(reflect.TypeOf(item)).Interface()
			// 等価判定をする
			if ok := reflect.DeepEqual(item, target); ok {
				return true
			}
		}
	}
	return false
}
