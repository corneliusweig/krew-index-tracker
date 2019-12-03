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
	"github.com/corneliusweig/krew-index-tracker/pkg/github"
	"github.com/corneliusweig/krew-index-tracker/pkg/repository"
	"github.com/corneliusweig/krew-index-tracker/pkg/repository/krew"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func SaveDownloadCountsToBigQuery(ctx context.Context, token string, isUpdateIndex bool) error {
	logrus.Infof("Determine repositories to inspect")
	repos, err := krew.NewRepositoryProvider(isUpdateIndex).List(ctx)
	if err != nil {
		return errors.Wrapf(err, "could not determine list of repositories")
	}

	logrus.Infof("Fetching repository download summaries")
	summaries, err := fetchSummaries(ctx, token, repos)
	if err != nil {
		return errors.Wrapf(err, "could not fetch repo summaries")
	}

	logrus.Infof("Uploading summaries to BigQuery")
	if err := bigquery.Upload(ctx, summaries); err != nil {
		return errors.Wrapf(err, "failed saving scraped data")
	}

	logrus.Infof("All good")
	return nil
}

func fetchSummaries(ctx context.Context, token string, handles []repository.Handle) ([]github.RepoSummary, error) {
	releases := github.NewReleaseFetcher(ctx, token)
	summaries := make([]github.RepoSummary, 0, len(handles))
	for _, h := range handles {
		summary, err := releases.Summary(h)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}
