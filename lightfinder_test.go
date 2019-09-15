package lightfinder

import (
	"testing"
)

type testDataItem struct {
	query   string
	links   []string
	results []string
}

var testData = []testDataItem{
	{"wikipedia", []string{"https://ru.wikipedia.org/wiki/Вики",
		"https://www.watcom.ru/",
		"https://www.rbc.ru/",
		"https://vk.com",
		"http://newsvolga.com/",
		"https://vz.ru/",
		"http://regions.ru/",
		"https://www.watcom.ru/",
		"https://www.rbc.ru/",
		"https://vk.com",
		"http://newsvolga.com/",
		"https://vz.ru/",
		"http://regions.ru/",
		"https://www.watcom.ru/",
		"https://www.rbc.ru/",
		"https://vk.com",
		"https://google.com",
		"http://newsvolga.com/",
		"https://vz.ru/",
		"http://regions.ru/",
		"https://www.watcom.ru/",
		"https://www.rbc.ru/",
		"https://vk.com",
		"https://google.com",
		"http://newsvolga.com/",
		"https://vz.ru/",
		"http://regions.ru/",
		"https://www.apple.com",
		"https://www.sports.ru/nba/?gr=www",
		"https://www.zebra.com/ru/ru.html",
		"https://intraservice.ru"},
		[]string{"https://ru.wikipedia.org/wiki/Вики"}},
	{"система подсчета посетителей", []string{"https://www.watcom.ru/", "https://ya.ru", "https://www.rbc.ru/"}, []string{"https://www.watcom.ru/"}},
	{"новости", []string{"https://www.watcom.ru/", "https://ya.ru", "https://www.rbc.ru/"}, []string{"https://www.watcom.ru/", "https://ya.ru", "https://www.rbc.ru/"}},
	{"Чак НорРис", []string{"https://google.com", "https://ru.wikipedia.org/wiki/%D0%9D%D0%BE%D1%80%D1%80%D0%B8%D1%81,_%D0%A7%D0%B0%D0%BA", "https://ru.wikipedia.org/wiki/%D0%9A%D1%80%D1%83%D1%82%D0%BE%D0%B9_%D0%A3%D0%BE%D0%BA%D0%B5%D1%80:_%D0%9F%D1%80%D0%B0%D0%B2%D0%BE%D1%81%D1%83%D0%B4%D0%B8%D0%B5_%D0%BF%D0%BE-%D1%82%D0%B5%D1%85%D0%B0%D1%81%D1%81%D0%BA%D0%B8"},
		[]string{"https://ru.wikipedia.org/wiki/%D0%9D%D0%BE%D1%80%D1%80%D0%B8%D1%81,_%D0%A7%D0%B0%D0%BA", "https://ru.wikipedia.org/wiki/%D0%9A%D1%80%D1%83%D1%82%D0%BE%D0%B9_%D0%A3%D0%BE%D0%BA%D0%B5%D1%80:_%D0%9F%D1%80%D0%B0%D0%B2%D0%BE%D1%81%D1%83%D0%B4%D0%B8%D0%B5_%D0%BF%D0%BE-%D1%82%D0%B5%D1%85%D0%B0%D1%81%D1%81%D0%BA%D0%B8"}},
}

// TestSingleQuerySearch - testing of SingleQuerySearch
func TestSingleQuerySearch(t *testing.T) {
	for _, tItem := range testData {
		sq := NewSingleQuery(tItem.query, tItem.links)
		realFind, err := sq.QuerySearch()
		if err != nil {
			t.Errorf("Can't make test, too much errors, %v", err)
		}
		if len(realFind) != len(tItem.results) {
			t.Errorf("got %v slice of resources, need to %v (query=%s)", realFind, tItem.results, tItem.query)
		} else {
			// detail comparison
			for _, realItem := range realFind {
				exists := false
				for _, expected := range tItem.results {
					if expected == realItem {
						exists = true
					}
				}
				if !exists {
					t.Errorf("wrong content, got %v slice of resources, need to %v  (query=%s)", realFind, tItem.results, tItem.query)
				}
			}
		}
	}
}

func BenchmarkSingleQuerySearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tItem := range testData {
			sq := NewSingleQuery(tItem.query, tItem.links)
			sq.QuerySearch()
		}
	}
}
