package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//get array of managers and their endpoints
func (r *Redfish) GetManagers() ([]string, error) {
	var mgrs OData
	var result = make([]string, 0)

	if r.AuthToken == nil || *r.AuthToken == "" {
		return result, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	response, err := r.httpRequest(r.Managers, "GET", nil, nil, false)
	if err != nil {
		return result, err
	}

	raw := response.Content
	if response.StatusCode != http.StatusOK {
		return result, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &mgrs)
	if err != nil {
		return result, err
	}

	if len(mgrs.Members) == 0 {
		return result, errors.New(fmt.Sprintf("BUG: Missing or empty Members attribute in Managers"))
	}

	for _, m := range mgrs.Members {
		result = append(result, *m.Id)
	}
	return result, nil
}

// get manager data for a particular account
func (r *Redfish) GetManagerData(managerEndpoint string) (*ManagerData, error) {
	var result ManagerData

	if r.AuthToken == nil || *r.AuthToken == "" {
		return nil, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	response, err := r.httpRequest(managerEndpoint, "GET", nil, nil, false)
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
	result.SelfEndpoint = &managerEndpoint
	return &result, nil
}

// map ID -> manager data
func (r *Redfish) MapManagersById() (map[string]*ManagerData, error) {
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
		if m.Id == nil {
			return result, errors.New(fmt.Sprintf("BUG: No Id found or Id is null in JSON data from %s", mgr))
		}
		result[*m.Id] = m
	}

	return result, nil
}

// map UUID -> manager data
func (r *Redfish) MapManagersByUuid() (map[string]*ManagerData, error) {
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
			return result, errors.New(fmt.Sprintf("BUG: No UUID found or UUID is null in JSON data from %s", mgr))
		}
		result[*m.UUID] = m
	}

	return result, nil
}
