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

package bigquery

import (
	"context"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/corneliusweig/krew-index-tracker/pkg/constants"
	"github.com/corneliusweig/krew-index-tracker/pkg/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func Upload(ctx context.Context, items []github.RepoSummary) error {
	client, err := bigquery.NewClient(ctx, constants.ProjectID)
	if err != nil {
		return errors.Wrapf(err, "failed to create bq client")
	}

	ds, err := ensureDataset(ctx, client)
	if err != nil {
		return errors.Wrapf(err, "failed to initialize dataset")
	}

	table, err := ensureTable(ctx, ds)
	if err != nil {
		return errors.Wrapf(err, "failed to initialize table")
	}

	return errors.Wrapf(table.Inserter().Put(ctx, items), "could not insert rows into BQ")
}

func ensureDataset(ctx context.Context, client *bigquery.Client) (*bigquery.Dataset, error) {
	// Creates the new BigQuery dataset.
	ds := client.Dataset(constants.BQDataset)

	if meta, _ := ds.Metadata(ctx); meta != nil {
		logrus.Infof("Dataset already exists")
		return ds, nil
	}

	if err := ds.Create(ctx, &bigquery.DatasetMetadata{
		Name:        constants.BQDataset,
		Description: "Download counts for all plugins in the centralized krew index",
		Location:    "US",
	}); err != nil {
		return nil, errors.Wrapf(err, "creating dataset")
	}
	logrus.Infof("Dataset created")

	return ds, nil
}

func ensureTable(ctx context.Context, ds *bigquery.Dataset) (*bigquery.Table, error) {
	schema, err := bigquery.InferSchema(github.RepoSummary{})
	if err != nil {
		return nil, errors.Wrapf(err, "could not infer schema for 'RepoSummary'")
	}
	logrus.Debugf("Schema looks good")

	table := ds.Table(constants.BQTableName)
	if meta, _ := table.Metadata(ctx); meta != nil {
		logrus.Infof("Found table with the same name")
	} else {
		meta = &bigquery.TableMetadata{
			Schema:           schema,
			TimePartitioning: &bigquery.TimePartitioning{Field: "CreatedAt", Expiration: 1000 * 24 * time.Hour},
		}
		if err := table.Create(ctx, meta); err != nil {
			return nil, errors.Wrapf(err, "could not create BQ table '%s'", constants.BQTableName)
		}
		logrus.Infof("Created table '%s'", constants.BQTableName)
	}

	return table, nil
}

func createTableName() string {
	return time.Now().Format("2006_01_02_150405")
}
