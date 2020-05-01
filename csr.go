package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (r *Redfish) fetchCSRHP(mgr *ManagerData) (string, error) {
	var csr string
	var oemHp ManagerDataOemHp
	var secsvc string
	var oemSSvc SecurityServiceDataOemHp
	var httpscertloc string
	var httpscert HTTPSCertDataOemHp

	// parse Oem section from JSON
	err := json.Unmarshal(mgr.Oem, &oemHp)
	if err != nil {
		return csr, err
	}

	// get SecurityService endpoint from .Oem.Hp.links.SecurityService
	if oemHp.Hp.Links.SecurityService.ID == nil {
		return csr, errors.New("BUG: .Hp.Links.SecurityService.Id not found or null")
	}
	secsvc = *oemHp.Hp.Links.SecurityService.ID

	if r.AuthToken == nil || *r.AuthToken == "" {
		return csr, errors.New("No authentication token found, is the session setup correctly?")
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
		return csr, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return csr, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &oemSSvc)
	if err != nil {
		return csr, err
	}

	if oemSSvc.Links.HTTPSCert.ID == nil {
		return csr, fmt.Errorf("BUG: .links.HttpsCert.Id not present or is null in data from %s", response.URL)
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
		}).Info("Requesting certficate signing request")
	}
	response, err = r.httpRequest(httpscertloc, "GET", nil, nil, false)
	if err != nil {
		return csr, err
	}

	raw = response.Content

	if response.StatusCode != http.StatusOK {
		return csr, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &httpscert)
	if err != nil {
		return csr, err
	}

	if httpscert.CSR == nil {
		// Note: We can't really distinguish between a running CSR generation or not.
		// If no CSR generation has started and no certificate was imported the API reports "CertificateSigningRequest": null,
		// whereas CertificateSigningRequest is not present when CSR generation is running but the JSON parser can't distinguish between both
		// situations
		return csr, fmt.Errorf("No CertificateSigningRequest found. Either CSR generation hasn't been started or is still running")
	}

	csr = *httpscert.CSR
	return csr, nil
}

func (r *Redfish) fetchCSRHPE(mgr *ManagerData) (string, error) {
	var csr string
	var oemHpe ManagerDataOemHpe
	var secsvc string
	var oemSSvc SecurityServiceDataOemHpe
	var httpscertloc string
	var httpscert HTTPSCertDataOemHpe

	// parse Oem section from JSON
	err := json.Unmarshal(mgr.Oem, &oemHpe)
	if err != nil {
		return csr, err
	}

	// get SecurityService endpoint from .Oem.Hp.links.SecurityService
	if oemHpe.Hpe.Links.SecurityService.ID == nil {
		return csr, errors.New("BUG: .Hpe.Links.SecurityService.Id not found or null")
	}
	secsvc = *oemHpe.Hpe.Links.SecurityService.ID

	if r.AuthToken == nil || *r.AuthToken == "" {
		return csr, errors.New("No authentication token found, is the session setup correctly?")
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
		return csr, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return csr, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &oemSSvc)
	if err != nil {
		return csr, err
	}

	if oemSSvc.Links.HTTPSCert.ID == nil {
		return csr, fmt.Errorf("BUG: .links.HttpsCert.Id not present or is null in data from %s", response.URL)
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
		}).Info("Requesting certficate signing request")
	}
	response, err = r.httpRequest(httpscertloc, "GET", nil, nil, false)
	if err != nil {
		return csr, err
	}

	raw = response.Content

	if response.StatusCode != http.StatusOK {
		return csr, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &httpscert)
	if err != nil {
		return csr, err
	}

	if httpscert.CSR == nil {
		// Note: We can't really distinguish between a running CSR generation or not.
		// If no CSR generation has started and no certificate was imported the API reports "CertificateSigningRequest": null,
		// whereas CertificateSigningRequest is not present when CSR generation is running but the JSON parser can't distinguish between both
		// situations
		return csr, fmt.Errorf("No CertificateSigningRequest found. Either CSR generation hasn't been started or is still running")
	}

	csr = *httpscert.CSR
	return csr, nil
}

