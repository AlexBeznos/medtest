package medtest

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
	"strings"
	"testing"
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
					File("testdata/sitemap.html")

				Convey("parse paths", func() {
					expected := "/path1, /path2"

					st.Parse(&conf)

					So(strings.Join(st.Paths, ", "), ShouldEqual, expected)
				})
			})
		})
	})
}
