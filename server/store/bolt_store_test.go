package store_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-registry/server/store"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

var _ = Describe("BoltStore", func() {
	var (
		err       error
		boltStore BoltStore
		dbFile    *os.File
		config    BoltConfig

		logger = boshlog.NewLogger(boshlog.LevelNone)
	)

	BeforeEach(func() {
		dbFile, err = ioutil.TempFile("", "test-bolt")
		Expect(err).ToNot(HaveOccurred())

		config = BoltConfig{
			DBFile: dbFile.Name(),
		}
		boltStore = NewBoltStore(config, logger)
	})

	Describe("Get", func() {
		It("returns the value if key exist", func() {
			err = boltStore.Save("fake-key", "fake-value")
			Expect(err).ToNot(HaveOccurred())

			value, found, err := boltStore.Get("fake-key")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(value).To(Equal("fake-value"))
		})

		It("returns false if key does not exist", func() {
			_, found, err := boltStore.Get("fake-key")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeFalse())
		})
	})

	Describe("Delete", func() {
		It("deletes the key if it exist", func() {
			err = boltStore.Save("fake-key", "fake-value")
			Expect(err).ToNot(HaveOccurred())

			err = boltStore.Delete("fake-key")
			Expect(err).ToNot(HaveOccurred())

			_, found, err := boltStore.Get("fake-key")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeFalse())
		})

		It("does not return error if key does not exist", func() {
			err = boltStore.Delete("fake-key")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Save", func() {
		It("stores the appropiate value when key does not exist", func() {
			err = boltStore.Save("fake-key", "fake-value")
			Expect(err).ToNot(HaveOccurred())

			value, found, err := boltStore.Get("fake-key")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(value).To(Equal("fake-value"))
		})

		It("updates the appropiate value when key already exist", func() {
			err = boltStore.Save("fake-key", "fake-value")
			Expect(err).ToNot(HaveOccurred())

			err = boltStore.Save("fake-key", "fake-new-value")
			Expect(err).ToNot(HaveOccurred())

			value, found, err := boltStore.Get("fake-key")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(value).To(Equal("fake-new-value"))
		})
	})
})
