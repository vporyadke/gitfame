// +build !solution

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

var cmd = &cobra.Command{
	Use:   "gitfame",
	Short: "",
	Long:  "Compute stats on user activity in a git repo",
	Run:   work,
}

var path, format, revision string
var committer bool
var extensionsFlag, languages, exclude, restrict []string
var checkExtensions int
var extensions map[string]int = make(map[string]int)

func init() {
	cmd.Flags().StringVar(&path, "repository", ".", "path to repo")
	cmd.Flags().StringVar(&format, "format", "tabular", "output format")
	cmd.Flags().StringVar(&revision, "revision", "HEAD", "revision")
	cmd.Flags().StringVar(&sortBy, "order-by", "lines", "sort key")
	cmd.Flags().BoolVar(&committer, "use-committer", false, "attribute files to committer")
	cmd.Flags().StringSliceVar(&extensionsFlag, "extensions", []string{}, "only count files with these extensions")
	cmd.Flags().StringSliceVar(&languages, "languages", []string{}, "only count these languages")
	cmd.Flags().StringSliceVar(&exclude, "exclude", []string{}, "exclude these patterns")
	cmd.Flags().StringSliceVar(&restrict, "restrict-to", []string{}, "restrict to these patterns")
}

// test

func prepareConfig() {
	if committer {
		author = "committer"
	}
	if len(extensionsFlag) > 0 {
		checkExtensions |= 1
		for _, e := range extensionsFlag {
			extensions[e] |= 1
		}
	}
	if len(languages) > 0 {
		checkExtensions |= 2
		for _, l := range languages {
			for _, e := range languageExtensions[l] {
				extensions[e] |= 2
			}
		}
	}
}

func work(cmd *cobra.Command, args []string) {
	prepareConfig()
	ls := exec.Command("git", "ls-tree", revision, "-r")
	ls.Dir = path
	out, err := ls.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(out))
		os.Exit(1)
	}
	rawList := strings.Split(string(out), "\n")
	for _, s := range rawList {
		if len(s) == 0 {
			continue
		}
		spaceCount := 0
		var pathStart int
		for i, c := range s {
			if unicode.IsSpace(c) {
				spaceCount++
			} else if spaceCount >= 3 {
				pathStart = i
				break
			}
		}
		filename := s[pathStart:]
		if check(filename) {
			calcStats(filename)
		}
	}
	sortStats()
	printStats()
}

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
