package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// GetAccounts - get array of accounts and their endpoints
func (r *Redfish) GetAccounts() ([]string, error) {
	var accsvc AccountService
	var accs OData
	var result = make([]string, 0)

	// check if vendor supports account management
	if r.Flavor == RedfishFlavorNotInitialized {
		err := r.GetVendorFlavor()
		if err != nil {
			return result, err
		}
	}
	if VendorCapabilities[r.FlavorString]&HasAccountService != HasAccountService {
		return result, errors.New("Account management is not support for this vendor")
	}

	if r.AuthToken == nil || *r.AuthToken == "" {
		return result, errors.New("No authentication token found, is the session setup correctly?")
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               r.AccountService,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting path to account service")
	}
	response, err := r.httpRequest(r.AccountService, "GET", nil, nil, false)
	if err != nil {
		return result, err
	}

	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return result, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &accsvc)
	if err != nil {
		return result, err
	}

	if accsvc.AccountsEndpoint == nil {
		return result, errors.New("BUG: No Accounts endpoint found")
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               *accsvc.AccountsEndpoint.ID,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting accounts")
	}
	response, err = r.httpRequest(*accsvc.AccountsEndpoint.ID, "GET", nil, nil, false)
	if err != nil {
		return result, err
	}

	raw = response.Content
	if response.StatusCode != http.StatusOK {
		return result, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &accs)
	if err != nil {
		return result, err
	}

	if len(accs.Members) == 0 {
		return result, fmt.Errorf("BUG: Missing or empty Members attribute in Accounts")
	}

	for _, a := range accs.Members {
		result = append(result, *a.ID)
	}
	return result, nil
}

// GetAccountData - get account data for a particular account
func (r *Redfish) GetAccountData(accountEndpoint string) (*AccountData, error) {
	var result AccountData

	// check if vendor supports account management
	if r.Flavor == RedfishFlavorNotInitialized {
		err := r.GetVendorFlavor()
		if err != nil {
			return nil, err
		}
	}
	if VendorCapabilities[r.FlavorString]&HasAccountService != HasAccountService {
		return nil, errors.New("Account management is not support for this vendor")
	}

	if r.AuthToken == nil || *r.AuthToken == "" {
		return nil, errors.New("No authentication token found, is the session setup correctly?")
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               accountEndpoint,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting account information")
	}
	response, err := r.httpRequest(accountEndpoint, "GET", nil, nil, false)
	if err != nil {
		return nil, err
	}

	// store unparsed content
	raw := response.Content

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(raw, &result)
	if err != nil {
		return nil, err
	}
	result.SelfEndpoint = &accountEndpoint
	return &result, nil
}

// MapAccountsByName - map username -> user data
func (r *Redfish) MapAccountsByName() (map[string]*AccountData, error) {
	var result = make(map[string]*AccountData)

	al, err := r.GetAccounts()
	if err != nil {
		return result, err
	}

	for _, acc := range al {
		a, err := r.GetAccountData(acc)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if a.UserName == nil {
			return result, errors.New("BUG: No UserName found or UserName is null")
		}
		// Note: some vendors like DELL/EMC use predefined number of accounts
		//       and report an empty UserName for unused accounts "slots"
		if *a.UserName == "" {
			if r.Verbose {
				log.WithFields(log.Fields{
					"hostname":      r.Hostname,
					"port":          r.Port,
					"timeout":       r.Timeout,
					"flavor":        r.Flavor,
					"flavor_string": r.FlavorString,
					"path":          *a.SelfEndpoint,
				}).Info("Discarding account because UserName field is empty")
			}
			continue
		}
		result[*a.UserName] = a
	}

	return result, nil
}

// MapAccountsByID - map ID -> user data
func (r *Redfish) MapAccountsByID() (map[string]*AccountData, error) {
	var result = make(map[string]*AccountData)

	al, err := r.GetAccounts()
	if err != nil {
		return result, err
	}

	for _, acc := range al {
		a, err := r.GetAccountData(acc)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if a.ID == nil {
			return result, errors.New("BUG: No Id found or Id is null")
		}
		result[*a.ID] = a
	}

	return result, nil
}

