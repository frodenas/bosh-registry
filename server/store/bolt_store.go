package store

import (
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/boltdb/bolt"
)

const boltStoreLogTag = "BoltRegistryStore"
const boltStoreFileMode = 0600
const boltStoreFileLockTimeout = 1
const boltStoreBucketName = "Registry"

type BoltStore struct {
	config BoltConfig
	logger boshlog.Logger
}

func NewBoltStore(
	config BoltConfig,
	logger boshlog.Logger,
) BoltStore {
	return BoltStore{
		config: config,
		logger: logger,
	}
}

func (s BoltStore) Delete(key string) error {
	db, err := s.openDB()
	if err != nil {
		return bosherr.WrapErrorf(err, "Deleting key '%s'", key)
	}
	defer db.Close()

	s.logger.Debug(boltStoreLogTag, "Deleting key '%s'", key)
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(boltStoreBucketName))
		if bucket != nil {
			return bucket.Delete([]byte(key))
		}
		return nil
	})
	if err != nil {
		return bosherr.WrapErrorf(err, "Deleting key '%s'", key)
	}

	return nil
}

func (s BoltStore) Get(key string) (string, bool, error) {
	db, err := s.openDB()
	if err != nil {
		return "", false, bosherr.WrapErrorf(err, "Reading key '%s'", key)
	}
	defer db.Close()

	var value []byte
	s.logger.Debug(boltStoreLogTag, "Reading key '%s'", key)
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(boltStoreBucketName))
		if bucket != nil {
			value = bucket.Get([]byte(key))
		}
		return nil
	})
	if value != nil {
		return string(value), true, nil
	}

	return "", false, nil
}

func (s BoltStore) Save(key string, value string) error {
	db, err := s.openDB()
	if err != nil {
		return bosherr.WrapErrorf(err, "Saving key '%s'", key)
	}
	defer db.Close()

	s.logger.Debug(boltStoreLogTag, "Saving key '%s'", key)
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(boltStoreBucketName))
		if err != nil {
			return bosherr.WrapErrorf(err, "Creating bucket '%s'", boltStoreBucketName)
		}
		return bucket.Put([]byte(key), []byte(value))
	})
	if err != nil {
		return bosherr.WrapErrorf(err, "Saving key '%s'", key)
	}

	return nil
}

func (s BoltStore) openDB() (db *bolt.DB, err error) {
	dbOptions := &bolt.Options{
		Timeout: boltStoreFileLockTimeout * time.Second,
	}
	db, err = bolt.Open(s.config.DBFile, boltStoreFileMode, dbOptions)
	if err != nil {
		return db, bosherr.WrapError(err, "Opening Bolt database")
	}

	return db, nil
}
