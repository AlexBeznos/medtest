package medtest

import (
	"log"
	"net/url"
	"strconv"
)

type Config struct {
	RootUrl     string
	SitemapPath string
}

func (c *Config) CombineUrl(path string, pageNumber int) string {
	u, err := url.Parse(c.RootUrl)
	if err != nil {
		log.Fatal(err)
	}
	pnumber := strconv.Itoa(pageNumber)

	u.Path = path
	query := u.Query()
	query.Set("page", pnumber)
	u.RawQuery = query.Encode()

	return u.String()
}

func (c *Config) PrepareSitemapUrl() string {
	u, err := url.Parse(c.RootUrl)
	if err != nil {
		log.Fatal(err)
	}

	u.Path = c.SitemapPath

	return u.String()
}

func (c *Config) GetParamFromPath(path string, name string) string {
	u, err := url.Parse(c.RootUrl)
	if err != nil {
		log.Fatal(err)
	}

	fullUrl, err := u.Parse(path)
	if err != nil {
		log.Fatal(err)
	}

	query := fullUrl.Query()
	param := query.Get(name)

	return param
}
