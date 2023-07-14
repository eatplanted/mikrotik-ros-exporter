package mikrotik

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Resource struct {
	CpuCount             float64 `json:"cpu-count,string"`
	CpuFrequency         float64 `json:"cpu-frequency,string"`
	CpuLoad              float64 `json:"cpu-load,string"`
	FreeHddSpace         float64 `json:"free-hdd-space,string"`
	FreeMemory           float64 `json:"free-memory,string"`
	TotalHddSpace        float64 `json:"total-hdd-space,string"`
	TotalMemory          float64 `json:"total-memory,string"`
	WriteSectSinceReboot float64 `json:"write-sect-since-reboot,string"`
	WriteSectTotal       float64 `json:"write-sect-total,string"`
}

func (c client) GetResource() (Resource, error) {
	resp, err := c.get("/system/resource")
	if err != nil {
		return Resource{}, err
	}

	if resp.StatusCode != 200 {
		errorMessage := fmt.Sprintf("received invalid status code: %d", resp.StatusCode)
		return Resource{}, errors.New(errorMessage)
	}

	defer resp.Body.Close()

	var resource Resource
	if err := json.NewDecoder(resp.Body).Decode(&resource); err != nil {
		return Resource{}, err
	}

	return resource, nil
}
