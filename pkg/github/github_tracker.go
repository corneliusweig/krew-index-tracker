/*
Copyright 2020 Cornelius Weig.

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

package github

import (
	"context"
	"net/http"
	"sync"

	api "github.com/google/go-github/v34/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/corneliusweig/krew-index-tracker/pkg/github/client"
	"github.com/corneliusweig/krew-index-tracker/pkg/github/repository"
	"github.com/corneliusweig/krew-index-tracker/pkg/github/repository/krew"
)

func SaveDownloadCountsToBigQuery(ctx context.Context, token string, isUpdateIndex bool) error {
	logrus.Infof("Determine repositories to inspect")
	repos, err := krew.NewRepositoryProvider(isUpdateIndex).List(ctx)
	if err != nil {
		return errors.Wrapf(err, "could not determine list of repositories")
	}

	logrus.Infof("Fetching repository download summaries")
	summaries := fetchSummaries(ctx, token, repos)

	logrus.Infof("Uploading summaries to BigQuery")
	if err := client.GithubBigQuery().Upload(ctx, summaries); err != nil {
		return errors.Wrapf(err, "failed saving scraped data")
	}

	return nil
}

func fetchSummaries(ctx context.Context, token string, handles []repository.Handle) []client.RepoSummary {
	var httpCli *http.Client
	if token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		httpCli = oauth2.NewClient(ctx, ts)
	}
	cli := api.NewClient(httpCli)

	var mu sync.Mutex // guards ret
	var ret []client.RepoSummary

	const poolSize = 20
	sem := make(chan struct{}, poolSize)
	for _, h := range handles {
		h := h
		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()
			summary, err := client.Summary(ctx, cli, h)
			if err != nil {
				logrus.Warnf("Could not fetch summary for %q: %v", h.PluginName, err)
				return
			}
			mu.Lock()
			defer mu.Unlock()
			ret = append(ret, summary)
		}()
	}
	for i := 0; i < poolSize; i++ {
		sem <- struct{}{}
	}
	return ret
}
