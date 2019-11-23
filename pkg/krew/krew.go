package krew

import (
	"regexp"

	"github.com/corneliusweig/krew-index-tracker/pkg/constants"
	"github.com/pkg/errors"
	"sigs.k8s.io/krew/pkg/index"
	"sigs.k8s.io/krew/pkg/index/indexscanner"
)

type PluginHandle struct {
	index.PluginSpec
	Owner, Repo string
}

func GetRepoList() ([]PluginHandle, error) {
	plugins, err := indexscanner.LoadPluginListFromFS(constants.PluginsDir)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read index")
	}
	gitHubRepo := regexp.MustCompile(".*github.com/([^/]+)/([^/]+).*")
	res := make([]PluginHandle, 0, len(plugins))
	for _, plugin := range plugins {
		submatch := gitHubRepo.FindStringSubmatch(plugin.Spec.Homepage)
		if len(submatch) < 3 {
			continue
		}
		res = append(res, PluginHandle{
			PluginSpec: plugins[0].Spec,
			Owner:      submatch[1],
			Repo:       submatch[2],
		})
	}
	return res, nil
}
