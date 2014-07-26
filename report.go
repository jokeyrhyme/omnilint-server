package main

import (
  "encoding/json"
)

type ReportItem struct {
  Message string `json:"message"`
  Source string `json:"source"`
  Severity int `json:"severity"`
  Type string `json:"type"`
  Line int `json:"line"`
  Column int `json:"column"`
}

type ReportTotals struct {
  Errors int `json:"errors"`
  Warnings int `json:"warnings"`
}

type Report struct {
  Totals ReportTotals `json:"totals"`
  Errors []ReportItem `json:"errors"`
}

func (r *Report) New() {
  r.Totals.Errors = 0;
  r.Errors = make([]ReportItem, 0)
}

func (r *Report) AddItem(item ReportItem) {
  r.Totals.Errors += 1
  r.Errors = append(r.Errors, item)
}

func (r *Report) ToJson() string {
  output, _ := json.Marshal(r)
  return string(output[:]);
}
