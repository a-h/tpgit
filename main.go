package main

import "flag"
import "github.com/a-h/tpgit/git"
import "fmt"
import "os"

var repo = flag.String("repo", "https://github.com/a-h/ver", "The repo to query for TargetProcess Ids")
var dryRun = flag.Bool("dryRun", true, "Set to true (default) to see what changes would be made.")

func main() {
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
		fmt.Println(entry.Body)
	}
}
