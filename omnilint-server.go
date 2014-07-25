package main

import (
  "os"
  "github.com/go-martini/martini"
  "github.com/yvasiyarov/gorelic"
)

func main() {
  var (
    newrelic_license = os.Getenv("NEWRELIC_LICENSE")
    newrelic_name = os.Getenv("NEWRELIC_NAME")
  )
  if newrelic_license != "" && newrelic_name != "" {
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
