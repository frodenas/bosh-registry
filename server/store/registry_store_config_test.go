package store_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-registry/server/store"
)

var _ = Describe("RegistryStoreConfig", func() {
	var (
		options RegistryStoreConfig

		validOptions = RegistryStoreConfig{
			Adapter: "adapter",
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

		It("returns error if Adapter is empty", func() {
			options.Adapter = ""

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Adapter"))
		})
	})
})
