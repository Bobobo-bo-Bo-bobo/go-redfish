package redfish

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// Login - Login to SessionEndpoint and get authentication token for this session
func (r *Redfish) Login() error {
	var sessions sessionServiceEndpoint

	if r.Username == "" || r.Password == "" {
		return fmt.Errorf("Both Username and Password must be set")
	}

	// Get session endpoint if not already defined by information from base endpoint .Links.Sessions
	// because some implementations (e.g. INSPUR) report SessionService endpoint but don't implement it.
	if r.Sessions == "" {
		if r.Verbose {
			log.WithFields(log.Fields{
				"hostname":           r.Hostname,
				"port":               r.Port,
				"timeout":            r.Timeout,
				"flavor":             r.Flavor,
				"flavor_string":      r.FlavorString,
				"path":               r.SessionService,
				"method":             "GET",
				"additional_headers": nil,
				"use_basic_auth":     true,
			}).Info("Requesting path to session service")
		}
		response, err := r.httpRequest(r.SessionService, "GET", nil, nil, true)
		if err != nil {
			return err
		}

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
		}

		raw := response.Content

		err = json.Unmarshal(raw, &sessions)
		if err != nil {
			return err
		}

		// check if management boards reports "ServiceEnabled" and if it does, check if is true
		if sessions.Enabled != nil {
			if !*sessions.Enabled {
				return fmt.Errorf("Session information from %s reports session service as disabled", response.URL)
			}
		}

		if sessions.Sessions == nil {
			return fmt.Errorf("BUG: No Sessions endpoint reported from %s", response.URL)
		}

		if sessions.Sessions.ID == nil {
			return fmt.Errorf("BUG: Malformed Sessions endpoint reported from %s: no @odata.id field found", response.URL)
		}

		r.Sessions = *sessions.Sessions.ID
	}

	jsonPayload := fmt.Sprintf("{ \"UserName\":\"%s\",\"Password\":\"%s\" }", r.Username, r.Password)
	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               r.Sessions,
			"method":             "POST",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Sending login data to session service")
	}
	if r.Debug {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               r.Sessions,
			"method":             "POST",
			"additional_headers": nil,
			"use_basic_auth":     false,
			"payload":            jsonPayload,
		}).Debug("Sending login data to session service")
	}
	response, err := r.httpRequest(r.Sessions, "POST", nil, strings.NewReader(jsonPayload), false)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		redfishError, err := r.ProcessError(response)
		if err != nil {
			return fmt.Errorf("HTTP POST for %s returned \"%s\" instead of \"200 OK\" or \"201 Created\"", response.URL, response.Status)
		}
		msg := r.GetErrorMessage(redfishError)
		if msg != "" {
			return fmt.Errorf("Login failed: %s", msg)
		}
		return fmt.Errorf("HTTP POST for %s returned \"%s\" instead of \"200 OK\" or \"201 Created\"", response.URL, response.Status)
	}

	token := response.Header.Get("x-auth-token")
	if token == "" {
		return fmt.Errorf("BUG: HTTP POST to SessionService endpoint %s returns OK but no X-Auth-Token in reply", response.URL)
	}
	r.AuthToken = &token

	session := response.Header.Get("location")
	if session == "" {
		return fmt.Errorf("BUG: HTTP POST to SessionService endpoint %s returns OK but has no Location in reply", response.URL)
	}

	// check if is a full URL
	if session[0] == '/' {
		if r.Port > 0 {
			session = fmt.Sprintf("https://%s:%d%s", r.Hostname, r.Port, session)
		} else {
			session = fmt.Sprintf("https://%s%s", r.Hostname, session)
		}
	}
	r.SessionLocation = &session

	return nil
}
