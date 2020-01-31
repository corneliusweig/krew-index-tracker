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

package uploader

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Entity struct {
	Id          string
	Description string
}

type SchemaGenerator func() (bigquery.Schema, error)

type Client struct {
	projectId string
	dataset   Entity
	table     Entity
	getSchema SchemaGenerator
}

func NewClient(project string, dataset, table Entity, getSchema SchemaGenerator) *Client {
	return &Client{
		projectId: project,
		dataset:   dataset,
		table:     table,
		getSchema: getSchema,
	}
}

func (c *Client) String() string {
	return fmt.Sprintf("%s.%s.%s", c.projectId, c.dataset.Id, c.table.Id)
}

func (c *Client) Upload(ctx context.Context, data interface{}) error {
	client, err := bigquery.NewClient(ctx, c.projectId)
	if err != nil {
		return errors.Wrapf(err, "failed to create bq client")
	}

	ds, err := c.ensureDataset(ctx, client)
	if err != nil {
		return errors.Wrapf(err, "failed to initialize dataset")
	}

	table, err := c.ensureTable(ctx, ds)
	if err != nil {
		return errors.Wrapf(err, "failed to initialize table")
	}

	return errors.Wrapf(table.Inserter().Put(ctx, data), "could not insert rows")
}

func (c *Client) ensureDataset(ctx context.Context, client *bigquery.Client) (*bigquery.Dataset, error) {
	// Creates the new BigQuery dataset.
	ds := client.Dataset(c.dataset.Id)

	if meta, _ := ds.Metadata(ctx); meta != nil {
		logrus.Infof("Dataset already exists")
		return ds, nil
	}

	if err := ds.Create(ctx, &bigquery.DatasetMetadata{
		Name:        c.dataset.Id,
		Description: c.dataset.Description,
		Location:    "US",
	}); err != nil {
		return nil, errors.Wrapf(err, "creating dataset")
	}
	logrus.Infof("Dataset created")

	return ds, nil
}

func (c *Client) ensureTable(ctx context.Context, ds *bigquery.Dataset) (*bigquery.Table, error) {
	table := ds.Table(c.table.Id)
	if meta, _ := table.Metadata(ctx); meta != nil {
		logrus.Infof("Found table with the same name")
		return table, nil
	}

	schema, err := c.getSchema()
	if err != nil {
		return nil, errors.Wrapf(err, "could not infer schema for %q", c)
	}
	logrus.Debugf("Schema looks good")

	meta := &bigquery.TableMetadata{
		Schema:           schema,
		Description:      c.table.Description,
		TimePartitioning: &bigquery.TimePartitioning{Field: "CreatedAt", Expiration: 1000 * 24 * time.Hour},
	}
	if err := table.Create(ctx, meta); err != nil {
		return nil, errors.Wrapf(err, "could not create BQ table %q", c)
	}
	logrus.Infof("Created table %q", c.table.Id)

	return table, nil
}
