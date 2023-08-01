package http

import (
	"net/url"
	"strings"
)

var currentConfig *HttpClientConfig

type HttpClientConfig struct {
	InsecureEndpoints []string `config:"insecure_endpoints"`
}

func Init(config *HttpClientConfig) {
	currentConfig = config
}

func IsInsecureEndpoint(endpoint string) (bool, error) {
	if currentConfig == nil {
		return false, nil
	}

	parsedEndpoint, err := url.Parse(endpoint)
	if err != nil {
		return false, err
	}
	for _, insecureEndpoint := range currentConfig.InsecureEndpoints {
		if strings.Contains(parsedEndpoint.Host, insecureEndpoint) {
			return true, nil
		}
	}
	return false, nil
}
