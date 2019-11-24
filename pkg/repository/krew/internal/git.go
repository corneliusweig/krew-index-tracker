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

package internal

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// UpdateAndCleanUntracked will fetch origin and set HEAD to origin/HEAD
// and also will create a pristine working directory by removing
// untracked files and directories.
func UpdateAndCleanUntracked(ctx context.Context, updateIndex bool, destinationPath string) error {
	if !updateIndex {
		logrus.Infof("Skipping index update")
		return nil
	}

	if err := git(ctx, destinationPath, "fetch", "origin", "master", "--verbose", "--depth", "1"); err != nil {
		return errors.Wrapf(err, "fetch index at %q failed", destinationPath)
	}

	if err := git(ctx, destinationPath, "reset", "--hard", "@{upstream}"); err != nil {
		return errors.Wrapf(err, "reset index at %q failed", destinationPath)
	}

	err := git(ctx, destinationPath, "clean", "-xfd")
	return errors.Wrapf(err, "clean index at %q failed", destinationPath)
}

func git(ctx context.Context, pwd string, args ...string) error {
	logrus.Infof("Going to run git %s", strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = pwd
	buf := bytes.Buffer{}
	var w io.Writer = &buf
	if logrus.InfoLevel < logrus.GetLevel() {
		w = io.MultiWriter(w, os.Stderr)
	}
	cmd.Stdout, cmd.Stderr = w, w
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "command execution failure, output=%q", buf.String())
	}
	return nil
}
