package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	"github.com/frodenas/bosh-registry/server"
	"github.com/frodenas/bosh-registry/server/store"
)

const mainLogTag = "main"

var (
	configFileOpt = flag.String("configFile", "", "Path to configuration file")
)

func main() {
	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr, os.Stderr)
	fs := boshsys.NewOsFileSystem(logger)

	defer logger.HandlePanic("Main")

	flag.Parse()

	config, err := NewConfigFromPath(*configFileOpt, fs)
	if err != nil {
		logger.Error(mainLogTag, "Loading config: %s", err.Error())
		os.Exit(1)
	}

	instanceHandler, err := createInstanceHandler(config, logger)
	if err != nil {
		logger.Error(mainLogTag, "Creating Registry Instance Handler: %s", err.Error())
		os.Exit(1)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	listener := server.NewListener(config.Server, instanceHandler, logger)
	errChan := listener.ListenAndServe()
	select {
	case err := <-errChan:
		if err != nil {
			logger.Error(mainLogTag, "Error occurred: %s", err.Error())
			os.Exit(1)
		}
	case sig := <-signals:
		logger.Debug(mainLogTag, "Exiting, received signal: %#v", sig)
		listener.Stop()
	}

	os.Exit(0)
}

func createInstanceHandler(config Config, logger boshlog.Logger) (*server.InstanceHandler, error) {
	store, err := store.NewStore(config.Store, logger)
	if err != nil {
		return nil, bosherr.WrapError(err, "Creating a Registry Store")
	}

	instanceHandler := server.NewInstanceHandler(config.Server, store, logger)

	return instanceHandler, nil
}
