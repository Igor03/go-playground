package httpx

import (
	"net/http"
	"net/url"
	"time"
)

type HttpxClientConfig struct {
	RawBaseUrl string
	Timeout    int
}

type httpxClient struct {
	clientConfig HttpxClientConfig
	baseUrl      *url.URL
	httpClient   *http.Client
}

func validateBaseUrl(hcc *HttpxClientConfig) (*url.URL, error) {
	url, err := url.Parse(hcc.RawBaseUrl)

	if err != nil {
		return nil, err
	}

	return url, nil

}

func (c *httpxClient) addPath(path string) (*url.URL, error) {
	url, err := c.baseUrl.Parse(path)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func New(hcc HttpxClientConfig) (*httpxClient, error) {
	hc := http.Client{
		Timeout: time.Duration(hcc.Timeout) * time.Second,
	}

	purl, err := validateBaseUrl(&hcc)

	if err != nil {
		return nil, err
	}

	return &httpxClient{
		clientConfig: hcc,
		httpClient:   &hc,
		baseUrl:      purl,
	}, nil
}
