package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	config_util "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"

	"github.com/prometheus/prometheus/storage/remote"
)

func initHTTPTransport(caFile string, keyFile, certFile string, insecure bool) (http.RoundTripper, error) {
	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}
	var tr http.RoundTripper = &http.Transport{}

	if insecure {
		tlsConfig.InsecureSkipVerify = insecure
	}

	caCertPool := x509.NewCertPool()
	if caFile != "" {
		caCert, err := os.ReadFile(filepath.Clean(caFile))
		if err != nil {
			return tr, err
		}
		caCertPool.AppendCertsFromPEM(caCert)
	}

	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return tr, err
		}

		tlsConfig.RootCAs = caCertPool
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	tr = &http.Transport{TLSClientConfig: tlsConfig}

	return tr, nil
}

func initRemoteWriteClient(baseURL *url.URL, timeout time.Duration, roundTripper http.RoundTripper, username, password string, headers map[string]string) (*remote.Client, error) {
	addressURL, err := url.Parse(baseURL.String())
	if err != nil {
		return nil, err
	}

	// build remote write client
	writeClient, err := remote.NewWriteClient("remote-write", &remote.ClientConfig{
		URL:              &config_util.URL{URL: addressURL},
		Timeout:          model.Duration(timeout),
		RetryOnRateLimit: true,
	})
	if err != nil {
		return nil, err
	}

	// set custom tls config from httpConfigFilePath
	// set custom headers to every request
	client, ok := writeClient.(*remote.Client)
	if !ok {
		return nil, err
	}

	if username != "" && password != "" {
		headers["Authorization"] = fmt.Sprintf("Basic %s", basicAuth(username, password))
	}

	client.Client.Transport = &setHeadersTransport{
		RoundTripper: roundTripper,
		headers:      headers,
	}

	return client, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

type setHeadersTransport struct {
	http.RoundTripper
	headers map[string]string
}

func (s *setHeadersTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range s.headers {
		req.Header.Set(key, value)
	}
	return s.RoundTripper.RoundTrip(req)
}
