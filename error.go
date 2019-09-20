package redfish

import (
	"encoding/json"
)

func (r *Redfish) ProcessError(response HttpResult) (*RedfishError, error) {
	var rerr RedfishError
	var err error

	err = json.Unmarshal(response.Content, &rerr)
	if err != nil {
		return nil, err
	}

	return &rerr, nil
}
