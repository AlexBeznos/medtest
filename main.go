package main

import (
  "encoding/json"
  "fmt"
  "log"
  "io/ioutil"

  "github.com/AlexBeznos/medtest/parser"
)

type FullPack struct {
  Tests []parser.QuestionsPage `json:"tests"`
}

func main() {

  conf := parser.Config{
    RootUrl: "https://www.med-test.in.ua",
    SitemapPath: "/uk/site-map",
  }
  sitemap := parser.Sitemap{}
  sitemap.Parse(&conf)

  pack := FullPack{}
  for _, path := range sitemap.Paths {
    fmt.Printf("Path: %s\n", path)

    fmt.Println("Parsing started...")

    qpage := parser.QuestionsPage{
      Path: path,
    }
    qpage.Parse(&conf)

    fmt.Println("Parsing finished!")

    pack.Tests = append(pack.Tests, qpage)
  }

  fmt.Println("Marshaling...")
  marshalled, _ := json.Marshal(pack)
  fmt.Println("Writing to data.json")
  err := ioutil.WriteFile("data.json", []byte(string(marshalled)), 0644)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Finish.")
}
