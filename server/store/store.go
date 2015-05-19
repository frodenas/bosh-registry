package store

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/mitchellh/mapstructure"
)

type Store interface {
	Delete(string) error
	Get(string) (string, bool, error)
	Save(string, string) error
}

func NewStore(
	config Config,
	logger boshlog.Logger,
) (Store, error) {
	switch {
	case config.Adapter == "bolt":
		boltConfig := BoltConfig{}
		err := mapstructure.Decode(config.Options, &boltConfig)
		if err != nil {
			return nil, bosherr.WrapError(err, "Decoding Bolt Registry Store configuration")
		}

		err = boltConfig.Validate()
		if err != nil {
			return nil, bosherr.WrapError(err, "Validating Bolt Registry Store configuration")
		}

		return NewBoltStore(boltConfig, logger), nil
	}

	return nil, bosherr.Errorf("Registry Store adapter '%s' not supported", config.Adapter)
}
