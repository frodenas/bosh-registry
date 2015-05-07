package store_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-registry/server/store"
)

var _ = Describe("BoltRegistryStoreConfig", func() {
	var (
		options BoltRegistryStoreConfig

		validOptions = BoltRegistryStoreConfig{
			DBFile: "fake-dbfile",
		}
	)

	Describe("Validate", func() {
		BeforeEach(func() {
			options = validOptions
		})

		It("does not return error if all fields are valid", func() {
			err := options.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if DBFile is empty", func() {
			options.DBFile = ""

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty DBFile"))
		})
	})
})
