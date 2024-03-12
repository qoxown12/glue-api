package model

type SmbStatus struct {
	Names       string   `json:"names"`
	Hostname    string   `json:"hostname"`
	IpAddress   string   `json:"ip_address"`
	ShareFolder string   `json:"folder_name"`
	SharePath   string   `json:"path"`
	Status      string   `json:"status"`
	State       string   `json:"state"`
	Port        []int    `json:"port"`
	Users       []string `json:"users"`
} //@name SmbStatus