func (r *Redfish) fetchCSRHuawei(mgr *ManagerData) (string, error) {
	var csr string
	var oemHuawei ManagerDataOemHuawei
	var secsvc string
	var oemSSvc SecurityServiceDataOemHuawei
	var httpscertloc string
	var httpscert HTTPSCertDataOemHuawei

	// parse Oem section from JSON
	err := json.Unmarshal(mgr.Oem, &oemHuawei)
	if err != nil {
		return csr, err
	}

	// get SecurityService endpoint from .Oem.Huawei.SecurityService
	if oemHuawei.Huawei.SecurityService.ID == nil {
		return csr, errors.New("BUG: .Huawei.SecurityService.Id not found or null")
	}
	secsvc = *oemHuawei.Huawei.SecurityService.ID

	if r.AuthToken == nil || *r.AuthToken == "" {
		return csr, errors.New("No authentication token found, is the session setup correctly?")
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
		return csr, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return csr, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &oemSSvc)
	if err != nil {
		return csr, err
	}

	if oemSSvc.Links.HTTPSCert.ID == nil {
		return csr, fmt.Errorf("BUG: .links.HttpsCert.Id not present or is null in data from %s", response.URL)
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
		}).Info("Requesting certficate signing request")
	}
	response, err = r.httpRequest(httpscertloc, "GET", nil, nil, false)
	if err != nil {
		return csr, err
	}

	raw = response.Content

	if response.StatusCode != http.StatusOK {
		return csr, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &httpscert)
	if err != nil {
		return csr, err
	}

	if httpscert.CSR == nil {
		// Note: We can't really distinguish between a running CSR generation or not.
		// If no CSR generation has started and no certificate was imported the API reports "CertificateSigningRequest": null,
		// whereas CertificateSigningRequest is not present when CSR generation is running but the JSON parser can't distinguish between both
		// situations
		return csr, fmt.Errorf("No CertificateSigningRequest found. Either CSR generation hasn't been started or is still running")
	}

	csr = *httpscert.CSR
	return csr, nil
}

func (r *Redfish) getCSRTargetHP(mgr *ManagerData) (string, error) {
	var csrTarget string
	var oemHp ManagerDataOemHp
	var secsvc string
	var oemSSvc SecurityServiceDataOemHp
	var httpscertloc string
	var httpscert HTTPSCertDataOemHp

	// parse Oem section from JSON
	err := json.Unmarshal(mgr.Oem, &oemHp)
	if err != nil {
		return csrTarget, err
	}

	// get SecurityService endpoint from .Oem.Hp.links.SecurityService
	if oemHp.Hp.Links.SecurityService.ID == nil {
		return csrTarget, errors.New("BUG: .Hp.Links.SecurityService.Id not found or null")
	}
	secsvc = *oemHp.Hp.Links.SecurityService.ID

	if r.AuthToken == nil || *r.AuthToken == "" {
		return csrTarget, errors.New("No authentication token found, is the session setup correctly?")
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
		return csrTarget, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return csrTarget, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &oemSSvc)
	if err != nil {
		return csrTarget, err
	}

	if oemSSvc.Links.HTTPSCert.ID == nil {
		return csrTarget, fmt.Errorf("BUG: .links.HttpsCert.Id not present or is null in data from %s", response.URL)
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
		}).Info("Requesting path to certificate signing request")
	}
	response, err = r.httpRequest(httpscertloc, "GET", nil, nil, false)

	if err != nil {
		return csrTarget, err
	}

	raw = response.Content

	if response.StatusCode != http.StatusOK {
		return csrTarget, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &httpscert)
	if err != nil {
		return csrTarget, err
	}

	if httpscert.Actions.GenerateCSR.Target == nil {
		return csrTarget, fmt.Errorf("BUG: .Actions.GenerateCSR.Target is not present or empty in JSON data from %s", response.URL)
	}

	csrTarget = *httpscert.Actions.GenerateCSR.Target
	return csrTarget, nil
}

