package redfish

// make a deep copy
func (r *Redfish) Clone() Redfish {
	var a string
	var sl string

	var cpy Redfish = Redfish{
		Hostname:        r.Hostname,
		Port:            r.Port,
		Username:        r.Username,
		Password:        r.Password,
		AuthToken:       nil,
		SessionLocation: nil,
		Timeout:         r.Timeout,
		InsecureSSL:     r.InsecureSSL,
		Debug:           r.Debug,
		Verbose:         r.Verbose,
		RawBaseContent:  r.RawBaseContent,
		AccountService:  r.AccountService,
		Chassis:         r.Chassis,
		Managers:        r.Managers,
		SessionService:  r.SessionService,
		Sessions:        r.Sessions,
		Systems:         r.Systems,
		Flavor:          r.Flavor,
		FlavorString:    r.FlavorString,
		initialised:     r.initialised,
	}

	if r.AuthToken != nil {
		a = *r.AuthToken
		cpy.AuthToken = &a
	}

	if r.SessionLocation != nil {
		sl = *r.SessionLocation
		cpy.SessionLocation = &sl
	}

	return cpy
}
