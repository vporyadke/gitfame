package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type User struct {
	commits, files map[string]struct{}
	lines          int
}

func newUser() User {
	commits := make(map[string]struct{})
	files := make(map[string]struct{})
	return User{commits, files, 0}
}

var users = make(map[string]User)

var author = "author"

func calcStats(filePath string) {
	blame := exec.Command("git", "blame", "--porcelain", revision, "--", filePath)
	blame.Dir = path
	out, err := blame.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(out))
		os.Exit(1)
	}
	if len(out) == 0 {
		// Можно было бы два раза позвать лог, или аккуратно читать хеш, но поверим, что не встертим @@@@
		log := exec.Command("git", "log", "--format=%H@@@@%an", revision, "--", filePath)
		log.Dir = path
		logOut, err := log.CombinedOutput()
		if err != nil {
			fmt.Println(err)
			fmt.Println(string(out))
			os.Exit(1)
		}
		last := strings.Split(strings.Split(string(logOut), "\n")[0], "@@@@")
		commit, username := last[0], last[1]
		var user User
		var ok bool
		if user, ok = users[username]; !ok {
			user = newUser()
		}
		user.files[filePath] = struct{}{}
		user.commits[commit] = struct{}{}
		users[username] = user
	}
	blameLines := strings.Split(string(out), "\n")
	var curCommit string
	var curLines int
	commits := make(map[string]string)
	for i := 0; i < len(blameLines); i++ {
		lineInfo := strings.Fields(blameLines[i])
		if len(lineInfo) == 4 {
			curCommit = lineInfo[0]
			curLines, _ = strconv.Atoi(lineInfo[3])
			if username, ok := commits[curCommit]; ok {
				user := users[username]
				user.lines += curLines
				users[username] = user
			} else {
			parse_fields:
				for {
					i++
					field := strings.Fields(blameLines[i])[0]
					switch field {
					case author:
						username := blameLines[i][len(field)+1:]
						var user User
						var ok bool
						if user, ok = users[username]; !ok {
							user = newUser()
						}
						user.commits[curCommit] = struct{}{}
						user.files[filePath] = struct{}{}
						user.lines += curLines
						users[username] = user
						commits[curCommit] = username
					case "filename":
						break parse_fields
					}
				}
			}
		}
		i++
	}
}
