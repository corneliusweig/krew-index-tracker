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

package client

import (
	"context"
	"time"

	api "github.com/google/go-github/v28/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/corneliusweig/krew-index-tracker/pkg/github/repository"
)

type ReleaseFetcher struct {
	ctx    context.Context
	client *api.Client
}

type RepoSummary struct {
	PluginName string           `json:"pluginName,omitempty"`
	Owner      string           `json:"owner,omitempty"`
	Repo       string           `json:"repo,omitempty"`
	CreatedAt  time.Time        `json:"created_at,omitempty"`
	Releases   []ReleaseSummary `json:"releases,omitempty"`
}

type ReleaseSummary struct {
	TagName       string         `json:"tag,omitempty"`
	PublishedAt   time.Time      `json:"published_at,omitempty"`
	ReleaseAssets []AssetSummary `json:"assets,omitempty"`
}

type AssetSummary struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	DownloadCount int    `json:"download_count"`
}

func NewReleaseFetcher(ctx context.Context, token string) *ReleaseFetcher {
	return &ReleaseFetcher{
		ctx:    ctx,
		client: getClient(ctx, token),
	}
}

func (rf *ReleaseFetcher) Summary(h repository.Handle) (RepoSummary, error) {
	logrus.Infof("Fetching summary for %s/%s", h.Owner, h.Repo)
	releases, _, err := rf.client.Repositories.ListReleases(rf.ctx, h.Owner, h.Repo, nil)
	if err != nil {
		return RepoSummary{}, errors.Wrapf(err, "listing releases of %s/%s", h.Owner, h.Repo)
	}
	return RepoSummary{
		PluginName: h.PluginName,
		Owner:      h.Owner,
		Repo:       h.Repo,
		CreatedAt:  time.Now(),
		Releases:   toReleaseSummaries(releases),
	}, nil
}

func toAssetSummaries(as []api.ReleaseAsset) (res []AssetSummary) {
	for _, asset := range as {
		res = append(res, AssetSummary{
			Name:          asset.GetName(),
			URL:           asset.GetBrowserDownloadURL(),
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
			PublishedAt:   r.GetPublishedAt().Time,
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
