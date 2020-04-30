package redfish

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// GoRedfishVersion - Library version
// Note: Be consistent with "Semantic Versioning 2.0.0" - see https://semver.org/
const GoRedfishVersion string = "1.2.1-20204030"
const _GoRedfishURL string = "https://git.ypbind.de/cgit/go-redfish/"

var userAgent = "go-redfish/" + GoRedfishVersion + " (" + _GoRedfishURL + ")"

// Error - Redfish error as defined by the standard
type Error struct {
	Error ErrorMessage `json:"error"`
}

// ErrorMessage - structure of each individual error messages
type ErrorMessage struct {
	Code                *string                    `json:"code"`
	Message             *string                    `json:"Message"`
	MessageExtendedInfo []ErrorMessageExtendedInfo `json:"@Message.ExtendedInfo"`
}

// ErrorMessageExtendedInfo - each individual error entry
type ErrorMessageExtendedInfo struct {
	MessageID         *string  `json:"MessageId"`
	Severity          *string  `json:"Severity"`
	Resolution        *string  `json:"Resolution"`
	Message           *string  `json:"Message"`
	MessageArgs       []string `json:"MessageArgs"`
	RelatedProperties []string `json:"RelatedProperties"`
}

// OData - Open Data procol struct
type OData struct {
	ID           *string `json:"@odata.id"`
	Type         *string `json:"@odata.type"`
	Context      *string `json:"@odata.context"`
	Members      []OData `json:"Members"`
	MembersCount int     `json:"Members@odata.count"`
}

type baseEndpoint struct {
	AccountService OData             `json:"AccountService"`
	Chassis        OData             `json:"Chassis"`
	Managers       OData             `json:"Managers"`
	SessionService OData             `json:"SessionService"`
	Systems        OData             `json:"Systems"`
	Links          baseEndpointLinks `json:"Links"`
}

type baseEndpointLinks struct {
	Sessions *OData `json:"Sessions"`
}

type sessionServiceEndpoint struct {
	Enabled        *bool  `json:"ServiceEnabled"`
	SessionTimeout int    `json:"SessionTimeout"`
	Sessions       *OData `json:"Sessions"`
}

// Status - Health status
type Status struct {
	State        *string `json:"State"`
	Health       *string `json:"Health"`
	HealthRollUp *string `json:"HealthRollUp"`
}

// SystemProcessorSummary - summary of system - CPU
type SystemProcessorSummary struct {
	Count  int     `json:"Count"`
	Status Status  `json:"Status"`
	Model  *string `json:"Model"`
}

// SystemMemorySummary - summary of system - memory
type SystemMemorySummary struct {
	TotalSystemMemoryGiB float64 `json:"TotalSystemMemoryGiB"`
	Status               Status  `json:"Status"`
}

// SystemActionsComputerReset - allowed action for computer reset
type SystemActionsComputerReset struct {
	Target          string   `json:"target"`
	ResetTypeValues []string `json:"ResetType@Redfish.AllowableValues"`
	ActionInfo      string   `json:"@Redfish.ActionInfo"`
}

// SystemActions - supported actions for computer reset
type SystemActions struct {
	ComputerReset *SystemActionsComputerReset `json:"#ComputerSystem.Reset"`
}

// ActionInfoParameter - informations about allowed actions
type ActionInfoParameter struct {
	Name            string   `json:"Name"`
	Required        bool     `json:"Required"`
	DataType        string   `json:"DataType"`
	AllowableValues []string `json:"AllowableValues"`
}

// SystemActionInfo - information about allowed actions
type SystemActionInfo struct {
	ODataContext string                `json:"@odata.context"`
	ODataID      string                `json:"@odata.id"`
	ODataType    string                `json:"@odata.type"`
	Parameters   []ActionInfoParameter `json:"Parameters"`
}

