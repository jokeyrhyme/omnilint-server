package main

import (
  "bytes"
	"bufio"
	"net/http"
	"os"
	"os/exec"
	"github.com/go-martini/martini"
	"github.com/yvasiyarov/gorelic"
)

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
			return php(res, req)
		}
		return 400, "Content-Type not currently supported"
	})

	m.Run()
}

func php(res http.ResponseWriter, req *http.Request) (int, string) {
	path, err := exec.LookPath("php")
	if err != nil {
		return 500, err.Error()
	}
	cmd := exec.Command(path, "-l")
	cmd.Stdin = bufio.NewReader(req.Body)
	var (
    out bytes.Buffer
  )
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
//		return 500, err.Error()
		return 500, out.String()
	}
	return 200, out.String()
}
