package githuburl

import (
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
)

var (
	gitHubRepo = regexp.MustCompile(".*github.com/([^/]+)/([^?/#]+).*")
)

// Parse returns the owner and repo for a GitHub URL.
func Parse(url string) (string, string, error) {
	submatch := gitHubRepo.FindStringSubmatch(url)
	if len(submatch) < 3 {
		return "", "", fmt.Errorf("'%s' is not a GitHub URL", url)
	}
	logrus.Debugf("%s -> %s/%s", url, submatch[1], submatch[2])
	return submatch[1], submatch[2], nil
}
