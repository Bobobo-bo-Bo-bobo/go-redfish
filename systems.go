package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

//get array of systems and their endpoints
func (r *Redfish) GetSystems() ([]string, error) {
	var systems OData
	var result = make([]string, 0)

	if r.AuthToken == nil || *r.AuthToken == "" {
		return result, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	response, err := r.httpRequest(r.Systems, "GET", nil, nil, false)
	if err != nil {
		return result, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return result, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &systems)
	if err != nil {
		return result, err
	}

	if len(systems.Members) == 0 {
		return result, errors.New("BUG: Array of system endpoints is empty")
	}

	for _, endpoint := range systems.Members {
		result = append(result, *endpoint.Id)
	}
	return result, nil
}

// get system data for a particular system
func (r *Redfish) GetSystemData(systemEndpoint string) (*SystemData, error) {
	var result SystemData

	if r.AuthToken == nil || *r.AuthToken == "" {
		return nil, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	response, err := r.httpRequest(systemEndpoint, "GET", nil, nil, false)
	if err != nil {
		return nil, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &result)
	if err != nil {
		return nil, err
	}

	result.SelfEndpoint = &systemEndpoint
	return &result, nil
}

// Map systems by ID
func (r *Redfish) MapSystemsById() (map[string]*SystemData, error) {
	var result = make(map[string]*SystemData)

	sysl, err := r.GetSystems()
	if err != nil {
		return result, nil
	}

	for _, sys := range sysl {
		s, err := r.GetSystemData(sys)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if s.Id == nil {
			return result, errors.New(fmt.Sprintf("BUG: No Id found for System at %s", sys))
		}

		result[*s.Id] = s
	}

	return result, nil
}

// Map systems by UUID
func (r *Redfish) MapSystemsByUuid() (map[string]*SystemData, error) {
	var result = make(map[string]*SystemData)

	sysl, err := r.GetSystems()
	if err != nil {
		return result, nil
	}

	for _, sys := range sysl {
		s, err := r.GetSystemData(sys)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if s.UUID == nil {
			return result, errors.New(fmt.Sprintf("BUG: No UUID found for System at %s", sys))
		}

		result[*s.UUID] = s
	}

	return result, nil
}

// Map systems by serial number
func (r *Redfish) MapSystemsBySerialNumber() (map[string]*SystemData, error) {
	var result = make(map[string]*SystemData)

	sysl, err := r.GetSystems()
	if err != nil {
		return result, nil
	}

	for _, sys := range sysl {
		s, err := r.GetSystemData(sys)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if s.SerialNumber == nil {
			return result, errors.New(fmt.Sprintf("BUG: No SerialNumber found for System at %s", sys))
		}

		result[*s.SerialNumber] = s
	}

	return result, nil
}

// get vendor specific "flavor"
func (r *Redfish) GetVendorFlavor() error {
	// get vendor "flavor" for vendor specific implementation details
	_sys, err := r.GetSystems()
	if err != nil {
		return err
	}
	// assuming every system has the same vendor, pick the first one to determine vendor flavor
	_sys0, err := r.GetSystemData(_sys[0])
	if _sys0.Manufacturer != nil {
		_manufacturer := strings.TrimSpace(strings.ToLower(*_sys0.Manufacturer))
		if _manufacturer == "hp" || _manufacturer == "hpe" {
			r.Flavor = REDFISH_HP
            r.FlavorString = "hp"
		} else if _manufacturer == "huawei" {
			r.Flavor = REDFISH_HUAWEI
            r.FlavorString = "huawei"
		} else if _manufacturer == "inspur" {
			r.Flavor = REDFISH_INSPUR
            r.FlavorString = "inspur"
		} else if _manufacturer == "supermicro" {
			r.Flavor = REDFISH_SUPERMICRO
            r.FlavorString = "supermicro"
		} else {
			r.Flavor = REDFISH_GENERAL
            r.FlavorString = "vanilla"
		}
	}

	return nil
}