// get endpoint of first free account slot
func (r *Redfish) dellGetFreeAccountSlot() (string, error) {
	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":      r.Hostname,
			"port":          r.Port,
			"timeout":       r.Timeout,
			"flavor":        r.Flavor,
			"flavor_string": r.FlavorString,
			"path":          r.AccountService,
		}).Info("Looking for unused account slot")
	}

	accountList, err := r.GetAccounts()
	if err != nil {
		return "", err
	}

	// Get account information to find the first unused account slot
	for slotIdx, accEndpt := range accountList {

		// Note: The first account slot is reserved and can't be modified
		//       ("Modifying the user configuration at index 1 is not allowed.")
		if slotIdx == 0 {
			continue
		}

		_acd, err := r.GetAccountData(accEndpt)
		if err != nil {
			return "", err
		}

		// should NEVER happen
		if _acd.UserName == nil {
			return "", errors.New("BUG: No UserName found or UserName is null")
		}

		if *_acd.UserName == "" {
			if r.Verbose {
				log.WithFields(log.Fields{
					"hostname":      r.Hostname,
					"port":          r.Port,
					"timeout":       r.Timeout,
					"flavor":        r.Flavor,
					"flavor_string": r.FlavorString,
					"path":          r.AccountService,
					"unused_slot":   accEndpt,
				}).Info("Found unused account slot")
			}
			return accEndpt, nil
		}
	}

	return "", nil
}

func (r *Redfish) dellAddAccount(acd AccountCreateData) error {
	var accountEnabled bool

	_unusedSlot, err := r.dellGetFreeAccountSlot()
	if err != nil {
		return err
	}

	if _unusedSlot == "" {
		return errors.New("No unused account slot found")
	}

	// Instead of adding an account we have to modify an existing
	// unused account slot.
	acd.Enabled = &accountEnabled
	return r.ModifyAccountByEndpoint(_unusedSlot, acd)
}

func (r *Redfish) hpBuildPrivilegeMap(flags uint) *AccountPrivilegeMapOemHp {
	var result AccountPrivilegeMapOemHp

	if flags&HpePrivilegeLogin == HpePrivilegeLogin {
		result.Login = true
	}

	if flags&HpePrivilegeRemoteConsole == HpePrivilegeRemoteConsole {
		result.RemoteConsole = true
	}

	if flags&HpePrivilegeUserConfig == HpePrivilegeUserConfig {
		result.UserConfig = true
	}

	if flags&HpePrivilegeVirtualMedia == HpePrivilegeVirtualMedia {
		result.VirtualMedia = true
	}

	if flags&HpePrivilegeVirtualPowerAndReset == HpePrivilegeVirtualPowerAndReset {
		result.VirtualPowerAndReset = true
	}

	if flags&HpePrivilegeIloConfig == HpePrivilegeIloConfig {
		result.ILOConfig = true
	}
	return &result
}

