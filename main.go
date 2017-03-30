package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"strconv"

	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/a-h/tpgit/backend"
	"github.com/a-h/tpgit/git"
	"github.com/a-h/tpgit/targetprocess"
)

var repoFlag = flag.String("repo", "", "The directory containing a git repo to query for TargetProcess ids in commit messages.")
var repoURLFlag = flag.String("repoURL", "", "The endpoint of the repo, used to construct the URL to the commits in TargetProcess e.g. https://bitbucket.com/org/repo/commits/ - the message will add the git hash to end of URL.")
var dryRun = flag.Bool("dryrun", true, "Set to true (default) to see what changes would be made.")
var url = flag.String("url", "", "Set to the root address of your TargetProcess account, e.g. https://example.tpondemand.com")
var username = flag.String("username", "", "Sets the username to use to authenticate against TargetProcess.")
var password = flag.String("password", "", "Sets the password to use to authenticate against TargetProcess.")
var maximumToAdd = flag.Int("max", 1, "Sets the maximum number of commits that the system will do in one run.")
var logFormat = flag.String("logformat", "json", "Set to json for JSON, or console for console friendly formatting.")
var quiet = flag.Bool("quiet", false, "Reduces log output.")

var backendFlag = flag.String("backend", "localfile", "Sets the backend to use to store the status of git entries.")

// Local File backened settings.
var localFileLocationFlag = flag.String("hashfile", "", "The name of the file to use to store the hashes, e.g. 'projectname.hashes'")

func main() {
	exitCode := run()
	defer os.Exit(exitCode)
}

func run() int {
	flag.Parse()

	if *logFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}
	logger := log.WithField("repo", *repoFlag)

	repo := *repoFlag
	if repo == "" {
		repo, _ = os.Getwd()
	}
	if repo == "" {
		logger.Errorf("repo flag missing")
		return -1
	}

	if *repoURLFlag == "" {
		logger.Errorf("repoURL flag missing")
		return -1
	}

	if *url == "" {
		logger.Errorf("url flag missing")
		return -1
	}

	if *username == "" || *password == "" {
		logger.Errorf("username or password flag missing")
		return -1
	}

	logger.Info("getting commit log from repo")
	commits, err := git.Log(repo)
	if err != nil {
		logger.Errorf("failed to get the commit log: %v\n", err)
		return -1
	}

	be, err := getBackend()
	if err != nil {
		logger.Errorf("failed to configure backend: %v", err)
		return -1
	}

	tp := targetprocess.NewAPI(*url, *username, *password)

	err = processCommmits(logger, commits, be, tp, *repoURLFlag)
	if err != nil {
		return -1
	}
	return 0
}

func getBackend() (Backend, error) {
	switch *backendFlag {
	case "localfile":
		filename := *localFileLocationFlag
		if filename == "" {
			return nil, errors.New("localfile backend: hashfile flag not set")
		}
		return backend.NewLocalFile(filename)
	}
	return nil, errors.New("backend not recognised")
}

type commenter interface {
	Comment(entityID int, message string) error
}

func processCommmits(logger *log.Entry, commits []git.Commit, be Backend, commenter commenter, commitURL string) error {
	commentsCreated := 0
	for _, entry := range commits {
		entryLogger := logger.WithField("hash", entry.Hash)

		processed, err := be.IsProcessed(entry.Hash)
		if err != nil {
			return err
		}
		if processed && !*quiet {
			entryLogger.Info("skipping already processed hash")
			continue
		}
		entryLogger = entryLogger.WithField("git_timestamp", entry.Timestamp)
		entryLogger = entryLogger.WithField("name", entry.Name)
		ids := extract(entry.Body)
		entryLogger = entryLogger.WithField("ids", ids)
		entryLogger = entryLogger.WithField("body", entry.Body)

		if !*quiet {
			entryLogger.Info("processing commit")
		}

		msg := fmt.Sprintf("Referenced in commit %v (%v) by %v:\n\n%s",
			commitURL+entry.Hash,
			entry.Date(),
			entry.Email,
			entry.Body)

		if !*dryRun && len(ids) > 0 {
			if commentsCreated > *maximumToAdd {
				if !*quiet {
					entryLogger.Infof("exceeded maximum of %d comments, not doing any more", *maximumToAdd)
				}
				continue
			}
			entryLogger.Info("adding comment to target process")
			err = addComments(commenter, ids, msg)
			commentsCreated += len(ids)
			if err != nil {
				entryLogger.Errorf("failed to write comment: %v", err)
			}
			entryLogger.Infof("written %d comments to TargetProces", commentsCreated)
		}

		if err = be.MarkProcessed(entry.Hash); err != nil {
			return err
		}
	}

	logger.Infof("cancelling lease")
	err := be.CancelLease()
	if err != nil {
		return fmt.Errorf("failed to cancel lease: %v", err)
	}

	logger.Infof("complete")
	return nil
}

func addComments(commenter commenter, ids []int, msg string) error {
	for _, id := range ids {
		err := commenter.Comment(id, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

var re = regexp.MustCompile(`(?i)(?:(?:TP)|(?:TP[\-\s\:]+)|(?:^#))(?P<id>\d+)`)

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
