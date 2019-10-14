package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (r *Redfish) getImportCertTargetHP(mgr *ManagerData) (string, error) {
	var certTarget string
	var oemHp ManagerDataOemHp
	var secsvc string
	var oemSSvc SecurityServiceDataOemHp
	var httpscertloc string
	var httpscert HTTPSCertDataOemHp

	// parse Oem section from JSON
	err := json.Unmarshal(mgr.Oem, &oemHp)
	if err != nil {
		return certTarget, err
	}

	// get SecurityService endpoint from .Oem.Hp.links.SecurityService
	if oemHp.Hp.Links.SecurityService.ID == nil {
		return certTarget, errors.New("BUG: .Hp.Links.SecurityService.Id not found or null")
	}
	secsvc = *oemHp.Hp.Links.SecurityService.ID

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
		return certTarget, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &oemSSvc)
	if err != nil {
		return certTarget, err
	}

	if oemSSvc.Links.HTTPSCert.ID == nil {
		return certTarget, fmt.Errorf("BUG: .links.HttpsCert.Id not present or is null in data from %s", response.URL)
	}

	httpscertloc = *oemSSvc.Links.HTTPSCert.ID
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
		return certTarget, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &httpscert)
	if err != nil {
		return certTarget, err
	}

	if httpscert.Actions.ImportCertificate.Target == nil {
		return certTarget, fmt.Errorf("BUG: .Actions.ImportCertificate.Target is not present or empty in JSON data from %s", response.URL)
	}

	certTarget = *httpscert.Actions.ImportCertificate.Target
	return certTarget, nil
}

func (r *Redfish) getImportCertTargetHPE(mgr *ManagerData) (string, error) {
	var certTarget string
	var oemHpe ManagerDataOemHpe
	var secsvc string
	var oemSSvc SecurityServiceDataOemHpe
	var httpscertloc string
	var httpscert HTTPSCertDataOemHpe

	// parse Oem section from JSON
	err := json.Unmarshal(mgr.Oem, &oemHpe)
	if err != nil {
		return certTarget, err
	}

	// get SecurityService endpoint from .Oem.Hp.links.SecurityService
	if oemHpe.Hpe.Links.SecurityService.ID == nil {
		return certTarget, errors.New("BUG: .Hpe.Links.SecurityService.Id not found or null")
	}
	secsvc = *oemHpe.Hpe.Links.SecurityService.ID

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
		return certTarget, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &oemSSvc)
	if err != nil {
		return certTarget, err
	}

	if oemSSvc.Links.HTTPSCert.ID == nil {
		return certTarget, fmt.Errorf("BUG: .links.HttpsCert.Id not present or is null in data from %s", response.URL)
	}

	httpscertloc = *oemSSvc.Links.HTTPSCert.ID
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
		return certTarget, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &httpscert)
	if err != nil {
		return certTarget, err
	}

	if httpscert.Actions.ImportCertificate.Target == nil {
		return certTarget, fmt.Errorf("BUG: .Actions.ImportCertificate.Target is not present or empty in JSON data from %s", response.URL)
	}

	certTarget = *httpscert.Actions.ImportCertificate.Target
	return certTarget, nil
}

func (r *Redfish) getImportCertTargetHuawei(mgr *ManagerData) (string, error) {
	var certTarget string
	var oemHuawei ManagerDataOemHuawei
	var secsvc string
	var oemSSvc SecurityServiceDataOemHuawei
	var httpscertloc string
	var httpscert HTTPSCertDataOemHuawei

	// parse Oem section from JSON
	err := json.Unmarshal(mgr.Oem, &oemHuawei)
	if err != nil {
		return certTarget, err
	}

	// get SecurityService endpoint from .Oem.Huawei.links.SecurityService
	if oemHuawei.Huawei.SecurityService.ID == nil {
		return certTarget, errors.New("BUG: .Huawei.SecurityService.Id not found or null")
	}
	secsvc = *oemHuawei.Huawei.SecurityService.ID

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
		return certTarget, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &oemSSvc)
	if err != nil {
		return certTarget, err
	}

	if oemSSvc.Links.HTTPSCert.ID == nil {
		return certTarget, fmt.Errorf("BUG: .links.HttpsCert.Id not present or is null in data from %s", response.URL)
	}

	httpscertloc = *oemSSvc.Links.HTTPSCert.ID

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
		return certTarget, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &httpscert)
	if err != nil {
		return certTarget, err
	}

	if httpscert.Actions.ImportCertificate.Target == nil {
		return certTarget, fmt.Errorf("BUG: .Actions.ImportCertificate.Target is not present or empty in JSON data from %s", response.URL)
	}

	certTarget = *httpscert.Actions.ImportCertificate.Target
	return certTarget, nil
}

// ImportCertificate - import certificate
func (r *Redfish) ImportCertificate(cert string) error {
	var certtarget string

	if r.AuthToken == nil || *r.AuthToken == "" {
		return errors.New("No authentication token found, is the session setup correctly?")
	}

	// set vendor flavor
	err := r.GetVendorFlavor()
	if err != nil {
		return err
	}

	// get list of Manager endpoint
	mgrList, err := r.GetManagers()
	if err != nil {
		return err
	}

	// pick the first entry
	mgr0, err := r.GetManagerData(mgrList[0])
	if err != nil {
		return err
	}

	// get endpoint SecurityService from Managers
	if r.Flavor == RedfishHP {
		certtarget, err = r.getImportCertTargetHP(mgr0)
		if err != nil {
			return err
		}

		// HP/HPE service processors (iLO) will reboot automatically
		// if the certificate has been imported successfully
	} else if r.Flavor == RedfishHPE {
		certtarget, err = r.getImportCertTargetHPE(mgr0)
		if err != nil {
			return err
		}

		// HP/HPE service processors (iLO) will reboot automatically
		// if the certificate has been imported successfully
	} else if r.Flavor == RedfishHuawei {
		certtarget, err = r.getImportCertTargetHuawei(mgr0)
		if err != nil {
			return err
		}

		// Reboot service processor to activate new certificate
		err = r.ResetSP()
		if err != nil {
			return err
		}
	} else if r.Flavor == RedfishInspur {
		return errors.New("Inspur management boards do not support certificate import")
	} else if r.Flavor == RedfishSuperMicro {
		return errors.New("SuperMicro management boards do not support certificate import")
	} else {
		return errors.New("Unable to get vendor for management board. If this vendor supports certificate import please file a feature request")
	}

	if certtarget == "" {
		return errors.New("BUG: Target for certificate import is not known")
	}

	// escape new lines
	rawcert := strings.Replace(cert, "\n", "\\n", -1)
	certPayload := fmt.Sprintf("{ \"Certificate\": \"%s\" }", rawcert)

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
			"payload":            certPayload,
		}).Debug("Importing SSL certificate")
	}
	response, err := r.httpRequest(certtarget, "POST", nil, strings.NewReader(certPayload), false)
	if err != nil {
		return err
	}
	// XXX: do we need to look at the content returned by HTTP POST ?

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP POST to %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	return nil
}
