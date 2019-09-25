package redfish

// HP/HPE: Oem data for Manager endpoint and SecurityService endpoint
type ManagerDataOemHpeLicense struct {
	Key    *string `json:"LicenseKey"`
	String *string `json:"LicenseString"`
	Type   *string `json:"LicenseType"`
	Expire *string `json:"LicenseExpire"`
}

type ManagerDataOemHpeFederationConfig struct {
	IPv6MulticastScope            *string `json:"IPv6MulticastScope"`
	MulticastAnnouncementInterval *int64  `json:"MulticastAnnouncementInterval"`
	MulticastDiscovery            *string `json:"MulticastDiscovery"`
	MulticastTimeToLive           *int64  `json:"MulticastTimeToLive"`
	ILOFederationManagement       *string `json:"iLOFederationManagement"`
}

type ManagerDataOemHpeFirmwareData struct {
	Date         *string `json:"Date"`
	DebugBuild   *bool   `json:"DebugBuild"`
	MajorVersion *uint64 `json:"MajorVersion"`
	MinorVersion *uint64 `json:"MinorVersion"`
	Time         *string `json:"Time"`
	Version      *string `json:"Version"`
}

type ManagerDataOemHpeFirmware struct {
	Current ManagerDataOemHpeFirmwareData `json:"Current"`
}

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

type ManagerDataOemHpe struct {
	Hpe *_managerDataOemHpe `json:"Hpe"`
}

type SecurityServiceDataOemHpeLinks struct {
	ESKM      OData `json:"ESKM"`
	HttpsCert OData `json:"HttpsCert"`
	SSO       OData `json:"SSO"`
}

type SecurityServiceDataOemHpe struct {
	Id    *string                        `json:"Id"`
	Type  *string                        `json:"Type"`
	Links SecurityServiceDataOemHpeLinks `json:"Links"`
}

type HttpsCertActionsOemHpe struct {
	GenerateCSR                LinkTargets  `json:"#HpeHttpsCert.GenerateCSR"`
	ImportCertificate          LinkTargets  `json:"#HpeHttpsCert.ImportCertificate"`
	X509CertificateInformation X509CertInfo `json:"X509CertificateInformation"`
}

type HttpsCertDataOemHpe struct {
	CSR     *string                `json:"CertificateSigningRequest"`
	Id      *string                `json:"Id"`
	Actions HttpsCertActionsOemHpe `json:"Actions"`
}

// HP(E) uses it's own privilege map instead of roles
type AccountPrivilegeMapOemHpe struct {
	Login                bool `json:"LoginPriv"`
	RemoteConsole        bool `json:"RemoteConsolePriv"`
	UserConfig           bool `json:"UserConfigPriv"`
	VirtualMedia         bool `json:"VirtualMediaPriv"`
	VirtualPowerAndReset bool `json:"VirtualPowerAndResetPriv"`
	ILOConfig            bool `json:"iLOConfigPriv"`
}
