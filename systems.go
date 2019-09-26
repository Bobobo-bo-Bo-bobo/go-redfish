package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

//get array of systems and their endpoints
func (r *Redfish) GetSystems() ([]string, error) {
	var systems OData
	var result = make([]string, 0)

	if r.AuthToken == nil || *r.AuthToken == "" {
		return result, errors.New("ERROR: No authentication token found, is the session setup correctly?")
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               r.Systems,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting available systems")
	}
	response, err := r.httpRequest(r.Systems, "GET", nil, nil, false)
	if err != nil {
		return result, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return result, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &systems)
	if err != nil {
		return result, err
	}

	if len(systems.Members) == 0 {
		return result, errors.New("BUG: Array of system endpoints is empty")
	}

	for _, endpoint := range systems.Members {
		result = append(result, *endpoint.Id)
	}
	return result, nil
}

// get system data for a particular system
func (r *Redfish) GetSystemData(systemEndpoint string) (*SystemData, error) {
	var result SystemData

	if r.AuthToken == nil || *r.AuthToken == "" {
		return nil, errors.New("ERROR: No authentication token found, is the session setup correctly?")
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               systemEndpoint,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting system information")
	}
	response, err := r.httpRequest(systemEndpoint, "GET", nil, nil, false)
	if err != nil {
		return nil, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &result)
	if err != nil {
		return nil, err
	}

	result.SelfEndpoint = &systemEndpoint
	return &result, nil
}

// Map systems by ID
func (r *Redfish) MapSystemsById() (map[string]*SystemData, error) {
	var result = make(map[string]*SystemData)

	sysl, err := r.GetSystems()
	if err != nil {
		return result, nil
	}

	for _, sys := range sysl {
		s, err := r.GetSystemData(sys)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if s.Id == nil {
			return result, errors.New(fmt.Sprintf("BUG: No Id found for System at %s", sys))
		}

		result[*s.Id] = s
	}

	return result, nil
}

// Map systems by UUID
func (r *Redfish) MapSystemsByUuid() (map[string]*SystemData, error) {
	var result = make(map[string]*SystemData)

	sysl, err := r.GetSystems()
	if err != nil {
		return result, nil
	}

	for _, sys := range sysl {
		s, err := r.GetSystemData(sys)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if s.UUID == nil {
			return result, errors.New(fmt.Sprintf("BUG: No UUID found for System at %s", sys))
		}

		result[*s.UUID] = s
	}

	return result, nil
}

// Map systems by serial number
func (r *Redfish) MapSystemsBySerialNumber() (map[string]*SystemData, error) {
	var result = make(map[string]*SystemData)

	sysl, err := r.GetSystems()
	if err != nil {
		return result, nil
	}

	for _, sys := range sysl {
		s, err := r.GetSystemData(sys)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if s.SerialNumber == nil {
			return result, errors.New(fmt.Sprintf("BUG: No SerialNumber found for System at %s", sys))
		}

		result[*s.SerialNumber] = s
	}

	return result, nil
}

// OEM HP or HPE can't be distinguished only by the Manufacurer field of the System endpoint
// because newer BIOS/iLO4 versions set the Manufacturer to "HPE" but still use Oem.Hp instead
// of Oem.Hpe for vendor specific data.
func (r *Redfish) distinguishHpHpeFlavors(sd *SystemData) (int, string, error) {
    // Because we are only interested in the key of the Oem dict, lets use a temporary
    // simple struct to avoid parsing the JSON data a second time
    type _ManagerDataOemHpOrHpe struct {
	Hp json.RawMessage `json:"Hp"`
	Hpe json.RawMessage `json:"Hpe"`
    }
    var _OemhpOrHpe _ManagerDataOemHpOrHpe

    // parse JSON and look at the Oem fields
    err := json.Unmarshal(sd.Oem, &_OemhpOrHpe)
    if err != nil {
	return REDFISH_FLAVOR_UNINITIALIZED, "<error>", err
    }

    if len(_OemhpOrHpe.Hp) == 0 && len(_OemhpOrHpe.Hpe) > 0 {
	return REDFISH_HPE, "hpe", nil
    }

    if len(_OemhpOrHpe.Hp) > 0 && len(_OemhpOrHpe.Hpe) == 0 {
	return REDFISH_HP, "hp", nil
    }

    if len(_OemhpOrHpe.Hp) == 0 &&  len(_OemhpOrHpe.Hpe) == 0 {
	return REDFISH_FLAVOR_UNINITIALIZED, "<bug>", errors.New("BUG: Manufacturer is hp or hpe but Oem.Hp and Oem.Hpe are both undefined")
    }
}

