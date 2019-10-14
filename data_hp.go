package redfish

// ManagerDataOemHpLicense - HP/HPE: Oem data for Manager endpoint and SecurityService endpoint
type ManagerDataOemHpLicense struct {
	Key    *string `json:"LicenseKey"`
	String *string `json:"LicenseString"`
	Type   *string `json:"LicenseType"`
	Expire *string `json:"LicenseExpire"`
}

// ManagerDataOemHpFederationConfig - same as ManagerDataOemHpeFederationConfig
type ManagerDataOemHpFederationConfig struct {
	IPv6MulticastScope            *string `json:"IPv6MulticastScope"`
	MulticastAnnouncementInterval *int64  `json:"MulticastAnnouncementInterval"`
	MulticastDiscovery            *string `json:"MulticastDiscovery"`
	MulticastTimeToLive           *int64  `json:"MulticastTimeToLive"`
	ILOFederationManagement       *string `json:"iLOFederationManagement"`
}

// ManagerDataOemHpFirmwareData - same as ManagerDataOemHpeFirmwareData
type ManagerDataOemHpFirmwareData struct {
	Date         *string `json:"Date"`
	DebugBuild   *bool   `json:"DebugBuild"`
	MajorVersion *uint64 `json:"MajorVersion"`
	MinorVersion *uint64 `json:"MinorVersion"`
	Time         *string `json:"Time"`
	Version      *string `json:"Version"`
}

// ManagerDataOemHpFirmware - same as ManagerDataOemHpeFirmware
type ManagerDataOemHpFirmware struct {
	Current ManagerDataOemHpFirmwareData `json:"Current"`
}

// ManagerDataOemHpLinks - same as ManagerDataOemHpeLinks
// NOTE: The result for HP/HPE are different depending if the HTTP header
// OData-Version is set or not. If OData-Version is _NOT_ set data are returned in
// .Oem.Hp.links with endpoints in "href". If OData-Version is set
// data are returned in .Oem.Hp.Links (note different case!) and endpoints are
// defined as @odata.id. We always set "OData-Version: 4.0"
type ManagerDataOemHpLinks struct {
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

type _managerDataOemHp struct {
	FederationConfig ManagerDataOemHpFederationConfig `json:"FederationConfig"`
	Firmware         ManagerDataOemHpFirmware         `json:"Firmware"`
	License          ManagerDataOemHpLicense          `json:"License"`
	Type             *string                          `json:"Type"`
	Links            ManagerDataOemHpLinks            `json:"Links"`
}

// ManagerDataOemHp - same as ManagerDataOemHpe
type ManagerDataOemHp struct {
	Hp *_managerDataOemHp `json:"Hp"`
}

// SecurityServiceDataOemHpLinks - same as SecurityServiceDataOemHpeLinks
type SecurityServiceDataOemHpLinks struct {
	ESKM      OData `json:"ESKM"`
	HTTPSCert OData `json:"HttpsCert"`
	SSO       OData `json:"SSO"`
}

// SecurityServiceDataOemHp - same as SecurityServiceDataOemHpe
type SecurityServiceDataOemHp struct {
	ID    *string                       `json:"Id"`
	Type  *string                       `json:"Type"`
	Links SecurityServiceDataOemHpLinks `json:"Links"`
}

// HTTPSCertActionsOemHp - same as HTTPSCertActionsOemHpe
type HTTPSCertActionsOemHp struct {
	GenerateCSR                LinkTargets  `json:"#HpHttpsCert.GenerateCSR"`
	ImportCertificate          LinkTargets  `json:"#HpHttpsCert.ImportCertificate"`
	X509CertificateInformation X509CertInfo `json:"X509CertificateInformation"`
}

// HTTPSCertDataOemHp - same as HTTPSCertDataOemHpe
type HTTPSCertDataOemHp struct {
	CSR     *string               `json:"CertificateSigningRequest"`
	ID      *string               `json:"Id"`
	Actions HTTPSCertActionsOemHp `json:"Actions"`
}

// AccountPrivilegeMapOemHp - HP(E) uses it's own privilege map instead of roles
type AccountPrivilegeMapOemHp struct {
	Login                bool `json:"LoginPriv"`
	RemoteConsole        bool `json:"RemoteConsolePriv"`
	UserConfig           bool `json:"UserConfigPriv"`
	VirtualMedia         bool `json:"VirtualMediaPriv"`
	VirtualPowerAndReset bool `json:"VirtualPowerAndResetPriv"`
	ILOConfig            bool `json:"iLOConfigPriv"`
}

// Bitset for privileges
const (
	HpePrivilegeNone       = 0
	HpePrivilegeLogin uint = 1 << iota
	HpePrivilegeRemoteConsole
	HpePrivilegeUserConfig
	HpePrivilegeVirtualMedia
	HpePrivilegeVirtualPowerAndReset
	HpePrivilegeIloConfig
)

// HPEPrivilegeMap - map privilege names to flags
var HPEPrivilegeMap = map[string]uint{
	"login":                HpePrivilegeLogin,
	"remoteconsole":        HpePrivilegeRemoteConsole,
	"userconfig":           HpePrivilegeUserConfig,
	"virtualmedia":         HpePrivilegeVirtualMedia,
	"virtualpowerandreset": HpePrivilegeVirtualPowerAndReset,
	"iloconfig":            HpePrivilegeIloConfig,
}

// HPEVirtualRoles - "Virtual" roles with predefined privilege map set
var HPEVirtualRoles = map[string]uint{
	"none":          0,
	"readonly":      HpePrivilegeLogin,
	"operator":      HpePrivilegeLogin | HpePrivilegeRemoteConsole | HpePrivilegeVirtualMedia | HpePrivilegeVirtualPowerAndReset,
	"administrator": HpePrivilegeLogin | HpePrivilegeRemoteConsole | HpePrivilegeUserConfig | HpePrivilegeVirtualMedia | HpePrivilegeVirtualPowerAndReset | HpePrivilegeIloConfig,
}
