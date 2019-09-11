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
	"strings"
)

// SearchItem - atomic item for search process
type SearchItem struct {
	query        string
	searcharea   string
	resourceHref string
	searchResult bool
}

// setQuery - set query value
func (si *SearchItem) setQuery(q string) *SearchItem {
	si.query = q
	return si
}

// setSearchArea - set searcharea value
func (si *SearchItem) setSearchArea(body string) *SearchItem {
	si.searcharea = body
	return si
}

// setResourceHref - set resourceHref value
func (si *SearchItem) setResourceHref(url string) *SearchItem {
	si.resourceHref = url
	return si
}

// setSearchResult - set searchResult value
func (si *SearchItem) setSearchResult(res bool) *SearchItem {
	si.searchResult = res
	return si
}

// ExecQuery - finds query string in the searchbody and write result
func (si *SearchItem) ExecQuery() *SearchItem {
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

// NewQueryItem - builder of SearchItem
func NewQueryItem(query, searchbody, resource string) *SearchItem {
	si := new(SearchItem)
	si.setQuery(query)
	si.setSearchArea(searchbody)
	si.setResourceHref(resource)
	return si
}
