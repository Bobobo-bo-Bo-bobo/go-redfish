package redfish

import (
	"encoding/json"
	"strings"
)

// ProcessError - processess error response from API
func (r *Redfish) ProcessError(response HTTPResult) (*Error, error) {
	var rerr Error
	var err error

	err = json.Unmarshal(response.Content, &rerr)
	if err != nil {
		return nil, err
	}

	return &rerr, nil
}

// GetErrorMessage - get error messages from error structure
func (r *Redfish) GetErrorMessage(rerr *Error) string {
	var result string
	var _list = make([]string, 0)

	if rerr.Error.Code != nil {
		// According to the API specificiation the error object can hold multiple entries (see https://redfish.dmtf.org/schemas/DSP0266_1.0.html#error-responses).
		for _, extinfo := range rerr.Error.MessageExtendedInfo {
			// On failure some vendors, like HP/HPE, don't set any Message, only MessageId. If there is no Message we return MessageId and hope for the best.
			if extinfo.Message != nil {
				_list = append(_list, *extinfo.Message)
			} else if extinfo.MessageID != nil {
				_list = append(_list, *extinfo.MessageID)
			}
		}
	} else {
		return result
	}

	result = strings.Join(_list, "; ")
	return result
}
