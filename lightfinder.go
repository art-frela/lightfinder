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
	"encoding/json"
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
	err          error
}

// execQuery - finds query string in the searchbody and write result
func (si *searchItem) execQuery() *searchItem {
	// TODO: filter stop word, using https://github.com/bbalet/stopwords
	// at next time ;-)
	// this is a simplest check for empty query or searchbody
	if si.query == "" || si.searcharea == "" {
		si.searchResult = false
		si.err = fmt.Errorf("empty query [%s] or searcharea [does not print, too much words], these must be filled", si.query)
		return si
	}
	si.searchResult = strings.Contains(si.searcharea, si.query)
	return si
}

// newSearchItem - builder of searchItem
func newSearchItem(query, searchbody, resource string) *searchItem {
	si := searchItem{
		query:        strings.ToLower(query),
		searcharea:   strings.ToLower(searchbody),
		resourceHref: resource,
	}
	return &si
}

// resource - simple structure for customize methods
type resource string

func (r resource) string() string {
	return string(r)
}

// getContent - fills content for resource, using http GET request
func (r resource) getContent() (string, error) {
	// simple validation
	if r == "" {
		err := fmt.Errorf("Empty request string")
		return "", err
	}
	bodybts, _, err := httpRequest(r.string())
	if err != nil {
		err = fmt.Errorf("Fail for get content from %s, %v", r.string(), err)
		return "", err
	}
	defer bodybts.Close()
	body, err := ioutil.ReadAll(bodybts)
	if err != nil {
		err = fmt.Errorf("Fail for read body content, from %s, %v", r.string(), err)
		return "", err
	}
	return string(body), nil
}

// newResources - adapter for slice of string to slice of resource
func newResources(items []string) []resource {
	resources := make([]resource, len(items))
	for i, item := range items {
		resources[i] = resource(item)
	}
	return resources
}

// httpRequest - common part of http request fot html search
func httpRequest(uri string) (rc io.ReadCloser, httpcode int, err error) {
	// set http client: timeout of request and switch off redirect
	c := http.Client{
		Timeout: requestTimeOut,
		// CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// 	return http.ErrUseLastResponse
		// },
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

// SingleQuery - single query for many resources
type SingleQuery struct {
	Search string   `json:"search"`
	Sites  []string `json:"sites"`
}

// QuerySearch - make search Query text at the Sites and returns Sites which contain Query
// If one resource responses httpcode<400 is eq success
func (sq *SingleQuery) QuerySearch() (containResources []string, err error) {
	// validate request
	if len(sq.Sites) == 0 {
		err = fmt.Errorf("Empty list of resources")
		return containResources, err
	}
	// define resources
	resources := newResources(sq.Sites)
	var wg sync.WaitGroup
	chResults := make(chan searchItem, len(sq.Sites))
	// process seasrch by resources
	for _, res := range resources {
		wg.Add(1)
		go func(ch chan searchItem, res resource) {
			content, err := res.getContent()
			if err != nil {
				result := searchItem{
					query:        sq.Search,
					resourceHref: res.string(),
					err:          err,
				}
				ch <- result
				wg.Done()
				return
			}
			sq.searchSingleResource(sq.Search, content, res.string(), &wg, ch)
		}(chResults, res)

	}
	go func() {
		wg.Wait()
		close(chResults)
	}()
	var errSearch searchErrors
	hasSuccessResponse := false
	for result := range chResults {
		//fmt.Printf(">>>GOT RESULT: %s\t%s\t%t\t%v\n", result.query, result.resourceHref, result.searchResult, result.err)
		if result.err != nil {
			errSearch.addError(result.resourceHref, result.err)
			continue
		}
		hasSuccessResponse = true
		if result.searchResult {
			containResources = append(containResources, result.resourceHref)
			continue
		}
		er := fmt.Errorf("resource does not contains [%s]", result.query)
		errSearch.addError(result.resourceHref, er)
	}
	if len(containResources) == 0 && !hasSuccessResponse { // for case when no one resource doesn't containr search string
		err = fmt.Errorf("Search result error, %s", errSearch.string())
	}
	return containResources, nil
}

// searchSingleResource - worker for search job
func (sq *SingleQuery) searchSingleResource(query, content, resource string, wg *sync.WaitGroup, out chan<- searchItem) {
	qi := newSearchItem(query, content, resource)
	qi.execQuery()
	out <- *qi
	wg.Done()
}

// NewSingleQuery - builder for SingleQuery
func NewSingleQuery(query string, links []string) *SingleQuery {
	sq := SingleQuery{
		Search: query,
		Sites:  links,
	}
	return &sq
}

// searchItemError - item of search error
type searchItemError struct {
	Resource string `json:"resource"`
	Error    string `json:"error"`
}

// searchErrors - slice of searchItemError
type searchErrors []searchItemError

// string - return string representation of search errors
func (se *searchErrors) string() string {
	bts, err := json.Marshal(*se)
	if err != nil {
		return fmt.Sprintf("%+v", *se)
	}
	return string(bts)
}

// addError - add error string to slice of search errors
func (se *searchErrors) addError(resource string, err error) {
	errItem := searchItemError{
		Resource: resource,
		Error:    err.Error(),
	}
	*se = append(*se, errItem)
}
