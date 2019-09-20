package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// Login to SessionEndpoint and get authentication token for this session
func (r *Redfish) Login() error {
	var sessions sessionServiceEndpoint

	if r.Username == "" || r.Password == "" {
		return errors.New(fmt.Sprintf("ERROR: Both Username and Password must be set"))
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
			return errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
		}

		raw := response.Content

		err = json.Unmarshal(raw, &sessions)
		if err != nil {
			return err
		}

		// check if management boards reports "ServiceEnabled" and if it does, check if is true
		if sessions.Enabled != nil {
			if !*sessions.Enabled {
				return errors.New(fmt.Sprintf("ERROR: Session information from %s reports session service as disabled\n", response.Url))
			}
		}

		if sessions.Sessions == nil {
			return errors.New(fmt.Sprintf("BUG: No Sessions endpoint reported from %s\n", response.Url))
		}

		if sessions.Sessions.Id == nil {
			return errors.New(fmt.Sprintf("BUG: Malformed Sessions endpoint reported from %s: no @odata.id field found\n", response.Url))
		}

		r.Sessions = *sessions.Sessions.Id
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
		redfish_error, err := r.ProcessError(response)
		if err != nil {
			return errors.New(fmt.Sprintf("ERROR: HTTP POST for %s returned \"%s\" instead of \"200 OK\" or \"201 Created\"", response.Url, response.Status))
		}
		if redfish_error.Error.Code != nil {
			// According to the API specificiation the error object can hold multiple entries (see https://redfish.dmtf.org/schemas/DSP0266_1.0.html#error-responses).
			// We will always pick the first and hope it provides suitable information.
			if len(redfish_error.Error.MessageExtendedInfo) > 0 {
				// On authentication failure some vendors like HP/HPE don't set any Message, only MessageId. If there is no Message we return MessageId and hope for the best.
				if redfish_error.Error.MessageExtendedInfo[0].Message != nil {
					return errors.New(fmt.Sprintf("ERROR: Login failed: %s\n", *redfish_error.Error.MessageExtendedInfo[0].Message))
				} else {
					if redfish_error.Error.MessageExtendedInfo[0].MessageId != nil {
						return errors.New(fmt.Sprintf("ERROR: Login failed: %s\n", *redfish_error.Error.MessageExtendedInfo[0].MessageId))
					}
					return errors.New(fmt.Sprintf("ERROR: HTTP POST for %s returned \"%s\" instead of \"200 OK\" or \"201 Created\"", response.Url, response.Status))
				}
			} else {
				return errors.New(fmt.Sprintf("ERROR: HTTP POST for %s returned \"%s\" instead of \"200 OK\" or \"201 Created\"", response.Url, response.Status))
			}
		} else {
			return errors.New(fmt.Sprintf("ERROR: HTTP POST for %s returned \"%s\" instead of \"200 OK\" or \"201 Created\"", response.Url, response.Status))
		}
	}

	token := response.Header.Get("x-auth-token")
	if token == "" {
		return errors.New(fmt.Sprintf("BUG: HTTP POST to SessionService endpoint %s returns OK but no X-Auth-Token in reply", response.Url))
	}
	r.AuthToken = &token

	session := response.Header.Get("location")
	if session == "" {
		return errors.New(fmt.Sprintf("BUG: HTTP POST to SessionService endpoint %s returns OK but has no Location in reply", response.Url))
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
