package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Login to SessionEndpoint and get authentication token for this session
func (r *Redfish) Login() error {
	var sessions sessionServiceEndpoint

	if r.Username == "" || r.Password == "" {
		return errors.New(fmt.Sprintf("ERROR: Both Username and Password must be set"))
	}

	response, err := r.httpRequest(r.SessionService, "GET", nil, nil, true)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("ERROR: HTTP POST for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
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

	jsonPayload := fmt.Sprintf("{ \"UserName\":\"%s\",\"Password\":\"%s\" }", r.Username, r.Password)
	response, err = r.httpRequest(*sessions.Sessions.Id, "POST", nil, strings.NewReader(jsonPayload), false)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("ERROR: HTTP POST for %s returned \"%s\" instead of \"200 OK\" or \"201 Created\"", response.Url, response.Status))
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
