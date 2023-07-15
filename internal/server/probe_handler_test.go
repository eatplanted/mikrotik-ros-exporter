package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/eatplanted/mikrotik-ros-exporter/internal/config"
)

func TestPrometheusTimeoutHTTPHeader(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
	}))
	defer testServer.Close()

	parsedUrl, err := url.Parse(testServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("/probe?target=%s&credential=default", parsedUrl.Host)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("X-Prometheus-Scrape-Timeout-Seconds", "2")

	recorder := httptest.NewRecorder()
	server := NewServer(config.Configuration{
		Credentials: map[string]config.Credential{
			"default": {},
		},
	})

	server.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("probe request handler returned wrong status code: %v, want %v", status, http.StatusOK)
	}
}

func TestTimeoutIsSetCorrectly(t *testing.T) {
	var testSuite = []struct {
		inConfigurationTimeout    float64
		inPrometheusScrapeTimeout string
		outTimeout                float64
	}{
		{0, "15", 14.5},
		{20, "15", 14.5},
		{5, "15", 5},
		{10, "", 10},
		{0, "", 119.5},
	}

	for _, test := range testSuite {
		request, _ := http.NewRequest("GET", "", nil)
		request.Header.Set("X-Prometheus-Scrape-Timeout-Seconds", test.inPrometheusScrapeTimeout)
		configuration := config.Configuration{
			Timeout: test.inConfigurationTimeout,
		}

		timeout, err := getTimeout(request, configuration)
		if err != nil {
			t.Error(err)
		}

		if timeout != test.outTimeout {
			t.Errorf("timeout is incorrect: %v, want %v", timeout, test.outTimeout)
		}
	}
}
