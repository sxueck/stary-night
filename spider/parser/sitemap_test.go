package parser

import (
	"lightning/spider"
	"log"
	"testing"
)

func TestExtractAllTitles(t *testing.T) {
	url := "https://blog.lingyf.com/sitemap.xml"
	content, err := spider.HttpRequestToBytes("GET", url)
	if err != nil {
		t.Failed()
	}

	sitemap, err := UnmarshalSitemap(content)
	if err != nil {
		t.Failed()
	}

	log.Printf("%s", sitemap)
	titles, err := ExtractSomeTitles(&sitemap.URL)
	if err != nil {
		t.Failed()
	}

	for _,v := range titles {
		log.Println(v)
	}
}
