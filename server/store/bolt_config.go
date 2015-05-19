package store

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type BoltConfig struct {
	DBFile string
}

func (c BoltConfig) Validate() error {
	if c.DBFile == "" {
		return bosherr.Error("Must provide a non-empty DBFile")
	}

	return nil
}
