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

func (rf *ReleaseFetcher) RepoSummary(owner, repo string) (RepoSummary, error) {
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
