package redfish

// ManagerDataOemHpeLicense - HP/HPE: Oem data for Manager endpoint and SecurityService endpoint
type ManagerDataOemHpeLicense struct {
	Key    *string `json:"LicenseKey"`
	String *string `json:"LicenseString"`
	Type   *string `json:"LicenseType"`
	Expire *string `json:"LicenseExpire"`
}

// ManagerDataOemHpeFederationConfig - HP/HPE federatin information
type ManagerDataOemHpeFederationConfig struct {
	IPv6MulticastScope            *string `json:"IPv6MulticastScope"`
	MulticastAnnouncementInterval *int64  `json:"MulticastAnnouncementInterval"`
	MulticastDiscovery            *string `json:"MulticastDiscovery"`
	MulticastTimeToLive           *int64  `json:"MulticastTimeToLive"`
	ILOFederationManagement       *string `json:"iLOFederationManagement"`
}

// ManagerDataOemHpeFirmwareData - HP/HPE firmware data
type ManagerDataOemHpeFirmwareData struct {
	Date         *string `json:"Date"`
	DebugBuild   *bool   `json:"DebugBuild"`
	MajorVersion *uint64 `json:"MajorVersion"`
	MinorVersion *uint64 `json:"MinorVersion"`
	Time         *string `json:"Time"`
	Version      *string `json:"Version"`
}

// ManagerDataOemHpeFirmware - current firmware of the management processor
type ManagerDataOemHpeFirmware struct {
	Current ManagerDataOemHpeFirmwareData `json:"Current"`
}

// ManagerDataOemHpeLinks - link targets for HPE vendor specific extensions
// NOTE: The result for HP/HPE are different depending if the HTTP header
// OData-Version is set or not. If OData-Version is _NOT_ set data are returned in
// .Oem.Hp.links with endpoints in "href". If OData-Version is set
// data are returned in .Oem.Hp.Links (note different case!) and endpoints are
// defined as @odata.id. We always set "OData-Version: 4.0"
type ManagerDataOemHpeLinks struct {
	ActiveHealthSystem   OData `json:"ActiveHealthSystem"`
	DateTimeService      OData `json:"DateTimeService"`
	EmbeddedMediaService OData `json:"EmbeddedMediaService"`
	FederationDispatch   OData `json:"FederationDispatch"`
	FederationGroups     OData `json:"FederationGroups"`
	FederationPeers      OData `json:"FederationPeers"`
	LicenseService       OData `json:"LicenseService"`
	SecurityService      OData `json:"SecurityService"`
	UpdateService        OData `json:"UpdateService"`
	VSPLogLocation       OData `json:"VSPLogLocation"`
}

type _managerDataOemHpe struct {
	FederationConfig ManagerDataOemHpeFederationConfig `json:"FederationConfig"`
	Firmware         ManagerDataOemHpeFirmware         `json:"Firmware"`
	License          ManagerDataOemHpeLicense          `json:"License"`
	Type             *string                           `json:"Type"`
	Links            ManagerDataOemHpeLinks            `json:"Links"`
}

// ManagerDataOemHpe - OEM data for HPE systems
type ManagerDataOemHpe struct {
	Hpe *_managerDataOemHpe `json:"Hpe"`
}

// SecurityServiceDataOemHpeLinks - HPE extension
type SecurityServiceDataOemHpeLinks struct {
	ESKM      OData `json:"ESKM"`
	HTTPSCert OData `json:"HttpsCert"`
	SSO       OData `json:"SSO"`
}

// SecurityServiceDataOemHpe - HPE extension
type SecurityServiceDataOemHpe struct {
	ID    *string                        `json:"Id"`
	Type  *string                        `json:"Type"`
	Links SecurityServiceDataOemHpeLinks `json:"Links"`
}

// HTTPSCertActionsOemHpe - OEM extension for certificate management
type HTTPSCertActionsOemHpe struct {
	GenerateCSR                LinkTargets  `json:"#HpeHttpsCert.GenerateCSR"`
	ImportCertificate          LinkTargets  `json:"#HpeHttpsCert.ImportCertificate"`
	X509CertificateInformation X509CertInfo `json:"X509CertificateInformation"`
}

// HTTPSCertDataOemHpe - HPE OEM extension
type HTTPSCertDataOemHpe struct {
	CSR     *string                `json:"CertificateSigningRequest"`
	ID      *string                `json:"Id"`
	Actions HTTPSCertActionsOemHpe `json:"Actions"`
}

// AccountPrivilegeMapOemHpe - HP(E) uses it's own privilege map instead of roles
type AccountPrivilegeMapOemHpe struct {
	Login                bool `json:"LoginPriv"`
	RemoteConsole        bool `json:"RemoteConsolePriv"`
	UserConfig           bool `json:"UserConfigPriv"`
	VirtualMedia         bool `json:"VirtualMediaPriv"`
	VirtualPowerAndReset bool `json:"VirtualPowerAndResetPriv"`
	ILOConfig            bool `json:"iLOConfigPriv"`
}
