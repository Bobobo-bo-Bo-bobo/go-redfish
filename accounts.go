package redfish

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

//get array of accounts and their endpoints
func (r *Redfish) GetAccounts() ([]string, error) {
	var accsvc AccountService
	var accs OData
	var result = make([]string, 0)

	// check if vendor supports account management
	if r.Flavor == REDFISH_FLAVOR_NOT_INITIALIZED {
		err := r.GetVendorFlavor()
		if err != nil {
			return result, err
		}
	}
	if VendorCapabilities[r.FlavorString]&HAS_ACCOUNTSERVICE != HAS_ACCOUNTSERVICE {
		return result, errors.New("ERROR: Account management is not support for this vendor")
	}

	if r.AuthToken == nil || *r.AuthToken == "" {
		return result, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
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
		return result, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
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
			"path":               *accsvc.AccountsEndpoint.Id,
			"method":             "GET",
			"additional_headers": nil,
			"use_basic_auth":     false,
		}).Info("Requesting accounts")
	}
	response, err = r.httpRequest(*accsvc.AccountsEndpoint.Id, "GET", nil, nil, false)
	if err != nil {
		return result, err
	}

	raw = response.Content
	if response.StatusCode != http.StatusOK {
		return result, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &accs)
	if err != nil {
		return result, err
	}

	if len(accs.Members) == 0 {
		return result, errors.New(fmt.Sprintf("BUG: Missing or empty Members attribute in Accounts"))
	}

	for _, a := range accs.Members {
		result = append(result, *a.Id)
	}
	return result, nil
}

// get account data for a particular account
func (r *Redfish) GetAccountData(accountEndpoint string) (*AccountData, error) {
	var result AccountData

	// check if vendor supports account management
	if r.Flavor == REDFISH_FLAVOR_NOT_INITIALIZED {
		err := r.GetVendorFlavor()
		if err != nil {
			return nil, err
		}
	}
	if VendorCapabilities[r.FlavorString]&HAS_ACCOUNTSERVICE != HAS_ACCOUNTSERVICE {
		return nil, errors.New("ERROR: Account management is not support for this vendor")
	}

	if r.AuthToken == nil || *r.AuthToken == "" {
		return nil, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
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
		return nil, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(raw, &result)
	if err != nil {
		return nil, err
	}
	result.SelfEndpoint = &accountEndpoint
	return &result, nil
}

// map username -> user data
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

// map ID -> user data
func (r *Redfish) MapAccountsById() (map[string]*AccountData, error) {
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
		if a.Id == nil {
			return result, errors.New("BUG: No Id found or Id is null")
		}
		result[*a.Id] = a
	}

	return result, nil
}

