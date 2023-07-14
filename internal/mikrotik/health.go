package mikrotik

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Health struct {
	Voltage     float64
	Temperature float64
}

type response struct {
	Id    string `json:".id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (c client) GetHealth() (Health, error) {
	resp, err := c.get("/system/health")
	if err != nil {
		return Health{}, err
	}

	if resp.StatusCode != 200 {
		errorMessage := fmt.Sprintf("received invalid status code: %d", resp.StatusCode)
		return Health{}, errors.New(errorMessage)
	}

	defer resp.Body.Close()
	var response []response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Health{}, err
	}

	var health Health

	for _, r := range response {
		switch r.Name {
		case "voltage":
			health.Voltage, err = strconv.ParseFloat(r.Value, 64)
			if err != nil {
				return Health{}, err
			}
		case "temperature":
			health.Temperature, err = strconv.ParseFloat(r.Value, 64)
			if err != nil {
				return Health{}, err
			}
		}
	}
	return health, nil
}