func (r *Redfish) getCSRTargetHPE(mgr *ManagerData) (string, error) {
	var csrTarget string
	var oemHpe ManagerDataOemHpe
	var secsvc string
	var oemSSvc SecurityServiceDataOemHpe
	var httpscertloc string
	var httpscert HTTPSCertDataOemHpe

	// parse Oem section from JSON
	err := json.Unmarshal(mgr.Oem, &oemHpe)
	if err != nil {
		return csrTarget, err
	}

	// get SecurityService endpoint from .Oem.Hp.links.SecurityService
	if oemHpe.Hpe.Links.SecurityService.ID == nil {
		return csrTarget, errors.New("BUG: .Hpe.Links.SecurityService.Id not found or null")
	}
	secsvc = *oemHpe.Hpe.Links.SecurityService.ID

	if r.AuthToken == nil || *r.AuthToken == "" {
		return csrTarget, errors.New("No authentication token found, is the session setup correctly?")
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
		return csrTarget, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return csrTarget, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &oemSSvc)
	if err != nil {
		return csrTarget, err
	}

	if oemSSvc.Links.HTTPSCert.ID == nil {
		return csrTarget, fmt.Errorf("BUG: .links.HttpsCert.Id not present or is null in data from %s", response.URL)
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
		}).Info("Requesting path to certificate signing request")
	}
	response, err = r.httpRequest(httpscertloc, "GET", nil, nil, false)

	if err != nil {
		return csrTarget, err
	}

	raw = response.Content

	if response.StatusCode != http.StatusOK {
		return csrTarget, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &httpscert)
	if err != nil {
		return csrTarget, err
	}

	if httpscert.Actions.GenerateCSR.Target == nil {
		return csrTarget, fmt.Errorf("BUG: .Actions.GenerateCSR.Target is not present or empty in JSON data from %s", response.URL)
	}

	csrTarget = *httpscert.Actions.GenerateCSR.Target
	return csrTarget, nil
}

func (r *Redfish) getCSRTargetHuawei(mgr *ManagerData) (string, error) {
	var csrTarget string
	var oemHuawei ManagerDataOemHuawei
	var secsvc string
	var oemSSvc SecurityServiceDataOemHuawei
	var httpscertloc string
	var httpscert HTTPSCertDataOemHuawei

	// parse Oem section from JSON
	err := json.Unmarshal(mgr.Oem, &oemHuawei)
	if err != nil {
		return csrTarget, err
	}

	// get SecurityService endpoint from .Oem.Huawei.SecurityService
	if oemHuawei.Huawei.SecurityService.ID == nil {
		return csrTarget, errors.New("BUG: .Huawei.SecurityService.Id not found or null")
	}
	secsvc = *oemHuawei.Huawei.SecurityService.ID

	if r.AuthToken == nil || *r.AuthToken == "" {
		return csrTarget, errors.New("No authentication token found, is the session setup correctly?")
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
		return csrTarget, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return csrTarget, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &oemSSvc)
	if err != nil {
		return csrTarget, err
	}

	if oemSSvc.Links.HTTPSCert.ID == nil {
		return csrTarget, fmt.Errorf("BUG: .links.HttpsCert.Id not present or is null in data from %s", response.URL)
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
		}).Info("Requesting path to certificate signing request")
	}
	response, err = r.httpRequest(httpscertloc, "GET", nil, nil, false)

	if err != nil {
		return csrTarget, err
	}

	raw = response.Content

	if response.StatusCode != http.StatusOK {
		return csrTarget, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &httpscert)
	if err != nil {
		return csrTarget, err
	}

	if httpscert.Actions.GenerateCSR.Target == nil {
		return csrTarget, fmt.Errorf("BUG: .Actions.GenerateCSR.Target is not present or empty in JSON data from %s", response.URL)
	}

	csrTarget = *httpscert.Actions.GenerateCSR.Target
	return csrTarget, nil
}

