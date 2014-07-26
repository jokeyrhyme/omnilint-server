package main

import (
  "os"
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
  m.Run()
}
