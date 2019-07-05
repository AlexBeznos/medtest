package medtest

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const sitemapItemsSelector = "#site-map a"

type Sitemap struct {
	Paths []string
}

// Public
func (s *Sitemap) Parse(conf *Config) error {
  var err error

	pageUrl := conf.PrepareSitemapUrl()

	// Request page
	res, err := http.Get(pageUrl)
	if err != nil {
    return BuildError("Sitemap parsing failed", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
    return BuildError("Sitemap parsing failed, page broken", err)
	}

	// Load HTML document to goquery
	doc, _ := goquery.NewDocumentFromReader(res.Body)

	// Find urls
	doc.Find(sitemapItemsSelector).Each(func(i int, anchor *goquery.Selection) {
		path, exist := anchor.Attr("href")

		if exist {
			s.Paths = append(s.Paths, path)
		}
	})

  return nil
}
