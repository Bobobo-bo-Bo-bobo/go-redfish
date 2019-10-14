package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GetChassis - get array of chassis and their endpoints
func (r *Redfish) GetChassis() ([]string, error) {
	var chassis OData
	var result = make([]string, 0)

	if r.AuthToken == nil || *r.AuthToken == "" {
		return result, errors.New("No authentication token found, is the session setup correctly?")
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               r.Chassis,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting installed chassis components")
	}
	response, err := r.httpRequest(r.Chassis, "GET", nil, nil, false)
	if err != nil {
		return result, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return result, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &chassis)
	if err != nil {
		return result, err
	}

	if len(chassis.Members) == 0 {
		return result, errors.New("BUG: Array of chassis endpoints is empty")
	}

	for _, endpoint := range chassis.Members {
		result = append(result, *endpoint.ID)
	}
	return result, nil
}

// GetChassisData - get chassis data for a particular chassis
func (r *Redfish) GetChassisData(chassisEndpoint string) (*ChassisData, error) {
	var result ChassisData

	if r.AuthToken == nil || *r.AuthToken == "" {
		return nil, errors.New("No authentication token found, is the session setup correctly?")
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               chassisEndpoint,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting chassis information")
	}
	response, err := r.httpRequest(chassisEndpoint, "GET", nil, nil, false)
	if err != nil {
		return nil, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &result)
	if err != nil {
		return nil, err
	}

	result.SelfEndpoint = &chassisEndpoint
	return &result, nil
}

// MapChassisByID - Map chassis by ID
func (r *Redfish) MapChassisByID() (map[string]*ChassisData, error) {
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
		if s.ID == nil {
			return result, fmt.Errorf("BUG: No Id found for Chassis at %s", chas)
		}

		result[*s.ID] = s
	}

	return result, nil
}

// GetPowerData - get power data from endpoint
func (r *Redfish) GetPowerData(powerEndpoint string) (*PowerData, error) {
	var result PowerData

	if r.AuthToken == nil || *r.AuthToken == "" {
		return nil, errors.New("No authentication token found, is the session setup correctly?")
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               powerEndpoint,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting power information")
	}
	response, err := r.httpRequest(powerEndpoint, "GET", nil, nil, false)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(response.Content, &result)
	if err != nil {
		return nil, err
	}

	result.SelfEndpoint = &powerEndpoint
	return &result, nil
}

// GetThermalData - get thermal data from endpoint
func (r *Redfish) GetThermalData(thermalEndpoint string) (*ThermalData, error) {
	var result ThermalData

	if r.AuthToken == nil || *r.AuthToken == "" {
		return nil, errors.New("No authentication token found, is the session setup correctly?")
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               thermalEndpoint,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting thermal information")
	}
	response, err := r.httpRequest(thermalEndpoint, "GET", nil, nil, false)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(response.Content, &result)
	if err != nil {
		return nil, err
	}

	result.SelfEndpoint = &thermalEndpoint
	return &result, nil
}
