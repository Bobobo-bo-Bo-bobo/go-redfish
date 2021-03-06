package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (r *Redfish) hpGetLicense(mgr *ManagerData) (*ManagerLicenseData, error) {
	var lic ManagerLicenseData
	var m ManagerDataOemHp

	err := json.Unmarshal(mgr.Oem, &m)
	if err != nil {
		return nil, err
	}

	// no license key means no license
	if m.Hp.License.Key == nil {
		return nil, nil
	}
	if *m.Hp.License.Key == "" {
		return nil, nil
	}

	lic.Name = "HP iLO license"
	if m.Hp.License.Expire != nil {
		lic.Expiration = *m.Hp.License.Expire
	}

	if m.Hp.License.Type != nil {
		lic.Type = *m.Hp.License.Type
	}

	if m.Hp.License.Key != nil {
		lic.License = *m.Hp.License.Key
	}

	return &lic, nil
}

func (r *Redfish) hpeGetLicense(mgr *ManagerData) (*ManagerLicenseData, error) {
	var lic ManagerLicenseData
	var m ManagerDataOemHpe

	err := json.Unmarshal(mgr.Oem, &m)
	if err != nil {
		return nil, err
	}

	// no license key means no license
	if m.Hpe.License.Key == nil {
		return nil, nil
	}
	if *m.Hpe.License.Key == "" {
		return nil, nil
	}

	lic.Name = "HPE iLO license"
	if m.Hpe.License.Expire != nil {
		lic.Expiration = *m.Hpe.License.Expire
	}

	if m.Hpe.License.Type != nil {
		lic.Type = *m.Hpe.License.Type
	}

	if m.Hpe.License.Key != nil {
		lic.License = *m.Hpe.License.Key
	}

	return &lic, nil
}

// GetLicense - get licenses of management board
func (r *Redfish) GetLicense(mgr *ManagerData) (*ManagerLicenseData, error) {
	if r.AuthToken == nil || *r.AuthToken == "" {
		return nil, errors.New("No authentication token found, is the session setup correctly?")
	}

	if r.Flavor == RedfishFlavorNotInitialized {
		err := r.GetVendorFlavor()
		if err != nil {
			return nil, err
		}
	}

	if r.Flavor == RedfishHP {
		return r.hpGetLicense(mgr)
	} else if r.Flavor == RedfishHPE {
		return r.hpeGetLicense(mgr)
	}

	return nil, errors.New("License operations are not supported for this vendor. If this vendor supports license operations please file a feature request")
}

func (r *Redfish) hpHpePrepareLicensePayload(l []byte) string {
	return fmt.Sprintf(`{ "LicenseKey": "%s" }`, string(l))
}

func (r *Redfish) hpSetLicense(mgr *ManagerData, l []byte) error {
	var m ManagerDataOemHp

	err := json.Unmarshal(mgr.Oem, &m)
	if err != nil {
		return err
	}

	// get LicenseService endpoint path from OEM data
	if m.Hp.Links.LicenseService.ID == nil || *m.Hp.Links.LicenseService.ID == "" {
		return fmt.Errorf("BUG: Expected LicenseService endpoint definition in .Oem.Hp.Links for vendor %s, but found none", r.FlavorString)
	}

	licensePayload := r.hpHpePrepareLicensePayload(l)

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               *m.Hp.Links.LicenseService.ID,
			"method":             "POST",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Uploading license")
	}
	if r.Debug {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               *m.Hp.Links.LicenseService.ID,
			"method":             "POST",
			"additional_headers": nil,
			"use_basic_auth":     false,
			"payload":            licensePayload,
		}).Debug("Uploading license")
	}

	response, err := r.httpRequest(*m.Hp.Links.LicenseService.ID, "POST", nil, strings.NewReader(licensePayload), false)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		redfishError, err := r.ProcessError(response)
		if err != nil {
			return fmt.Errorf("License installation returned \"%s\" instead of \"200 OK\" or \"201 Created\" and didn't return an error object", response.Status)
		}
		msg := r.GetErrorMessage(redfishError)
		if msg != "" {
			return fmt.Errorf("License installation failed: %s", msg)
		}
		return fmt.Errorf("License installation returned \"%s\" instead of \"200 OK\" or \"201 Created\" and didn't return an error object", response.Status)
	}
	return nil
}

func (r *Redfish) hpeSetLicense(mgr *ManagerData, l []byte) error {
	var m ManagerDataOemHpe

	err := json.Unmarshal(mgr.Oem, &m)
	if err != nil {
		return err
	}

	// get LicenseService endpoint path from OEM data
	if m.Hpe.Links.LicenseService.ID == nil || *m.Hpe.Links.LicenseService.ID == "" {
		return fmt.Errorf("BUG: Expected LicenseService endpoint definition in .Oem.Hpe.Links for vendor %s, but found none", r.FlavorString)
	}

	licensePayload := r.hpHpePrepareLicensePayload(l)

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               *m.Hpe.Links.LicenseService.ID,
			"method":             "POST",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Uploading license")
	}
	if r.Debug {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               *m.Hpe.Links.LicenseService.ID,
			"method":             "POST",
			"additional_headers": nil,
			"use_basic_auth":     false,
			"payload":            licensePayload,
		}).Debug("Uploading license")
	}

	response, err := r.httpRequest(*m.Hpe.Links.LicenseService.ID, "POST", nil, strings.NewReader(licensePayload), false)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		redfishError, err := r.ProcessError(response)
		if err != nil {
			return fmt.Errorf("License installation returned \"%s\" instead of \"200 OK\" or \"201 Created\" and didn't return an error object", response.Status)
		}
		msg := r.GetErrorMessage(redfishError)
		if msg != "" {
			return fmt.Errorf("License installation failed: %s", msg)
		}
		return fmt.Errorf("License installation returned \"%s\" instead of \"200 OK\" or \"201 Created\" and didn't return an error object", response.Status)
	}

	return nil
}

// AddLicense - add license to management board
func (r *Redfish) AddLicense(mgr *ManagerData, l []byte) error {
	if r.AuthToken == nil || *r.AuthToken == "" {
		return errors.New("No authentication token found, is the session setup correctly?")
	}

	if r.Flavor == RedfishFlavorNotInitialized {
		err := r.GetVendorFlavor()
		if err != nil {
			return err
		}
	}

	if r.Flavor == RedfishHP {
		return r.hpSetLicense(mgr, l)
	} else if r.Flavor == RedfishHPE {
		return r.hpeSetLicense(mgr, l)
	}

	return errors.New("License operations are not supported for this vendor. If this vendor supports license operations please file a feature request")

}
