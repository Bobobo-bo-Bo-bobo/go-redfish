package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GetManagers - get array of managers and their endpoints
func (r *Redfish) GetManagers() ([]string, error) {
	var mgrs OData
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
			"path":               r.Managers,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting user accounts")
	}
	response, err := r.httpRequest(r.Managers, "GET", nil, nil, false)
	if err != nil {
		return result, err
	}

	raw := response.Content
	if response.StatusCode != http.StatusOK {
		return result, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &mgrs)
	if err != nil {
		return result, err
	}

	if len(mgrs.Members) == 0 {
		return result, fmt.Errorf("BUG: Missing or empty Members attribute in Managers")
	}

	for _, m := range mgrs.Members {
		result = append(result, *m.ID)
	}
	return result, nil
}

// GetManagerData - get manager data for an particular account
func (r *Redfish) GetManagerData(managerEndpoint string) (*ManagerData, error) {
	var result ManagerData

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
			"path":               managerEndpoint,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting information for user")
	}
	response, err := r.httpRequest(managerEndpoint, "GET", nil, nil, false)
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
	result.SelfEndpoint = &managerEndpoint
	return &result, nil
}

// MapManagersByID - map ID -> manager data
func (r *Redfish) MapManagersByID() (map[string]*ManagerData, error) {
	var result = make(map[string]*ManagerData)

	ml, err := r.GetManagers()
	if err != nil {
		return result, err
	}

	for _, mgr := range ml {
		m, err := r.GetManagerData(mgr)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if m.ID == nil {
			return result, fmt.Errorf("BUG: No Id found or Id is null in JSON data from %s", mgr)
		}
		result[*m.ID] = m
	}

	return result, nil
}

// MapManagersByUUID - map UUID -> manager data
func (r *Redfish) MapManagersByUUID() (map[string]*ManagerData, error) {
	var result = make(map[string]*ManagerData)

	ml, err := r.GetManagers()
	if err != nil {
		return result, err
	}

	for _, mgr := range ml {
		m, err := r.GetManagerData(mgr)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if m.UUID == nil {
			return result, fmt.Errorf("BUG: No UUID found or UUID is null in JSON data from %s", mgr)
		}
		result[*m.UUID] = m
	}

	return result, nil
}
