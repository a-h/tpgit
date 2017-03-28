package main

import (
	"flag"
	"os"
	"regexp"

	"strconv"

	"fmt"

	"github.com/a-h/tpgit/git"
)

var repo = flag.String("repo", "https://github.com/a-h/ver", "The repo to query for TargetProcess ids in commit messages.")
var dryRun = flag.Bool("dryRun", true, "Set to true (default) to see what changes would be made.")

func main() {
	flag.Parse()

	r, err := git.Clone(*repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to clone the repo: %v\n", err)
		os.Exit(-1)
	}

	log, err := r.Log()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get the log: %v\n", err)
		os.Exit(-1)
	}

	for _, entry := range log {
		fmt.Printf("%d - %s\n", extract(entry.Body), entry.Body)
		fmt.Println("----------------------------------------------")
	}
}

var re = regexp.MustCompile(`(?i)(?:(?:TP)?(?:\-|\s+|\:|^#))(?P<id>\d+)`)

func extract(message string) []int {
	ids := []int{}

	for _, m := range re.FindAllStringSubmatch(message, -1) {
		if len(m) > 1 {
			sm := m[1] // The first captured group.
			id, err := strconv.Atoi(sm)
			if err != nil {
				continue
			}
			if !contains(ids, id) {
				ids = append(ids, id)
			}
		}
	}
	return ids
}

func contains(values []int, value int) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}
