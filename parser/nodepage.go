package parser

import (
	"log"
  "strconv"
  "strings"
	"net/http"

  "github.com/PuerkitoBio/goquery"
)

type QuestionsPage struct {
  Path string `json:"path"`
  Questions []Question `json:"questions"`
  Index []string `json:"index"`
  NumberOfPages int `json:"numberOfPages"`
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

// Public
func (qpage *QuestionsPage) Parse(conf *Config) {
  qpage.getMetaData(conf)

  chQuestions := make(chan Question)
  chFinished := make(chan bool)

  // Kick off the parsing
  for i := 1; i <= qpage.NumberOfPages; i++ {
    url := conf.CombineUrl(qpage.Path, i)
    go parseQuestions(url, chQuestions, chFinished)
  }

  // Subscription to parsed questions
  for c := 0; c < qpage.NumberOfPages; {
		select {
		case question := <-chQuestions:
			qpage.Questions = append(qpage.Questions, question)
		case <-chFinished:
			c++
		}
	}
}

// Private
func (qpage *QuestionsPage) getMetaData(conf *Config) {
  url := conf.CombineUrl(qpage.Path, 1)

  // Request page
  res, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()

  if res.StatusCode != 200 {
    log.Fatalf("Status code error: %d, %s\n%s", res.StatusCode, res.Status, url)
  }

  // Load HTML document to goquery
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }

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
}

func parseQuestions(url string, chQuestions chan Question, chFinish chan bool) {
  defer func() {
    chFinish <- true
  }()

  // Request page
  res, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }

  defer res.Body.Close()

  if res.StatusCode != 200 {
    log.Fatalf("Status code error: %d, %s\n%s", res.StatusCode, res.Status, url)
  }

  // Load HTML document to goquery
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }

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
func (qpage *QuestionsPage) parseBreadcrumb (i int, crumbItem *goquery.Selection) {
  var indexItem string

  if crumbItem.HasClass("active") {
    indexItem = crumbItem.Find("span span").Text()
  } else {
    indexItem = crumbItem.Find("a span").Text()
  }

  qpage.Index = append(qpage.Index, indexItem)
}

func(a *Answer) fillFromSelection(item *goquery.Selection) {
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
