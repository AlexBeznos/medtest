package main

import (
  "encoding/json"
  "fmt"
  "log"
  "io/ioutil"

  "github.com/AlexBeznos/med-test-parser/parser"
)

func main() {
  test := parser.Test{
    Path: "/uk/node/list/tests/html/10",
  }
  fmt.Println("Parsing started...")
  test.Parse()
  fmt.Println("Parsing finished!")

  // Temp solution
  fmt.Println("Marshaling...")
  marshalled, _ := json.Marshal(test)
  fmt.Println("Writing to data.json")
  err := ioutil.WriteFile("data.json", []byte(string(marshalled)), 0644)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Finish.")
}
