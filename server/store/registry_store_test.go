package store_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-registry/server/store"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
)

var _ = Describe("RegistryStore", func() {
	var (
		config = RegistryStoreConfig{}
		logger = boshlog.NewLogger(boshlog.LevelNone)
	)

	Describe("NewRegistryStore", func() {
		It("returns error if Adapter is not supported", func() {
			config.Adapter = "fake-adapter"

			_, err := NewRegistryStore(config, logger)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Registry Store adapter 'fake-adapter' not supported"))
		})

		Context("when adapter is bolt", func() {
			BeforeEach(func() {
				config.Adapter = "bolt"
				config.Options = nil
			})

			It("does not return error if bolt configuration is not valid", func() {
				boltRegistryStoreConfig := map[string]interface{}{
					"DBFile": "fake-dbfile",
				}
				config.Options = boltRegistryStoreConfig
				_, err := NewRegistryStore(config, logger)
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns error if bolt configuration is not valid", func() {
				_, err := NewRegistryStore(config, logger)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Validating Bolt Registry Store configuration"))
			})
		})
	})
})
