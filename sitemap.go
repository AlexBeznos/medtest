package medtest

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const sitemapItemsSelector = "#site-map a"

type Sitemap struct {
	Paths []string
}

// Public
func (s *Sitemap) Parse(conf *Config) {
	pageUrl := conf.PrepareSitemapUrl()

	// Request page
	res, err := http.Get(pageUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Status code error: %d, %s\n%s", res.StatusCode, res.Status, pageUrl)
	}

	// Load HTML document to goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find urls
	doc.Find(sitemapItemsSelector).Each(func(i int, anchor *goquery.Selection) {
		path, exist := anchor.Attr("href")

		if exist {
			s.Paths = append(s.Paths, path)
		}
	})
}
