package redfish

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Logout - Logout from SessionEndpoint and delete authentication token for this session
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
		return fmt.Errorf("BUG: X-Auth-Token set (value: %s) but no SessionLocation for this session found", *r.AuthToken)
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               *r.SessionLocation,
			"method":             "DELETE",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Removing session authentication")
	}
	response, err := r.httpRequest(*r.SessionLocation, "DELETE", nil, nil, false)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP DELETE for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	r.AuthToken = nil
	r.SessionLocation = nil

	return nil
}
