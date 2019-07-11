package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Initialise Redfish basic data
func (r *Redfish) Initialise() error {
	var base baseEndpoint
	var raw []byte

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
		}).Info("Requesting basic information")
	}
	response, err := r.httpRequest("/redfish/v1/", "GET", nil, nil, false)
	if err != nil {
		return err
	}

	raw = response.Content
	r.RawBaseContent = string(raw)

	// Some managementboards (e.g. IBM/Lenovo) will redirect to a different webserver running on a different port.
	// To avoid futher problems for non-GET methods we will parse the new location and set the port accordig to the
	// Location header.
	if response.StatusCode == http.StatusMovedPermanently || response.StatusCode == http.StatusFound {

		location := response.Header.Get("Location")

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
				"status_code":        response.StatusCode,
				"status":             response.Status,
				"location":           location,
			}).Info("HTTP request redirected by the server")
		}

		// Note: Although RFC 2616 specify "The new permanent URI SHOULD be given by the Location field in the response."
		//       we will barf because we have no way to obtain the redirect URL.
		if location == "" {
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
				"status_code":        response.StatusCode,
				"status":             response.Status,
			}).Fatal("HTTP request returned 3xx status but didn't set new Location header")
		}

		new_loc, err := url.Parse(location)
		if err != nil {
			return err
		}

		// XXX: We don't expect to be redirected to a new _server_, protocol (HTTPS is mandatory
		//      for Redfish) or path (/redfish/v1 is the mandatory path for Redfish API accesS),
		//      so we choose to ignore everything else except the port.
		_host, _port, _ := net.SplitHostPort(new_loc.Host)
		if _port != "" {
			new_port, err := net.LookupPort("tcp", _port)
			if err != nil {
				return err
			}

			r.Port = new_port
			if r.Verbose {
				log.WithFields(log.Fields{
					"hostname":      r.Hostname,
					"port":          r.Port,
					"timeout":       r.Timeout,
					"flavor":        r.Flavor,
					"flavor_string": r.FlavorString,
				}).Info("Port configuration has been updated")
			}
		}

		// At least warn if the redirect points to another host when verbosity is requested
		if r.Verbose {
			new_host := strings.ToLower(_host)
			if _host == "" {
				new_host = strings.ToLower(new_loc.Host)
			}

			if new_host != strings.ToLower(r.Hostname) {
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
					"status_code":        response.StatusCode,
					"status":             response.Status,
					"location":           location,
				}).Warning("Ignoring redirect to new server as indicated by the Location header sent by the server")
			}
		}

		// Re-request base information from new location
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
			}).Info("Rerequesting basic information")
		}
		response, err := r.httpRequest("/redfish/v1/", "GET", nil, nil, false)
		if err != nil {
			return err
		}

		raw = response.Content
		r.RawBaseContent = string(raw)
	} else if response.StatusCode != http.StatusOK {
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
