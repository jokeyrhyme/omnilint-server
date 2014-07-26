package main

import (
  "bytes"
  "encoding/json"
  "net/http"
  "os/exec"
  "regexp"
  "strconv"
  "strings"
)

type CodeSnifferReport struct {
  Totals ReportTotals `json:"totals"`
  Files struct {
    STDIN struct {
      Errors int `json:"errors"`
      Warnings int `json:"warnings"`
      Messages []ReportItem `json:"messages"`
    } `json:"STDIN"`
  } `json:"files"`
}

func (r *CodeSnifferReport) ToJson() string {
  output, _ := json.Marshal(r)
  return string(output[:]);
}

func ParseSyntax(code string) ([]ReportItem, error) {
  // borrowing liberally from:
  // https://github.com/AtomLinter/linter-php/blob/master/lib/linter-php.coffee
	path, err := exec.LookPath("php")
	if err != nil {
		return make([]ReportItem, 0), err
	}
	cmd := exec.Command(path, "-l", "-n", "-d display_errors=On", "-d log_errors=Off")
	cmd.Stdin = strings.NewReader(code)
	var (
    out bytes.Buffer
  )
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
    var (
      matchingLines [][]string
      results []ReportItem
      result ReportItem
    )
    re := regexp.MustCompile(`(Parse|Fatal) (?P<error>error):(\s*(?P<type>parse|syntax) error,?)?\s*(?P<message>(unexpected '(?P<near>[^']+)')?.*) in .*? on line (?P<line>\d+)`);
    matchingLines = re.FindAllStringSubmatch(out.String(), -1);
    results = make([]ReportItem, 0)
    for _, match := range matchingLines {
      result = ReportItem{}
      result.Source = match[1] // Parse|Fatal
      result.Message = match[5] // <message>
      if match[8] != "" {
        result.Line, _ = strconv.Atoi(match[8]) // <line>
      }
      results = append(results, result)
    }
		return results, nil
	}
	return make([]ReportItem, 0), nil
}

func CodeSniffer(code string) ([]ReportItem, error) {
	path, err := exec.LookPath("phpcs")
	if err != nil {
		return make([]ReportItem, 0), err
	}
	cmd := exec.Command(path, "--report=json")
	cmd.Stdin = strings.NewReader(code)
	var (
    out bytes.Buffer
  )
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
    var (
      report CodeSnifferReport
    )
    err = json.Unmarshal([]byte(out.String()), &report)
    if err != nil {
      return make([]ReportItem, 0), err
    }
    return report.Files.STDIN.Messages, nil
	}
	return make([]ReportItem, 0), nil
}

func CheckPhp(res http.ResponseWriter, code string) (int, string) {
  var (
    someResults, results []ReportItem
    err error
  )
  results = make([]ReportItem, 0)

  someResults, err = ParseSyntax(code)
	if err != nil {
		return 500, err.Error()
	}
  results = append(results, someResults...);

  someResults, err = CodeSniffer(code)
	if err != nil {
		return 500, err.Error()
	}
  results = append(results, someResults...);

  if len(results) > 0 {
    res.Header().Set("Content-Type", "application/json")
    report := new(Report)
    for _, item := range results {
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
