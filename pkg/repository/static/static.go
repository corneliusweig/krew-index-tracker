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

package static

import (
	"bufio"
	"context"
	"os"

	"github.com/corneliusweig/krew-index-tracker/pkg/repository"
	"github.com/corneliusweig/krew-index-tracker/pkg/repository/githuburl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Provider struct {
	filename string
}

var _ repository.Provider = Provider{}

// NewStaticRepositoryProvider creates a repository provider which is configured by a static
// list of github repositories. It expects one github URL per line in the file.
func NewStaticRepositoryProvider(filename string) repository.Provider {
	return &Provider{filename: filename}
}

// List reads in a static list of github repositories with one entry per line.
func (p Provider) List(ctx context.Context) ([]repository.Handle, error) {
	repoList, err := os.Open(p.filename)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open repository list at %s", p.filename)
	}
	defer repoList.Close()

	lines := bufio.NewScanner(repoList)
	var res []repository.Handle
	for lines.Scan() {
		owner, repo, err := githuburl.Parse(lines.Text())
		if err != nil {
			logrus.Infof("Skipping repository plugin: %s", err)
			continue
		}
		// todo(corneliusweig): the static provider does not handle PluginName
		res = append(res, repository.Handle{Owner: owner, Repo: repo})
	}

	if err := lines.Err(); err != nil {
		return nil, errors.Wrapf(err, "could not scan repository list")
	}

	return res, nil
}
