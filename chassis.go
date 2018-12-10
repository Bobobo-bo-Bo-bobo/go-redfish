package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//get array of chassis and their endpoints
func (r *Redfish) GetChassis() ([]string, error) {
	var chassis OData
	var result = make([]string, 0)

	if r.AuthToken == nil || *r.AuthToken == "" {
		return result, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	response, err := r.httpRequest(r.Chassis, "GET", nil, nil, false)
	if err != nil {
		return result, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return result, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &chassis)
	if err != nil {
		return result, err
	}

	if len(chassis.Members) == 0 {
		return result, errors.New("BUG: Array of chassis endpoints is empty")
	}

	for _, endpoint := range chassis.Members {
		result = append(result, *endpoint.Id)
	}
	return result, nil
}

// get chassis data for a particular chassis
func (r *Redfish) GetChassisData(chassisEndpoint string) (*ChassisData, error) {
	var result ChassisData

	if r.AuthToken == nil || *r.AuthToken == "" {
		return nil, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	response, err := r.httpRequest(chassisEndpoint, "GET", nil, nil, false)
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

	result.SelfEndpoint = &chassisEndpoint
	return &result, nil
}

// Map chassis by ID
func (r *Redfish) MapChassisById() (map[string]*ChassisData, error) {
	var result = make(map[string]*ChassisData)

	chasl, err := r.GetChassis()
	if err != nil {
		return result, nil
	}

	for _, chas := range chasl {
		s, err := r.GetChassisData(chas)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if s.Id == nil {
			return result, errors.New(fmt.Sprintf("BUG: No Id found for Chassis at %s", chas))
		}

		result[*s.Id] = s
	}

	return result, nil
}

// get Power data from
func (r *Redfish) GetPowerData(powerEndpoint string) (*PowerData, error) {
	var result PowerData

	response, err := r.httpRequest(powerEndpoint, "GET", nil, nil, false)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(response.Content, &result)
	if err != nil {
		return nil, err
	}

	result.SelfEndpoint = &powerEndpoint
	return &result, nil
}

// get Thermal data from
func (r *Redfish) GetThermalData(thermalEndpoint string) (*ThermalData, error) {
	var result ThermalData

	response, err := r.httpRequest(thermalEndpoint, "GET", nil, nil, false)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(response.Content, &result)
	if err != nil {
		return nil, err
	}

	result.SelfEndpoint = &thermalEndpoint
	return &result, nil
}
