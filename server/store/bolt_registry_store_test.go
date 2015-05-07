package store_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-registry/server/store"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
)

var _ = Describe("BoltRegistryStore", func() {
	var (
		err               error
		boltRegistryStore BoltRegistryStore
		dbFile            *os.File
		config            BoltRegistryStoreConfig

		logger = boshlog.NewLogger(boshlog.LevelNone)
	)

	BeforeEach(func() {
		dbFile, err = ioutil.TempFile("", "test-bolt")
		Expect(err).ToNot(HaveOccurred())

		config = BoltRegistryStoreConfig{
			DBFile: dbFile.Name(),
		}
		boltRegistryStore = NewBoltRegistryStore(config, logger)
	})

	Describe("Get", func() {
		It("returns the value if key exist", func() {
			err = boltRegistryStore.Save("fake-key", "fake-value")
			Expect(err).ToNot(HaveOccurred())

			value, found, err := boltRegistryStore.Get("fake-key")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(value).To(Equal("fake-value"))
		})

		It("returns false if key does not exist", func() {
			_, found, err := boltRegistryStore.Get("fake-key")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeFalse())
		})
	})

	Describe("Delete", func() {
		It("deletes the key if it exist", func() {
			err = boltRegistryStore.Save("fake-key", "fake-value")
			Expect(err).ToNot(HaveOccurred())

			err = boltRegistryStore.Delete("fake-key")
			Expect(err).ToNot(HaveOccurred())

			_, found, err := boltRegistryStore.Get("fake-key")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeFalse())
		})

		It("does not return error if key does not exist", func() {
			err = boltRegistryStore.Delete("fake-key")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Save", func() {
		It("stores the appropiate value when key does not exist", func() {
			err = boltRegistryStore.Save("fake-key", "fake-value")
			Expect(err).ToNot(HaveOccurred())

			value, found, err := boltRegistryStore.Get("fake-key")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(value).To(Equal("fake-value"))
		})

		It("updates the appropiate value when key already exist", func() {
			err = boltRegistryStore.Save("fake-key", "fake-value")
			Expect(err).ToNot(HaveOccurred())

			err = boltRegistryStore.Save("fake-key", "fake-new-value")
			Expect(err).ToNot(HaveOccurred())

			value, found, err := boltRegistryStore.Get("fake-key")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(value).To(Equal("fake-new-value"))
		})
	})
})