func (r *Redfish) makeCSRPayloadHP(csr CSRData) string {
	// Note: HPE uses the same format as HP
	var csrstr string

	if csr.C == "" {
		csr.C = "XX"
	}
	if csr.S == "" {
		csr.S = "-"
	}
	if csr.L == "" {
		csr.L = "-"
	}
	if csr.O == "" {
		csr.O = "-"
	}
	if csr.OU == "" {
		csr.OU = "-"
	}

	csrstr += fmt.Sprintf("\"Country\": \"%s\", ", csr.C)
	csrstr += fmt.Sprintf("\"State\": \"%s\", ", csr.S)
	csrstr += fmt.Sprintf("\"City\": \"%s\", ", csr.L)
	csrstr += fmt.Sprintf("\"OrgName\": \"%s\", ", csr.O)
	csrstr += fmt.Sprintf("\"OrgUnit\": \"%s\", ", csr.OU)

	if csr.CN != "" {
		csrstr += fmt.Sprintf("\"CommonName\": \"%s\" ", csr.CN)
	} else {
		csrstr += fmt.Sprintf("\"CommonName\": \"%s\" ", r.Hostname)
	}

	csrstr = "{ " + csrstr + " } "
	return csrstr
}

func (r *Redfish) makeCSRPayloadVanilla(csr CSRData) string {
	var csrstr string

	if csr.C != "" {
		csrstr += fmt.Sprintf("\"Country\": \"%s\", ", csr.C)
	} else {
		csrstr += "XX"
	}

	if csr.S != "" {
		csrstr += fmt.Sprintf("\"State\": \"%s\", ", csr.S)
	}

	if csr.L != "" {
		csrstr += fmt.Sprintf("\"City\": \"%s\", ", csr.L)
	}

	if csr.O != "" {
		csrstr += fmt.Sprintf("\"OrgName\": \"%s\", ", csr.O)
	}

	if csr.OU != "" {
		csrstr += fmt.Sprintf("\"OrgUnit\": \"%s\", ", csr.OU)
	}

	if csr.CN != "" {
		csrstr += fmt.Sprintf("\"CommonName\": \"%s\" ", csr.CN)
	} else {
		csrstr += fmt.Sprintf("\"CommonName\": \"%s\" ", r.Hostname)
	}

	csrstr = "{ " + csrstr + " } "
	return csrstr
}

func (r *Redfish) makeCSRPayload(csr CSRData) string {
	var csrstr string

	if r.Flavor == RedfishHP || r.Flavor == RedfishHPE {
		csrstr = r.makeCSRPayloadHP(csr)
	} else {
		csrstr = r.makeCSRPayloadVanilla(csr)
	}

	return csrstr
}

func (r *Redfish) validateCSRData(csr CSRData) error {
	switch r.Flavor {
	case RedfishDell:
		return nil
	case RedfishHP:
		if csr.C == "" || csr.CN == "" || csr.O == "" || csr.OU == "" || csr.L == "" || csr.S == "" {
			return fmt.Errorf("HP requires C, CN, O, OU, L and S to be set")
		}
	case RedfishHPE:
		if csr.C == "" || csr.CN == "" || csr.O == "" || csr.OU == "" || csr.L == "" || csr.S == "" {
			return fmt.Errorf("HP requires C, CN, O, OU, L and S to be set")
		}
	case RedfishHuawei:
		// Huawei: Doesn't accept / as part of any field - see Issue#11
		if strings.Index(csr.C, "/") != -1 {
			return fmt.Errorf("Huwaei doesn't accept / as part of any field")
		}
		if strings.Index(csr.S, "/") != -1 {
			return fmt.Errorf("Huwaei doesn't accept / as part of any field")
		}
		if strings.Index(csr.L, "/") != -1 {
			return fmt.Errorf("Huwaei doesn't accept / as part of any field")
		}
		if strings.Index(csr.O, "/") != -1 {
			return fmt.Errorf("Huwaei doesn't accept / as part of any field")
		}
		if strings.Index(csr.OU, "/") != -1 {
			return fmt.Errorf("Huwaei doesn't accept / as part of any field")
		}
		if strings.Index(csr.CN, "/") != -1 {
			return fmt.Errorf("Huwaei doesn't accept / as part of any field")
		}

	case RedfishInspur:
		return nil
	case RedfishSuperMicro:
		return nil
	default:
		return nil
	}
	return nil
}

