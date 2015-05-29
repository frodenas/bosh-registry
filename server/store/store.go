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
		if err := mapstructure.Decode(config.Options, &boltConfig); err != nil {
			return nil, bosherr.WrapError(err, "Decoding Bolt Registry Store configuration")
		}

		if err := boltConfig.Validate(); err != nil {
			return nil, bosherr.WrapError(err, "Validating Bolt Registry Store configuration")
		}

		return NewBoltStore(boltConfig, logger), nil
	}

	return nil, bosherr.Errorf("Registry Store adapter '%s' not supported", config.Adapter)
}
