package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//get array of roles and their endpoints
func (r *Redfish) GetRoles() ([]string, error) {
	var accsvc AccountService
	var roles OData
	var result = make([]string, 0)

	if r.AuthToken == nil || *r.AuthToken == "" {
		return result, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	response, err := r.httpRequest(r.AccountService, "GET", nil, nil, false)
	if err != nil {
		return result, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return result, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &accsvc)
	if err != nil {
		return result, err
	}

	// Some managementboards (e.g. HPE iLO) don't use roles but an internal ("Oem") privilege map instead
	if accsvc.RolesEndpoint == nil {
		return result, nil
	}

	response, err = r.httpRequest(*accsvc.RolesEndpoint.Id, "GET", nil, nil, false)
	if err != nil {
		return result, err
	}
	raw = response.Content

	if response.StatusCode != http.StatusOK {
		return result, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &roles)
	if err != nil {
		return result, err
	}

	if len(roles.Members) == 0 {
		return result, errors.New(fmt.Sprintf("BUG: Missing or empty Members attribute in Roles"))
	}

	for _, r := range roles.Members {
		result = append(result, *r.Id)
	}
	return result, nil
}

// get role data for a particular role
func (r *Redfish) GetRoleData(roleEndpoint string) (*RoleData, error) {
	var result RoleData

	if r.AuthToken == nil || *r.AuthToken == "" {
		return nil, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	response, err := r.httpRequest(roleEndpoint, "GET", nil, nil, false)
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

	result.SelfEndpoint = &roleEndpoint
	return &result, nil
}

// map roles by name
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
			return result, errors.New("ERROR: No Name found or Name is null")
		}

		result[*rl.Name] = rl
	}

	return result, nil
}

// map roles by ID
func (r *Redfish) MapRolesById() (map[string]*RoleData, error) {
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
		if rl.Id == nil {
			return result, errors.New("ERROR: No Id found or Id is null")
		}

		result[*rl.Id] = rl
	}

	return result, nil
}
