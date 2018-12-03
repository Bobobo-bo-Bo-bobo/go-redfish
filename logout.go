package redfish

import (
	"errors"
	"fmt"
	"net/http"
)

// Logout from SessionEndpoint and delete authentication token for this session
func (r *Redfish) Logout() error {

	if r.AuthToken == nil {
		// do nothing for Logout when we don't even have an authentication token
		return nil
	}
	if *r.AuthToken == "" {
		// do nothing for Logout when we don't even have an authentication token
		return nil
	}

	if r.SessionLocation == nil || *r.SessionLocation == "" {
		return errors.New(fmt.Sprintf("BUG: X-Auth-Token set (value: %s) but no SessionLocation for this session found\n", *r.AuthToken))
	}

	response, err := r.httpRequest(*r.SessionLocation, "DELETE", nil, nil, false)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("ERROR: HTTP DELETE for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	r.AuthToken = nil
	r.SessionLocation = nil

	return nil
}
