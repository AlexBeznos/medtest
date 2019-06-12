package main

import (
  "fmt"
  "log"
  "net/http"
  "encoding/json"
  "strings"

  "github.com/PuerkitoBio/goquery"
)

type Answer struct {
  Text string `json:"text"`
  IsCorrect bool `json:"isCorrect"`
}
type Question struct {
  Text string `json:"text"`
  Answers []Answer `json:"answers"`
}

func getAnswer(ansItem *goquery.Selection) Answer {
  var text string
  var isCorrect bool

  switch isCorrect := ansItem.HasClass("alert-success"); isCorrect {
    case true:
      text = ansItem.Find("strong").Text()
    case false:
      text = ansItem.Text()
  }

  return Answer{
    strings.TrimSpace(text),
    isCorrect,
  } 
}

func main() {
  rootUrl := "https://www.med-test.in.ua/uk/node/list/tests/html/10"
  var quiz []Question

  // Request page
  res, err := http.Get(rootUrl)

  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()

  if res.StatusCode != 200 {
    log.Fatalf("Status code error: %d, %s\n%s", res.StatusCode, res.Status, rootUrl)
  }

  // Load HTML document to goquery
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }

  // Find question blocks
  doc.Find("#tests-content .container .row .test .panel").Each(func (i int, testItem *goquery.Selection) {
    questionText := testItem.Find(".ask a").Text()
    var answers []Answer

    testItem.Find(".answer .list-group-item").Each(func (i int, ansItem *goquery.Selection) {
      answers = append(answers, getAnswer(ansItem))
    })

    quiz = append(quiz, Question{
      strings.TrimSpace(questionText),
      answers,
    })
  })

  marshalled, _ := json.Marshal(quiz)
  fmt.Println(string(marshalled))
}
