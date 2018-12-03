package redfish

// Supermicro: Oem data for Manager.Actions endpoint
type ManagerActionsDataOemSupermicro struct {
	Oem _managerActionsDataOemSupermicro `json:"Oem"`
}

type _managerActionsDataOemSupermicro struct {
	ManagerReset LinkTargets `json:"#Manager.Reset"`
}
