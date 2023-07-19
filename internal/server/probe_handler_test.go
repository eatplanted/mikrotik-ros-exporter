package server

import (
	"encoding/json"
	"fmt"
	"github.com/eatplanted/mikrotik-ros-exporter/internal/mikrotik"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/eatplanted/mikrotik-ros-exporter/internal/config"
)

func init() {
	logrus.SetOutput(ioutil.Discard)
}

func TestPrometheusTimeoutHTTPHeader(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
	}))
	defer testServer.Close()

	url := fmt.Sprintf("/probe?target=%s&credential=default", testServer.URL)
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

func TestSkipTLSVerifyHTTPHeader_SetTrue(t *testing.T) {
	testServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rest/system/health" {
			json.NewEncoder(w).Encode([]interface{}{
				map[string]interface{}{
					".id":   "*E",
					"name":  "temperature",
					"type":  "C",
					"value": "49",
				},
			})
		}

		if r.URL.Path == "/rest/system/resource" {
			json.NewEncoder(w).Encode(mikrotik.Resource{
				CpuCount: 4,
			})
		}
	}))

	defer testServer.Close()

	url := fmt.Sprintf("/probe?target=%s&credential=default&skip_tls_verify=true", testServer.URL)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	server := NewServer(config.Configuration{
		Credentials: map[string]config.Credential{
			"default": {},
		},
	})

	server.ServeHTTP(recorder, request)

	body, err := io.ReadAll(recorder.Body)
	if !strings.Contains(string(body), "mikrotik_probe_success 1") {
		t.Errorf("probe request handler returned wrong status code: %s, want %s", body, "mikrotik_probe_success 1")
	}

	if !strings.Contains(string(body), "mikrotik_system_health_temperature 49") {
		t.Errorf("probe request handler returned wrong status code: %s, want %s", body, "mikrotik_system_health_temperature 49")
	}

	if !strings.Contains(string(body), "mikrotik_system_resource_cpu_count 4") {
		t.Errorf("probe request handler returned wrong status code: %s, want %s", body, "mikrotik_system_resource_cpu_count 4")
	}
}

func TestSkipTLSVerifyHTTPHeader_SetFalse(t *testing.T) {
	testServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	defer testServer.Close()

	url := fmt.Sprintf("/probe?target=%s&credential=default&skip_tls_verify=false", testServer.URL)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	server := NewServer(config.Configuration{
		Credentials: map[string]config.Credential{
			"default": {},
		},
	})

	server.ServeHTTP(recorder, request)

	body, err := io.ReadAll(recorder.Body)
	if !strings.Contains(string(body), "mikrotik_probe_success 0\n") {
		t.Errorf("probe request handler returned wrong status code: %s, want %s", body, "mikrotik_probe_success 0")
	}
}