// SystemData - System information
type SystemData struct {
	UUID               *string                 `json:"UUID"`
	Status             Status                  `json:"Status"`
	SerialNumber       *string                 `json:"SerialNumber"`
	ProcessorSummary   *SystemProcessorSummary `json:"ProcessorSummary"`
	Processors         *OData                  `json:"Processors"`
	PowerState         *string                 `json:"Powerstate"`
	Name               *string                 `json:"Name"`
	Model              *string                 `json:"Model"`
	MemorySummary      *SystemMemorySummary    `json:"MemorySummary"`
	Memory             *OData                  `json:"Memory"`
	Manufacturer       *string                 `json:"Manufacturer"`
	LogServices        *OData                  `json:"LogServices"`
	ID                 *string                 `json:"Id"`
	EthernetInterfaces *OData                  `json:"EthernetInterfaces"`
	BIOSVersion        *string                 `json:"BiosVersion"`
	BIOS               *OData                  `json:"Bios"`
	Actions            *SystemActions          `json:"Actions"`
	Oem                json.RawMessage         `json:"Oem"`
	SelfEndpoint       *string
	// map normalized (converted to lowercase) to supported reset types
	allowedResetTypes map[string]string
	// name of the reset type property, usually "ResetType", but may vary (e.g. when specified otherwise in @Redfish.ActionInfo)
	resetTypeProperty string
}

// AccountService - Account handling
type AccountService struct {
	AccountsEndpoint *OData `json:"Accounts"`
	RolesEndpoint    *OData `json:"Roles"`
}

// AccountData - individual accounts
type AccountData struct {
	ID       *string `json:"Id"`
	Name     *string `json:"Name"`
	UserName *string `json:"UserName"`
	Password *string `json:"Password"`
	RoleID   *string `json:"RoleId"`
	Enabled  *bool   `json:"Enabled"`
	Locked   *bool   `json:"Locked"`

	SelfEndpoint *string
}

// RoleData - individual roles
type RoleData struct {
	ID                 *string  `json:"Id"`
	Name               *string  `json:"Name"`
	IsPredefined       *bool    `json:"IsPredefined"`
	Description        *string  `json:"Description"`
	AssignedPrivileges []string `json:"AssignedPrivileges"`
	//    OemPrivileges   []string    `json:"OemPrivileges"`
	SelfEndpoint *string
}

// ChassisData - Chassis information
type ChassisData struct {
	ID           *string         `json:"Id"`
	Name         *string         `json:"Name"`
	ChassisType  *string         `json:"ChassisType"`
	Manufacturer *string         `json:"Manufacturer"`
	Model        *string         `json:"Model"`
	SerialNumber *string         `json:"SerialNumber"`
	PartNumber   *string         `json:"PartNumber"`
	AssetTag     *string         `json:"AssetTag"`
	IndicatorLED *string         `json:"IndicatorLED"`
	Status       Status          `json:"Status"`
	Oem          json.RawMessage `json:"Oem"`
	Thermal      *OData          `json:"Thermal"`
	Power        *OData          `json:"Power"`

	SelfEndpoint *string
}

// TemperatureData - temperature readings
type TemperatureData struct {
	ODataID                   *string `json:"@odata.id"`
	MemberID                  *string `json:"MemberId"`
	SensorNumber              *int    `json:"SensorNumber"`
	Name                      *string `json:"Name"`
	ReadingCelsius            *int    `json:"ReadingCelsius"`
	LowerThresholdNonCritical *int    `json:"LowerThresholdNonCritical"`
	LowerThresholdCritical    *int    `json:"LowerThresholdCritical"`
	LowerThresholdFatal       *int    `json:"LowerThresholdFatal"`
	UpperThresholdNonCritical *int    `json:"UpperThresholdNonCritical"`
	UpperThresholdCritical    *int    `json:"UpperThresholdCritical"`
	UpperThresholdFatal       *int    `json:"UpperThresholdFatal"`
	MinReadingRangeTemp       *int    `json:"MinReadingRangeTemp"`
	MaxReadingRangeTemp       *int    `json:"MaxReadingRangeTemp"`
	Status                    Status  `json:"Status"`
}

