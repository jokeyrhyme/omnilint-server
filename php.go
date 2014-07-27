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

type codeSnifferReport struct {
	Totals reportTotals `json:"totals"`
	Files  struct {
		STDIN struct {
			Errors   int          `json:"errors"`
			Warnings int          `json:"warnings"`
			Messages []reportItem `json:"messages"`
		} `json:"STDIN"`
	} `json:"files"`
}

func parseSyntax(code string) ([]reportItem, error) {
	// borrowing liberally from:
	// https://github.com/AtomLinter/linter-php/blob/master/lib/linter-php.coffee
	path, err := exec.LookPath("php")
	if err != nil {
		return make([]reportItem, 0), err
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
			results       []reportItem
			result        reportItem
		)
		re := regexp.MustCompile(`(Parse|Fatal) (?P<error>error):(\s*(?P<type>parse|syntax) error,?)?\s*(?P<message>(unexpected '(?P<near>[^']+)')?.*) in .*? on line (?P<line>\d+)`)
		matchingLines = re.FindAllStringSubmatch(out.String(), -1)
		results = make([]reportItem, 0)
		for _, match := range matchingLines {
			result = reportItem{}
			result.Source = match[1]  // Parse|Fatal
			result.Message = match[5] // <message>
			if match[8] != "" {
				result.Line, _ = strconv.Atoi(match[8]) // <line>
			}
			results = append(results, result)
		}
		return results, nil
	}
	return make([]reportItem, 0), nil
}

func codeSniffer(code string, options map[string]string) ([]reportItem, error) {
	path, err := exec.LookPath("phpcs")
	if err != nil {
		return make([]reportItem, 0), err
	}
	if options["phpcs.standard"] != "" {
		re := regexp.MustCompile(`\W`)
		options["phpcs.standard"] = re.ReplaceAllLiteralString(options["phpcs.standard"], "")
	} else {
		options["phpcs.standard"] = "PSR2"
	}
	cmd := exec.Command(path, "--report=json", "--standard="+options["phpcs.standard"])
	cmd.Stdin = strings.NewReader(code)
	var (
		out bytes.Buffer
	)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		var (
			report codeSnifferReport
		)
		err = json.Unmarshal([]byte(out.String()), &report)
		if err != nil {
			return make([]reportItem, 0), err
		}
		return report.Files.STDIN.Messages, nil
	}
	return make([]reportItem, 0), nil
}

func checkPhp(res http.ResponseWriter, code string, params map[string]string) (int, string) {
	var (
		someResults, results []reportItem
		err                  error
	)
	results = make([]reportItem, 0)

	someResults, err = parseSyntax(code)
	if err != nil {
		return 500, err.Error()
	}
	results = append(results, someResults...)

	someResults, err = codeSniffer(code, params)
	if err != nil {
		return 500, err.Error()
	}
	results = append(results, someResults...)

	if len(results) > 0 {
		res.Header().Set("Content-Type", "application/json")
		report := new(report)
		for _, item := range results {
			report.AddItem(item)
		}
		var body []byte
		body, err = json.Marshal(report)
		if err != nil {
			return 500, err.Error()
		}
		return 200, string(body[:])
	}
	return 204, ""
}
