package model

// IscsiServiceCreate model info
// @Description Iscsi Service daemon 구조체
type IscsiServiceCreate struct {
	Service_Type string `json:"service_type"`
	Service_Id   string `json:"service_id"`
	Service_Name string `json:"service_name"`
	Placement    struct {
		Hosts []string `json:"hosts"`
	} `json:"placement"`
	Spec struct {
		Api_Password    string `json:"api_password"`
		Api_User        string `json:"api_user"`
		Api_Port        int16  `json:"api_port"`
		Pool            string `json:"pool"`
		Trusted_Ip_List string `json:"trusted_ip_list"`
	}
} //@name IscsiServiceCreate

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
