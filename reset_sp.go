package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func (r *Redfish) getManagerResetTarget_Supermicro(mgr *ManagerData) (string, error) {
	var actions_sm ManagerActionsDataOemSupermicro
	var target string

	err := json.Unmarshal(mgr.Actions, &actions_sm)
	if err != nil {
		return target, err
	}

	if actions_sm.Oem.ManagerReset.Target == nil || *actions_sm.Oem.ManagerReset.Target == "" {
		return target, errors.New(fmt.Sprintf("ERROR: No ManagerReset.Target found in data or ManagerReset.Target is null"))
	}

	return *actions_sm.Oem.ManagerReset.Target, nil
}

func (r *Redfish) getManagerResetTarget_Vanilla(mgr *ManagerData) (string, error) {
	var actions_sm ManagerActionsData
	var target string

	err := json.Unmarshal(mgr.Actions, &actions_sm)
	if err != nil {
		return target, err
	}

	if actions_sm.ManagerReset.Target == nil || *actions_sm.ManagerReset.Target == "" {
		return target, errors.New(fmt.Sprintf("ERROR: No ManagerReset.Target found in data or ManagerReset.Target is null"))
	}

	return *actions_sm.ManagerReset.Target, nil
}

func (r *Redfish) getManagerResetTarget(mgr *ManagerData) (string, error) {
	var err error
	var sp_reset_target string

	if r.Flavor == REDFISH_SUPERMICRO {
		sp_reset_target, err = r.getManagerResetTarget_Supermicro(mgr)
	} else {
		sp_reset_target, err = r.getManagerResetTarget_Vanilla(mgr)
	}
	if err != nil {
		return sp_reset_target, err
	}

	return sp_reset_target, nil
}

func (r *Redfish) ResetSP() error {
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

	sp_reset_target, err := r.getManagerResetTarget(mgr0)
	if err != nil {
		return err
	}

	sp_reset_payload := "{ \"ResetType\": \"ForceRestart\" }"
	response, err := r.httpRequest(sp_reset_target, "POST", nil, strings.NewReader(sp_reset_payload), false)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("ERROR: HTTP POST to %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	return nil
}
