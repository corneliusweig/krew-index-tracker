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
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/krew/pkg/index/indexscanner"

	"github.com/corneliusweig/krew-index-tracker/pkg/github/constants"
	"github.com/corneliusweig/krew-index-tracker/pkg/github/repository"
	"github.com/corneliusweig/krew-index-tracker/pkg/github/repository/githuburl"
	"github.com/corneliusweig/krew-index-tracker/pkg/github/repository/krew/internal"
)

type IndexRepositoryProvider struct {
	updateIndex bool
}

var _ repository.Provider = IndexRepositoryProvider{}

func NewRepositoryProvider(isUpdateIndex bool) repository.Provider {
	return &IndexRepositoryProvider{updateIndex: isUpdateIndex}
}

func (k IndexRepositoryProvider) List(ctx context.Context) ([]repository.Handle, error) {
	logrus.Debugf("Updating krew index")
	if err := internal.UpdateAndCleanUntracked(ctx, k.updateIndex, constants.IndexDir); err != nil {
		return nil, errors.Wrapf(err, "could not update krew index")
	}

	return getRepoList()
}

func getRepoList() ([]repository.Handle, error) {
	logrus.Infof("Reading repo list")
	plugins, err := indexscanner.LoadPluginListFromFS(constants.PluginsDir)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read index")
	}

	res := make([]repository.Handle, 0, len(plugins))
	for _, plugin := range plugins {
		homepage := resolveWellKnownDeviations(plugin.Spec.Homepage)
		owner, repo, err := githuburl.Parse(homepage)
		if err != nil {
			logrus.Infof("Skipping repository plugin: %s", err)
			continue
		}
		res = append(res, repository.Handle{Owner: owner, Repo: repo, PluginName: plugin.GetName()})
	}
	return res, nil
}

var knownHomePages = map[string]string{
	`https://sigs.k8s.io/krew`:                                   "kubernetes-sigs/krew",
	`https://krew.sigs.k8s.io/`:                                  "kubernetes-sigs/krew",
	`https://kubernetes.github.io/ingress-nginx/kubectl-plugin/`: "kubernetes/ingress-nginx",
	`https://kudo.dev/`:                                          "kudobuilder/kudo",
	`https://kubevirt.io`:                                        "kubevirt/kubectl-virt-plugin",
	`https://popeyecli.io`:                                       "derailed/popeye",
	`https://soluble-ai.github.io/kubetap/`:                      "soluble-ai/kubetap",
	"https://schemahero.io":                                      "schemahero/schemahero",
}

func resolveWellKnownDeviations(homepage string) string {
	repo, ok := knownHomePages[homepage]
	if ok {
		return "github.com/" + repo
	}
	return homepage
}
