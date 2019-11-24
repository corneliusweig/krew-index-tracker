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

package github

import (
	"context"
	"time"

	api "github.com/google/go-github/v28/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type ReleaseFetcher struct {
	ctx    context.Context
	client *api.Client
}

type RepoSummary struct {
	Owner    string           `json:"owner,omitempty"`
	Repo     string           `json:"repo,omitempty"`
	Releases []ReleaseSummary `json:"releases,omitempty"`
}

type ReleaseSummary struct {
	TagName       string         `json:"tag,omitempty"`
	Name          string         `json:"name,omitempty"`
	PublishedAt   time.Time      `json:"published_at,omitempty"`
	User          string         `json:"user,omitempty"`
	ReleaseAssets []AssetSummary `json:"assets,omitempty"`
}

type AssetSummary struct {
	Name          string `json:"name"`
	State         string `json:"state"`
	DownloadCount int    `json:"download_count"`
}

func NewReleaseFetcher(ctx context.Context, token string) *ReleaseFetcher {
	return &ReleaseFetcher{
		ctx:    ctx,
		client: getClient(ctx, token),
	}
}

func (rf *ReleaseFetcher) Summary(owner, repo string) (RepoSummary, error) {
	logrus.Infof("Fetching summary for %s/%s", owner, repo)
	releases, _, err := rf.client.Repositories.ListReleases(rf.ctx, owner, repo, nil)
	if err != nil {
		return RepoSummary{}, errors.Wrapf(err, "listing releases of %s/%s", owner, repo)
	}
	return RepoSummary{
		Owner:    owner,
		Repo:     repo,
		Releases: toReleaseSummaries(releases),
	}, nil
}

func toAssetSummaries(as []api.ReleaseAsset) (res []AssetSummary) {
	for _, asset := range as {
		res = append(res, AssetSummary{
			Name:          asset.GetName(),
			State:         asset.GetState(),
			DownloadCount: asset.GetDownloadCount(),
		})
	}
	return
}

func toReleaseSummaries(rs []*api.RepositoryRelease) (res []ReleaseSummary) {
	for _, r := range rs {
		if r.GetPrerelease() || r.GetDraft() {
			logrus.Debugf("Skipping release %s", r.GetName())
			continue
		}
		res = append(res, ReleaseSummary{
			TagName:       r.GetTagName(),
			Name:          r.GetName(),
			PublishedAt:   r.GetPublishedAt().Time,
			User:          r.GetAuthor().GetLogin(),
			ReleaseAssets: toAssetSummaries(r.Assets),
		})
	}
	return
}

func getClient(ctx context.Context, token string) *api.Client {
	if len(token) == 0 {
		return api.NewClient(nil)
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return api.NewClient(tc)
}
