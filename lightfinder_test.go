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
	{"wiki", []string{"https://ru.wikipedia.org/wiki/Вики",
		"https://www.watcom.ru/",
		"https://www.rbc.ru/",
		"https://vk.com",
		"https://google.com",
		"http://newsvolga.com/",
		"https://vz.ru/",
		"http://regions.ru/",
		"https://news.google.com/?hl=ru&gl=RU&ceid=RU:ru"},
		[]string{"https://ru.wikipedia.org/wiki/Вики"}},
	{"система подсчета посетителей", []string{"https://www.watcom.ru/", "https://ya.ru", "https://www.rbc.ru/"}, []string{"https://www.watcom.ru/"}},
	{"новости", []string{"https://www.watcom.ru/", "https://ya.ru", "https://www.rbc.ru/"}, []string{"https://www.watcom.ru/", "https://ya.ru", "https://www.rbc.ru/"}},
}

// TestSingleQuerySearch - testing of SingleQuerySearch
func TestSingleQuerySearch(t *testing.T) {
	for _, tItem := range testData {
		realFind := SingleQuerySearch(tItem.query, tItem.links)
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
			SingleQuerySearch(tItem.query, tItem.links)
		}
	}
}
