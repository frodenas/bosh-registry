package registry_test

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-registry/client"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	"github.com/frodenas/bosh-registry/server/listener/fakes"
)

var _ = Describe("Client", func() {
	var (
		err error

		instanceHandler *fakes.FakeInstanceHandler
		mux             *http.ServeMux
		registryServer  *httptest.Server

		options        ClientOptions
		logger         boshlog.Logger
		registryClient Client

		instanceID           string
		expectedAgentSet     AgentSettings
		expectedAgentSetJSON []byte
	)

	BeforeEach(func() {
		logger = boshlog.NewLogger(boshlog.LevelNone)
		instanceHandler = fakes.NewFakeInstanceHandler("fake-username", "fake-password")
		mux = http.NewServeMux()
		mux.HandleFunc("/", instanceHandler.HandleFunc)

		instanceID = "fake-instance-id"
		expectedAgentSet = AgentSettings{AgentID: "fake-agent-id"}
		expectedAgentSetJSON, err = json.Marshal(expectedAgentSet)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("when using http", func() {
		BeforeEach(func() {
			registryServer = httptest.NewServer(mux)
			serverURL, err := url.Parse(registryServer.URL)
			Expect(err).ToNot(HaveOccurred())
			serverHost, serverPortString, err := net.SplitHostPort(serverURL.Host)
			Expect(err).ToNot(HaveOccurred())
			serverPort, err := strconv.ParseInt(serverPortString, 10, 64)
			Expect(err).ToNot(HaveOccurred())

			options = ClientOptions{
				Protocol: serverURL.Scheme,
				Host:     serverHost,
				Port:     int(serverPort),
				Username: "fake-username",
				Password: "fake-password",
			}
			registryClient = NewClient(options, logger)
		})

		AfterEach(func() {
			registryServer.Close()
		})

		Describe("Delete", func() {
			Context("when settings for the instance exist in the registry", func() {
				BeforeEach(func() {
					instanceHandler.InstanceSettings = expectedAgentSetJSON
				})

				It("deletes settings in the registry", func() {
					err = registryClient.Delete(instanceID)
					Expect(err).ToNot(HaveOccurred())
					Expect(instanceHandler.InstanceSettings).To(Equal([]byte{}))
				})
			})

			Context("when settings for instance does not exist", func() {
				It("should not return an error", func() {
					Expect(instanceHandler.InstanceSettings).To(Equal([]byte{}))
					err = registryClient.Delete(instanceID)
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})

		Describe("Endpoint", func() {
			It("returns the BOSH Registry endpoint", func() {
				endpoint := registryClient.Endpoint()
				Expect(endpoint).To(Equal(fmt.Sprintf("%s://%s:%d", options.Protocol, options.Host, options.Port)))
			})
		})

		Describe("EndpointWithCredentials", func() {
			It("returns the BOSH Registry endpoint with credentials", func() {
				endpoint := registryClient.EndpointWithCredentials()
				Expect(endpoint).To(Equal(fmt.Sprintf("%s://%s:%s@%s:%d", options.Protocol, options.Username, options.Password, options.Host, options.Port)))
			})
		})

		Describe("Fetch", func() {
			Context("when settings for the instance exist in the registry", func() {
				BeforeEach(func() {
					instanceHandler.InstanceSettings = expectedAgentSetJSON
				})

				It("fetches settings from the registry", func() {
					agentSet, err := registryClient.Fetch(instanceID)
					Expect(err).ToNot(HaveOccurred())
					Expect(agentSet).To(Equal(expectedAgentSet))
				})
			})

			Context("when settings for instance does not exist", func() {
				It("returns an error", func() {
					Expect(instanceHandler.InstanceSettings).To(Equal([]byte{}))
					agentSet, err := registryClient.Fetch(instanceID)
					Expect(err).To(HaveOccurred())
					Expect(agentSet).To(Equal(AgentSettings{}))
				})
			})
		})

		Describe("Update", func() {
			It("updates settings in the registry", func() {
				Expect(instanceHandler.InstanceSettings).To(Equal([]byte{}))
				err := registryClient.Update(instanceID, expectedAgentSet)
				Expect(err).ToNot(HaveOccurred())
				Expect(instanceHandler.InstanceSettings).To(Equal(expectedAgentSetJSON))
			})
		})
	})

	Context("when using https", func() {
		BeforeEach(func() {
			registryServer = httptest.NewTLSServer(mux)
			serverURL, err := url.Parse(registryServer.URL)
			Expect(err).ToNot(HaveOccurred())
			serverHost, serverPortString, err := net.SplitHostPort(serverURL.Host)
			Expect(err).ToNot(HaveOccurred())
			serverPort, err := strconv.ParseInt(serverPortString, 10, 64)
			Expect(err).ToNot(HaveOccurred())

			options = ClientOptions{
				Protocol: serverURL.Scheme,
				Host:     serverHost,
				Port:     int(serverPort),
				Username: "fake-username",
				Password: "fake-password",
				TLS: ClientTLSOptions{
					InsecureSkipVerify: true,
					CertFile:           "../test/assets/public.pem",
					KeyFile:            "../test/assets/private.pem",
					CACertFile:         "../test/assets/ca.pem",
				},
			}
			registryClient = NewClient(options, logger)
		})

		AfterEach(func() {
			registryServer.Close()
		})

		Describe("Delete", func() {
			It("deletes settings in the registry", func() {
				instanceHandler.InstanceSettings = expectedAgentSetJSON
				err = registryClient.Delete(instanceID)
				Expect(err).ToNot(HaveOccurred())
				Expect(instanceHandler.InstanceSettings).To(Equal([]byte{}))
			})
		})

		Describe("Endpoint", func() {
			It("returns the BOSH Registry endpoint", func() {
				endpoint := registryClient.Endpoint()
				Expect(endpoint).To(Equal(fmt.Sprintf("%s://%s:%d", options.Protocol, options.Host, options.Port)))
			})
		})

		Describe("EndpointWithCredentials", func() {
			It("returns the BOSH Registry endpoint with credentials", func() {
				endpoint := registryClient.EndpointWithCredentials()
				Expect(endpoint).To(Equal(fmt.Sprintf("%s://%s:%s@%s:%d", options.Protocol, options.Username, options.Password, options.Host, options.Port)))
			})
		})

		Describe("Fetch", func() {
			It("fetches settings from the registry", func() {
				instanceHandler.InstanceSettings = expectedAgentSetJSON
				agentSet, err := registryClient.Fetch(instanceID)
				Expect(err).ToNot(HaveOccurred())
				Expect(agentSet).To(Equal(expectedAgentSet))
			})
		})

		Describe("Update", func() {
			It("updates settings in the registry", func() {
				Expect(instanceHandler.InstanceSettings).To(Equal([]byte{}))
				err := registryClient.Update(instanceID, expectedAgentSet)
				Expect(err).ToNot(HaveOccurred())
				Expect(instanceHandler.InstanceSettings).To(Equal(expectedAgentSetJSON))
			})
		})
	})
})
