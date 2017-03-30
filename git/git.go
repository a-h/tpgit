package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Commit is the data stored within a git log output.
type Commit struct {
	Hash string `json:"hash"`
	Body string `json:"body"`
	// Name is the author name.
	Name      string `json:"name"`
	Email     string `json:"email"`
	Timestamp int64  `json:"timestamp"`
}

// Date Converts the Unix timestamp into a Time.
func (c Commit) Date() time.Time {
	return time.Unix(c.Timestamp, 0)
}

var lineSeparator = ":7e7dd4cbeda4c5f65b46e9d55ac526f63fa9a7c9:\n"
var separator = ":ec0c7bc17e1ef95b57f47e6ee9f63f54ac187325:"

// Log gets the git log of the repository.
func Log(directory string) ([]Commit, error) {
	logfmt := "--pretty=format:" +
		"%H" + separator + // Hash
		"%B" + separator + // Subject
		"%aN" + separator + // Author Name
		"%aE" + separator + // Author Email
		"%ad" + separator + // Date
		"%at" + lineSeparator // Timestamp

	cmd := exec.Command("git", "--no-pager", "log", "--first-parent", "master", "--reverse", logfmt)
	cmd.Dir = directory
	out, err := cmd.CombinedOutput()

	if err != nil {
		return []Commit{}, fmt.Errorf("failed to get the log of %s with err '%v' and message '%s'", directory, err, string(out))
	}

	return parseLogOutput(string(out))
}

func parseLogOutput(output string) ([]Commit, error) {
	commits := []Commit{}
	for _, line := range strings.Split(output, lineSeparator) {
		if line == "" {
			break
		}
		c, err := parseLogLine(line)
		if err != nil {
			return commits, err
		}
		commits = append(commits, c)
	}
	return commits, nil
}

func parseLogLine(line string) (Commit, error) {
	lineParts := strings.Split(line, separator)

	if len(lineParts) != 6 {
		return Commit{}, fmt.Errorf("failed to parse log line '%s', unexpected number of commit parts found", line)
	}

	ts, err := strconv.ParseInt(lineParts[5], 10, 64)

	if err != nil {
		return Commit{}, fmt.Errorf("failed to parse timestamp value of '%s' for line '%s' with err %v", lineParts[5], line, err)
	}

	return Commit{
		Hash:      strings.TrimSpace(lineParts[0]),
		Body:      strings.TrimSpace(lineParts[1]),
		Name:      strings.TrimSpace(lineParts[2]),
		Email:     strings.TrimSpace(lineParts[3]),
		Timestamp: ts,
	}, nil
}
