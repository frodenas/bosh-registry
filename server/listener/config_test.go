package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-registry/server/listener"
)

var _ = Describe("Config", func() {
	var (
		options Config

		validOptions = Config{
			Protocol: "http",
			Address:  "fake-host",
			Port:     25777,
			Username: "fake-username",
			Password: "fake-password",
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

		It("returns error if Protocol is empty", func() {
			options.Protocol = ""

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a valid Protocol"))
		})

		It("returns error if Protocol is not valid", func() {
			options.Protocol = "protocol"

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a valid Protocol"))
		})

		It("returns error if Address is empty", func() {
			options.Address = ""

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Address"))
		})

		It("returns error if Port is empty", func() {
			options.Port = 0

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Port"))
		})

		It("returns error if Username is empty", func() {
			options.Username = ""

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Username"))
		})

		It("returns error if Password is empty", func() {
			options.Password = ""

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Password"))
		})
	})
})

var _ = Describe("TLSConfig", func() {
	var (
		options Config

		validOptions = Config{
			Protocol: "https",
			Address:  "fake-host",
			Port:     5555,
			Username: "fake-username",
			Password: "fake-password",
			TLS: TLSConfig{
				CertFile: "fake-certificate",
				KeyFile:  "fake-key",
			},
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

		It("returns error if CertFile is empty", func() {
			options.TLS.CertFile = ""

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty CertFile"))
		})

		It("returns error if KeyFile is empty", func() {
			options.TLS.KeyFile = ""

			err := options.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty KeyFile"))
		})
	})
})
