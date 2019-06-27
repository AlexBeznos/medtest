package parser

import (
	"log"
  "strconv"
  "strings"
	"net/url"
	"net/http"

  "github.com/PuerkitoBio/goquery"
)

const lastPageLinkSelector = ".navigation .pagination .last a"
const breadcrumbsSelector = "ul.breadcrumb li"
const questionItemsSelection = "#tests-content .container .row .test .panel"
const answerItemsSelection = ".answer .list-group-item"

// Public
func (t *Test) Parse() {
  t.getMetaData()
  chQuestions := make(chan Question)
  chFinished := make(chan bool)

  // Kick off the parsing
  for i := 1; i <= t.NumberOfPages; i++ {
    go parseQuestions(i, t.Path, chQuestions, chFinished)
  }

  // Subscription to parsed questions
  for c := 0; c < t.NumberOfPages; {
		select {
		case question := <-chQuestions:
			t.Questions = append(t.Questions, question)
		case <-chFinished:
			c++
		}
	}
}

// Private
func (t *Test) getMetaData() {
  url := prepareUrl(1, t.Path)

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
    numberOfPages := getNumberOfPagesFromUrl(href)
    t.NumberOfPages = numberOfPages
  } else {
    t.NumberOfPages = 1
  }

  // Get breadcrumbs
  doc.Find(breadcrumbsSelector).Each(t.parseBreadcrumb)
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


// Utils
func (t *Test) parseBreadcrumb (i int, crumbItem *goquery.Selection) {
  var indexItem string

  if crumbItem.HasClass("active") {
    indexItem = crumbItem.Find("span span").Text()
  } else {
    indexItem = crumbItem.Find("a span").Text()
  }

  t.Index = append(t.Index, indexItem)
}

func parseQuestions(pageNumber int, path string, chQuestions chan Question, chFinish chan bool) {
  defer func() {
    chFinish <- true
  }()
  url := prepareUrl(pageNumber, path)

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

func prepareUrl(pageNumber int, path string) string {
  u, err := url.Parse(RootUrl)
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

func getNumberOfPagesFromUrl(href string) int {
  u, err := url.Parse(RootUrl)
  if err != nil {
    log.Fatal(err)
  }

  fullUrl, err := u.Parse(href)
  if err != nil {
    log.Fatal(err)
  }

  query := fullUrl.Query()
  pageNumber := query.Get("page")
  num, err := strconv.Atoi(pageNumber)
  if err != nil {
    log.Fatal(err)
  }

  return num
}
