package main

import (
  "bytes"
  "encoding/json"
  "io"
  "os/exec"
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

func ParseSyntax(code io.Reader) (string, error) {
	path, err := exec.LookPath("php")
	if err != nil {
		return "", err
	}
	cmd := exec.Command(path, "-l")
	cmd.Stdin = code
	var (
    out bytes.Buffer
  )
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		return out.String(), nil
	}
	return "", nil
}

func CodeSniffer(code io.Reader) ([]ReportItem, error) {
	path, err := exec.LookPath("phpcs")
	if err != nil {
		return make([]ReportItem, 0), err
	}
	cmd := exec.Command(path, "--report=json")
	cmd.Stdin = code
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
