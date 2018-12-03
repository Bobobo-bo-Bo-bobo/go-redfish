package redfish

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const UserAgent string = "go-redfish/1.0.0"

type OData struct {
	Id           *string `json:"@odata.id"`
	Type         *string `json:"@odata.type"`
	Context      *string `json:"@odata.context"`
	Members      []OData `json:"Members"`
	MembersCount int     `json:"Members@odata.count"`
}

type baseEndpoint struct {
	AccountService OData `json:"AccountService"`
	Chassis        OData `json:"Chassis"`
	Managers       OData `json:"Managers"`
	SessionService OData `json:"SessionService"`
	Systems        OData `json:"Systems"`
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
	Count  int    `json:"Count"`
	Status Status `json:"Status"`
}

type SystemData struct {
	UUID         *string `json:"UUID"`
	Status       Status  `json:"Status"`
	SerialNumber *string `json:"SerialNumber"`
	//ProcessorSummary  *SystemProcessorSummary  `json:"ProcessorSummary"`
	Processors *OData  `json:"Processors"`
	PowerState *string `json:"Powerstate"`
	Name       *string `json:"Name"`
	Model      *string `json:"Model"`
	//MemorySummary *SystemMemorySummary    `json:"MemorySummary"`
	Memory             *OData  `json:"Memory"`
	Manufacturer       *string `json:"Manufacturer"`
	LogServices        *OData  `json:"LogServices"`
	Id                 *string `json:"Id"`
	EthernetInterfaces *OData  `json:"EthernetInterfaces"`
	BIOSVersion        *string `json:"BiosVersion"`
	BIOS               *OData  `json:"Bios"`
	// Actions
	SelfEndpoint *string
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

const (
	REDFISH_GENERAL uint = iota
	REDFISH_HP
	REDFISH_HUAWEI
	REDFISH_INSPUR
	REDFISH_LENOVO
	REDFISH_SUPERMICRO
)

// service processor capabilities
const (
    HAS_ACCOUNTSERVICE uint = 1 << iota // has AccountService endpoint
    HAS_SECURITYSERVICE
    HAS_ACCOUNT_ROLES
)

// map capabilities by vendor
var VendorCapabilities = map[string]uint{
    "hp": HAS_ACCOUNTSERVICE | HAS_SECURITYSERVICE,
    "huawei": HAS_ACCOUNTSERVICE | HAS_SECURITYSERVICE | HAS_ACCOUNT_ROLES,
    "inspur": 0,
    "supermicro": HAS_ACCOUNTSERVICE | HAS_ACCOUNT_ROLES,
    "vanilla": HAS_ACCOUNTSERVICE | HAS_SECURITYSERVICE | HAS_ACCOUNT_ROLES,
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

	httpRequest(string, string, *map[string]string, io.Reader, bool) (HttpResult, error)
	getCSRTarget_HP(*ManagerData) (string, error)
	getCSRTarget_Huawei(*ManagerData) (string, error)
	makeCSRPayload(CSRData) string
	makeCSRPayload_HP(CSRData) string
	makeCSRPayload_Vanilla(CSRData) string
	fetchCSR_HP(*ManagerData) (string, error)
	fetchCSR_Huawei(*ManagerData) (string, error)
	getImportCertTarget_HP(*ManagerData) (string, error)
	getImportCertTarget_Huawei(*ManagerData) (string, error)
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
	Flavor uint
    FlavorString    string
}
