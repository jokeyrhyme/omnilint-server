package main

import (
	"io/ioutil"
  "log"
	"net/http"
	"os"
  "strings"
  "reflect"
	"github.com/go-martini/martini"
  "github.com/martini-contrib/cors"
	"github.com/yvasiyarov/gorelic"
)

func toMapStringString(m map[string][]string) map[string]string {
	out := make(map[string]string)
	for k, v := range m {
		if len(v) > 0 {
			out[k] = v[0]
		}
	}
	return out
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

  var logger *log.Logger
  logger = m.Injector.Get(reflect.TypeOf(logger)).Interface().(*log.Logger)

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
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return 400, "Invalid request body"
		}
		if contentType == "application/x-php" {
			qsa := toMapStringString(req.URL.Query())
			return checkPhp(res, string(body[:]), qsa)
		}
		return 415, "Content-Type not currently supported"
	})

  var corsOrigins = os.Getenv("CORS_ORIGINS")
  if corsOrigins != "" {
    logger.Println("activating CORS: " + corsOrigins)
    m.Use(cors.Allow(&cors.Options{
      AllowOrigins: strings.Split(corsOrigins, ","),
      AllowMethods: []string{"GET", "POST"},
      AllowHeaders: []string{"Origin", "Content-Type"},
    }))
  }

	m.Run()
}
