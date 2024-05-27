package model

type SmbNormalStatus struct {
	Names       string   `json:"names"`
	Hostname    string   `json:"hostname"`
	IpAddress   string   `json:"ip_address"`
	ShareFolder string   `json:"folder_name"`
	SharePath   string   `json:"path"`
	Status      string   `json:"status"`
	State       string   `json:"state"`
	Port        []int    `json:"port"`
	FsName      string   `json:"fs_name"`
	VolumePath  string   `json:"volume_path"`
	Users       []string `json:"users"`
} //@name SmbNormalStatus

type SmbStatus interface{} //@name SmbStatus