// AddAccount - Add account
func (r *Redfish) AddAccount(acd AccountCreateData) error {
	var acsd AccountService
	var accep string
	var payload string
	var _flags uint
	var found bool

	if r.AuthToken == nil || *r.AuthToken == "" {
		return errors.New("No authentication token found, is the session setup correctly?")
	}

	if r.Flavor == RedfishFlavorNotInitialized {
		err := r.GetVendorFlavor()
		if err != nil {
			return err
		}
	}

	// check if vendor supports account management
	if VendorCapabilities[r.FlavorString]&HasAccountService != HasAccountService {
		return errors.New("Account management is not support for this vendor")
	}

	// Note: DELL/EMC iDRAC uses a hardcoded, predefined number of account slots
	//       and as a consequence only support GET and HEAD on the "usual" endpoints
	if r.Flavor == RedfishDell {
		return r.dellAddAccount(acd)
	}

	// get Accounts endpoint from AccountService
	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               r.AccountService,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting path to account service")
	}
	response, err := r.httpRequest(r.AccountService, "GET", nil, nil, false)
	if err != nil {
		return nil
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	err = json.Unmarshal(response.Content, &acsd)
	if err != nil {
		return err
	}

	if acsd.AccountsEndpoint == nil {
		return fmt.Errorf("No Accounts endpoint found in response from %s", response.URL)
	}

	if acsd.AccountsEndpoint.ID == nil || *acsd.AccountsEndpoint.ID == "" {
		return fmt.Errorf("BUG: Defined Accounts endpoint from %s does not have @odata.id value", response.URL)
	}

	accep = *acsd.AccountsEndpoint.ID

	if r.Flavor == RedfishHP || r.Flavor == RedfishHPE {
		if acd.UserName == "" || acd.Password == "" {
			return errors.New("Required field(s) missing")
		}

		// OemHpPrivilegeMap is an INTERNAL map but it MUST be exported to be accessed by json.Marshal
		if acd.OemHpPrivilegeMap != nil {
			log.WithFields(log.Fields{
				"hostname":          r.Hostname,
				"port":              r.Port,
				"timeout":           r.Timeout,
				"flavor":            r.Flavor,
				"flavor_string":     r.FlavorString,
				"role":              acd.Role,
				"oemhpprivilegemap": *acd.OemHpPrivilegeMap,
				"hpeprivileges":     acd.HPEPrivileges,
			}).Warning("Internal field OemHpPrivilegeMap is set, discarding it's content")
			acd.OemHpPrivilegeMap = nil
		}

		if acd.Role != "" {
			// map "roles" to privileges
			virtualRole := strings.TrimSpace(strings.ToLower(acd.Role))
			_flags, found = HPEVirtualRoles[virtualRole]
			if !found {
				return fmt.Errorf("Unknown role %s", acd.Role)
			}

			// If additional privileges are set we add them too
			_flags |= acd.HPEPrivileges

			acd.OemHpPrivilegeMap = r.hpBuildPrivilegeMap(_flags)
		}

		rawPrivPayload, err := json.Marshal(*acd.OemHpPrivilegeMap)
		if err != nil {
			return err
		}

		payload = fmt.Sprintf("{ \"UserName\": \"%s\", \"Password\": \"%s\", \"Oem\":{ \"Hp\":{ \"LoginName\": \"%s\", \"Privileges\": %s }}}", acd.UserName, acd.Password, acd.UserName, string(rawPrivPayload))
	} else {
		if acd.UserName == "" || acd.Password == "" || acd.Role == "" {
			return errors.New("Required field(s) missing")
		}

		// check of requested role exists, role Names are _NOT_ unique (e.g. Supermicro report all names as "User Role") but Id is
		rmap, err := r.MapRolesByID()
		if err != nil {
			return err
		}

		_, found := rmap[acd.Role]
		if !found {
			return fmt.Errorf("Requested role %s not found", acd.Role)
		}

		payload = fmt.Sprintf("{ \"UserName\": \"%s\", \"Password\": \"%s\", \"RoleId\": \"%s\" }", acd.UserName, acd.Password, acd.Role)
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               accep,
			"method":             "POST",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Adding account")
	}
	if r.Debug {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               accep,
			"method":             "POST",
			"additional_headers": nil,
			"use_basic_auth":     false,
			"payload":            payload,
		}).Debug("Adding account")
	}
	response, err = r.httpRequest(accep, "POST", nil, strings.NewReader(payload), false)
	if err != nil {
		return err
	}

	// some vendors like Supermicro imposes limits on fields like password and return HTTP 400 - Bad Request
	if response.StatusCode == http.StatusBadRequest {
		rerr, err := r.ProcessError(response)
		if err != nil {
			return fmt.Errorf("Operation failed, returned \"%s\" and no error information", response.Status)
		}

		//
		// For instance Supermicro responds for creation with passwords exceeding the maximal password length with:
		// {
		//   "error": {
		// 	"code": "Base.v1_4_0.GeneralError",
		// 	"Message": "A general error has occurred. See ExtendedInfo for more information.",
		// 	"@Message.ExtendedInfo": [
		// 	  {
		// 		"MessageId": "Base.v1_4_0.PropertyValueFormatError",
		// 		"Severity": "Warning",
		// 		"Resolution": "Correct the value for the property in the request body and resubmit the request if the operation failed.",
		// 		"Message": "The value this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long. for the property Password is of a different format than the property can accept.",
		// 		"MessageArgs": [
		// 		  "this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.",
		// 		  "Password"
		// 		],
		// 		"RelatedProperties": [
		// 		  "Password"
		// 		]
		// 	  }
		// 	]
		//   }
		// }
		//
		errmsg := r.GetErrorMessage(rerr)
		if errmsg != "" {
			return fmt.Errorf("%s", errmsg)
		}
		return fmt.Errorf("Operation failed, returned \"%s\" and no error information", response.Status)
	}

	// any other error ? (HTTP 400 has been handled above)
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("Operation failed, returned \"%s\" and no error information", response.Status)
	}
	return nil
}

