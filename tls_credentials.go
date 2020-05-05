package orlop

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/switch-bit/orlop/log"
	"google.golang.org/grpc/credentials"
)

// NewServerTLSConfig returns a new tls.VaultConfig from the given configuration input
func NewServerTLSConfig(cfg HasTLSConfig, vault HasVaultConfig) (*tls.Config, error) {
	config := &tls.Config{
		ClientAuth: cfg.GetClientAuth(),
		MinVersion: tls.VersionTLS12,
	}

	if !strSliceContains(config.NextProtos, "http/1.1") {
		config.NextProtos = append(config.NextProtos, "http/1.1")
	}

	if !cfg.GetEnabled() {
		log.Trace("tls disabled")
		return config, nil
	}

	if !strSliceContains(config.NextProtos, "h2") {
		// Enable HTTP/2
		config.NextProtos = append([]string{"h2"}, config.NextProtos...)
	}

	var err error
	t := CloneTLSConfig(cfg)
	t.Cert.Secret, t.Key.Secret, err = GenerateCertificates(t.Generate, vault)
	if err != nil {
		return nil, err
	}

	certPEMBlock, err := LoadKey(t.GetCert(), vault, "certificate")
	if err == nil {
		log.Trace("certificate loaded")

		keyPEMBlock, err := LoadKey(t.GetKey(), vault, "private_key")
		if err != nil {
			log.WithError(err).Error("error loading private key")
			return nil, err
		}

		if t.GetRootCA().GetEnabled() {
			rootcaPEMBlock, err := LoadKey(t.GetRootCA(), vault, "issuing_ca")
			if err == nil {
				config.ClientCAs = x509.NewCertPool()

				if !config.ClientCAs.AppendCertsFromPEM(rootcaPEMBlock) {
					log.WithError(err).Error("failed to append client CA certificates")
					return nil, fmt.Errorf("tls: failed to append client CA certificates")
				}
			} else {
				log.WithError(err).Error("error loading CA")
			}
		}

		c, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
		if err != nil {
			log.WithError(err).Error("error creating key pair")
			return nil, err
		}

		config.Certificates = append(config.Certificates, c)
	} else {
		log.WithError(err).Error("error loading certificate")
	}

	return config, nil
}

// NewClientTLSConfig returns a new tls.VaultConfig from the given configuration input
func NewClientTLSConfig(cfg HasTLSConfig, vault HasVaultConfig) (*tls.Config, error) {
	if !cfg.GetEnabled() {
		log.Trace("tls disabled")
		return &tls.Config{}, nil
	}

	config := &tls.Config{
		ServerName:         cfg.GetOverride(),
		InsecureSkipVerify: cfg.GetInsecure(),
		MinVersion:         tls.VersionTLS12,
	}

	var err error
	t := CloneTLSConfig(cfg)
	t.Cert.Secret, t.Key.Secret, err = GenerateCertificates(t.Generate, vault)
	if err != nil {
		return nil, err
	}

	certPEMBlock, err := LoadKey(t.GetCert(), vault, "certificate")
	if err == nil {
		keyPEMBlock, err := LoadKey(t.GetKey(), vault, "private_key")
		if err != nil {
			log.WithError(err).Error("error loading private key")
			return nil, err
		}

		config.RootCAs = x509.NewCertPool()

		if !config.RootCAs.AppendCertsFromPEM(certPEMBlock) {
			log.WithError(err).Error("failed to append certificates")
			return nil, fmt.Errorf("tls: failed to append certificates")
		}

		if t.GetRootCA().GetEnabled() {
			rootcaPEMBlock, err := LoadKey(t.GetRootCA(), vault, "issuing_ca")
			if err == nil {
				if !config.RootCAs.AppendCertsFromPEM(rootcaPEMBlock) {
					log.WithError(err).Error("failed to append root certificates")
					return nil, fmt.Errorf("tls: failed to append root certificates")
				}
			} else {
				log.WithError(err).Error("error loading CA")
			}
		}

		c, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
		if err != nil {
			log.WithError(err).Error("error creating key pair")
			return nil, err
		}

		config.Certificates = append(config.Certificates, c)
	} else {
		log.WithError(err).Error("error loading certificate")
	}

	return config, nil
}

// NewClientTLSCredentials returns new TLS credentials based on the provided configuration
func NewClientTLSCredentials(cfg HasTLSConfig, vault HasVaultConfig) (credentials.TransportCredentials, error) {
	t, err := NewClientTLSConfig(cfg, vault)
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(t), nil
}

// NewServerTLSCredentials returns new TLS credentials based on the provided configuration
func NewServerTLSCredentials(cfg HasTLSConfig, vault HasVaultConfig) (credentials.TransportCredentials, error) {
	t, err := NewServerTLSConfig(cfg, vault)
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(t), nil
}

// GenerateCertificates generates certificates by accessing the configured endpoint in vault
func GenerateCertificates(g HasCertGenerationConfig, vault HasVaultConfig) (cert string, key string, err error) {
	// If Vault not enabled or certificate generation not enabled, just return
	if !vault.GetEnabled() || !g.GetEnabled() {
		return
	}

	// Connect to Vault
	client, err := NewVault(vault)
	if err != nil {
		return "", "", err
	}

	params := map[string]interface{}{
		"common_name": g.GetCommonName(),
		"format":      "pem_bundle",
	}
	if len(g.GetAltNames()) > 0 {
		params["alt_names"] = g.GetAltNames()
	}
	if g.GetTTL().Seconds() > 60 {
		params["ttl"] = g.GetTTL().String()
	}

	// Write the params to the path to generate the certificate
	secret, err := client.Write(g.GetPath(), params)
	if err != nil {
		return "", "", err
	}

	// Set the generated certificate and private key as secrets
	cert = decodeCertInfo(secret.Data, "certificate")
	key = decodeCertInfo(secret.Data, "private_key")

	return
}

func decodeCertInfo(data map[string]interface{}, key string) string {
	if d, ok := data[key]; ok {
		if s, ok := d.(string); ok {
			return base64.StdEncoding.EncodeToString([]byte(s))
		}
	}

	return ""
}
