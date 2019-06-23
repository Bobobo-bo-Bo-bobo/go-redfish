package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Initialise Redfish basic data
func (r *Redfish) Initialise() error {
	var base baseEndpoint

	if r.Debug || r.Verbose {
		// Logging setup
		var log_fmt *log.TextFormatter = new(log.TextFormatter)
		log_fmt.FullTimestamp = true
		log_fmt.TimestampFormat = time.RFC3339
		log.SetFormatter(log_fmt)
	}

	if r.Debug {
		log.SetLevel(log.DebugLevel)
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               "/redfish/v1/",
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Rquesting basic information")
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
	// some systems don't have the (mandatory!) AccountService endpoint (e.g. LENOVO)
	if base.AccountService.Id != nil {
		r.AccountService = *base.AccountService.Id
	}

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