func (r *Redfish) dellDeleteAccount(endpoint string) error {
	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               endpoint,
			"method":             "PATCH",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Releasing DELL/EMC account slot")
	}
	if r.Debug {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               endpoint,
			"method":             "PATCH",
			"additional_headers": nil,
			"use_basic_auth":     false,
			"payload":            DELLEmptyAccountSlot,
		}).Info("Releasing DELL/EMC account slot")
	}

	response, err := r.httpRequest(endpoint, "PATCH", nil, strings.NewReader(DELLEmptyAccountSlot), false)
	if r.Debug {
	}

	if response.StatusCode != http.StatusOK {
		// TODO: Check error object
		return fmt.Errorf("HTTP PATCH for %s returned \"%s\"", endpoint, response.Status)
	}
	return err
}

// DeleteAccount - delete an account
func (r *Redfish) DeleteAccount(u string) error {
	if r.AuthToken == nil || *r.AuthToken == "" {
		return errors.New("No authentication token found, is the session setup correctly?")
	}

	// check if vendor supports account management
	if r.Flavor == RedfishFlavorNotInitialized {
		err := r.GetVendorFlavor()
		if err != nil {
			return err
		}
	}
	if VendorCapabilities[r.FlavorString]&HasAccountService != HasAccountService {
		return errors.New("Account management is not support for this vendor")
	}

	// get endpoint for account to delete
	amap, err := r.MapAccountsByName()
	if err != nil {
		return err
	}

	adata, found := amap[u]
	if !found {
		return fmt.Errorf("Account %s not found", u)
	}

	if adata.SelfEndpoint == nil || *adata.SelfEndpoint == "" {
		return fmt.Errorf("BUG: SelfEndpoint not set or empty in account data for %s", u)
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               *adata.SelfEndpoint,
			"method":             "DELETE",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Deleting account")
	}

	// Note: DELL/EMC only
	if r.Flavor == RedfishDell {
		return r.dellDeleteAccount(*adata.SelfEndpoint)
	}

	response, err := r.httpRequest(*adata.SelfEndpoint, "DELETE", nil, nil, false)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP DELETE for %s returned \"%s\" instead of \"200 OK\"", response.URL, response.Status)
	}

	return nil
}

// ChangePassword - change account password
func (r *Redfish) ChangePassword(u string, p string) error {
	var payload string

	if u == "" {
		return errors.New("Username is empty")
	}

	if p == "" {
		return fmt.Errorf("Password for %s is empty", u)
	}

	if r.AuthToken == nil || *r.AuthToken == "" {
		return errors.New("No authentication token found, is the session setup correctly?")
	}

	// check if vendor supports account management
	if r.Flavor == RedfishFlavorNotInitialized {
		err := r.GetVendorFlavor()
		if err != nil {
			return err
		}
	}
	if VendorCapabilities[r.FlavorString]&HasAccountService != HasAccountService {
		return errors.New("Account management is not support for this vendor")
	}

	// check if the account exists
	amap, err := r.MapAccountsByName()
	if err != nil {
	}

	adata, found := amap[u]
	if !found {
		return fmt.Errorf("Account %s not found", u)
	}

	if adata.SelfEndpoint == nil || *adata.SelfEndpoint == "" {
		return fmt.Errorf("BUG: SelfEndpoint not set or empty in account data for %s", u)
	}
	payload = fmt.Sprintf("{ \"Password\": \"%s\" }", p)

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               *adata.SelfEndpoint,
			"method":             "PATCH",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Changing account password")
	}
	if r.Debug {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               *adata.SelfEndpoint,
			"method":             "PATCH",
			"additional_headers": nil,
			"use_basic_auth":     false,
			"payload":            payload,
		}).Debug("Changing account password")
	}
	response, err := r.httpRequest(*adata.SelfEndpoint, "PATCH", nil, strings.NewReader(payload), false)
	if err != nil {
		return err
	}

	// some vendors like Supermicro imposes limits on fields like password and return HTTP 400 - Bad Request
	if response.StatusCode == http.StatusBadRequest {
		rerr, err := r.ProcessError(response)
		if err != nil {
			return fmt.Errorf("Operation failed, returned \"%s\" and no error information", response.Status)
		}

		//
		// For instance Supermicro responds for creation with passwords exceeding the maximal password length with:
		// {
		//   "error": {
		// 	"code": "Base.v1_4_0.GeneralError",
		// 	"Message": "A general error has occurred. See ExtendedInfo for more information.",
		// 	"@Message.ExtendedInfo": [
		// 	  {
		// 		"MessageId": "Base.v1_4_0.PropertyValueFormatError",
		// 		"Severity": "Warning",
		// 		"Resolution": "Correct the value for the property in the request body and resubmit the request if the operation failed.",
		// 		"Message": "The value this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long. for the property Password is of a different format than the property can accept.",
		// 		"MessageArgs": [
		// 		  "this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.this_password_is_too_long.",
		// 		  "Password"
		// 		],
		// 		"RelatedProperties": [
		// 		  "Password"
		// 		]
		// 	  }
		// 	]
		//   }
		// }
		//
		errmsg := r.GetErrorMessage(rerr)
		if errmsg != "" {
			return fmt.Errorf("%s", errmsg)
		}
		return fmt.Errorf("Operation failed, returned \"%s\" and no error information", response.Status)
	}

	// any other error ? (HTTP 400 has been handled above)
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("HTTP POST for %s returned \"%s\" instead of \"200 OK\" or \"201 Created\"", response.URL, response.Status)
	}
	return nil
}

