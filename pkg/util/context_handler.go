package util

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func ContextWithCtrlCHandler(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGPIPE)

	go func() {
		<-sigs
		signal.Stop(sigs)
		cancel()
		logrus.Infof("Aborted.")
	}()

	return ctx
}
