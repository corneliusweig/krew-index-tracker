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

package client

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Homebrew struct {
	url string
}

type HomebrewAPITarget struct {
	Count int `json:"krew"`
}

type HomebrewAPIInstalls struct {
	Aggregate30d  HomebrewAPITarget `json:"30d,omitempty"`
	Aggregate90d  HomebrewAPITarget `json:"90d,omitempty"`
	Aggregate365d HomebrewAPITarget `json:"365d,omitempty"`
}

type HomebrewAPIAnalytics struct {
	Installs          HomebrewAPIInstalls `json:"install"`
	InstallsOnRequest HomebrewAPIInstalls `json:"install_on_request"`
	InstallErrors     HomebrewAPIInstalls `json:"build_error"`
}

type HomebrewAPIResponse struct {
	Analytics HomebrewAPIAnalytics `json:"analytics,omitempty"`
}

type HomebrewStats struct {
	CreatedAt             time.Time
	Installs30d           int `bigquery:"installs_30d"`
	Installs90d           int `bigquery:"installs_90d"`
	Installs365d          int `bigquery:"installs_365d"`
	InstallsOnRequest30d  int `bigquery:"installs_on_request_30d"`
	InstallsOnRequest90d  int `bigquery:"installs_on_request_90d"`
	InstallsOnRequest365d int `bigquery:"installs_on_request_365d"`
	BuildErrors30d        int `bigquery:"build_errors_30d"`
}

func NewHomebrew(url string) *Homebrew {
	return &Homebrew{url: url}
}

func (h *Homebrew) FetchAnalytics(ctx context.Context) (HomebrewStats, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", h.url, nil)
	if err != nil {
		return HomebrewStats{}, errors.Wrapf(err, "error creating GET request for homebrew")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return HomebrewStats{}, errors.Wrapf(err, "error requesting analytics data from homebrew")
	}
	defer res.Body.Close()

	if res.StatusCode < 200 && 300 <= res.StatusCode {
		return HomebrewStats{}, errors.New("fetching homebrew analytics failed")
	}

	var body HomebrewAPIResponse
	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		return HomebrewStats{}, errors.Wrapf(err, "error parsing response from homebrew")
	}

	return stats(&body), nil
}

func stats(r *HomebrewAPIResponse) HomebrewStats {
	return HomebrewStats{
		CreatedAt:             time.Now(),
		Installs30d:           r.Analytics.Installs.Aggregate30d.Count,
		Installs90d:           r.Analytics.Installs.Aggregate90d.Count,
		Installs365d:          r.Analytics.Installs.Aggregate365d.Count,
		InstallsOnRequest30d:  r.Analytics.InstallsOnRequest.Aggregate30d.Count,
		InstallsOnRequest90d:  r.Analytics.InstallsOnRequest.Aggregate90d.Count,
		InstallsOnRequest365d: r.Analytics.InstallsOnRequest.Aggregate365d.Count,
		BuildErrors30d:        r.Analytics.InstallErrors.Aggregate30d.Count,
	}
}
