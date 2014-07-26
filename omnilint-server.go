package main

import (
	"bufio"
  "encoding/json"
	"net/http"
	"os"
	"github.com/go-martini/martini"
	"github.com/yvasiyarov/gorelic"
)

type JsonSerializable interface {
  ToJson() string
}

func CheckPhp(res http.ResponseWriter, req *http.Request) (int, string) {
  // result, err := ParseSyntax(bufio.NewReader(req.Body))
  result, err := CodeSniffer(bufio.NewReader(req.Body))
	if err != nil {
		return 500, err.Error()
	}
  if len(result) > 0 {
    res.Header().Set("Content-Type", "application/json")
    report := new(Report)
    for _, item := range result {
      report.AddItem(item)
    }
    var body []byte
    body, err = json.Marshal(report)
    if err != nil {
      return 500, err.Error()
    }
    return 200, string(body[:]);
  }
	return 204, ""
}

func main() {
	var (
		newrelicLicense = os.Getenv("NEWRELIC_LICENSE")
		newrelicName    = os.Getenv("NEWRELIC_NAME")
	)
	if newrelicLicense != "" && newrelicName != "" {
		agent := gorelic.NewAgent()
		agent.Verbose = true
		agent.NewrelicLicense = os.Getenv("NEWRELIC_LICENSE")
		agent.NewrelicName = os.Getenv("NEWRELIC_NAME")
		agent.Run()
	}

	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
	})

	m.Post("**", func(res http.ResponseWriter, req *http.Request) (int, string) {
		var (
			contentTypeHeader []string
			contentType       string
		)
		contentTypeHeader = req.Header["Content-Type"]
		if len(contentTypeHeader) > 0 {
			contentType = contentTypeHeader[0]
		} else {
			return 400, "Content-Type header is mandatory"
		}
		if contentType == "application/x-php" {
			return CheckPhp(res, req)
		}
		return 415, "Content-Type not currently supported"
	})

	m.Run()
}
