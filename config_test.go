package medtest

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConfigMethods(t *testing.T) {
	conf := Config{
		RootUrl:     "https://www.med-test.in.ua/",
		SitemapPath: "/uk/site-map",
	}

	Convey("Given config with root url and sitemap path", t, func() {
		Convey("#GetParamFromPath", func() {
			Convey("When param exist in path it should return result", func() {
				param := conf.GetParamFromPath("/uk/site-map?page=10", "page")

				So(param, ShouldEqual, "10")
			})

			Convey("When param not exists it should return empty string", func() {
				param := conf.GetParamFromPath("/uk/site-map?page=10", "limit")

				So(param, ShouldEqual, "")
			})
		})

		Convey("#PrepareSitemapUrl", func() {
			Convey("it return full sitemap url", func() {
				url := conf.PrepareSitemapUrl()

				So(url, ShouldEqual, "https://www.med-test.in.ua/uk/site-map")
			})
		})

		Convey("#CombineUrl", func() {
			Convey("return full url with required page number", func() {
				number := 125
				path := "uk/node/path"
				url := conf.CombineUrl(path, number)
				expected := fmt.Sprintf("%s%s?page=%d", conf.RootUrl, path, number)

				So(url, ShouldEqual, expected)
			})
		})
	})
}
