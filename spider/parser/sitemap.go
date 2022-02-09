package parser

import (
	"encoding/xml"
	"fmt"
	"lightning/spider"
)

// Custom sitemap.xml Rule

type SiteMapXML struct {
	XMLName xml.Name `xml:"urlset"`
	//Text    string   `xml:",chardata"`
	Xmlns string     `xml:"xmlns,attr"`
	Xhtml string     `xml:"xhtml,attr"`
	URL   SiteMapURL `xml:"url"`
}

type SiteMapURL []struct {
	//Text    string `xml:",chardata"`
	Loc     string `xml:"loc"`
	Lastmod string `xml:"lastmod"`
}

func UnmarshalSitemap(xmlBody *[]byte) (*SiteMapXML, error) {
	sml := &SiteMapXML{}
	err := xml.Unmarshal(*xmlBody, sml)
	if err != nil {
		return nil, fmt.Errorf("XML unmarshal parsing error : %s", err)
	}

	return sml, nil
}

func ExtractSomeTitles(urls *SiteMapURL) ([]string, error) {
	var title []string
	for _, v := range *urls {
		u := v.Loc
		t := spider.ObtainSiteTitle(u)

		if len(t) < 15 {
			continue
		}

		title = append(title, t)
	}

	return title, nil
}
