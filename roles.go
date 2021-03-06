package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GetRoles - get array of roles and their endpoints
func (r *Redfish) GetRoles() ([]string, error) {
	var accsvc AccountService
	var roles OData
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
			"path":               r.AccountService,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting path for account roles")
	}
	response, err := r.httpRequest(r.AccountService, "GET", nil, nil, false)
	if err != nil {
		return result, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return result, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &accsvc)
	if err != nil {
		return result, err
	}

	// Some managementboards (e.g. HPE iLO) don't use roles but an internal ("Oem") privilege map instead
	if accsvc.RolesEndpoint == nil {
		return result, nil
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               *accsvc.RolesEndpoint.ID,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting account roles")
	}
	response, err = r.httpRequest(*accsvc.RolesEndpoint.ID, "GET", nil, nil, false)
	if err != nil {
		return result, err
	}
	raw = response.Content

	if response.StatusCode != http.StatusOK {
		return result, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &roles)
	if err != nil {
		return result, err
	}

	if len(roles.Members) == 0 {
		return result, fmt.Errorf("BUG: Missing or empty Members attribute in Roles")
	}

	for _, r := range roles.Members {
		result = append(result, *r.ID)
	}
	return result, nil
}

// GetRoleData - get role data for a particular role
func (r *Redfish) GetRoleData(roleEndpoint string) (*RoleData, error) {
	var result RoleData

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
			"path":               roleEndpoint,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting role information")
	}
	response, err := r.httpRequest(roleEndpoint, "GET", nil, nil, false)
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

	result.SelfEndpoint = &roleEndpoint
	return &result, nil
}

// MapRolesByName - map roles by name
func (r *Redfish) MapRolesByName() (map[string]*RoleData, error) {
	var result = make(map[string]*RoleData)

	rll, err := r.GetRoles()
	if err != nil {
		return result, err
	}

	for _, ro := range rll {
		rl, err := r.GetRoleData(ro)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if rl.Name == nil {
			return result, errors.New("No Name found or Name is null")
		}

		result[*rl.Name] = rl
	}

	return result, nil
}

// MapRolesByID - map roles by ID
func (r *Redfish) MapRolesByID() (map[string]*RoleData, error) {
	var result = make(map[string]*RoleData)

	rll, err := r.GetRoles()
	if err != nil {
		return result, err
	}

	for _, ro := range rll {
		rl, err := r.GetRoleData(ro)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if rl.ID == nil {
			return result, errors.New("No Id found or Id is null")
		}

		result[*rl.ID] = rl
	}

	return result, nil
}
