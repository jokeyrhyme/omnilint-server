package main

import (
  "io/ioutil"
	"net/http"
	"os"
	"github.com/go-martini/martini"
	"github.com/yvasiyarov/gorelic"
)

type JsonSerializable interface {
  ToJson() string
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
    body, err := ioutil.ReadAll(req.Body);
    if err != nil {
      return 400, "Invalid request body"
    }
		if contentType == "application/x-php" {
			return CheckPhp(res, string(body[:]))
		}
		return 415, "Content-Type not currently supported"
	})

	m.Run()
}
