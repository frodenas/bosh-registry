package main_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	fakesys "github.com/cloudfoundry/bosh-agent/system/fakes"

	. "github.com/frodenas/bosh-registry/main"

	"github.com/frodenas/bosh-registry/server/listener"
	"github.com/frodenas/bosh-registry/server/store"
)

var validServerConfig = server.Config{
	Protocol: "http",
	Address:  "fake-host",
	Port:     5555,
	Username: "fake-username",
	Password: "fake-password",
}

var validStoreOptions = store.Config{
	Adapter: "fake-adapter",
}

var validConfig = Config{
	Server: validServerConfig,
	Store:  validStoreOptions,
}

var _ = Describe("NewConfigFromPath", func() {
	var (
		fs *fakesys.FakeFileSystem
	)

	BeforeEach(func() {
		fs = fakesys.NewFakeFileSystem()
	})

	It("returns error if config is not valid", func() {
		err := fs.WriteFileString("/config.json", "{}")
		Expect(err).ToNot(HaveOccurred())

		_, err = NewConfigFromPath("/config.json", fs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Validating config"))
	})

	It("returns error if file contains invalid json", func() {
		err := fs.WriteFileString("/config.json", "-")
		Expect(err).ToNot(HaveOccurred())

		_, err = NewConfigFromPath("/config.json", fs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Unmarshalling config"))
	})

	It("returns error if file cannot be read", func() {
		err := fs.WriteFileString("/config.json", "{}")
		Expect(err).ToNot(HaveOccurred())

		fs.ReadFileError = errors.New("fake-read-err")

		_, err = NewConfigFromPath("/config.json", fs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("fake-read-err"))
	})
})

var _ = Describe("Config", func() {
	var (
		config Config
	)

	Describe("Validate", func() {
		BeforeEach(func() {
			config = validConfig
		})

		It("does not return error if all server and store sections are valid", func() {
			err := config.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if server section is not valid", func() {
			config.Server.Address = ""

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validating Server configuration"))
		})

		It("returns error if store section is not valid", func() {
			config.Store.Adapter = ""

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validating Store configuration"))
		})
	})
})