func (r *Redfish) AddAccount(acd AccountCreateData) error {
	var acsd AccountService
	var accep string
	var payload string
	var rerr RedfishError

	if r.AuthToken == nil || *r.AuthToken == "" {
		return errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	if r.Flavor == REDFISH_FLAVOR_NOT_INITIALIZED {
		err := r.GetVendorFlavor()
		if err != nil {
			return err
		}
	}

	// check if vendor supports account management
	if VendorCapabilities[r.FlavorString]&HAS_ACCOUNTSERVICE != HAS_ACCOUNTSERVICE {
		return errors.New("ERROR: Account management is not support for this vendor")
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
		return errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	err = json.Unmarshal(response.Content, &acsd)
	if err != nil {
		return err
	}

	if acsd.AccountsEndpoint == nil {
		return errors.New(fmt.Sprintf("ERROR: No Accounts endpoint found in response from %s", response.Url))
	}

	if acsd.AccountsEndpoint.Id == nil || *acsd.AccountsEndpoint.Id == "" {
		return errors.New(fmt.Sprintf("BUG: Defined Accounts endpoint from %s does not have @odata.id value", response.Url))
	}

	accep = *acsd.AccountsEndpoint.Id

	if r.Flavor == REDFISH_HP {
		// XXX: Use Oem specific privilege map
	} else {
		if acd.UserName == "" || acd.Password == "" || acd.Role == "" {
			return errors.New("ERROR: Required field(s) missing")
		}

		// check of requested role exists, role Names are _NOT_ unique (e.g. Supermicro report all names as "User Role") but Id is
		rmap, err := r.MapRolesById()
		if err != nil {
			return err
		}

		_, found := rmap[acd.Role]
		if !found {
			return errors.New(fmt.Sprintf("ERROR: Requested role %s not found", acd.Role))
		}

		payload = fmt.Sprintf("{ \"UserName\": \"%s\", \"Password\": \"%s\", \"RoleId\": \"%s\" }", acd.UserName, acd.Password, acd.Role)
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
			err = json.Unmarshal(response.Content, &rerr)
			if err != nil {
				return errors.New(fmt.Sprintf("ERROR: HTTP POST for %s returned \"%s\" and no error information", response.Url, response.Status))
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
			errmsg := ""
			if len(rerr.Error.MessageExtendedInfo) > 0 {
				for _, e := range rerr.Error.MessageExtendedInfo {
					if e.Message != nil || *e.Message != "" {
						if errmsg == "" {
							errmsg += *e.Message
						} else {
							errmsg += "; " + *e.Message
						}
					}
				}
			} else {
				if rerr.Error.Message != nil || *rerr.Error.Message != "" {
					errmsg = *rerr.Error.Message
				} else {
					errmsg = fmt.Sprintf("HTTP POST for %s returned \"%s\" and error information but error information neither contains @Message.ExtendedInfo nor Message", response.Url, response.Status)
				}
			}
			return errors.New(fmt.Sprintf("ERROR: %s", errmsg))
		}

		// any other error ? (HTTP 400 has been handled above)
		if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusBadRequest {
			return errors.New(fmt.Sprintf("ERROR: HTTP POST for %s returned \"%s\" instead of \"200 OK\" or \"201 Created\"", response.Url, response.Status))
		}
	}
	return nil
}

func (r *Redfish) DeleteAccount(u string) error {
	if r.AuthToken == nil || *r.AuthToken == "" {
		return errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	// check if vendor supports account management
	if r.Flavor == REDFISH_FLAVOR_NOT_INITIALIZED {
		err := r.GetVendorFlavor()
		if err != nil {
			return err
		}
	}
	if VendorCapabilities[r.FlavorString]&HAS_ACCOUNTSERVICE != HAS_ACCOUNTSERVICE {
		return errors.New("ERROR: Account management is not support for this vendor")
	}

	// get endpoint for account to delete
	amap, err := r.MapAccountsByName()
	if err != nil {
		return err
	}

	adata, found := amap[u]
	if !found {
		return errors.New(fmt.Sprintf("ERROR: Account %s not found", u))
	}

	if adata.SelfEndpoint == nil || *adata.SelfEndpoint == "" {
		return errors.New(fmt.Sprintf("BUG: SelfEndpoint not set or empty in account data for %s", u))
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
	response, err := r.httpRequest(*adata.SelfEndpoint, "DELETE", nil, nil, false)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("ERROR: HTTP DELETE for %s returned \"%s\" instead of \"200 OK\"", response.Url, response.Status))
	}

	return nil
}

func (r *Redfish) ChangePassword(u string, p string) error {
	var payload string
	var rerr RedfishError

	if u == "" {
		return errors.New("ERROR: Username is empty")
	}

	if p == "" {
		return errors.New(fmt.Sprintf("ERROR: Password for %s is empty", u))
	}

	if r.AuthToken == nil || *r.AuthToken == "" {
		return errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	// check if vendor supports account management
	if r.Flavor == REDFISH_FLAVOR_NOT_INITIALIZED {
		err := r.GetVendorFlavor()
		if err != nil {
			return err
		}
	}
	if VendorCapabilities[r.FlavorString]&HAS_ACCOUNTSERVICE != HAS_ACCOUNTSERVICE {
		return errors.New("ERROR: Account management is not support for this vendor")
	}

	// check if the account exists
	amap, err := r.MapAccountsByName()
	if err != nil {
	}

	adata, found := amap[u]
	if !found {
		return errors.New(fmt.Sprintf("ERROR: Account %s not found", u))
	}

	if adata.SelfEndpoint == nil || *adata.SelfEndpoint == "" {
		return errors.New(fmt.Sprintf("BUG: SelfEndpoint not set or empty in account data for %s", u))
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
		err = json.Unmarshal(response.Content, &rerr)
		if err != nil {
			return errors.New(fmt.Sprintf("ERROR: HTTP POST for %s returned \"%s\" and no error information", response.Url, response.Status))
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
		errmsg := ""
		if len(rerr.Error.MessageExtendedInfo) > 0 {
			for _, e := range rerr.Error.MessageExtendedInfo {
				if e.Message != nil || *e.Message != "" {
					if errmsg == "" {
						errmsg += *e.Message
					} else {
						errmsg += "; " + *e.Message
					}
				}
			}
		} else {
			if rerr.Error.Message != nil || *rerr.Error.Message != "" {
				errmsg = *rerr.Error.Message
			} else {
				errmsg = fmt.Sprintf("HTTP POST for %s returned \"%s\" and error information but error information neither contains @Message.ExtendedInfo nor Message", response.Url, response.Status)
			}
		}
		return errors.New(fmt.Sprintf("ERROR: %s", errmsg))
	}

	// any other error ? (HTTP 400 has been handled above)
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusBadRequest {
		return errors.New(fmt.Sprintf("ERROR: HTTP POST for %s returned \"%s\" instead of \"200 OK\" or \"201 Created\"", response.Url, response.Status))
	}
	return nil
}

func (r *Redfish) makeAccountCreateModifyPayload(acd AccountCreateData) (string, error) {
	var payload string

	// handle HP(E) PrivilegeMap
	if r.Flavor == REDFISH_HP {
		//FIXME: handle HP(E) PrivilegeMap
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

func (r *Redfish) ModifyAccount(u string, acd AccountCreateData) error {
	var rerr RedfishError

	if r.AuthToken == nil || *r.AuthToken == "" {
		return errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	// check if vendor supports account management
	if r.Flavor == REDFISH_FLAVOR_NOT_INITIALIZED {
		err := r.GetVendorFlavor()
		if err != nil {
			return err
		}
	}
	if VendorCapabilities[r.FlavorString]&HAS_ACCOUNTSERVICE != HAS_ACCOUNTSERVICE {
		return errors.New("ERROR: Account management is not support for this vendor")
	}

	// get endpoint for account to modify/check if account with this name already exists
	umap, err := r.MapAccountsByName()
	if err != nil {
		return err
	}

	udata, found := umap[u]
	if !found {
		return errors.New(fmt.Sprintf("ERROR: User %s not found", u))
	}
	if udata.SelfEndpoint == nil || *udata.SelfEndpoint == "" {
		return errors.New(fmt.Sprintf("BUG: SelfEndpoint is not set or empty for user %s", u))
	}

	if r.Flavor == REDFISH_HP {
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
			err = json.Unmarshal(response.Content, &rerr)
			if err != nil {
				return errors.New(fmt.Sprintf("ERROR: HTTP POST for %s returned \"%s\" and no error information", response.Url, response.Status))
			}
			errmsg := ""
			if len(rerr.Error.MessageExtendedInfo) > 0 {
				for _, e := range rerr.Error.MessageExtendedInfo {
					if e.Message != nil || *e.Message != "" {
						if errmsg == "" {
							errmsg += *e.Message
						} else {
							errmsg += "; " + *e.Message
						}
					}
				}
			} else {
				if rerr.Error.Message != nil || *rerr.Error.Message != "" {
					errmsg = *rerr.Error.Message
				} else {
					errmsg = fmt.Sprintf("HTTP POST for %s returned \"%s\" and error information but error information neither contains @Message.ExtendedInfo nor Message", response.Url, response.Status)
				}
			}
			return errors.New(fmt.Sprintf("ERROR: %s", errmsg))
		}

		// any other error ? (HTTP 400 has been handled above)
		if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusBadRequest {
			return errors.New(fmt.Sprintf("ERROR: HTTP POST for %s returned \"%s\" instead of \"200 OK\" or \"201 Created\"", response.Url, response.Status))
		}
	}
	return nil
}
