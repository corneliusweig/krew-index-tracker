/*
Copyright 2019 Cornelius Weig.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
