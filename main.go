package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"

	"strconv"

	"strings"

	"bufio"

	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/a-h/tpgit/git"
	"github.com/a-h/tpgit/targetprocess"
)

var repo = flag.String("repo", "https://github.com/a-h/ver", "The repo to query for TargetProcess ids in commit messages.")
var dryRun = flag.Bool("dryRun", true, "Set to true (default) to see what changes would be made.")
var url = flag.String("url", "", "Set to the root address of your TargetProcess account, e.g. https://example.tpondemand.com")
var username = flag.String("username", "", "Sets the username to use to authenticate against TargetProcess.")
var password = flag.String("password", "", "Sets the password to use to authenticate against TargetProcess.")
var maximumToAdd = flag.Int("max", 1, "Sets the maximum number of commits that the system will do in one run.")

func main() {
	exitCode := run()
	defer os.Exit(exitCode)
}

func run() int {
	flag.Parse()

	log.SetFormatter(&log.JSONFormatter{})
	logger := log.WithField("repo", *repo)

	if *url == "" {
		logger.Errorf("please provide a TargetProcess URL to access")
		return -1
	}

	if *username == "" || *password == "" {
		logger.Errorf("please provide TargetProcess authentication details via the -username and -password command line args")
		return -1
	}

	dir, err := os.Getwd()
	if err != nil {
		logger.Errorf("couldn't get working directory: %v", err)
	}
	fn := path.Join(dir, getFileNameFromRepo(*repo))
	logger.Debug("hash filename is %v", fn)

	lockfileName := fn + ".lock"
	if _, err := os.Stat(lockfileName); !os.IsNotExist(err) {
		logger.Errorf("lock file %v is present for repo %v. not starting.\n", lockfileName, *repo)
		return -1
	}
	err = ioutil.WriteFile(lockfileName, []byte("lock"), 0666)
	if err != nil {
		logger.Errorf("error creating lock file %v: %v", lockfileName, err)
		return -1
	}
	defer os.Remove(lockfileName)

	hashes, err := loadHashes(fn)
	if !os.IsNotExist(err) {
		logger.Errorf("failed to load previous hashes from file %v: %v\n", fn, err)
		return -1
	}

	logger.Info("getting commit log from repo")
	commitlog, err := git.Log(*repo)
	if err != nil {
		logger.Errorf("failed to get the commit log: %v\n", err)
		return -1
	}

	commentsCreated := 0
	for _, entry := range commitlog {
		entryLogger := logger.WithField("hash", entry.Hash)

		if _, ok := hashes[entry.Hash]; ok {
			entryLogger.Info("skipping already processed hash")
			continue
		}
		entryLogger = entryLogger.WithField("git_timestamp", entry.Timestamp)
		entryLogger = entryLogger.WithField("name", entry.Name)
		ids := extract(entry.Body)
		entryLogger = entryLogger.WithField("ids", ids)
		entryLogger = entryLogger.WithField("body", entry.Body)
		entryLogger.Info("processing commit")

		msg := fmt.Sprintf("Referenced in commit %v (%v) by %v:\n\n%s",
			getCommitURL(*repo, entry.Hash),
			entry.Date(),
			entry.Email,
			entry.Body)

		entryLogger.Info("adding comment to target process")

		tp := targetprocess.NewAPI(*url, *username, *password)

		if !*dryRun {
			err = addCommentToTargetProcess(tp, ids, msg)
			commentsCreated++
			if err != nil {
				entryLogger.Errorf("failed to write comment: %v", err)
			}
			entryLogger.Infof("written %d comments to TargetProces")
			if commentsCreated > *maximumToAdd {
				entryLogger.Infof("exceeded maximum of %d comments, not doing any more")
				continue
			}
		}

		hashes[entry.Hash] = true
	}

	logger.Infof("writing %d hashes to %v", len(hashes), fn)
	err = saveHashes(fn, hashes)
	if err != nil {
		logger.Errorf("failed to write hashes to file: %v", err)
		return -1
	}

	logger.Infof("complete")

	return 0
}

func addCommentToTargetProcess(tp targetprocess.API, ids []int, msg string) error {
	for _, id := range ids {
		_, err := tp.Comment(id, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func getCommitURL(repo string, hash string) string {
	repo = strings.TrimSuffix(repo, ".git")
	return repo + "/commits/" + hash
}

func getFileNameFromRepo(repo string) string {
	return strings.NewReplacer("/", "", ":", "", ".", "-").Replace(repo)
}

func saveHashes(fileName string, hashes map[string]bool) error {
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		return err
	}

	for k := range hashes {
		_, err := file.WriteString(k + "\n")
		if err != nil {
			return fmt.Errorf("failed to write hash to file: %v", err)
		}
	}
	return file.Sync()
}

func loadHashes(fileName string) (map[string]bool, error) {
	op := make(map[string]bool)

	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		return op, err
	}

	r := bufio.NewReader(file)

	var line string
	for {
		line, err = r.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line != "" {
			op[line] = true
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
