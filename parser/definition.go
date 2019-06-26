package parser

const RootUrl = "https://www.med-test.in.ua"

type Test struct {
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