// FanData - fan readings
type FanData struct {
	ODataID                   *string         `json:"@odata.id"`
	MemberID                  *string         `json:"MemberId"`
	SensorNumber              *int            `json:"SensorNumber"`
	FanName                   *string         `json:"FanName"`
	Name                      *string         `json:"Name"`
	Reading                   *int            `json:"Reading"`
	LowerThresholdNonCritical *int            `json:"LowerThresholdNonCritical"`
	LowerThresholdCritical    *int            `json:"LowerThresholdCritical"`
	LowerThresholdFatal       *int            `json:"LowerThresholdFatal"`
	UpperThresholdNonCritical *int            `json:"UpperThresholdNonCritical"`
	UpperThresholdCritical    *int            `json:"UpperThresholdCritical"`
	UpperThresholdFatal       *int            `json:"UpperThresholdFatal"`
	MinReadingRange           *int            `json:"MinReadingRange"`
	MaxReadingRange           *int            `json:"MaxReadingRange"`
	Status                    Status          `json:"Status"`
	ReadingUnits              *string         `json:"ReadingUnits"`
	PartNumber                *string         `json:"PartNumber"`
	PhysicalContext           *string         `json:"PhysicalContext"`
	Oem                       json.RawMessage `json:"Oem"`
}

// ThermalData - thermal data
type ThermalData struct {
	ODataID      *string           `json:"@odata.id"`
	ID           *string           `json:"Id"`
	Status       Status            `json:"Status"`
	Temperatures []TemperatureData `json:"Temperatures"`
	Fans         []FanData         `json:"Fans"`
	SelfEndpoint *string
}

// PowerMetricsData - current power data/metrics
type PowerMetricsData struct {
	MinConsumedWatts     *int `json:"MinConsumedWatts"`
	MaxConsumedWatts     *int `json:"MaxConsumedWatts"`
	AverageConsumedWatts *int `json:"AverageConsumedWatts"`
	IntervalInMin        *int `json:"IntervalInMin"`
}

// PowerLimitData - defined power limits
type PowerLimitData struct {
	LimitInWatts   *int    `json:"LimitInWatts"`
	LimitException *string `json:"LimitException"`
}

// PowerControlData - power control information
type PowerControlData struct {
	ID                 *string `json:"@odata.id"`
	MemberID           *string `json:"MemberId"`
	Name               *string `json:"Name"`
	PowerConsumedWatts *int
	PowerMetrics       PowerMetricsData `json:"PowerMetrics"`
	PowerLimit         PowerLimitData   `json:"PowerLimit"`
	Status             Status           `json:"Status"`
	Oem                json.RawMessage  `json:"Oem"`
}

// VoltageData - voltage information
type VoltageData struct {
	ODataID                   *string  `json:"@odata.id"`
	MemberID                  *string  `json:"MemberId"`
	Name                      *string  `json:"Name"`
	SensorNumber              *int     `json:"SensorNumber"`
	Status                    Status   `json:"Status"`
	ReadingVolts              *float64 `json:"ReadingVolts"`
	UpperThresholdNonCritical *float64 `json:"UpperThresholdNonCritical"`
	UpperThresholdCritical    *float64 `json:"UpperThresholdCritical"`
	UpperThresholdFatal       *float64 `json:"UpperThresholdFatal"`
	LowerThresholdNonCritical *float64 `json:"LowerThresholdNonCritical"`
	LowerThresholdCritical    *float64 `json:"LowerThresholdCritical"`
	LowerThresholdFatal       *float64 `json:"LowerThresholdFatal"`
	MinReadingRange           *float64 `json:"MinReadingRange"`
	MaxReadingRange           *float64 `json:"MaxReadingRange"`
	PhysicalContext           *string  `json:"PhysicalContext"`
}

// PSUData - power supply information
type PSUData struct {
	ODataID              *string         `json:"@odata.id"`
	MemberID             *string         `json:"MemberId"`
	Name                 *string         `json:"Name"`
	Status               Status          `json:"Status"`
	PowerSupplyType      *string         `json:"PowerSupplyType"`
	LineInputVoltage     *int            `json:"LineInputVoltage"`
	LineInputVoltageType *string         `json:"LineInputVoltageType"`
	LastPowerOutputWatts *int            `json:"LastPowerOutputWatts"`
	PowerCapacityWatts   *int            `json:"PowerCapacityWatts"`
	Model                *string         `json:"Model"`
	FirmwareVersion      *string         `json:"FirmwareVersion"`
	SerialNumber         *string         `json:"SerialNumber"`
	Manufacturer         *string         `json:"Manufacturer"`
	PartNumber           *string         `json:"PartNumber"`
	Oem                  json.RawMessage `json:"Oem"`
	Redundancy           []OData         `json:"Redundancy"`
}

