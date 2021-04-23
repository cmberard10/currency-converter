package http

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

//Service ...
type Service struct {
	HTTPClient *http.Client
	Config     *Config
}

//Config holds the db connection information
type Config struct {
	ClientTimeoutInSeconds                 int `required:"true" split_words:"true" default:"30"`
}

//NewHTTPDAO ....
func NewHTTPDAO(config *Config) *Service {
	s := &Service{
		Config: config,
	}
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,

		ExpectContinueTimeout: 10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS11,
		},
	}

	c := &http.Client{
		Timeout:   time.Second * time.Duration(config.ClientTimeoutInSeconds),
		Transport: transport,
	}

	s.HTTPClient = c
	return s
}
