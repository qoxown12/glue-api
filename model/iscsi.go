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
	Api_Password  string `yaml:"api_password"`
	Api_User      string `yaml:"api_user"`
	Api_Port      int    `yaml:"api_port"`
	Pool          string `yaml:"pool"`
	TrustedIpList string `yaml:"trusted_ip_list"`
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

type IscsiCommon interface{} // @name IscsiCommon

type GlueUrl struct {
	ActiveName string `json:"active_name"`
}
type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

type IscsiTargetCreate struct {
	Target_Iqn  string    `json:"target_iqn"`
	Portals     []Portals `json:"portals"`
	Disks       []Disks   `json:"disks"`
	Clients     []Clients `json:"clients"`
	Groups      []Groups  `json:"groups"`
	Acl_Enabled bool      `json:"acl_enabled"`
	Auth        Auth      `json:"auth"`
} //@name IscsiTargetCreate
type IscsiTargetUpdate struct {
	New_Target_Iqn string    `json:"new_target_iqn"`
	Portals        []Portals `json:"portals"`
	Disks          []Disks   `json:"disks"`
	Clients        []Clients `json:"clients"`
	Groups         []Groups  `json:"groups"`
	Acl_Enabled    bool      `json:"acl_enabled"`
	Auth           Auth      `json:"auth"`
} //@name IscsiTargetUpdate
type Portals struct {
	Host string `json:"host"`
	Ip   string `json:"ip"`
}
type Disks struct {
	Pool      string   `json:"pool"`
	Image     string   `json:"image"`
	Controls  struct{} `json:"controls"`
	Backstore string   `json:"backstore"`
	Lun       int      `json:"lun"`
}
type Clients struct {
}
type Groups struct {
}
type Auth struct {
	User            string `json:"user"`
	Password        string `json:"password"`
	Mutual_User     string `json:"mutual_user"`
	Mutual_Password string `json:"mutual_password"`
} //@name Auth

type Iscsihosts []struct {
	Placement struct {
		Hosts []string `json:"hosts"`
	} `json:"placement"`
}
