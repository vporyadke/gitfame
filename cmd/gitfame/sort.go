package main

import (
	"fmt"
	"os"
	"sort"
)

type UserOutput struct {
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
	Name    string `json:"name"`
}

var sortBy string

var stats []UserOutput

func byLines(i, j int) bool {
	if stats[i].Lines != stats[j].Lines {
		return stats[i].Lines > stats[j].Lines
	}
	if stats[i].Commits != stats[j].Commits {
		return stats[i].Commits > stats[j].Commits
	}
	if stats[i].Files != stats[j].Files {
		return stats[i].Files > stats[j].Files
	}
	return stats[i].Name < stats[j].Name
}

func byCommits(i, j int) bool {
	if stats[i].Commits != stats[j].Commits {
		return stats[i].Commits > stats[j].Commits
	}
	if stats[i].Lines != stats[j].Lines {
		return stats[i].Lines > stats[j].Lines
	}
	if stats[i].Files != stats[j].Files {
		return stats[i].Files > stats[j].Files
	}
	return stats[i].Name < stats[j].Name
}

func byFiles(i, j int) bool {
	if stats[i].Files != stats[j].Files {
		return stats[i].Files > stats[j].Files
	}
	if stats[i].Lines != stats[j].Lines {
		return stats[i].Lines > stats[j].Lines
	}
	if stats[i].Commits != stats[j].Commits {
		return stats[i].Commits > stats[j].Commits
	}
	return stats[i].Name < stats[j].Name
}

func sortStats() {
	for name, data := range users {
		stats = append(stats, UserOutput{Name: name, Commits: len(data.commits), Files: len(data.files), Lines: data.lines})
	}
	switch sortBy {
	case "lines":
		sort.Slice(stats, byLines)
	case "commits":
		sort.Slice(stats, byCommits)
	case "files":
		sort.Slice(stats, byFiles)
	default:
		fmt.Printf("unknown sort option: %v\n", sortBy)
		os.Exit(1)
	}
}
