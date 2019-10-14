package redfish

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (r *Redfish) getManagerResetTargetSupermicro(mgr *ManagerData) (string, error) {
	var actionsSm ManagerActionsDataOemSupermicro
	var target string

	err := json.Unmarshal(mgr.Actions, &actionsSm)
	if err != nil {
		return target, err
	}

	if actionsSm.Oem.ManagerReset.Target == nil || *actionsSm.Oem.ManagerReset.Target == "" {
		return target, fmt.Errorf("No ManagerReset.Target found in data or ManagerReset.Target is null")
	}

	return *actionsSm.Oem.ManagerReset.Target, nil
}

func (r *Redfish) getManagerResetTargetVanilla(mgr *ManagerData) (string, error) {
	var actionsSm ManagerActionsData
	var target string

	err := json.Unmarshal(mgr.Actions, &actionsSm)
	if err != nil {
		return target, err
	}

	if actionsSm.ManagerReset.Target == nil || *actionsSm.ManagerReset.Target == "" {
		return target, fmt.Errorf("No ManagerReset.Target found in data or ManagerReset.Target is null")
	}

	return *actionsSm.ManagerReset.Target, nil
}

func (r *Redfish) getManagerResetTarget(mgr *ManagerData) (string, error) {
	var err error
	var spResetTarget string

	if r.Flavor == RedfishSuperMicro {
		spResetTarget, err = r.getManagerResetTargetSupermicro(mgr)
	} else {
		spResetTarget, err = r.getManagerResetTargetVanilla(mgr)
	}
	if err != nil {
		return spResetTarget, err
	}

	return spResetTarget, nil
}

// ResetSP - reset service processor
func (r *Redfish) ResetSP() error {
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

	spResetTarget, err := r.getManagerResetTarget(mgr0)
	if err != nil {
		return err
	}

	spResetPayload := "{ \"ResetType\": \"ForceRestart\" }"
	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               spResetTarget,
			"method":             "POST",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting service processor restart")
	}
	if r.Debug {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               spResetTarget,
			"method":             "POST",
			"additional_headers": nil,
			"use_basic_auth":     false,
			"payload":            spResetPayload,
		}).Debug("Requesting service processor restart")
	}
	response, err := r.httpRequest(spResetTarget, "POST", nil, strings.NewReader(spResetPayload), false)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP POST to %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	return nil
}
