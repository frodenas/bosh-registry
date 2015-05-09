package server

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
)

const listenerLogTag = "RegistryServerListener"

type Listener struct {
	config   Config
	handler  *InstanceHandler
	logger   boshlog.Logger
	listener net.Listener
}

func NewListener(
	config Config,
	handler *InstanceHandler,
	logger boshlog.Logger,
) Listener {
	return Listener{
		config:  config,
		handler: handler,
		logger:  logger,
	}
}

func (l *Listener) ListenAndServe() <-chan error {
	errChan := make(chan error, 1)

	tcpListener, err := net.ListenTCP(
		"tcp",
		&net.TCPAddr{
			IP:   net.ParseIP(l.config.Address),
			Port: l.config.Port,
		},
	)
	if err != nil {
		errChan <- bosherr.WrapError(err, "Starting Registry TCP Listener")
		return errChan
	}

	if l.config.Protocol == "https" {
		certificates, err := tls.LoadX509KeyPair(l.config.TLS.CertFile, l.config.TLS.KeyFile)
		if err != nil {
			errChan <- bosherr.WrapError(err, "Loading X509 Key Pair")
			return errChan
		}

		certPool := x509.NewCertPool()
		if l.config.TLS.CACertFile != "" {
			caCert, err := ioutil.ReadFile(l.config.TLS.CACertFile)
			if err != nil {
				errChan <- bosherr.WrapError(err, "Loading CA certificate")
				return errChan
			}

			if !certPool.AppendCertsFromPEM(caCert) {
				errChan <- bosherr.WrapError(err, "Invalid CA Certificate")
				return errChan
			}
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{certificates},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    certPool,
			MinVersion:   tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			},
			PreferServerCipherSuites: true,
			SessionTicketsDisabled:   true,
		}

		l.listener = tls.NewListener(tcpListener, tlsConfig)
	} else {
		l.listener = tcpListener
	}

	httpServer := http.Server{}
	mux := http.NewServeMux()
	mux.HandleFunc("/instances/", l.handler.HandleFunc)
	httpServer.Handler = mux

	l.logger.Debug(listenerLogTag, "Starting Registry Server at %s://%s:%d", l.config.Protocol, l.config.Address, l.config.Port)
	go func() {
		err := httpServer.Serve(l.listener)
		errChan <- err
	}()

	return errChan
}

func (l *Listener) Stop() {
	l.logger.Debug(listenerLogTag, "Stopping Registry Server")
	l.listener.Close()
}