func (r *Redfish) makeAccountCreateModifyPayload(acd AccountCreateData) (string, error) {
	var payload string
	var _flags uint
	var found bool

	// handle HP(E) PrivilegeMap
	if r.Flavor == RedfishHP || r.Flavor == RedfishHPE {
		// OemHpPrivilegeMap is an INTERNAL map but it MUST be exported to be accessed by json.Marshal
		if acd.OemHpPrivilegeMap != nil {
			log.WithFields(log.Fields{
				"hostname":          r.Hostname,
				"port":              r.Port,
				"timeout":           r.Timeout,
				"flavor":            r.Flavor,
				"flavor_string":     r.FlavorString,
				"role":              acd.Role,
				"oemhpprivilegemap": *acd.OemHpPrivilegeMap,
				"hpeprivileges":     acd.HPEPrivileges,
			}).Warning("Internal field OemHpPrivilegeMap is set, discarding it's content")
			acd.OemHpPrivilegeMap = nil
		}

		if acd.Role != "" {
			// map "roles" to privileges
			virtualRole := strings.TrimSpace(strings.ToLower(acd.Role))
			_flags, found = HPEVirtualRoles[virtualRole]
			if !found {
				return "", fmt.Errorf("Unknown role %s", acd.Role)
			}

			// If additional privileges are set we add them too
			_flags |= acd.HPEPrivileges

			acd.OemHpPrivilegeMap = r.hpBuildPrivilegeMap(_flags)
		}

		rawPrivPayload, err := json.Marshal(*acd.OemHpPrivilegeMap)
		if err != nil {
			return "", err
		}

		payload = fmt.Sprintf("{ \"UserName\": \"%s\", \"Password\": \"%s\", \"Oem\":{ \"Hp\":{ \"LoginName\": \"%s\", \"Privileges\": %s }}}", acd.UserName, acd.Password, acd.UserName, string(rawPrivPayload))
	} else {
		// force exclustion of privilege map for non-HP(E) systems
		acd.OemHpPrivilegeMap = nil
		raw, err := json.Marshal(acd)
		if err != nil {
			return payload, err
		}
		payload = string(raw)
	}
	return payload, nil
}