// PowerData - power data
type PowerData struct {
	OdataID       *string            `json:"@odata.id"`
	Context       *string            `json:"@odata.context"`
	ID            *string            `json:"Id"`
	PowerControl  []PowerControlData `json:"PowerControl"`
	Voltages      []VoltageData      `json:"Voltages"`
	PowerSupplies []PSUData          `json:"PowerSupplies"`
	SelfEndpoint  *string
}

// ManagerLicenseData - license data for management board
type ManagerLicenseData struct {
	Name       string
	Expiration string
	Type       string
	License    string
}

// ManagerActionsData - list of allowed actions of the management processor
type ManagerActionsData struct {
	ManagerReset LinkTargets `json:"#Manager.Reset"`
}

// ManagerData - information about the management processor
type ManagerData struct {
	ID              *string         `json:"Id"`
	Name            *string         `json:"Name"`
	ManagerType     *string         `json:"ManagerType"`
	UUID            *string         `json:"UUID"`
	Status          Status          `json:"Status"`
	FirmwareVersion *string         `json:"FirmwareVersion"`
	Oem             json.RawMessage `json:"Oem"`
	Actions         json.RawMessage `json:"Actions"` // may contain vendor specific data and endpoints

	/* futher data
	   VirtualMedia
	   SerialConsole
	   NetworkProtocol
	   GraphicalConsole
	   FirmwareVersion
	   EthernetInterfaces
	   Actions
	*/

	SelfEndpoint *string
}

// X509CertInfo - X509 certificate information
type X509CertInfo struct {
	Issuer         *string `json:"Issuer"`
	SerialNumber   *string `json:"SerialNumber"`
	Subject        *string `json:"Subject"`
	ValidNotAfter  *string `json:"ValidNotAfter"`
	ValidNotBefore *string `json:"ValidNotBefore"`
}

// LinkTargets - available link targets
type LinkTargets struct {
	Target     *string `json:"target"`
	ActionInfo *string `json:"@Redfish.ActionInfo"`
}

// CSRData - data for CSR subject
type CSRData struct {
	C  string // Country
	S  string // State or province
	L  string // Locality or city
	O  string // Organisation
	OU string // Organisational unit
	CN string // Common name
}

// AccountCreateData - data for account creation
type AccountCreateData struct {
	UserName string `json:",omitempty"`
	Password string `json:",omitempty"`
	// for service processors supporting roles
	Role string `json:"RoleId,omitempty"`

	Enabled *bool `json:",omitempty"`
	Locked  *bool `json:",omitempty"`

	// for HP(E) iLO which supports Oem specific
	// Note: OemHpPrivilegeMap is an _internal_ struct but must be exported for json.Marshal !
	//       Don't use this structm use HPEPrivileges instead
	OemHpPrivilegeMap *AccountPrivilegeMapOemHp `json:",omitempty"`
	HPEPrivileges     uint
}

// Redfish vendor flavors
const (
	RedfishFlavorNotInitialized uint = iota
	RedfishGeneral
	RedfishHP
	RedfishHPE
	RedfishHuawei
	RedfishInspur
	RedfishLenovo
	RedfishSuperMicro
	RedfishDell
)

// service processor capabilities
const (
	HasAccountService uint = 1 << iota // has AccountService endpoint
	HasSecurityService
	HasAccountRoles
	HasChassis
	HasLicense
)