// GenCSR - generate CSR
func (r *Redfish) GenCSR(csr CSRData) error {
	var csrstr string
	var gencsrtarget string

	if r.AuthToken == nil || *r.AuthToken == "" {
		return errors.New("No authentication token found, is the session setup correctly?")
	}

	// set vendor flavor
	err := r.GetVendorFlavor()
	if err != nil {
		return err
	}

	err = r.validateCSRData(csr)
	if err != nil {
		return err
	}

	csrstr = r.makeCSRPayload(csr)

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
		gencsrtarget, err = r.getCSRTargetHP(mgr0)
		if err != nil {
			return err
		}
	} else if r.Flavor == RedfishHPE {
		gencsrtarget, err = r.getCSRTargetHPE(mgr0)
		if err != nil {
			return err
		}
	} else if r.Flavor == RedfishHuawei {
		gencsrtarget, err = r.getCSRTargetHuawei(mgr0)
		if err != nil {
			return err
		}
	} else if r.Flavor == RedfishInspur {
		return errors.New("Inspur management boards do not support CSR generation")
	} else if r.Flavor == RedfishSuperMicro {
		return errors.New("SuperMicro management boards do not support CSR generation")
	} else {
		return errors.New("Unable to get vendor for management board. If this vendor supports CSR generation please file a feature request")
	}

	if gencsrtarget == "" {
		return errors.New("BUG: CSR generation target is not known")
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               gencsrtarget,
			"method":             "POST",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting CSR generation")
	}
	if r.Debug {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               gencsrtarget,
			"method":             "POST",
			"additional_headers": nil,
			"use_basic_auth":     false,
			"payload":            csrstr,
		}).Debug("Requesting CSR generation")
	}
	response, err := r.httpRequest(gencsrtarget, "POST", nil, strings.NewReader(csrstr), false)
	if err != nil {
		return err
	}
	// XXX: do we need to look at the content returned by HTTP POST ?

	switch response.StatusCode {
	case http.StatusOK:
		fallthrough
	case http.StatusCreated:
		fallthrough
	case http.StatusAccepted:
	default:
		return fmt.Errorf("HTTP POST to %s returned \"%s\" instead of \"200 OK\", \"201 Created\" or \"202 Accepted\"", response.URL, response.Status)
	}

	return nil
}

// FetchCSR - fetch CSR
func (r *Redfish) FetchCSR() (string, error) {
	var csrstr string

	// set vendor flavor
	err := r.GetVendorFlavor()
	if err != nil {
		return csrstr, err
	}

	// get list of Manager endpoint
	mgrList, err := r.GetManagers()
	if err != nil {
		return csrstr, err
	}

	// pick the first entry
	mgr0, err := r.GetManagerData(mgrList[0])
	if err != nil {
		return csrstr, err
	}

	// get endpoint SecurityService from Managers
	if r.Flavor == RedfishHP {
		csrstr, err = r.fetchCSRHP(mgr0)
		if err != nil {
			return csrstr, err
		}
	} else if r.Flavor == RedfishHPE {
		csrstr, err = r.fetchCSRHPE(mgr0)
		if err != nil {
			return csrstr, err
		}
	} else if r.Flavor == RedfishHuawei {
		csrstr, err = r.fetchCSRHuawei(mgr0)
		if err != nil {
			return csrstr, err
		}
	} else if r.Flavor == RedfishInspur {
		return csrstr, errors.New("Inspur management boards do not support CSR generation")
	} else if r.Flavor == RedfishSuperMicro {
		return csrstr, errors.New("SuperMicro management boards do not support CSR generation")
	} else {
		return csrstr, errors.New("Unable to get vendor for management board. If this vendor supports CSR generation please file a feature request")
	}

	// convert "raw" string (new lines escaped as \n) to real string (new lines are new lines)
	return strings.Replace(csrstr, "\\n", "\n", -1), nil
}
