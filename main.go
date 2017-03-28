package main

import (
	"flag"
	"io"
	"os"
	"regexp"

	"strconv"

	"strings"

	"bufio"

	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/a-h/tpgit/git"
)

var repo = flag.String("repo", "https://github.com/a-h/ver", "The repo to query for TargetProcess ids in commit messages.")
var dryRun = flag.Bool("dryRun", true, "Set to true (default) to see what changes would be made.")

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()

	log.SetFormatter(&log.JSONFormatter{})
	logger := log.WithField("repo", *repo)

	fn := getFileNameFromRepo(*repo)
	lockfileName := "." + fn
	if _, err := os.Stat(lockfileName); !os.IsNotExist(err) {
		logger.Errorf("lock file %v is present for repo %v. not starting.\n", lockfileName, *repo)
		return -1
	}
	err := ioutil.WriteFile(lockfileName, []byte{}, 0666)
	defer os.Remove(lockfileName)

	previousHashes, err := loadHashes(fn)
	if !os.IsNotExist(err) {
		logger.Errorf("failed to load previous hashes from file %v with error: %v\n", fn, err)
		return -1
	}

	logger.Info("cloning repo")
	r, err := git.Clone(*repo)
	if err != nil {
		logger.Errorf("failed to clone the repo: %v\n", err)
		return -1
	}

	logger.Info("getting commit log from repo")
	commitlog, err := r.Log()
	if err != nil {
		logger.Errorf("failed to get the commit log: %v\n", err)
		return -1
	}

	for _, entry := range commitlog {
		entryLogger := logger.WithField("hash", entry.Hash)

		if _, ok := previousHashes[entry.Hash]; ok {
			entryLogger.Info("skipping already processed hash")
			continue
		}
		entryLogger = entryLogger.WithField("git_timestamp", entry.Timestamp)
		entryLogger = entryLogger.WithField("name", entry.Name)
		entryLogger = entryLogger.WithField("ids", extract(entry.Body))
		entryLogger = entryLogger.WithField("body", entry.Body)
		entryLogger.Info("processing commit")
	}

	return 0
}

func getFileNameFromRepo(repo string) string {
	return strings.NewReplacer("/", "-", ":", "-", ".", "-").Replace(repo)
}

func loadHashes(fileName string) (map[string]bool, error) {
	op := make(map[string]bool)

	file, err := os.Open(fileName)
	if err != nil {
		return op, err
	}
	defer file.Close()

	r := bufio.NewReader(file)

	var line string
	for {
		line, err = r.ReadString('\n')
		op[strings.TrimSpace(line)] = true
		if err != nil {
			break
		}
	}

	if err != io.EOF {
		return op, err
	}

	return op, nil
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
