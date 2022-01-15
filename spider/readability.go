package spider

import (
	"fmt"
	"github.com/go-shiori/go-readability"
	"golang.org/x/net/html"
	"strings"
)

// ExtraHtmlSourceCode also need a set of headers for specific sites
func ExtraHtmlSourceCode(url string, options ...*Option) (string, []byte, error) {
	br, err := HttpRequestToBytes("GET", url, options...)
	if err != nil {
		return "", nil, fmt.Errorf("error or null fetch result : %s", err)
	}

	var article = readability.Article{}
	doc, _ := html.Parse(strings.NewReader(string(*br)))

	article, err = readabilityContent(doc)
	if err != nil {
		return "", nil, err
	}

	return article.Title, []byte(article.Content), nil
}

func readabilityContent(doc *html.Node) (readability.Article, error) {
	ps := &readability.Parser{}
	return ps.ParseDocument(doc, nil)
}
