package mikrotik

import (
	"fmt"
	"net/http"
	"time"
)

type Configuration struct {
	Timeout  float64
	Address  string
	Port     int
	Username string
	Password string
}

type Client interface {
	GetHealth() (Health, error)
	GetResource() (Resource, error)
}

type client struct {
	configuration Configuration
	httpClient    http.Client
}

func NewClient(configuration Configuration) Client {
	return &client{
		configuration: configuration,
		httpClient: http.Client{
			Timeout: time.Duration(configuration.Timeout) * time.Second,
		},
	}
}

func (c client) get(path string) (*http.Response, error) {
	url := c.buildURL(path)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(c.configuration.Username, c.configuration.Password)

	return c.httpClient.Do(request)
}

func (c client) buildURL(path string) string {
	return fmt.Sprintf("http://%s/rest%s", c.configuration.Address, path)
}
