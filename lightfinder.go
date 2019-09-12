/*Package lightfinder is a very simple text finder of query string at the specified web resources

Author: Artem Karpov mailto:art.frela@gmail.com
Date: 2019-09-11
Subject: Geekbrains Go course, 2nd qrt "Essentials of Golang"

Task: 1. Напишите функцию, которая будет получать на вход строку с поисковым запросом (string) и массив ссылок на страницы,
по которым стоит произвести поиск ([]string). Результатом работы функции должен быть массив строк со ссылками на страницы,
на которых обнаружен поисковый запрос. Функция должна искать точное соответствие фразе в тексте ответа от сервера по каждой из ссылок.
*/
package lightfinder

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	requestTimeOut = 30 * time.Second
)

// searchItem - atomic item for search process
type searchItem struct {
	query        string
	searcharea   string
	resourceHref string
	searchResult bool
}

// setQuery - set query value
func (si *searchItem) setQuery(q string) *searchItem {
	si.query = q
	return si
}

// setSearchArea - set searcharea value
func (si *searchItem) setSearchArea(body string) *searchItem {
	si.searcharea = body
	return si
}

// setResourceHref - set resourceHref value
func (si *searchItem) setResourceHref(url string) *searchItem {
	si.resourceHref = url
	return si
}

// setSearchResult - set searchResult value
func (si *searchItem) setSearchResult(res bool) *searchItem {
	si.searchResult = res
	return si
}

// execQuery - finds query string in the searchbody and write result
func (si *searchItem) execQuery() *searchItem {
	// TODO: filter stop word, using https://github.com/bbalet/stopwords
	// at next time ;-)
	// this is a simplest check for empty query or searchbody
	if si.query == "" || si.searcharea == "" {
		si.searchResult = false
		return si
	}
	si.setSearchResult(strings.Contains(si.searcharea, si.query))
	return si
}

// newQueryItem - builder of searchItem
func newQueryItem(query, searchbody, resource string) *searchItem {
	si := new(searchItem)
	si.setQuery(strings.ToLower(query))
	si.setSearchArea(strings.ToLower(searchbody))
	si.setResourceHref(resource)
	return si
}

// resource - simple structure for customize methods
type resource string

func newResources(items []string) []resource {
	resources := make([]resource, len(items))
	for i, item := range items {
		resources[i] = resource(item)
	}
	return resources
}

func (r resource) string() string {
	return string(r)
}

// getContent - fills content for resource, using http GET request
func (r resource) getContent() string {
	// simple validation
	if r == "" {
		return ""
	}
	bodybts, _, err := httpRequest(r.string())
	if err != nil {
		return ""
	}
	body, err := ioutil.ReadAll(bodybts)
	if err != nil {
		return ""
	}
	return string(body)
}

// httpRequest - common part of http request
func httpRequest(uri string) (rc io.Reader, httpcode int, err error) {
	// set http client: timeout of request and switch off redirect
	c := http.Client{
		Timeout: requestTimeOut,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	// make encoded url string, for cyrillic symbols and other
	u, err := url.Parse(uri)
	if err != nil {
		return
	}
	q := u.Query()
	u.RawQuery = q.Encode()
	// make request
	request, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return
	}
	request.Header.Set("Accept", "text/html")
	request.Header.Set("User-Agent", "SampleGoClient/1.0")

	httpData, err := c.Do(request)
	if err != nil {
		return
	}
	httpcode = httpData.StatusCode
	if httpData.StatusCode != http.StatusOK {
		err = fmt.Errorf("some error at time http.request, request=%s; httpcode=%d", u.String(), httpData.StatusCode)
		return
	}
	rc = httpData.Body
	return
}

// SingleQuerySearch - executes simple text search for single query at the many resources
func SingleQuerySearch(q string, links []string) (containResources []string) {
	resources := newResources(links)
	var wg sync.WaitGroup
	chResults := make(chan string, len(links))
	for _, res := range resources {
		wg.Add(1)
		go searchSingleResource(q, res.getContent(), res.string(), &wg, chResults)
	}
	go func() {
		wg.Wait()
		close(chResults)
	}()
	for existResource := range chResults {
		containResources = append(containResources, existResource)
	}
	return
}

// searchSingleResource - worker for search job
func searchSingleResource(query, content, resource string, wg *sync.WaitGroup, out chan<- string) {
	qi := newQueryItem(query, content, resource)
	qi.execQuery()
	if qi.searchResult {
		out <- resource
	}
	wg.Done()
}
