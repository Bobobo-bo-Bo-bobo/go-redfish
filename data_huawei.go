package redfish

// HTTPSCertActionsOemHuawei - Huawei: Oem data for Manager endpoint and SecurityService endpoint
type HTTPSCertActionsOemHuawei struct {
	GenerateCSR                LinkTargets  `json:"#HpHttpsCert.GenerateCSR"`
	ImportCertificate          LinkTargets  `json:"#HttpsCert.ImportServerCertificate"`
	X509CertificateInformation X509CertInfo `json:"X509CertificateInformation"`
}

// HTTPSCertDataOemHuawei - OEM specific extension
type HTTPSCertDataOemHuawei struct {
	CSR     *string                   `json:"CertificateSigningRequest"`
	ID      *string                   `json:"Id"`
	Actions HTTPSCertActionsOemHuawei `json:"Actions"`
}

// ManagerDataOemHuaweiLoginRule - Huawei specific login rules
type ManagerDataOemHuaweiLoginRule struct {
	MemberID    *string `json:"MemberId"`
	RuleEnabled bool    `json:"RuleEnabled"`
	StartTime   *string `json:"StartTime"`
	EndTime     *string `json:"EndTime"`
	IP          *string `json:"IP"`
	Mac         *string `json:"Mac"`
}

// SecurityServiceDataOemHuaweiLinks - OEM specific extension
type SecurityServiceDataOemHuaweiLinks struct {
	HTTPSCert OData `json:"HttpsCert"`
}

// SecurityServiceDataOemHuawei - OEM specific extension
type SecurityServiceDataOemHuawei struct {
	ID    *string                           `json:"Id"`
	Name  *string                           `json:"Name"`
	Links SecurityServiceDataOemHuaweiLinks `json:"Links"`
}

type _managerDataOemHuawei struct {
	BMCUpTime       *string                         `json:"BMCUpTime"`
	ProductUniqueID *string                         `json:"ProductUniqueID"`
	PlatformType    *string                         `json:"PlatformType"`
	LoginRule       []ManagerDataOemHuaweiLoginRule `json:"LoginRule"`

	SecurityService OData `json:"SecurityService"`
	SnmpService     OData `json:"SnmpService"`
	SMTPService     OData `json:"SmtpService"`
	SyslogService   OData `json:"SyslogService"`
	KvmService      OData `json:"KvmService"`
	NtpService      OData `json:"NtpService"`
}

// ManagerDataOemHuawei - OEM specific extension
type ManagerDataOemHuawei struct {
	Huawei _managerDataOemHuawei `json:"Huawei"`
}
