package medtest

import (
	"net/url"
	"strconv"
  "errors"
  "fmt"
)

type Config struct {
	RootUrl     string
	SitemapPath string
}

func (c *Config) CombineUrl(path string, pageNumber int) string {
	u, _ := url.Parse(c.RootUrl)
	pnumber := strconv.Itoa(pageNumber)

	u.Path = path
	query := u.Query()
	query.Set("page", pnumber)
	u.RawQuery = query.Encode()

	return u.String()
}

func (c *Config) PrepareSitemapUrl() string {
	u, _ := url.Parse(c.RootUrl)
	u.Path = c.SitemapPath

	return u.String()
}

func (c *Config) GetParamFromPath(path string, name string) string {
	u, _ := url.Parse(c.RootUrl)

	fullUrl, _ := u.Parse(path)

	query := fullUrl.Query()
	param := query.Get(name)

	return param
}

func BuildError(msg string, err error) error {
  var result error

  if err != nil {
    result = errors.New(fmt.Sprintf("%s. Error is: %s", msg, err.Error()))
  } else {
    result = errors.New(msg)
  }

  return result
}
