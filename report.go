package main

type reportItem struct {
	Message  string `json:"message"`
	Source   string `json:"source"`
	Severity int    `json:"severity"`
	Type     string `json:"type"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
}

func (i *reportItem) New() {
	i.Severity, i.Line, i.Column = 0, 0, 0
	i.Message, i.Source = "", ""
	i.Type = "ERROR"
}

type reportTotals struct {
	Errors   int `json:"errors"`
	Warnings int `json:"warnings"`
}

type report struct {
	Totals reportTotals `json:"totals"`
	Errors []reportItem `json:"errors"`
}

func (r *report) New() {
	r.Totals.Errors = 0
	r.Errors = make([]reportItem, 0)
}

func (r *report) AddItem(item reportItem) {
	r.Totals.Errors++
	r.Errors = append(r.Errors, item)
}
