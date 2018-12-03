package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func (r *Redfish) fetchCSR_HP(mgr *ManagerData) (string, error) {
	var csr string
	var oemHp ManagerDataOemHp
	var secsvc string
	var oemSSvc SecurityServiceDataOemHp
	var httpscertloc string
	var httpscert HttpsCertDataOemHp

	// parse Oem section from JSON
	err := json.Unmarshal(mgr.Oem, &oemHp)
	if err != nil {
		return csr, err
	}

	// get SecurityService endpoint from .Oem.Hp.links.SecurityService
	if oemHp.Hp.Links.SecurityService.Id == nil {
		return csr, errors.New("BUG: .Hp.Links.SecurityService.Id not found or null")
	} else {
		secsvc = *oemHp.Hp.Links.SecurityService.Id
	}

	if r.AuthToken == nil || *r.AuthToken == "" {
		return csr, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	response, err := r.httpRequest(secsvc, "GET", nil, nil, false)
	if err != nil {
		return csr, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return csr, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &oemSSvc)
	if err != nil {
		return csr, err
	}

	if oemSSvc.Links.HttpsCert.Id == nil {
		return csr, errors.New(fmt.Sprintf("BUG: .links.HttpsCert.Id not present or is null in data from %s", response.Url))
	}

	httpscertloc = *oemSSvc.Links.HttpsCert.Id

	response, err = r.httpRequest(httpscertloc, "GET", nil, nil, false)
	if err != nil {
		return csr, err
	}

	raw = response.Content

	if response.StatusCode != http.StatusOK {
		return csr, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
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
		return csr, errors.New(fmt.Sprintf("ERROR: No CertificateSigningRequest found. Either CSR generation hasn't been started or is still running"))
	}

	csr = *httpscert.CSR
	return csr, nil
}

func (r *Redfish) fetchCSR_Huawei(mgr *ManagerData) (string, error) {
	var csr string
	var oemHuawei ManagerDataOemHuawei
	var secsvc string
	var oemSSvc SecurityServiceDataOemHuawei
	var httpscertloc string
	var httpscert HttpsCertDataOemHuawei

	// parse Oem section from JSON
	err := json.Unmarshal(mgr.Oem, &oemHuawei)
	if err != nil {
		return csr, err
	}

	// get SecurityService endpoint from .Oem.Huawei.SecurityService
	if oemHuawei.Huawei.SecurityService.Id == nil {
		return csr, errors.New("BUG: .Huawei.SecurityService.Id not found or null")
	} else {
		secsvc = *oemHuawei.Huawei.SecurityService.Id
	}

	if r.AuthToken == nil || *r.AuthToken == "" {
		return csr, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	response, err := r.httpRequest(secsvc, "GET", nil, nil, false)
	if err != nil {
		return csr, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return csr, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &oemSSvc)
	if err != nil {
		return csr, err
	}

	if oemSSvc.Links.HttpsCert.Id == nil {
		return csr, errors.New(fmt.Sprintf("BUG: .links.HttpsCert.Id not present or is null in data from %s", response.Url))
	}

	httpscertloc = *oemSSvc.Links.HttpsCert.Id

	response, err = r.httpRequest(httpscertloc, "GET", nil, nil, false)
	if err != nil {
		return csr, err
	}

	raw = response.Content

	if response.StatusCode != http.StatusOK {
		return csr, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
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
		return csr, errors.New(fmt.Sprintf("ERROR: No CertificateSigningRequest found. Either CSR generation hasn't been started or is still running"))
	}

	csr = *httpscert.CSR
	return csr, nil
}

func (r *Redfish) getCSRTarget_HP(mgr *ManagerData) (string, error) {
	var csrTarget string
	var oemHp ManagerDataOemHp
	var secsvc string
	var oemSSvc SecurityServiceDataOemHp
	var httpscertloc string
	var httpscert HttpsCertDataOemHp

	// parse Oem section from JSON
	err := json.Unmarshal(mgr.Oem, &oemHp)
	if err != nil {
		return csrTarget, err
	}

	// get SecurityService endpoint from .Oem.Hp.links.SecurityService
	if oemHp.Hp.Links.SecurityService.Id == nil {
		return csrTarget, errors.New("BUG: .Hp.Links.SecurityService.Id not found or null")
	} else {
		secsvc = *oemHp.Hp.Links.SecurityService.Id
	}

	if r.AuthToken == nil || *r.AuthToken == "" {
		return csrTarget, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	response, err := r.httpRequest(secsvc, "GET", nil, nil, false)
	if err != nil {
		return csrTarget, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return csrTarget, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &oemSSvc)
	if err != nil {
		return csrTarget, err
	}

	if oemSSvc.Links.HttpsCert.Id == nil {
		return csrTarget, errors.New(fmt.Sprintf("BUG: .links.HttpsCert.Id not present or is null in data from %s", response.Url))
	}

	httpscertloc = *oemSSvc.Links.HttpsCert.Id

	response, err = r.httpRequest(httpscertloc, "GET", nil, nil, false)

	if err != nil {
		return csrTarget, err
	}

	raw = response.Content

	if response.StatusCode != http.StatusOK {
		return csrTarget, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &httpscert)
	if err != nil {
		return csrTarget, err
	}

	if httpscert.Actions.GenerateCSR.Target == nil {
		return csrTarget, errors.New(fmt.Sprintf("BUG: .Actions.GenerateCSR.Target is not present or empty in JSON data from %s", response.Url))
	}

	csrTarget = *httpscert.Actions.GenerateCSR.Target
	return csrTarget, nil
}

func (r *Redfish) getCSRTarget_Huawei(mgr *ManagerData) (string, error) {
	var csrTarget string
	var oemHuawei ManagerDataOemHuawei
	var secsvc string
	var oemSSvc SecurityServiceDataOemHuawei
	var httpscertloc string
	var httpscert HttpsCertDataOemHuawei

	// parse Oem section from JSON
	err := json.Unmarshal(mgr.Oem, &oemHuawei)
	if err != nil {
		return csrTarget, err
	}

	// get SecurityService endpoint from .Oem.Huawei.SecurityService
	if oemHuawei.Huawei.SecurityService.Id == nil {
		return csrTarget, errors.New("BUG: .Huawei.SecurityService.Id not found or null")
	} else {
		secsvc = *oemHuawei.Huawei.SecurityService.Id
	}

	if r.AuthToken == nil || *r.AuthToken == "" {
		return csrTarget, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	response, err := r.httpRequest(secsvc, "GET", nil, nil, false)
	if err != nil {
		return csrTarget, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return csrTarget, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &oemSSvc)
	if err != nil {
		return csrTarget, err
	}

	if oemSSvc.Links.HttpsCert.Id == nil {
		return csrTarget, errors.New(fmt.Sprintf("BUG: .links.HttpsCert.Id not present or is null in data from %s", response.Url))
	}

	httpscertloc = *oemSSvc.Links.HttpsCert.Id

	response, err = r.httpRequest(httpscertloc, "GET", nil, nil, false)

	if err != nil {
		return csrTarget, err
	}

	raw = response.Content

	if response.StatusCode != http.StatusOK {
		return csrTarget, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &httpscert)
	if err != nil {
		return csrTarget, err
	}

	if httpscert.Actions.GenerateCSR.Target == nil {
		return csrTarget, errors.New(fmt.Sprintf("BUG: .Actions.GenerateCSR.Target is not present or empty in JSON data from %s", response.Url))
	}

	csrTarget = *httpscert.Actions.GenerateCSR.Target
	return csrTarget, nil
}

func (r *Redfish) makeCSRPayload_HP(csr CSRData) string {
	var csrstr string = ""

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

func (r *Redfish) makeCSRPayload_Vanilla(csr CSRData) string {
	var csrstr string = ""

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

	if r.Flavor == REDFISH_HP {
		csrstr = r.makeCSRPayload_HP(csr)
	} else {
		csrstr = r.makeCSRPayload_Vanilla(csr)
	}

	return csrstr
}

func (r *Redfish) GenCSR(csr CSRData) error {
	var csrstr string = ""
	var gencsrtarget string = ""

	if r.AuthToken == nil || *r.AuthToken == "" {
		return errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	// set vendor flavor
	err := r.GetVendorFlavor()
	if err != nil {
		return err
	}

	csrstr = r.makeCSRPayload(csr)

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
		gencsrtarget, err = r.getCSRTarget_HP(mgr0)
		if err != nil {
			return err
		}
	} else if r.Flavor == REDFISH_HUAWEI {
		gencsrtarget, err = r.getCSRTarget_Huawei(mgr0)
		if err != nil {
			return err
		}
	} else if r.Flavor == REDFISH_INSPUR {
		return errors.New("ERROR: Inspur management boards do not support CSR generation")
	} else if r.Flavor == REDFISH_SUPERMICRO {
		return errors.New("ERROR: SuperMicro management boards do not support CSR generation")
	} else {
		return errors.New("ERROR: Unable to get vendor for management board. If this vendor supports CSR generation please file a feature request")
	}

	if gencsrtarget == "" {
		return errors.New("BUG: CSR generation target is not known")
	}

	response, err := r.httpRequest(gencsrtarget, "POST", nil, strings.NewReader(csrstr), false)
	if err != nil {
		return err
	}
	// XXX: do we need to look at the content returned by HTTP POST ?

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("ERROR: HTTP POST to %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	return nil
}

func (r *Redfish) FetchCSR() (string, error) {
	var csrstr string = ""

	// set vendor flavor
	err := r.GetVendorFlavor()
	if err != nil {
		return csrstr, err
	}

	// get list of Manager endpoint
	mgr_list, err := r.GetManagers()
	if err != nil {
		return csrstr, err
	}

	// pick the first entry
	mgr0, err := r.GetManagerData(mgr_list[0])
	if err != nil {
		return csrstr, err
	}

	// get endpoint SecurityService from Managers
	if r.Flavor == REDFISH_HP {
		csrstr, err = r.fetchCSR_HP(mgr0)
		if err != nil {
			return csrstr, err
		}
	} else if r.Flavor == REDFISH_HUAWEI {
		csrstr, err = r.fetchCSR_Huawei(mgr0)
		if err != nil {
			return csrstr, err
		}
	} else if r.Flavor == REDFISH_INSPUR {
		return csrstr, errors.New("ERROR: Inspur management boards do not support CSR generation")
	} else if r.Flavor == REDFISH_SUPERMICRO {
		return csrstr, errors.New("ERROR: SuperMicro management boards do not support CSR generation")
	} else {
		return csrstr, errors.New("ERROR: Unable to get vendor for management board. If this vendor supports CSR generation please file a feature request")
	}

	// convert "raw" string (new lines escaped as \n) to real string (new lines are new lines)
	return strings.Replace(csrstr, "\\n", "\n", -1), nil
}