// get vendor specific "flavor"
func (r *Redfish) GetVendorFlavor() error {
	// get vendor "flavor" for vendor specific implementation details
	_sys, err := r.GetSystems()
	if err != nil {
		return err
	}
	// assuming every system has the same vendor, pick the first one to determine vendor flavor
	_sys0, err := r.GetSystemData(_sys[0])
	if _sys0.Manufacturer != nil {
		_manufacturer := strings.TrimSpace(strings.ToLower(*_sys0.Manufacturer))
		if r.Debug {
			log.WithFields(log.Fields{
				"hostname":      r.Hostname,
				"port":          r.Port,
				"timeout":       r.Timeout,
				"flavor":        r.Flavor,
				"flavor_string": r.FlavorString,
				"manufacturer":  *_sys0.Manufacturer,
				"_manufacturer": _manufacturer,
			}).Debug("Identifying vendor flavor")
		}
		if _manufacturer == "hp" || _manufacturer == "_hpe" {
		    r.Flavor, r.FlavorString, err := distinguishHpHpeFlavors(_sys0)
		    if err != nil {
			return err
		    }
		} else if _manufacturer == "huawei" {
			r.Flavor = REDFISH_HUAWEI
			r.FlavorString = "huawei"
		} else if _manufacturer == "inspur" {
			r.Flavor = REDFISH_INSPUR
			r.FlavorString = "inspur"
		} else if _manufacturer == "supermicro" {
			r.Flavor = REDFISH_SUPERMICRO
			r.FlavorString = "supermicro"
		} else if _manufacturer == "dell inc." {
			r.Flavor = REDFISH_DELL
			r.FlavorString = "dell"
		} else if _manufacturer == "ibm" {
			r.Flavor = REDFISH_LENOVO
			r.FlavorString = "lenovo"
		} else {
			r.Flavor = REDFISH_GENERAL
			r.FlavorString = "vanilla"
		}
	}

	return nil
}

