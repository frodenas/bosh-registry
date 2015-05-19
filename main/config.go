package main

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	"github.com/frodenas/bosh-registry/server"
	"github.com/frodenas/bosh-registry/server/store"
)

type Config struct {
	Server server.Config `json:"server,omitempty"`
	Store  store.Config  `json:"store,omitempty"`
}

func NewConfigFromPath(path string, fs boshsys.FileSystem) (Config, error) {
	var config Config

	bytes, err := fs.ReadFile(path)
	if err != nil {
		return config, bosherr.WrapErrorf(err, "Reading config file '%s'", path)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, bosherr.WrapError(err, "Unmarshalling config")
	}

	err = config.Validate()
	if err != nil {
		return config, bosherr.WrapError(err, "Validating config")
	}

	return config, nil
}

func (c Config) Validate() error {
	err := c.Server.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating Server configuration")
	}

	err = c.Store.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating Store configuration")
	}

	return nil
}
