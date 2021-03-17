package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
)

type statWriter interface {
	write(entry UserOutput)
	flush()
}

type tabularWriter struct {
	*tabwriter.Writer
}

func newTabular() tabularWriter {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 1, 8, 1, ' ', 0)
	fmt.Fprintln(w, "Name\tLines\tCommits\tFiles")
	return tabularWriter{w}
}

func (w tabularWriter) write(entry UserOutput) {
	fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", entry.Name, entry.Lines, entry.Commits, entry.Files)
}

func (w tabularWriter) flush() {
	w.Flush()
}

type csvWriter struct {
	*csv.Writer
}

func newCSV() csvWriter {
	w := csv.NewWriter(os.Stdout)
	headers := []string{"Name", "Lines", "Commits", "Files"}
	_ = w.Write(headers)
	return csvWriter{w}
}

func (w csvWriter) write(entry UserOutput) {
	_ = w.Write([]string{entry.Name, strconv.Itoa(entry.Lines), strconv.Itoa(entry.Commits), strconv.Itoa(entry.Files)})
}

func (w csvWriter) flush() {
	w.Flush()
}

type jsonWriter struct{}

func (jsonWriter) write(entry UserOutput) {
	s, _ := json.Marshal(entry)
	fmt.Println(string(s))
}

func (jsonWriter) flush() {}

func printStats() {
	var w statWriter
	switch format {
	case "tabular":
		w = newTabular()
	case "csv":
		w = newCSV()
	case "json":
		s, _ := json.Marshal(stats)
		fmt.Println(string(s))
		return
	case "json-lines":
		w = jsonWriter{}
	default:
		fmt.Printf("unknown format option:%v\n", format)
	}
	for _, entry := range stats {
		w.write(entry)
	}
	w.flush()
}
