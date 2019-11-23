package bigquery

import (
	"context"
	"log"
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
		log.Fatalf("Failed to create client: %v", err)
	}

	ds, err := ensureDataset(ctx, client)
	if err != nil {
		return err
	}

	table, err := ensureTable(ctx, ds)
	if err != nil {
		return err
	}

	return errors.Wrapf(table.Inserter().Put(ctx, items), "could not insert rows into BQ")
}

func ensureDataset(ctx context.Context, client *bigquery.Client) (*bigquery.Dataset, error) {
	// Creates the new BigQuery dataset.
	ds := client.Dataset(constants.BQDataset)

	if meta, _ := ds.Metadata(ctx); meta != nil {
		logrus.Debugf("Dataset already exists")
		return ds, nil
	}

	if err := ds.Create(ctx, &bigquery.DatasetMetadata{
		Name:        constants.BQDataset,
		Description: "Download counts for all plugins in the centralized krew index",
		Location:    "US",
	}); err != nil {
		return nil, errors.Wrapf(err, "creating dataset")
	}
	logrus.Debugf("Dataset created")

	return ds, nil
}

func ensureTable(ctx context.Context, ds *bigquery.Dataset) (*bigquery.Table, error) {
	tableName := createTableName()

	schema, err := bigquery.InferSchema(github.RepoSummary{})
	if err != nil {
		return nil, errors.Wrapf(err, "could not infer schema for 'RepoSummary'")
	}
	logrus.Debugf("Schema looks good")

	table := ds.Table(tableName)
	if meta, _ := table.Metadata(ctx); meta != nil {
		logrus.Infof("Found table with the same name")
	} else {
		if err := table.Create(ctx, &bigquery.TableMetadata{Schema: schema}); err != nil {
			return nil, errors.Wrapf(err, "could not create bq table")
		}
		logrus.Infof("Created table '%s'", tableName)
	}

	return table, nil
}

func createTableName() string {
	return time.Now().Format("2006_01_02_150405")
}
