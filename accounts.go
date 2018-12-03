package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//get array of accounts and their endpoints
func (r *Redfish) GetAccounts() ([]string, error) {
	var accsvc AccountService
	var accs OData
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

	if accsvc.AccountsEndpoint == nil {
		return result, errors.New("BUG: No Accounts endpoint found")
	}

	response, err = r.httpRequest(*accsvc.AccountsEndpoint.Id, "GET", nil, nil, false)
	if err != nil {
		return result, err
	}

	raw = response.Content
	if response.StatusCode != http.StatusOK {
		return result, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &accs)
	if err != nil {
		return result, err
	}

	if len(accs.Members) == 0 {
		return result, errors.New(fmt.Sprintf("BUG: Missing or empty Members attribute in Accounts"))
	}

	for _, a := range accs.Members {
		result = append(result, *a.Id)
	}
	return result, nil
}

// get account data for a particular account
func (r *Redfish) GetAccountData(accountEndpoint string) (*AccountData, error) {
	var result AccountData

	if r.AuthToken == nil || *r.AuthToken == "" {
		return nil, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	response, err := r.httpRequest(accountEndpoint, "GET", nil, nil, false)
	if err != nil {
		return nil, err
	}

	// store unparsed content
	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &result)
	if err != nil {
		return nil, err
	}
	result.SelfEndpoint = &accountEndpoint
	return &result, nil
}

// map username -> user data
func (r *Redfish) MapAccountsByName() (map[string]*AccountData, error) {
	var result = make(map[string]*AccountData)

	al, err := r.GetAccounts()
	if err != nil {
		return result, err
	}

	for _, acc := range al {
		a, err := r.GetAccountData(acc)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if a.UserName == nil {
			return result, errors.New("BUG: No UserName found or UserName is null")
		}
		result[*a.UserName] = a
	}

	return result, nil
}

// map ID -> user data
func (r *Redfish) MapAccountsById() (map[string]*AccountData, error) {
	var result = make(map[string]*AccountData)

	al, err := r.GetAccounts()
	if err != nil {
		return result, err
	}

	for _, acc := range al {
		a, err := r.GetAccountData(acc)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if a.Id == nil {
			return result, errors.New("BUG: No Id found or Id is null")
		}
		result[*a.Id] = a
	}

	return result, nil
}