// set reset type map to map normalized state to supported variable value
func (r *Redfish) setAllowedResetTypes(sd *SystemData) error {
	if sd.Actions == nil {
		return errors.New(fmt.Sprintf("BUG: SystemData object don't define an Actions key"))
	}
	if sd.Actions.ComputerReset == nil {
		return errors.New(fmt.Sprintf("BUG: SystemData.Actions don't define a #ComputerSystem.Reset key"))
	}
	if sd.Actions.ComputerReset.Target == "" {
		return errors.New(fmt.Sprintf("BUG: SystemData.Actions.#ComputerSystem.Reset don't define a target key"))
	}

	if sd.Actions.ComputerReset.ActionInfo != "" {
		// TODO: Fetch information from ActionInfo URL
		if r.Verbose {
			log.WithFields(log.Fields{
				"hostname":           r.Hostname,
				"port":               r.Port,
				"timeout":            r.Timeout,
				"flavor":             r.Flavor,
				"flavor_string":      r.FlavorString,
				"path":               sd.Actions.ComputerReset.ActionInfo,
				"method":             "GET",
				"additional_headers": nil,
				"use_basic_auth":     false,
			}).Info("Requesting valid actions for system reset")
		}
		result, err := r.httpRequest(sd.Actions.ComputerReset.ActionInfo, "GET", nil, nil, false)
		if err != nil {
			return err
		}
		if result.StatusCode != 200 {
			return errors.New(fmt.Sprintf("ERROR: HTTP GET Request to %s returned %d - %s (expected 200)", sd.Actions.ComputerReset.ActionInfo, result.StatusCode, result.Status))
		}

		var sai SystemActionInfo
		raw := result.Content

		err = json.Unmarshal(raw, &sai)
		if err != nil {
			return err
		}

		// XXX: Assuming ActionInfo field "Parameters" for reset only contains one entry and this it contains the name
		//      of the field
		if len(sai.Parameters) == 0 {
			return errors.New("BUG: ActionInfo for system reset is either not defined or empty")
		}
		if sai.Parameters[0].Name == "" {
			return errors.New("BUG: ActionInfo.Parameters[0] don't have required field Name (or it is empty)")
		}
		if len(sai.Parameters[0].AllowableValues) == 0 {
			return errors.New(fmt.Sprintf("BUG: List of supported reset types in ActionInfo is not defined or empty"))
		}
		sd.allowedResetTypes = make(map[string]string)
		for _, t := range sai.Parameters[0].AllowableValues {
			_t := strings.ToLower(t)
			sd.allowedResetTypes[_t] = t
		}

		sd.resetTypeProperty = sai.Parameters[0].Name

	} else {
		if len(sd.Actions.ComputerReset.ResetTypeValues) == 0 {
			return errors.New(fmt.Sprintf("BUG: List of supported reset types is not defined or empty"))
		}

		sd.allowedResetTypes = make(map[string]string)
		for _, t := range sd.Actions.ComputerReset.ResetTypeValues {
			_t := strings.ToLower(t)
			sd.allowedResetTypes[_t] = t
		}
		// XXX: Is there a way to check the name of the reset type (is it always ResetType ?) ?
		//      For instance HPE define an extra key "AvailableActions" containing "PropertyName" for "Reset" action
		//      According to https://redfish.dmtf.org/schemas/DSP0266_1.7.0.html#allowable-values the name is the
		//      part before "@Redfish.AllowableValues" but this sucks for JSON parsing in Go!
		sd.resetTypeProperty = "ResetType"
	}
	return nil
}

// set power state
func (r *Redfish) SetSystemPowerState(sd *SystemData, state string) error {
	// do we already know the supported reset types?
	if len(sd.allowedResetTypes) == 0 {
		err := r.setAllowedResetTypes(sd)
		if err != nil {
			return err
		}
	}

	_state := strings.TrimSpace(strings.ToLower(state))
	resetType, found := sd.allowedResetTypes[_state]
	if found {
		// build payload
		payload := fmt.Sprintf("{ \"%s\": \"%s\" }", sd.resetTypeProperty, resetType)
		if r.Verbose {
			log.WithFields(log.Fields{
				"hostname":           r.Hostname,
				"port":               r.Port,
				"timeout":            r.Timeout,
				"flavor":             r.Flavor,
				"flavor_string":      r.FlavorString,
				"path":               sd.Actions.ComputerReset.Target,
				"method":             "POST",
				"additional_headers": nil,
				"use_basic_auth":     false,
			}).Info("Setting new system power state")
		}
		if r.Debug {
			log.WithFields(log.Fields{
				"hostname":           r.Hostname,
				"port":               r.Port,
				"timeout":            r.Timeout,
				"flavor":             r.Flavor,
				"flavor_string":      r.FlavorString,
				"path":               sd.Actions.ComputerReset.Target,
				"method":             "POST",
				"additional_headers": nil,
				"use_basic_auth":     false,
				"payload":            payload,
			}).Debug("Setting new system power state")
		}
		result, err := r.httpRequest(sd.Actions.ComputerReset.Target, "POST", nil, strings.NewReader(payload), false)
		if err != nil {
			return err
		}
		// DTMF Redfish schema definition defines the list of return codes following a POST operation
		// (see https://redfish.dmtf.org/schemas/DSP0266_1.7.0.html#post-action-a-id-post-action-a-)
		if result.StatusCode != 200 && result.StatusCode != 202 && result.StatusCode != 204 {
			return errors.New(fmt.Sprintf("ERROR: HTTP POST to %s returns HTTP status %d - %s (expect 200, 202 or 204)", sd.Actions.ComputerReset.Target, result.StatusCode, result.Status))
		}
	} else {
		return errors.New("ERROR: Requested PowerState is not supported for this system")
	}
	return nil
}
