package redfish

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Note: Be consistent with "Semantic Versioning 2.0.0" - see https://semver.org/
const GoRedfishVersion string = "1.2.1-2019.10.13"
const _GoRedfishUrl string = "https://git.ypbind.de/cgit/go-redfish/"

var UserAgent string = "go-redfish/" + GoRedfishVersion + " (" + _GoRedfishUrl + ")"

type RedfishError struct {
	Error RedfishErrorMessage `json:"error"`
}

type RedfishErrorMessage struct {
	Code                *string                           `json:"code"`
	Message             *string                           `json:"Message"`
	MessageExtendedInfo []RedfishErrorMessageExtendedInfo `json:"@Message.ExtendedInfo"`
}

type RedfishErrorMessageExtendedInfo struct {
	MessageId         *string  `json:"MessageId"`
	Severity          *string  `json:"Severity"`
	Resolution        *string  `json:"Resolution"`
	Message           *string  `json:"Message"`
	MessageArgs       []string `json:"MessageArgs"`
	RelatedProperties []string `json:"RelatedProperties"`
}

type OData struct {
	Id           *string `json:"@odata.id"`
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

type Status struct {
	State        *string `json:"State"`
	Health       *string `json:"Health"`
	HealthRollUp *string `json:"HealthRollUp"`
}

type SystemProcessorSummary struct {
	Count  int     `json:"Count"`
	Status Status  `json:"Status"`
	Model  *string `json:"Model"`
}

type SystemMemorySummary struct {
	TotalSystemMemoryGiB float64 `json:"TotalSystemMemoryGiB"`
	Status               Status  `json:"Status"`
}

type SystemActionsComputerReset struct {
	Target          string   `json:"target"`
	ResetTypeValues []string `json:"ResetType@Redfish.AllowableValues"`
	ActionInfo      string   `json:"@Redfish.ActionInfo"`
}

type SystemActions struct {
	ComputerReset *SystemActionsComputerReset `json:"#ComputerSystem.Reset"`
}

type ActionInfoParameter struct {
	Name            string   `json:"Name"`
	Required        bool     `json:"Required"`
	DataType        string   `json:"DataType"`
	AllowableValues []string `json:"AllowableValues"`
}

type SystemActionInfo struct {
	ODataContext string                `json:"@odata.context"`
	ODataId      string                `json:"@odata.id"`
	ODataType    string                `json:"@odata.type"`
	Parameters   []ActionInfoParameter `json:"Parameters"`
}

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
	Id                 *string                 `json:"Id"`
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

type AccountService struct {
	AccountsEndpoint *OData `json:"Accounts"`
	RolesEndpoint    *OData `json:"Roles"`
}

type AccountData struct {
	Id       *string `json:"Id"`
	Name     *string `json:"Name"`
	UserName *string `json:"UserName"`
	Password *string `json:"Password"`
	RoleId   *string `json:"RoleId"`
	Enabled  *bool   `json:"Enabled"`
	Locked   *bool   `json:"Locked"`

	SelfEndpoint *string
}

type RoleData struct {
	Id                 *string  `json:"Id"`
	Name               *string  `json:"Name"`
	IsPredefined       *bool    `json:"IsPredefined"`
	Description        *string  `json:"Description"`
	AssignedPrivileges []string `json:"AssignedPrivileges"`
	//    OemPrivileges   []string    `json:"OemPrivileges"`
	SelfEndpoint *string
}

type ChassisData struct {
	Id           *string         `json:"Id"`
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

type TemperatureData struct {
	ODataId                   *string `json:"@odata.id"`
	MemberId                  *string `json:"MemberId"`
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

type FanData struct {
	ODataId                   *string         `json:"@odata.id"`
	MemberId                  *string         `json:"MemberId"`
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

type ThermalData struct {
	ODataId      *string           `json:"@odata.id"`
	Id           *string           `json:"Id"`
	Status       Status            `json:"Status"`
	Temperatures []TemperatureData `json:"Temperatures"`
	Fans         []FanData         `json:"Fans"`
	SelfEndpoint *string
}

type PowerMetricsData struct {
	MinConsumedWatts     *int `json:"MinConsumedWatts"`
	MaxConsumedWatts     *int `json:"MaxConsumedWatts"`
	AverageConsumedWatts *int `json:"AverageConsumedWatts"`
	IntervalInMin        *int `json:"IntervalInMin"`
}

type PowerLimitData struct {
	LimitInWatts   *int    `json:"LimitInWatts"`
	LimitException *string `json:"LimitException"`
}

type PowerControlData struct {
	Id                 *string `json:"@odata.id"`
	MemberId           *string `json:"MemberId"`
	Name               *string `json:"Name"`
	PowerConsumedWatts *int
	PowerMetrics       PowerMetricsData `json:"PowerMetrics"`
	PowerLimit         PowerLimitData   `json:"PowerLimit"`
	Status             Status           `json:"Status"`
	Oem                json.RawMessage  `json:"Oem"`
}

type VoltageData struct {
	ODataId                   *string  `json:"@odata.id"`
	MemberId                  *string  `json:"MemberId"`
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

type PSUData struct {
	ODataId              *string         `json:"@odata.id"`
	MemberId             *string         `json:"MemberId"`
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

type PowerData struct {
	OdataId       *string            `json:"@odata.id"`
	Context       *string            `json:"@odata.context"`
	Id            *string            `json:"Id"`
	PowerControl  []PowerControlData `json:"PowerControl"`
	Voltages      []VoltageData      `json:"Voltages"`
	PowerSupplies []PSUData          `json:"PowerSupplies"`
	SelfEndpoint  *string
}

type ManagerLicenseData struct {
	Name       string
	Expiration string
	Type       string
	License    string
}

type ManagerActionsData struct {
	ManagerReset LinkTargets `json:"#Manager.Reset"`
}

type ManagerData struct {
	Id              *string         `json:"Id"`
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

type X509CertInfo struct {
	Issuer         *string `json:"Issuer"`
	SerialNumber   *string `json:"SerialNumber"`
	Subject        *string `json:"Subject"`
	ValidNotAfter  *string `json:"ValidNotAfter"`
	ValidNotBefore *string `json:"ValidNotBefore"`
}

type LinkTargets struct {
	Target     *string `json:"target"`
	ActionInfo *string `json:"@Redfish.ActionInfo"`
}

// data for CSR subject
type CSRData struct {
	C  string // Country
	S  string // State or province
	L  string // Locality or city
	O  string // Organisation
	OU string // Organisational unit
	CN string // Common name
}

// data for account creation
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

const (
	REDFISH_FLAVOR_NOT_INITIALIZED uint = iota
	REDFISH_GENERAL
	REDFISH_HP
	REDFISH_HPE
	REDFISH_HUAWEI
	REDFISH_INSPUR
	REDFISH_LENOVO
	REDFISH_SUPERMICRO
	REDFISH_DELL
)

// service processor capabilities
const (
	HAS_ACCOUNTSERVICE uint = 1 << iota // has AccountService endpoint
	HAS_SECURITYSERVICE
	HAS_ACCOUNT_ROLES
	HAS_CHASSIS
	HAS_LICENSE
)

// map capabilities by vendor
var VendorCapabilities = map[string]uint{
	"hp":         HAS_ACCOUNTSERVICE | HAS_SECURITYSERVICE | HAS_CHASSIS | HAS_LICENSE,
	"hpe":        HAS_ACCOUNTSERVICE | HAS_SECURITYSERVICE | HAS_CHASSIS | HAS_LICENSE,
	"huawei":     HAS_ACCOUNTSERVICE | HAS_SECURITYSERVICE | HAS_ACCOUNT_ROLES | HAS_CHASSIS,
	"inspur":     0,
	"supermicro": HAS_ACCOUNTSERVICE | HAS_ACCOUNT_ROLES | HAS_CHASSIS,
	"dell":       HAS_ACCOUNTSERVICE | HAS_ACCOUNT_ROLES | HAS_CHASSIS,
	"lenovo":     HAS_CHASSIS,
	"vanilla":    HAS_ACCOUNTSERVICE | HAS_SECURITYSERVICE | HAS_ACCOUNT_ROLES | HAS_CHASSIS,
	"":           HAS_ACCOUNTSERVICE | HAS_SECURITYSERVICE | HAS_ACCOUNT_ROLES | HAS_CHASSIS,
}

type HttpResult struct {
	Url        string
	StatusCode int
	Status     string
	Header     http.Header
	Content    []byte
}

type BaseRedfish interface {
	Initialize() error
	Login() error
	Logout() error
	GetSystems() ([]string, error)
	GetSystemData(string) (*SystemData, error)
	MapSystensById() (map[string]*SystemData, error)
	MapSystemsByUuid() (map[string]*SystemData, error)
	MapSystemsBySerialNumber() (map[string]*SystemData, error)
	GetAccounts() ([]string, error)
	GetAccountData(string) (*AccountData, error)
	MapAccountsByName() (map[string]*AccountData, error)
	MapAccountsById() (map[string]*AccountData, error)
	GetRoles() ([]string, error)
	GetRoleData(string) (*AccountData, error)
	MapRolesByName() (map[string]*RoleData, error)
	MapRolesById() (map[string]*RoleData, error)
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
	ProcessError(HttpResult) (*RedfishError, error)
	GetLicense(*ManagerData) (*ManagerLicenseData, error)
	GetErrorMessage(*RedfishError) string
	IsInitialised() bool
	Clone() *Redfish

	httpRequest(string, string, *map[string]string, io.Reader, bool) (HttpResult, error)
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
