package model

// NfsClusterInfoList model info
// @Description Glue NFS Cluster 상세정보 및 리스트 구조체
type NfsClusterInfoList []string //@name NfsClusterInfoList

type NfsClusterList interface{} //@name NfsClusterList
// NfsExportDetailed model info
// @Description Glue NFS Export 상세정보 구조체
type NfsExportDetailed []struct {
	AccessType string `json:"access_type"`
	Clients    []any  `json:"clients"`
	ClusterID  string `json:"cluster_id"`
	ExportID   int    `json:"export_id"`
	Fsal       struct {
		FsName string `json:"fs_name"`
		Name   string `json:"name"`
		UserID string `json:"user_id"`
	} `json:"fsal"`
	Path          string   `json:"path"`
	Protocols     []int    `json:"protocols"`
	Pseudo        string   `json:"pseudo"`
	SecurityLabel bool     `json:"security_label"`
	Squash        string   `json:"squash"`
	Transports    []string `json:"transports"`
} //@name NfsExportDetailed

// NfsExportCreate model info
// @Description Glue NFS Export 생성 구조체
type NfsExportCreate struct {
	AccessType    string   `json:"access_type"`
	Fsal          Fsal     `json:"fsal"`
	Path          string   `json:"path"`
	Protocols     []int    `json:"protocols"`
	Pseudo        string   `json:"pseudo"`
	SecurityLabel bool     `json:"security_label"`
	Squash        string   `json:"squash"`
	Transports    []string `json:"transports"`
} //@name NfsExportCreate

type Fsal struct {
	Name          string `json:"name"`
	FsName        string `json:"fs_name"`
	SecLabelXattr string `json:"sec_label_xattr"`
}

// NfsExportUpdate model info
// @Description Glue NFS Export 수정 구조체
type NfsExportUpdate struct {
	AccessType    string   `json:"access_type"`
	ExportID      int      `json:"export_id"`
	Fsal          Fsal     `json:"fsal"`
	Path          string   `json:"path"`
	Protocols     []int    `json:"protocols"`
	Pseudo        string   `json:"pseudo"`
	SecurityLabel bool     `json:"security_label"`
	Squash        string   `json:"squash"`
	Transports    []string `json:"transports"`
} //@name NfsExportUpdate

type NfsClusterCreateCount struct {
	ServiceType string            `yaml:"service_type"`
	ServiceID   string            `yaml:"service_id"`
	Placement   NfsPlacementCount `yaml:"placement"`
	Spec        NfsSpec           `yaml:"spec"`
}
type NfsClusterCreate struct {
	ServiceType string       `yaml:"service_type"`
	ServiceID   string       `yaml:"service_id"`
	Placement   NfsPlacement `yaml:"placement"`
	Spec        NfsSpec      `yaml:"spec"`
}
type NfsPlacementCount struct {
	Count int      `yaml:"count"`
	Hosts []string `yaml:"hosts"`
}
type NfsPlacement struct {
	Hosts []string `yaml:"hosts"`
}
type NfsSpec struct {
	Port int `yaml:"port"`
}
type NfsIngress struct {
	ServiceType string         `yaml:"service_type"`
	ServiceID   string         `yaml:"service_id"`
	Placement   NfsPlacement   `yaml:"placement"`
	Spec        NfsIngressSpec `yaml:"spec"`
}
type NfsIngressSpec struct {
	BackendService           string   `yaml:"backend_service"`
	VirtualIp                string   `yaml:"virtual_ip"`
	FrontendPort             int      `yaml:"frontend_port"`
	MonitorPort              int      `yaml:"monitor_port"`
	VirtualInterfaceNetworks []string `yaml:"virtual_interface_networks"`
	UseKeepalivedMulticast   bool     `yaml:"use_keepalived_multicast"`
}
