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
	"cloud.google.com/go/bigquery"
	"github.com/corneliusweig/krew-index-tracker/pkg/globals"
	"github.com/corneliusweig/krew-index-tracker/pkg/uploader"
)

func HomebrewBigQuery() *uploader.Client {
	return uploader.NewClient(
		globals.ProjectID,
		uploader.Entity{
			Id:          globals.BQDataset,
			Description: "Download counts for all plugins in the centralized krew index",
		},
		uploader.Entity{
			Id:          "homebrew_krew",
			Description: "Download counts for krew view homebrew",
		},
		func() (bigquery.Schema, error) {
			return bigquery.InferSchema(HomebrewStats{})
		},
	)
}
