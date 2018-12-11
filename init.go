package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Initialise Redfish basic data
func (r *Redfish) Initialise() error {
	var base baseEndpoint

	if r.Debug {
		r.logger = log.New(os.Stderr, "", log.Lshortfile|log.Ldate|log.Ltime|log.Lmicroseconds)
	}

	response, err := r.httpRequest("/redfish/v1/", "GET", nil, nil, false)
	if err != nil {
		return err
	}

	raw := response.Content
	r.RawBaseContent = string(raw)

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &base)
	if err != nil {
		return err
	}

	// extract required endpoints
	if base.AccountService.Id == nil {
		return errors.New(fmt.Sprintf("BUG: No AccountService endpoint found in base configuration from %s", response.Url))
	}
	r.AccountService = *base.AccountService.Id

	if base.Chassis.Id == nil {
		return errors.New(fmt.Sprintf("BUG: No Chassis endpoint found in base configuration from %s", response.Url))
	}
	r.Chassis = *base.Chassis.Id

	if base.Managers.Id == nil {
		return errors.New(fmt.Sprintf("BUG: No Managers endpoint found in base configuration from %s", response.Url))
	}
	r.Managers = *base.Managers.Id

	if base.SessionService.Id == nil {
		return errors.New(fmt.Sprintf("BUG: No SessionService endpoint found in base configuration from %s", response.Url))
	}
	r.SessionService = *base.SessionService.Id

	// Get session endpoint from .Links.Sessions because some implementations (e.g. INSPUR) report SessionService endpoint but don't implement it.
	if base.Links.Sessions != nil {
		if base.Links.Sessions.Id != nil && *base.Links.Sessions.Id != "" {
			r.Sessions = *base.Links.Sessions.Id
		}
	}

	if base.Systems.Id == nil {
		return errors.New(fmt.Sprintf("BUG: No Systems endpoint found in base configuration from %s", response.Url))
	}
	r.Systems = *base.Systems.Id

	return nil
}
