package medtest

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"strings"
	"testing"
)

func TestNodepageMethods(t *testing.T) {
	type Questions struct {
		Questions []Question `json:"questions"`
	}

	conf := Config{
		RootUrl:     "https://www.med-test.in.ua",
		SitemapPath: "/uk/site-map",
	}
	originPath := "/uk/node/list/tests/html/293"

	Convey("Given config with root url and sitemap path", t, func() {
		Convey("when test have only one page", func() {
			Convey("parse only one page", func() {
				qpage := QuestionsPage{
					Path: originPath,
				}
				expectedNumberOfPages := 1
				expectedIndex := "Krok M, Diagnostics laboratory, Archive, Infections, 2016"
				file, _ := ioutil.ReadFile("testdata/nodepage/without_pagination_result.json")
				expectedQuestions := Questions{}
				_ = json.Unmarshal([]byte(file), &expectedQuestions)

				defer gock.Off()
				gock.New(conf.RootUrl).
					Get(qpage.Path).
					MatchParams(map[string]string{
						"page": "1",
					}).
					Persist().
					Reply(200).
					File("testdata/nodepage/without_pagination.html")

				qpage.Parse(&conf)

				So(qpage.NumberOfPages, ShouldEqual, expectedNumberOfPages)
				So(strings.Join(qpage.Index, ", "), ShouldEqual, expectedIndex)
				So(Questions{qpage.Questions}, ShouldResemble, expectedQuestions)
			})
		})

		Convey("when test have more than one page", func() {
			Convey("parse all one pages", func() {
				qpage := QuestionsPage{
					Path: originPath,
				}
				expectedNumberOfPages := 2
				expectedIndex := "Krok M, Diagnostics laboratory, Archive, Infections, 2016"
				file, _ := ioutil.ReadFile("testdata/nodepage/full_test_result.json")
				expectedQuestions := Questions{}
				_ = json.Unmarshal([]byte(file), &expectedQuestions)

				defer gock.Off()
				gock.New(conf.RootUrl).
					Get(qpage.Path).
					MatchParams(map[string]string{
						"page": "1",
					}).
					Persist().
					Reply(200).
					File("testdata/nodepage/full_test.html")

				defer gock.Off()
				gock.New(conf.RootUrl).
					Get(qpage.Path).
					MatchParams(map[string]string{
						"page": "2",
					}).
					Persist().
					Reply(200).
					File("testdata/nodepage/full_test_2.html")

				qpage.Parse(&conf)

				So(qpage.NumberOfPages, ShouldEqual, expectedNumberOfPages)
				So(strings.Join(qpage.Index, ", "), ShouldEqual, expectedIndex)

				for i := range expectedQuestions.Questions {
					So(qpage.Questions, ShouldContain, expectedQuestions.Questions[i])
				}
			})
		})
	})
}
