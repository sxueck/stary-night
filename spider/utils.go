package spider

func ObtainSiteTitle(url string) string {
	title, _, err := ExtraHtmlSourceCode(url)
	if err != nil {
		// if it is unknown, the title is stored in the database
		// first and then periodically retrieved later
		return "UnknownTitle"
	}

	return title
}

func GetAListOfWebSiteArticles() {

}