// VendorCapabilities - map capabilities by vendor
var VendorCapabilities = map[string]uint{
	"hp":         HasAccountService | HasSecurityService | HasChassis | HasLicense,
	"hpe":        HasAccountService | HasSecurityService | HasChassis | HasLicense,
	"huawei":     HasAccountService | HasSecurityService | HasAccountRoles | HasChassis,
	"inspur":     0,
	"supermicro": HasAccountService | HasAccountRoles | HasChassis,
	"dell":       HasAccountService | HasAccountRoles | HasChassis,
	"lenovo":     HasChassis,
	"vanilla":    HasAccountService | HasSecurityService | HasAccountRoles | HasChassis,
	"":           HasAccountService | HasSecurityService | HasAccountRoles | HasChassis,
}

// HTTPResult - result of the http_request calls
type HTTPResult struct {
	URL        string
	StatusCode int
	Status     string
	Header     http.Header
	Content    []byte
}

// BaseRedfish - interface definition
type BaseRedfish interface {
	Initialize() error
	Login() error
	Logout() error
	GetSystems() ([]string, error)
	GetSystemData(string) (*SystemData, error)
	MapSystensByID() (map[string]*SystemData, error)
	MapSystemsByUUID() (map[string]*SystemData, error)
	MapSystemsBySerialNumber() (map[string]*SystemData, error)
	GetAccounts() ([]string, error)
	GetAccountData(string) (*AccountData, error)
	MapAccountsByName() (map[string]*AccountData, error)
	MapAccountsByID() (map[string]*AccountData, error)
	GetRoles() ([]string, error)
	GetRoleData(string) (*AccountData, error)
	MapRolesByName() (map[string]*RoleData, error)
	MapRolesByID() (map[string]*RoleData, error)
	GenCSR(CSRData) error
	FetchCSR() (string, error)
	ImportCertificate(string) error
	ResetSP() error
	GetVendorFlavor() error
	AddAccount(AccountCreateData) error
	ModifyAccount(string, AccountCreateData) error
	DeleteAccount(string) error
	ChangePassword(string, string) error
	SetSystemPowerState(*SystemData, string) error
	ProcessError(HTTPResult) (*Error, error)
	GetLicense(*ManagerData) (*ManagerLicenseData, error)
	GetErrorMessage(*Error) string
	IsInitialised() bool
	Clone() *Redfish
	GetManagers() ([]string, error)
	GetManagerData(string) (*ManagerData, error)
	MapManagersByID() (map[string]*ManagerData, error)
	MapManagersByUUID() (map[string]*ManagerData, error)

	httpRequest(string, string, *map[string]string, io.Reader, bool) (HTTPResult, error)
	getCSRTarget_HP(*ManagerData) (string, error)
	getCSRTarget_HPE(*ManagerData) (string, error)
	getCSRTarget_Huawei(*ManagerData) (string, error)
	makeCSRPayload(CSRData) string
	makeCSRPayload_HP(CSRData) string
	makeCSRPayload_Vanilla(CSRData) string
	fetchCSR_HP(*ManagerData) (string, error)
	fetchCSR_HPE(*ManagerData) (string, error)
	fetchCSR_Huawei(*ManagerData) (string, error)
	getImportCertTarget_HP(*ManagerData) (string, error)
	getImportCertTarget_HPE(*ManagerData) (string, error)
	getImportCertTarget_Huawei(*ManagerData) (string, error)
	makeAccountCreateModifyPayload(AccountCreateData) (string, error)
	setAllowedResetTypes(*SystemData) error
	hpGetLicense(*ManagerData) (*ManagerLicenseData, error)
	hpeGetLicense(*ManagerData) (*ManagerLicenseData, error)
	hpHpePrepareLicensePayload([]byte) string
}

// Redfish - object to access Redfish API
type Redfish struct {
	Hostname        string
	Port            int
	Username        string
	Password        string
	AuthToken       *string
	SessionLocation *string
	Timeout         time.Duration
	InsecureSSL     bool
	Debug           bool
	Verbose         bool
	RawBaseContent  string

	// endpoints
	AccountService string
	Chassis        string
	Managers       string
	SessionService string
	Sessions       string
	Systems        string

	// Vendor "flavor"
	Flavor       uint
	FlavorString string

	initialised bool
}
