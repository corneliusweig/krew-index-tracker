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

	api "github.com/google/go-github/v34/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/corneliusweig/krew-index-tracker/pkg/github/repository"
)

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

func Summary(ctx context.Context, cli *api.Client, h repository.Handle) (RepoSummary, error) {
	logrus.Infof("Fetching summary for %s/%s", h.Owner, h.Repo)
	releases, _, err := cli.Repositories.ListReleases(ctx, h.Owner, h.Repo, nil)
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

func toAssetSummaries(as []*api.ReleaseAsset) (res []AssetSummary) {
	for _, asset := range as {
		res = append(res, AssetSummary{
			Name:          asset.GetName(),
			URL:           asset.GetBrowserDownloadURL(),
			DownloadCount: asset.GetDownloadCount(),
		})
	}
	return
}
