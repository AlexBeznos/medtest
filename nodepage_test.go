package medtest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
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

		Convey("when request for metadata failed", func() {
			Convey("return error", func() {
				qpage := QuestionsPage{
					Path: originPath,
				}
				defer gock.Off()
				gock.New(conf.RootUrl).
					Get(qpage.Path).
					MatchParams(map[string]string{
						"page": "1",
					}).
					Persist().
					Reply(500)
				expectedError := errors.New("Page broken. URL: https://www.med-test.in.ua/uk/node/list/tests/html/293?page=1")

				err := qpage.Parse(&conf)

				So(err, ShouldResemble, expectedError)
			})
		})

		Convey("when root path and path empty", func() {
			Convey("return error", func() {
				qpage := QuestionsPage{
					Path: "",
				}
				config := Config{
					RootUrl: "",
				}
				expectedError := errors.New("Request failed. Error is: Get ?page=1: unsupported protocol scheme \"\"")

				err := qpage.Parse(&config)

				So(err, ShouldResemble, expectedError)
			})
		})

		Convey("when one of the pages failed", func() {
			Convey("return error", func() {
				qpage := QuestionsPage{
					Path: originPath,
				}
				expectedNumberOfPages := 2
				expectedIndex := "Krok M, Diagnostics laboratory, Archive, Infections, 2016"
				expectedError := errors.New("Page broken. URL: https://www.med-test.in.ua/uk/node/list/tests/html/293?page=2")

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
					Reply(500)

				err := qpage.Parse(&conf)

				So(qpage.NumberOfPages, ShouldEqual, expectedNumberOfPages)
				So(strings.Join(qpage.Index, ", "), ShouldEqual, expectedIndex)
				So(err, ShouldResemble, expectedError)
				So(len(qpage.Questions), ShouldEqual, 0)
			})
		})
	})
}
