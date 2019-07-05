package medtest

import (
	"errors"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestSitemapMethods(t *testing.T) {
	conf := Config{
		RootUrl:     "https://www.med-test.in.ua",
		SitemapPath: "/uk/site-map",
	}

	Convey("Given config with root url and sitemap path", t, func() {
		st := Sitemap{}

		Convey("#Parse", func() {
			Convey("when page successfully loaded", func() {
				defer gock.Off()
				gock.New(conf.RootUrl).
					Get(conf.SitemapPath).
					Reply(200).
					File("testdata/sitemap/valid.html")

				Convey("parse paths", func() {
					expected := "/path1, /path2"

					st.Parse(&conf)

					So(strings.Join(st.Paths, ", "), ShouldEqual, expected)
				})
			})

			Convey("when page can't be loaded", func() {
				defer gock.Off()
				gock.New(conf.RootUrl).
					Get(conf.SitemapPath).
					Reply(500)

				Convey("return error", func() {
					err := st.Parse(&conf)

					So(err, ShouldResemble, errors.New("Sitemap parsing failed, page broken"))
				})
			})

			Convey("when url fucked up", func() {
				Convey("return error", func() {
					config := Config{
						RootUrl:     "",
						SitemapPath: "",
					}
					err := st.Parse(&config)

					So(err, ShouldResemble, errors.New("Sitemap parsing failed. Error is: Get : unsupported protocol scheme \"\""))
				})
			})
		})
	})
}
