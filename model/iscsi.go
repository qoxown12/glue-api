package model

// IscsiServiceCreate model info
// @Description Iscsi Service daemon 구조체
type IscsiServiceCreate struct {
	Service_Type string    `yaml:"service_type"`
	Service_Id   string    `yaml:"service_id"`
	Placement    Placement `yaml:"placement"`
	Spec         Spec      `yaml:"spec"`
} //@name IscsiServiceCreate
type Spec struct {
	Api_Password string `yaml:"api_password"`
	Api_User     string `yaml:"api_user"`
	Api_Port     int    `yaml:"api_port"`
	Pool         string `yaml:"pool"`
}
type Placement struct {
	Hosts []string `yaml:"hosts"`
}

// IscsiServiceCreateCount model info
// @Description Iscsi Service daemon 구조체
type IscsiServiceCreateCount struct {
	Service_Type string         `yaml:"service_type"`
	Service_Id   string         `yaml:"service_id"`
	Placement    PlacementCount `yaml:"placement"`
	Spec         Spec           `yaml:"spec"`
} //@name IscsiServiceCreateCount
type PlacementCount struct {
	Count int      `yaml:"count"`
	Hosts []string `yaml:"hosts"`
}
type Place struct {
	Hosts []string `json:"hosts"`
}
type IscsiService []struct {
	Placement   Place  `json:"placement"`
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	ServiceType string `json:"service_type"`
	Spec        struct {
		APIPassword string `json:"api_password"`
		APIPort     int    `json:"api_port"`
		APIUser     string `json:"api_user"`
		Pool        string `json:"pool"`
	} `json:"spec"`
	Status struct {
		Ports   []int `json:"ports"`
		Running int   `json:"running"`
		Size    int   `json:"size"`
	} `json:"status"`
} //@name IscsiService

// IscsiTargetList model info
// @Description Iscsi Target List 구조체
type IscsiTargetList struct {
	// Created       string `json:"created"`
	// DiscoveryAuth struct {
	// 	MutualPassword                  string `json:"mutual_password"`
	// 	MutualPasswordEncryptionEnabled bool   `json:"mutual_password_encryption_enabled"`
	// 	MutualUsername                  string `json:"mutual_username"`
	// 	Password                        string `json:"password"`
	// 	PasswordEncryptionEnabled       bool   `json:"password_encryption_enabled"`
	// 	Username                        string `json:"username"`
	// } `json:"discovery_auth"`
	// Disks struct {
	// } `json:"disks"`
	// Epoch    int `json:"epoch"`
	// Gateways struct {
	// 	Gwvm struct {
	// 		ActiveLuns int    `json:"active_luns"`
	// 		Created    string `json:"created"`
	// 		Updated    string `json:"updated"`
	// 	} `json:"gwvm"`
	// } `json:"gateways"`
	Targets interface {
	} `json:"targets"`
	// Updated string `json:"updated"`
	// Version int    `json:"version"`
} //@name IscsiTargetList

// IscsiDiskList model info
// @Description Iscsi Disk List 구조체
type IscsiDiskList struct {
	Disks interface{} `json:"disks"`
} //@name IscsiDiskList

// IscsiDiscoveryInfo model info
// @Description Iscsi Discovery 정보 구조체
type IscsiDiscoveryInfo struct {
	Discovery_auth struct {
		Username        string `json:"username"`
		Password        string `json:"password"`
		Mutual_username string `json:"mutual_username"`
		Mutual_password string `json:"mutual_password"`
	} `json:"discovery_auth"`
} // @name IscsiDiscoveryInfo
