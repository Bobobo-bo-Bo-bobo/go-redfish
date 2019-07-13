package redfish

// HP/HPE: Oem data for Manager endpoint and SecurityService endpoint
type ManagerDataOemHpLicense struct {
	Key    *string `json:"LicenseKey"`
	String *string `json:"LicenseString"`
	Type   *string `json:"LicenseType"`
}

type ManagerDataOemHpFederationConfig struct {
	IPv6MulticastScope            *string `json:"IPv6MulticastScope"`
	MulticastAnnouncementInterval *int64  `json:"MulticastAnnouncementInterval"`
	MulticastDiscovery            *string `json:"MulticastDiscovery"`
	MulticastTimeToLive           *int64  `json:"MulticastTimeToLive"`
	ILOFederationManagement       *string `json:"iLOFederationManagement"`
}

type ManagerDataOemHpFirmwareData struct {
	Date         *string `json:"Date"`
	DebugBuild   *bool   `json:"DebugBuild"`
	MajorVersion *uint64 `json:"MajorVersion"`
	MinorVersion *uint64 `json:"MinorVersion"`
	Time         *string `json:"Time"`
	Version      *string `json:"Version"`
}

type ManagerDataOemHpFirmware struct {
	Current ManagerDataOemHpFirmwareData `json:"Current"`
}

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

type ManagerDataOemHp struct {
	Hp _managerDataOemHp `json:"Hp"`
}

type SecurityServiceDataOemHpLinks struct {
	ESKM      OData `json:"ESKM"`
	HttpsCert OData `json:"HttpsCert"`
	SSO       OData `json:"SSO"`
}

type SecurityServiceDataOemHp struct {
	Id    *string                       `json:"Id"`
	Type  *string                       `json:"Type"`
	Links SecurityServiceDataOemHpLinks `json:"Links"`
}

type HttpsCertActionsOemHp struct {
	GenerateCSR                LinkTargets  `json:"#HpHttpsCert.GenerateCSR"`
	ImportCertificate          LinkTargets  `json:"#HpHttpsCert.ImportCertificate"`
	X509CertificateInformation X509CertInfo `json:"X509CertificateInformation"`
}

type HttpsCertDataOemHp struct {
	CSR     *string               `json:"CertificateSigningRequest"`
	Id      *string               `json:"Id"`
	Actions HttpsCertActionsOemHp `json:"Actions"`
}

// HP(E) uses it's own privilege map instead of roles
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
	HPE_PRIVILEGE_NONE       = 0
	HPE_PRIVILEGE_LOGIN uint = 1 << iota
	HPE_PRIVILEGE_REMOTECONSOLE
	HPE_PRIVILEGE_USERCONFIG
	HPE_PRIVILEGE_VIRTUALMEDIA
	HPE_PRIVILEGE_VIRTUALPOWER_AND_RESET
	HPE_PRIVILEGE_ILOCONFIG
)

// map privilege names to flags
var HPEPrivilegeMap = map[string]uint{
	"login":                HPE_PRIVILEGE_LOGIN,
	"remoteconsole":        HPE_PRIVILEGE_REMOTECONSOLE,
	"userconfig":           HPE_PRIVILEGE_USERCONFIG,
	"virtualmedia":         HPE_PRIVILEGE_VIRTUALMEDIA,
	"virtualpowerandreset": HPE_PRIVILEGE_VIRTUALPOWER_AND_RESET,
	"iloconfig":            HPE_PRIVILEGE_ILOCONFIG,
}

// "Virtual" roles with predefined privilege map set
var HPEVirtualRoles = map[string]uint{
	"none":          0,
	"readonly":      HPE_PRIVILEGE_LOGIN,
	"operator":      HPE_PRIVILEGE_LOGIN | HPE_PRIVILEGE_REMOTECONSOLE | HPE_PRIVILEGE_VIRTUALMEDIA | HPE_PRIVILEGE_VIRTUALPOWER_AND_RESET,
	"administrator": HPE_PRIVILEGE_LOGIN | HPE_PRIVILEGE_REMOTECONSOLE | HPE_PRIVILEGE_USERCONFIG | HPE_PRIVILEGE_VIRTUALMEDIA | HPE_PRIVILEGE_VIRTUALPOWER_AND_RESET | HPE_PRIVILEGE_ILOCONFIG,
}
