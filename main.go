/*
Copyright 2019 Cornelius Weig

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

package main

import (
	"context"

	"github.com/corneliusweig/krew-index-tracker/pkg/bigquery"
	"github.com/corneliusweig/krew-index-tracker/pkg/constants"
	"github.com/corneliusweig/krew-index-tracker/pkg/git"
	"github.com/corneliusweig/krew-index-tracker/pkg/github"
	"github.com/corneliusweig/krew-index-tracker/pkg/krew"
	"github.com/corneliusweig/krew-index-tracker/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	token         string
	isUpdateIndex bool
)

var rootCmd = &cobra.Command{
	Use:     "krew-index-tracker",
	Example: "krew-index-tracker",
	Short:   "Generate a markdown changelog of merged pull requests since last release",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := git.UpdateAndCleanUntracked(isUpdateIndex, constants.IndexDir); err != nil {
			logrus.Fatal(err)
		}

		repos, err := krew.GetRepoList()
		if err != nil {
			logrus.Fatal(err)
		}

		ctx := util.ContextWithCtrlCHandler(context.Background())
		releaseFetcher := github.NewReleaseFetcher(ctx, token)

		var summaries []github.RepoSummary
		for _, repo := range repos {
			summary, err := releaseFetcher.RepoSummary(repo.Owner, repo.Repo)
			if err != nil {
				logrus.Warn(err)
				continue
			}
			summaries = append(summaries, summary)
		}

		if err := bigquery.Upload(ctx, summaries); err != nil {
			logrus.Error(err)
		}
		logrus.Debugf("All good")
	},
}

func main() {
	rootCmd.Flags().StringVar(&token, "token", "", "Specify personal Github Token if you are hitting a rate limit anonymously. https://github.com/settings/tokens")
	rootCmd.Flags().BoolVar(&isUpdateIndex, "update-index", false, "Call git to ensure that the index is up to date")
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
