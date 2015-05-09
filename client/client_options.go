package registry

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
)

// ClientOptions are the options used to create a BOSH Registry client.
type ClientOptions struct {
	// BOSH Registry protocol
	Protocol string `json:"protocol,omitempty"`

	// BOSH Registry hostname
	Host string `json:"host,omitempty"`

	// BOSH Registry port
	Port int `json:"port,omitempty"`

	// BOSH Registry username
	Username string `json:"username,omitempty"`

	// BOSH Registry password
	Password string `json:"password,omitempty"`

	// BOSH Registry TLS options (only when using protocol https)
	TLS TLSConfig `json:"tls,omitempty"`
}

// TLSConfig are the TLS options used to create a BOSH Registry client.
type TLSConfig struct {
	// If the Client must skip the verification of the server certificates
	InsecureSkipVerify bool `json:"insecure_skip_verify,omitempty"`

	// Certificate file (PEM format)
	CertFile string `json:"certfile,omitempty"`

	// Private key file (PEM format)
	KeyFile string `json:"keyfile,omitempty"`

	// Roor CA certificate file (PEM format)
	CACertFile string `json:"cacertfile,omitempty"`
}

// Validate validates the Client options.
func (o ClientOptions) Validate() error {
	if o.Protocol == "" {
		return bosherr.Error("Must provide a non-empty Protocol")
	}

	if o.Host == "" {
		return bosherr.Error("Must provide a non-empty Host")
	}

	if o.Port == 0 {
		return bosherr.Error("Must provide a non-empty Port")
	}

	if o.Username == "" {
		return bosherr.Error("Must provide a non-empty Username")
	}

	if o.Password == "" {
		return bosherr.Error("Must provide a non-empty Password")
	}

	if o.Protocol == "https" {
		err := o.TLS.Validate()
		if err != nil {
			return bosherr.WrapError(err, "Validating TLS configuration")
		}
	}

	return nil
}

// Validate validates the TLS options.
func (o TLSConfig) Validate() error {
	if o.CertFile == "" {
		return bosherr.Error("Must provide a non-empty CertFile")
	}

	if o.KeyFile == "" {
		return bosherr.Error("Must provide a non-empty KeyFile")
	}

	return nil
}