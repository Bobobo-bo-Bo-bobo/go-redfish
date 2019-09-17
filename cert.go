package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (r *Redfish) getImportCertTarget_HP(mgr *ManagerData) (string, error) {
	var certTarget string
	var oemHp ManagerDataOemHp
	var secsvc string
	var oemSSvc SecurityServiceDataOemHp
	var httpscertloc string
	var httpscert HttpsCertDataOemHp

	// parse Oem section from JSON
	err := json.Unmarshal(mgr.Oem, &oemHp)
	if err != nil {
		return certTarget, err
	}

	// Newer systems, e.g. iLO5+ use Oem.Hpe instead of Oem.Hp as part of the HP/HPE split on november, 1st 2015
	// XXX: Is there a more elegant solution to handle Oem.Hp and Oem.Hpe which essentially provide the same information
	//      (at least for the certificate handling) ?

	// NOTE: Hp and Hpe are mutually exclusive !
	if oemHp.Hp == nil && oemHp.Hpe == nil {
		return certTarget, errors.New("BUG: Neither .Oem.Hp nor .Oem.Hpe are found")
	}
	if oemHp.Hp != nil && oemHp.Hpe != nil {
		return certTarget, errors.New("BUG: Both .Oem.Hp and .Oem.Hpe are found")
	}

	// Point .Hpe to .Hp and continue processing
	if oemHp.Hpe != nil {
		oemHp.Hp = oemHp.Hpe
		oemHp.Hpe = nil
	}

	// get SecurityService endpoint from .Oem.Hp.links.SecurityService
	if oemHp.Hp.Links.SecurityService.Id == nil {
		return certTarget, errors.New("BUG: .Hp.Links.SecurityService.Id not found or null")
	} else {
		secsvc = *oemHp.Hp.Links.SecurityService.Id
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               secsvc,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting path to security service")
	}
	response, err := r.httpRequest(secsvc, "GET", nil, nil, false)
	if err != nil {
		return certTarget, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return certTarget, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &oemSSvc)
	if err != nil {
		return certTarget, err
	}

	if oemSSvc.Links.HttpsCert.Id == nil {
		return certTarget, errors.New(fmt.Sprintf("BUG: .links.HttpsCert.Id not present or is null in data from %s", response.Url))
	}

	httpscertloc = *oemSSvc.Links.HttpsCert.Id
	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               httpscertloc,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting path for SSL certificate import")
	}
	response, err = r.httpRequest(httpscertloc, "GET", nil, nil, false)
	if err != nil {
		return certTarget, err
	}

	raw = response.Content

	if response.StatusCode != http.StatusOK {
		return certTarget, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &httpscert)
	if err != nil {
		return certTarget, err
	}

	if httpscert.Actions.ImportCertificate.Target == nil {
		return certTarget, errors.New(fmt.Sprintf("BUG: .Actions.ImportCertificate.Target is not present or empty in JSON data from %s", response.Url))
	}

	certTarget = *httpscert.Actions.ImportCertificate.Target
	return certTarget, nil
}

func (r *Redfish) getImportCertTarget_Huawei(mgr *ManagerData) (string, error) {
	var certTarget string
	var oemHuawei ManagerDataOemHuawei
	var secsvc string
	var oemSSvc SecurityServiceDataOemHuawei
	var httpscertloc string
	var httpscert HttpsCertDataOemHuawei

	// parse Oem section from JSON
	err := json.Unmarshal(mgr.Oem, &oemHuawei)
	if err != nil {
		return certTarget, err
	}

	// get SecurityService endpoint from .Oem.Huawei.links.SecurityService
	if oemHuawei.Huawei.SecurityService.Id == nil {
		return certTarget, errors.New("BUG: .Huawei.SecurityService.Id not found or null")
	} else {
		secsvc = *oemHuawei.Huawei.SecurityService.Id
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               secsvc,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting path to security service")
	}
	response, err := r.httpRequest(secsvc, "GET", nil, nil, false)
	if err != nil {
		return certTarget, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return certTarget, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &oemSSvc)
	if err != nil {
		return certTarget, err
	}

	if oemSSvc.Links.HttpsCert.Id == nil {
		return certTarget, errors.New(fmt.Sprintf("BUG: .links.HttpsCert.Id not present or is null in data from %s", response.Url))
	}

	httpscertloc = *oemSSvc.Links.HttpsCert.Id

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               httpscertloc,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting path for SSL certificate import")
	}
	response, err = r.httpRequest(httpscertloc, "GET", nil, nil, false)
	if err != nil {
		return certTarget, err
	}

	raw = response.Content

	if response.StatusCode != http.StatusOK {
		return certTarget, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &httpscert)
	if err != nil {
		return certTarget, err
	}

	if httpscert.Actions.ImportCertificate.Target == nil {
		return certTarget, errors.New(fmt.Sprintf("BUG: .Actions.ImportCertificate.Target is not present or empty in JSON data from %s", response.Url))
	}

	certTarget = *httpscert.Actions.ImportCertificate.Target
	return certTarget, nil
}

func (r *Redfish) ImportCertificate(cert string) error {
	var certtarget string = ""

	if r.AuthToken == nil || *r.AuthToken == "" {
		return errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	// set vendor flavor
	err := r.GetVendorFlavor()
	if err != nil {
		return err
	}

	// get list of Manager endpoint
	mgr_list, err := r.GetManagers()
	if err != nil {
		return err
	}

	// pick the first entry
	mgr0, err := r.GetManagerData(mgr_list[0])
	if err != nil {
		return err
	}

	// get endpoint SecurityService from Managers
	if r.Flavor == REDFISH_HP {
		certtarget, err = r.getImportCertTarget_HP(mgr0)
		if err != nil {
			return err
		}

		// HP/HPE service processors (iLO) will reboot automatically
		// if the certificate has been imported successfully
	} else if r.Flavor == REDFISH_HUAWEI {
		certtarget, err = r.getImportCertTarget_Huawei(mgr0)
		if err != nil {
			return err
		}

		// Reboot service processor to activate new certificate
		err = r.ResetSP()
		if err != nil {
			return err
		}
	} else if r.Flavor == REDFISH_INSPUR {
		return errors.New("ERROR: Inspur management boards do not support certificate import")
	} else if r.Flavor == REDFISH_SUPERMICRO {
		return errors.New("ERROR: SuperMicro management boards do not support certificate import")
	} else {
		return errors.New("ERROR: Unable to get vendor for management board. If this vendor supports certificate import please file a feature request")
	}

	if certtarget == "" {
		return errors.New("BUG: Target for certificate import is not known")
	}

	// escape new lines
	rawcert := strings.Replace(cert, "\n", "\\n", -1)
	cert_payload := fmt.Sprintf("{ \"Certificate\": \"%s\" }", rawcert)

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               certtarget,
			"method":             "POST",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Importing SSL certificate")
	}
	if r.Debug {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               certtarget,
			"method":             "POST",
			"additional_headers": nil,
			"use_basic_auth":     false,
			"payload":            cert_payload,
		}).Debug("Importing SSL certificate")
	}
	response, err := r.httpRequest(certtarget, "POST", nil, strings.NewReader(cert_payload), false)
	if err != nil {
		return err
	}
	// XXX: do we need to look at the content returned by HTTP POST ?

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("ERROR: HTTP POST to %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	return nil
}
