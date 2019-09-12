package main

import (
	"flag"
	"fmt"
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
	r := finder.SingleQuerySearch(*query, wwwlist)
	if len(r) > 0 {
		fmt.Printf("Text [%s] contains the %d resources\n", *query, len(r))
		for i, ir := range r {
			fmt.Println("\t", i+1, ir)
		}
		return
	}
	fmt.Printf("Text [%s] does not contain any resources %v\n", *query, wwwlist)
}
