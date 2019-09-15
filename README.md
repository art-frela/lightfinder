# lightfinder - Go package for simplest search text at the web sites and provide links those contains query string

(homework#1.1 for [Geekbrains Go course](https://geekbrains.ru/geek_university/golang), 2nd qrt "Essentials of Golang")

## Purpose

Get query string and list of resource, exec search and return list of resource those contain query.  

## Install

`go get github.com/art-frela/lightfinder`

## Example

```golang
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	finder "github.com/art-frela/lightfinder"
)

func main() {
	// flags set and Parse
	query := flag.String("q", "Чак Норрис", "Query string for search at the web sites")
	list := flag.String("w", "https://google.com;https://ru.wikipedia.org/wiki/Норрис,_Чак;https://ru.wikipedia.org/wiki/Крутой_Уокер:_Правосудие_по-техасски", "List of websites semicolon separated")
	flag.Parse()
	// split list to slice of links
	wwwlist := strings.Split(*list, ";")
	sq := finder.NewSingleQuery(*query, wwwlist)
	r, err := sq.QuerySearch()
	if err != nil {
		log.Println("search error", err)
		return
	}
	if len(r) > 0 {
		fmt.Printf("Text [%s] contains %d resources\n", *query, len(r))
		for i, ir := range r {
			fmt.Println("\t", i+1, ir)
		}
		return
	}
	fmt.Printf("Text [%s] does not contain any resources %v\n", *query, wwwlist)
}
```
