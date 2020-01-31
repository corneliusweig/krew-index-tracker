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

package homebrew

import (
	"context"

	"github.com/corneliusweig/krew-index-tracker/pkg/homebrew/client"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func SaveAnalyticsToBigQuery(ctx context.Context) error {
	fetcher := client.NewHomebrew("https://formulae.brew.sh/api/formula/krew.json")

	logrus.Infof("Fetching homebrew download statistics")
	homebrewStats, err := fetcher.FetchAnalytics(ctx)
	if err != nil {
		return errors.Wrapf(err, "could not fetch homebrew analytics")
	}

	logrus.Infof("Uploading summaries to BigQuery")
	if err := client.HomebrewBigQuery().Upload(ctx, homebrewStats); err != nil {
		return errors.Wrapf(err, "failed saving scraped data")
	}

	return nil
}
