package server_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-registry/server"

	"github.com/frodenas/bosh-registry/server/store/fakes"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
)

var _ = Describe("InstanceHandler", func() {
	var (
		err              error
		logger           boshlog.Logger
		responseRecorder *httptest.ResponseRecorder
		request          *http.Request
		registryStore    *fakes.FakeRegistryStore
		instanceHandler  *InstanceHandler

		config = Config{
			Protocol: "http",
			Address:  "fake-host",
			Port:     5555,
			Username: "fake-username",
			Password: "fake-password",
		}
	)

	BeforeEach(func() {
		registryStore = &fakes.FakeRegistryStore{}
		logger = boshlog.NewLogger(boshlog.LevelNone)
		instanceHandler = NewInstanceHandler(config, registryStore, logger)
	})

	Describe("HandleFunc", func() {
		BeforeEach(func() {
			responseRecorder = httptest.NewRecorder()
		})

		It("returns a Not Found error if path does not contain instanceID", func() {
			request, err = http.NewRequest("GET", "/instances", nil)
			Expect(err).NotTo(HaveOccurred())

			instanceHandler.HandleFunc(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
			Expect(responseRecorder.Body.String()).To(ContainSubstring("not_found"))
		})
	})

	Describe("HandleGet", func() {
		BeforeEach(func() {
			responseRecorder = httptest.NewRecorder()
		})

		It("returns the instance settings", func() {
			registryStore.GetFound = true
			registryStore.GetValue = "fake-instance-settings"

			request, err = http.NewRequest("GET", "/instances/fake-instance-id/settings", nil)
			Expect(err).NotTo(HaveOccurred())

			instanceHandler.HandleFunc(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			Expect(responseRecorder.Body.String()).To(ContainSubstring("fake-instance-settings"))
			Expect(registryStore.GetCalled).To(BeTrue())
		})

		It("returns a Not Found error if instance settings have not been found", func() {
			request, err = http.NewRequest("GET", "/instances/fake-instance-id/settings", nil)
			Expect(err).NotTo(HaveOccurred())

			instanceHandler.HandleFunc(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
			Expect(responseRecorder.Body.String()).To(ContainSubstring("not_found"))
			Expect(registryStore.GetCalled).To(BeTrue())
		})

		It("returns a Bad request error if registry store returns an error", func() {
			registryStore.GetErr = errors.New("fake-registry-store-error")

			request, err = http.NewRequest("GET", "/instances/fake-instance-id/settings", nil)
			Expect(err).NotTo(HaveOccurred())

			instanceHandler.HandleFunc(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			Expect(responseRecorder.Body.String()).To(ContainSubstring("error"))
			Expect(registryStore.GetCalled).To(BeTrue())
		})
	})

	Describe("HandlePut", func() {
		BeforeEach(func() {
			responseRecorder = httptest.NewRecorder()
		})

		It("returns an OK status if instance settings have been saved", func() {
			request, err = http.NewRequest("PUT", "/instances/fake-instance-id/settings", bytes.NewReader([]byte("fake-instance-settings")))
			request.SetBasicAuth("fake-username", "fake-password")
			Expect(err).NotTo(HaveOccurred())

			instanceHandler.HandleFunc(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			Expect(registryStore.SaveCalled).To(BeTrue())
		})

		It("returns a Bad request error if registry store returns an error", func() {
			registryStore.SaveErr = errors.New("fake-registry-store-error")

			request, err = http.NewRequest("PUT", "/instances/fake-instance-id/settings", bytes.NewReader([]byte("fake-instance-settings")))
			request.SetBasicAuth("fake-username", "fake-password")
			Expect(err).NotTo(HaveOccurred())

			instanceHandler.HandleFunc(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			Expect(responseRecorder.Body.String()).To(ContainSubstring("error"))
			Expect(registryStore.SaveCalled).To(BeTrue())
		})

		It("returns an Unauthorized error if request does not contain credentials", func() {
			request, err = http.NewRequest("PUT", "/instances/fake-instance-id/settings", bytes.NewReader([]byte("fake-instance-settings")))
			Expect(err).NotTo(HaveOccurred())

			instanceHandler.HandleFunc(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
			Expect(responseRecorder.HeaderMap).To(HaveKey("Www-Authenticate"))
		})
	})

	Describe("HandleDelete", func() {
		BeforeEach(func() {
			responseRecorder = httptest.NewRecorder()
		})

		It("returns an OK status if instance settings have been deleted", func() {
			registryStore.GetFound = true
			registryStore.GetValue = "fake-instance-settings"

			request, err = http.NewRequest("DELETE", "/instances/fake-instance-id/settings", nil)
			request.SetBasicAuth("fake-username", "fake-password")
			Expect(err).NotTo(HaveOccurred())

			instanceHandler.HandleFunc(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			Expect(registryStore.DeleteCalled).To(BeTrue())
		})

		It("returns a Bad request error if registry store returns an error", func() {
			registryStore.DeleteErr = errors.New("fake-registry-store-error")

			request, err = http.NewRequest("DELETE", "/instances/fake-instance-id/settings", nil)
			request.SetBasicAuth("fake-username", "fake-password")
			Expect(err).NotTo(HaveOccurred())

			instanceHandler.HandleFunc(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			Expect(responseRecorder.Body.String()).To(ContainSubstring("error"))
			Expect(registryStore.DeleteCalled).To(BeTrue())
		})

		It("returns an Unauthorized error if request does not contain credentials", func() {
			request, err = http.NewRequest("DELETE", "/instances/fake-instance-id/settings", nil)
			Expect(err).NotTo(HaveOccurred())

			instanceHandler.HandleFunc(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
			Expect(responseRecorder.HeaderMap).To(HaveKey("Www-Authenticate"))
		})
	})
})
