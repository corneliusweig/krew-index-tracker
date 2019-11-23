package git

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
