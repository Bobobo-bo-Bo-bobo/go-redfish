package redfish

// make a deep copy
func (r *Redfish) Clone() *Redfish {
	var a *string = new(string)
	var sl *string = new(string)

	var cpy *Redfish = new(Redfish)

	if cpy != nil {
		cpy.Hostname = r.Hostname
		cpy.Port = r.Port
		cpy.Username = r.Username
		cpy.Password = r.Password
		cpy.AuthToken = nil
		cpy.SessionLocation = nil
		cpy.Timeout = r.Timeout
		cpy.InsecureSSL = r.InsecureSSL
		cpy.Debug = r.Debug
		cpy.Verbose = r.Verbose
		cpy.RawBaseContent = r.RawBaseContent
		cpy.AccountService = r.AccountService
		cpy.Chassis = r.Chassis
		cpy.Managers = r.Managers
		cpy.SessionService = r.SessionService
		cpy.Sessions = r.Sessions
		cpy.Systems = r.Systems
		cpy.Flavor = r.Flavor
		cpy.FlavorString = r.FlavorString
		cpy.initialised = r.initialised

		if r.AuthToken != nil {
			*a = *r.AuthToken
			cpy.AuthToken = a
		}

		if r.SessionLocation != nil {
			*sl = *r.SessionLocation
			cpy.SessionLocation = sl
		}
	}

	return cpy
}
