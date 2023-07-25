package mikrotik

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Interface struct {
	Id          string  `json:".id"`
	Name        string  `json:"name"`
	Disabled    bool    `json:"disabled,string"`
	Running     bool    `json:"running,string"`
	RxByte      float64 `json:"rx-byte,float64,string"`
	RxDrop      float64 `json:"rx-drop,float64,string"`
	RxError     float64 `json:"rx-error,float64,string"`
	RxPacket    float64 `json:"rx-packet,float64,string"`
	TxByte      float64 `json:"tx-byte,float64,string"`
	TxDrop      float64 `json:"tx-drop,float64,string"`
	TxError     float64 `json:"tx-error,float64,string"`
	TxPacket    float64 `json:"tx-packet,float64,string"`
	TxQueueDrop float64 `json:"tx-queue-drop,float64,string"`
	Type        string  `json:"type"`
}

func (i *Interface) IsActive() bool {
	return i.Running && !i.Disabled
}

func (c *client) GetInterfaces() ([]Interface, error) {
	resp, err := c.get("/interface")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		errorMessage := fmt.Sprintf("received invalid status code: %d", resp.StatusCode)
		return nil, errors.New(errorMessage)
	}

	defer resp.Body.Close()

	var interfaces []Interface
	if err := json.NewDecoder(resp.Body).Decode(&interfaces); err != nil {
		return nil, err
	}

	return interfaces, nil
}
