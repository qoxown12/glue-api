package model

type NvmeOfGatewayName []struct {
	Hostname    string `json:"hostname"`
	Daemon_name string `json:"daemon_name"`
}
type NvmeOfSubSystemList struct {
	ErrorMessage string `json:"error_message"`
	Subsystems   []struct {
		Nqn            string `json:"nqn"`
		SerialNumber   string `json:"serial_number"`
		ModelNumber    string `json:"model_number"`
		MinCntlid      int    `json:"min_cntlid"`
		MaxCntlid      int    `json:"max_cntlid"`
		NamespaceCount int    `json:"namespace_count"`
		Subtype        string `json:"subtype"`
		EnableHa       bool   `json:"enable_ha"`
	} `json:"subsystems"`
	Status int `json:"status"`
} //@name NvmeOfSubSystemList

type NvmeOfNameSpaceList struct {
	ErrorMessage string `json:"error_message"`
	SubsystemNqn string `json:"subsystem_nqn"`
	Namespaces   []struct {
		Nsid               int    `json:"nsid"`
		BdevName           string `json:"bdev_name"`
		RbdImageName       string `json:"rbd_image_name"`
		RbdPoolName        string `json:"rbd_pool_name"`
		BlockSize          int    `json:"block_size"`
		RbdImageSize       string `json:"rbd_image_size"`
		UUID               string `json:"uuid"`
		LoadBalancingGroup int    `json:"load_balancing_group"`
		RwIosPerSecond     string `json:"rw_ios_per_second"`
		RwMbytesPerSecond  string `json:"rw_mbytes_per_second"`
		RMbytesPerSecond   string `json:"r_mbytes_per_second"`
		WMbytesPerSecond   string `json:"w_mbytes_per_second"`
	} `json:"namespaces"`
	Status int `json:"status"`
} //@name NvmeOfNameSpaceList

type NvmeOfTargetVerify struct {
	Genctr  int `json:"genctr"`
	Records []struct {
		Trtype  string `json:"trtype"`
		Adrfam  string `json:"adrfam"`
		Subtype string `json:"subtype"`
		Treq    string `json:"treq"`
		Portid  int    `json:"portid"`
		Trsvcid string `json:"trsvcid"`
		Subnqn  string `json:"subnqn"`
		Traddr  string `json:"traddr"`
		Eflags  string `json:"eflags"`
		Sectype string `json:"sectype"`
	} `json:"records"`
} //@name NvmeOfTargetVerify

type NvmeOfList struct {
	Devices []struct {
		HostNQN    string `json:"HostNQN"`
		HostID     string `json:"HostID"`
		Subsystems []struct {
			Subsystem    string `json:"Subsystem"`
			SubsystemNQN string `json:"SubsystemNQN"`
			Controllers  []struct {
				Controller   string `json:"Controller"`
				SerialNumber string `json:"SerialNumber"`
				ModelNumber  string `json:"ModelNumber"`
				Firmware     string `json:"Firmware"`
				Transport    string `json:"Transport"`
				Address      string `json:"Address"`
				Namespaces   []any  `json:"Namespaces"`
				Paths        []struct {
					Path     string `json:"Path"`
					ANAState string `json:"ANAState"`
				} `json:"Paths"`
			} `json:"Controllers"`
			Namespaces []struct {
				NameSpace    string `json:"NameSpace"`
				Generic      string `json:"Generic"`
				Nsid         int    `json:"NSID"`
				UsedBytes    int64  `json:"UsedBytes"`
				MaximumLBA   int    `json:"MaximumLBA"`
				PhysicalSize int64  `json:"PhysicalSize"`
				SectorSize   int    `json:"SectorSize"`
			} `json:"Namespaces"`
		} `json:"Subsystems"`
	} `json:"Devices"`
} // @name NvmeOfList

type HostNvmeOfList struct {
	Hostname string `json:"Hostname"`
	Detail   NvmeOfList
} // @name HostNvmeOfList

type NvmeOfPath struct {
	Hostname string `json:"Hostname"`
	Devices  []struct {
		NameSpace    int    `json:"NameSpace"`
		DevicePath   string `json:"DevicePath"`
		GenericPath  string `json:"GenericPath"`
		Firmware     string `json:"Firmware"`
		ModelNumber  string `json:"ModelNumber"`
		SerialNumber string `json:"SerialNumber"`
		UsedBytes    int64  `json:"UsedBytes"`
		MaximumLBA   int    `json:"MaximumLBA"`
		PhysicalSize int64  `json:"PhysicalSize"`
		SectorSize   int    `json:"SectorSize"`
	} `json:"Devices"`
}
