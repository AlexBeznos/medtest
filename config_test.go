package medtest

import (
	"fmt"
  "errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConfigMethods(t *testing.T) {
	conf := Config{
		RootUrl:     "https://www.med-test.in.ua/",
		SitemapPath: "/uk/site-map",
	}

	Convey("Given config with root url and sitemap path", t, func() {
		Convey("#CombineUrl", func() {
			Convey("return full url with required page number", func() {
				number := 125
				path := "uk/node/path"
        url := conf.CombineUrl(path, number)
				expected := fmt.Sprintf("%s%s?page=%d", conf.RootUrl, path, number)

				So(url, ShouldEqual, expected)
			})
		})

		Convey("#PrepareSitemapUrl", func() {
			Convey("it return full sitemap url", func() {
				url := conf.PrepareSitemapUrl()

				So(url, ShouldEqual, "https://www.med-test.in.ua/uk/site-map")
			})
		})

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
	})

  Convey("#BuildError", t, func() {
    Convey("when error provided", func() {
      Convey("return error with provided error inside", func() {
        err := errors.New("some shit happened")
        msg := "Hello world"
        result := errors.New("Hello world. Error is: some shit happened")

        So(BuildError(msg, err), ShouldResemble, result)
      })
    })

    Convey("when error not provided", func() {
      Convey("return error with message only", func() {
        msg := "Hello world"
        result := errors.New("Hello world")

        So(BuildError(msg, nil), ShouldResemble, result)
      })
    })
  })
}
