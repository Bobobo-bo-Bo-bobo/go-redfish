package redfish

// Huawei: Oem data for Manager endpoint and SecurityService endpoint
type HttpsCertActionsOemHuawei struct {
	GenerateCSR                LinkTargets  `json:"#HpHttpsCert.GenerateCSR"`
	ImportCertificate          LinkTargets  `json:"#HttpsCert.ImportServerCertificate"`
	X509CertificateInformation X509CertInfo `json:"X509CertificateInformation"`
}

type HttpsCertDataOemHuawei struct {
	CSR     *string                   `json:"CertificateSigningRequest"`
	Id      *string                   `json:"Id"`
	Actions HttpsCertActionsOemHuawei `json:"Actions"`
}

type ManagerDataOemHuaweiLoginRule struct {
	MemberId    *string `json:"MemberId"`
	RuleEnabled bool    `json:"RuleEnabled"`
	StartTime   *string `json:"StartTime"`
	EndTime     *string `json:"EndTime"`
	IP          *string `json:"IP"`
	Mac         *string `json:"Mac"`
}
type SecurityServiceDataOemHuaweiLinks struct {
	HttpsCert OData `json:"HttpsCert"`
}

type SecurityServiceDataOemHuawei struct {
	Id    *string                           `json:"Id"`
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
	SmtpService     OData `json:"SmtpService"`
	SyslogService   OData `json:"SyslogService"`
	KvmService      OData `json:"KvmService"`
	NtpService      OData `json:"NtpService"`
}

type ManagerDataOemHuawei struct {
	Huawei _managerDataOemHuawei `json:"Huawei"`
}
