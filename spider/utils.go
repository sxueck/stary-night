package spider

import (
	"lightning/config"
	"log"
	"strings"
)

func ObtainSiteTitle(url string) string {
	title, _, err := ExtraHtmlSourceCode(url)
	if err != nil {
		log.Println("failed to get the website name directly")
		if strings.Contains(err.Error(), "timeout") {
			title, _, err = ExtraHtmlSourceCode(url, WithProxyAddress(config.Cfg.ProxyAddr))
			if err == nil {
				return title
			}
			log.Println(err)
		}
		// if it is unknown, the title is stored in the database
		// first and then periodically retrieved later
		return "UnknownTitle"
	}

	return title
}

func GetAListOfWebSiteArticles() {

}
