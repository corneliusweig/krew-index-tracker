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

package krew

import (
	"regexp"

	"github.com/corneliusweig/krew-index-tracker/pkg/constants"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/krew/pkg/index"
	"sigs.k8s.io/krew/pkg/index/indexscanner"
)

var (
	gitHubRepo = regexp.MustCompile(".*github.com/([^/]+)/([^/]+).*")
)

type PluginHandle struct {
	index.PluginSpec
	Owner, Repo string
}

func GetRepoList() ([]PluginHandle, error) {
	logrus.Debugf("Reading repo list")

	plugins, err := indexscanner.LoadPluginListFromFS(constants.PluginsDir)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read index")
	}
	res := make([]PluginHandle, 0, len(plugins))
	for _, plugin := range plugins {
		homepage := plugin.Spec.Homepage
		submatch := gitHubRepo.FindStringSubmatch(homepage)
		if len(submatch) < 3 {
			logrus.Infof("Skipping repository '%s'", homepage)
			continue
		}
		logrus.Debugf("%s -> %s/%s", homepage, submatch[1], submatch[2])
		res = append(res, PluginHandle{
			PluginSpec: plugins[0].Spec,
			Owner:      submatch[1],
			Repo:       submatch[2],
		})
	}
	return res, nil
}
