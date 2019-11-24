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

package tracker

import (
	"context"

	"github.com/corneliusweig/krew-index-tracker/pkg/bigquery"
	"github.com/corneliusweig/krew-index-tracker/pkg/constants"
	"github.com/corneliusweig/krew-index-tracker/pkg/git"
	"github.com/corneliusweig/krew-index-tracker/pkg/github"
	"github.com/corneliusweig/krew-index-tracker/pkg/krew"
	"github.com/sirupsen/logrus"
)

func SaveDownloadCountsToBigQuery(ctx context.Context, token string, isUpdateIndex bool) {
	logrus.Debugf("U")
	if err := git.UpdateAndCleanUntracked(ctx, isUpdateIndex, constants.IndexDir); err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("Reading repo list")
	repos, err := krew.GetRepoList()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("Fetching repo download summaries")
	summaries := fetchSummaries(ctx, token, repos)

	logrus.Infof("Uploading summaries to BigQuery")
	if err := bigquery.Upload(ctx, summaries); err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("All good")
}

func fetchSummaries(ctx context.Context, token string, repos []krew.PluginHandle) []github.RepoSummary {
	releaseFetcher := github.NewReleaseFetcher(ctx, token)
	summaries := make([]github.RepoSummary, 0, len(repos))
	for _, repo := range repos {
		summary, err := releaseFetcher.RepoSummary(repo.Owner, repo.Repo)
		if err != nil {
			logrus.Warn(err)
			continue
		}
		summaries = append(summaries, summary)
	}
	return summaries
}
