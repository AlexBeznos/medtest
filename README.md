# MEDTEST

Typical usage of `medtest` package should be like this:

```golang
// main.go
// Initialize config
conf := medtest.Config{
  RootUrl: "https://www.med-test.in.ua",
  SitemapPath: "/uk/site-map",
}

// Initialize sitemap struct
sitemap := medtest.Sitemap{}
err := sitemap.Parse(&conf) // parse sitemap
if err != nil {
  log.Fatal(err)
}

// Initialize some kind of struct | slice | array which will contain all questions
pack := FullPack{}

// Result of sitemap parsing will be stored inside `Paths` field.
for _, path := range sitemap.Paths { // Here we iterating through all paths
  fmt.Printf("Path: %s\n", path)

  fmt.Println("Parsing started...")

  // Initialize `QuestionsPage` struct with path which should be parsed
  qpage := medtest.QuestionsPage{
    Path: path,
  }
  
  // parse page itself
  err := qpage.Parse(&conf)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Parsing finished!")

  // push parsed page into agregation method
  pack.Tests = append(pack.Tests, qpage)
}

fmt.Println("Marshaling...")

// Convert struct into json
marshalled, _ := json.Marshal(pack)
fmt.Println("Writing to data.json")
// write json into file
err := ioutil.WriteFile("data.json", []byte(string(marshalled)), 0644)
if err != nil {
  log.Fatal(err)
}
fmt.Println("Finish.")
```
