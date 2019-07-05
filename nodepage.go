package medtest

import (
	"log"
	"net/http"
	"strconv"
  "fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type QuestionsPage struct {
	Path          string     `json:"path"`
	Questions     []Question `json:"questions"`
	Index         []string   `json:"index"`
	NumberOfPages int        `json:"numberOfPages"`
}
type Question struct {
	Text    string   `json:"text"`
	Answers []Answer `json:"answers"`
}

type Answer struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
}

const lastPageLinkSelector = ".navigation .pagination .last a"
const breadcrumbsSelector = "ul.breadcrumb li"
const questionItemsSelection = "#tests-content .container .row .test .panel"
const answerItemsSelection = ".answer .list-group-item"

func (qpage *QuestionsPage) Parse(conf *Config) error {
  err := qpage.getMetaData(conf)
  if err != nil {
    return err
  }

	chQuestions := make(chan Question)
	chFinished := make(chan bool)
  chErrors := make(chan error)

	// Kick off the parsing
	for i := 1; i <= qpage.NumberOfPages; i++ {
		url := conf.CombineUrl(qpage.Path, i)
		go parseQuestions(url, chQuestions, chErrors, chFinished)
	}

	// Subscription to parsed questions
	for c := 0; c < qpage.NumberOfPages; {
		select {
		case question := <-chQuestions:
			qpage.Questions = append(qpage.Questions, question)
    case err := <-chErrors:
      qpage.Questions = make([]Question, 0)
      return err
		case <-chFinished:
			c++
		}
	}

  return nil
}

func (qpage *QuestionsPage) getMetaData(conf *Config) error {
	url := conf.CombineUrl(qpage.Path, 1)

	// Request page
	res, err := http.Get(url)
	if err != nil {
		return BuildError("Request failed", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
    msg := fmt.Sprintf("Page broken. URL: %s", url)
		return BuildError(msg, nil)
	}

	// Load HTML document to goquery
	doc, _ := goquery.NewDocumentFromReader(res.Body)

	// Is page have pagination
	if doc.Find(lastPageLinkSelector).Length() != 0 {
		lastPageLink := doc.Find(lastPageLinkSelector)
		href, _ := lastPageLink.Attr("href")
		pageParam := conf.GetParamFromPath(href, "page")
		numberOfPages, err := strconv.Atoi(pageParam)
		if err != nil {
			log.Fatal(err)
		}
		qpage.NumberOfPages = numberOfPages
	} else {
		qpage.NumberOfPages = 1
	}

	// Get breadcrumbs
	doc.Find(breadcrumbsSelector).Each(qpage.parseBreadcrumb)

  return nil
}

func parseQuestions(url string, chQuestions chan Question, chErrors chan error, chFinish chan bool) {
	defer func() {
		chFinish <- true
	}()

	// Request page
	res, err := http.Get(url)
	if err != nil {
		chErrors <- BuildError("Request failed", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
    msg := fmt.Sprintf("Page broken. URL: %s", url)
    chErrors <- BuildError(msg, nil)
	}

	// Load HTML document to goquery
	doc, _ := goquery.NewDocumentFromReader(res.Body)

	// Find question blocks
	doc.Find(questionItemsSelection).Each(func(i int, testItem *goquery.Selection) {
		qtext := testItem.Find(".ask a").Text()
		question := Question{
			Text: strings.TrimSpace(qtext),
		}

		testItem.Find(answerItemsSelection).Each(func(i int, item *goquery.Selection) {
			answer := Answer{}
			answer.fillFromSelection(item)
			question.Answers = append(question.Answers, answer)
		})

		chQuestions <- question
	})
}

// Utils
func (qpage *QuestionsPage) parseBreadcrumb(i int, crumbItem *goquery.Selection) {
	var indexItem string

	if crumbItem.HasClass("active") {
		indexItem = crumbItem.Find("span span").Text()
	} else {
		indexItem = crumbItem.Find("a span").Text()
	}

	qpage.Index = append(qpage.Index, indexItem)
}

func (a *Answer) fillFromSelection(item *goquery.Selection) {
	var text string
	isCorrect := item.HasClass("alert-success")

	if isCorrect {
		text = item.Find("strong").Text()
	} else {
		text = item.Text()
	}

	a.Text = strings.TrimSpace(text)
	a.IsCorrect = isCorrect
}