// ModifyAccount - modify an account
func (r *Redfish) ModifyAccount(u string, acd AccountCreateData) error {
	if r.AuthToken == nil || *r.AuthToken == "" {
		return errors.New("No authentication token found, is the session setup correctly?")
	}

	// check if vendor supports account management
	if r.Flavor == RedfishFlavorNotInitialized {
		err := r.GetVendorFlavor()
		if err != nil {
			return err
		}
	}
	if VendorCapabilities[r.FlavorString]&HasAccountService != HasAccountService {
		return errors.New("Account management is not support for this vendor")
	}

	// get endpoint for account to modify/check if account with this name already exists
	umap, err := r.MapAccountsByName()
	if err != nil {
		return err
	}

	udata, found := umap[u]
	if !found {
		return fmt.Errorf("User %s not found", u)
	}
	if udata.SelfEndpoint == nil || *udata.SelfEndpoint == "" {
		return fmt.Errorf("BUG: SelfEndpoint is not set or empty for user %s", u)
	}

	payload, err := r.makeAccountCreateModifyPayload(acd)
	if err != nil {
		return err
	}

	if r.Verbose {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               *udata.SelfEndpoint,
			"method":             "PATCH",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Modifying account")
	}
	if r.Debug {
		log.WithFields(log.Fields{
			"hostname":           r.Hostname,
			"port":               r.Port,
			"timeout":            r.Timeout,
			"flavor":             r.Flavor,
			"flavor_string":      r.FlavorString,
			"path":               *udata.SelfEndpoint,
			"method":             "PATCH",
			"additional_headers": nil,
			"use_basic_auth":     false,
			"payload":            payload,
		}).Debug("Modifying account")
	}
	response, err := r.httpRequest(*udata.SelfEndpoint, "PATCH", nil, strings.NewReader(payload), false)
	if err != nil {
		return err
	}

	// some vendors like Supermicro imposes limits on fields like password and return HTTP 400 - Bad Request
	if response.StatusCode == http.StatusBadRequest {
		rerr, err := r.ProcessError(response)
		if err != nil {
			return fmt.Errorf("Operation failed, returned \"%s\" and no error information", response.Status)
		}

		errmsg := r.GetErrorMessage(rerr)
		if errmsg != "" {
			return fmt.Errorf("%s", errmsg)
		}
		return fmt.Errorf("Operation failed, returned \"%s\" and no error information", response.Status)
	}

	// any other error ? (HTTP 400 has been handled above)
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("HTTP POST for %s returned \"%s\" instead of \"200 OK\" or \"201 Created\"", response.URL, response.Status)
	}
	return nil
}

// ModifyAccountByEndpoint - modify account by it's endpoint
func (r *Redfish) ModifyAccountByEndpoint(endpoint string, acd AccountCreateData) error {
	if r.AuthToken == nil || *r.AuthToken == "" {
		return errors.New("No authentication token found, is the session setup correctly?")
	}

	// check if vendor supports account management
	if r.Flavor == RedfishFlavorNotInitialized {
		err := r.GetVendorFlavor()
		if err != nil {
			return err
		}
	}
	if VendorCapabilities[r.FlavorString]&HasAccountService != HasAccountService {
		return errors.New("Account management is not support for this vendor")
	}

	if r.Flavor == RedfishHP || r.Flavor == RedfishHPE {
		// XXX: Use Oem specific privilege map
	} else {

		payload, err := r.makeAccountCreateModifyPayload(acd)
		if err != nil {
			return err
		}

		if r.Verbose {
			log.WithFields(log.Fields{
				"hostname":           r.Hostname,
				"port":               r.Port,
				"timeout":            r.Timeout,
				"flavor":             r.Flavor,
				"flavor_string":      r.FlavorString,
				"path":               endpoint,
				"method":             "PATCH",
				"additional_headers": nil,
				"use_basic_auth":     false,
			}).Info("Modifying account")
		}
		if r.Debug {
			log.WithFields(log.Fields{
				"hostname":           r.Hostname,
				"port":               r.Port,
				"timeout":            r.Timeout,
				"flavor":             r.Flavor,
				"flavor_string":      r.FlavorString,
				"path":               endpoint,
				"method":             "PATCH",
				"additional_headers": nil,
				"use_basic_auth":     false,
				"payload":            payload,
			}).Debug("Modifying account")
		}
		response, err := r.httpRequest(endpoint, "PATCH", nil, strings.NewReader(payload), false)
		if err != nil {
			return err
		}

		// some vendors like Supermicro imposes limits on fields like password and return HTTP 400 - Bad Request
		if response.StatusCode == http.StatusBadRequest {
			rerr, err := r.ProcessError(response)
			if err != nil {
				return fmt.Errorf("Operation failed, returned \"%s\" and no error information", response.Status)
			}

			errmsg := r.GetErrorMessage(rerr)
			if errmsg != "" {
				return fmt.Errorf("%s", errmsg)
			}
			return fmt.Errorf("Operation failed, returned \"%s\" and no error information", response.Status)
		}
	}
	return nil
}
