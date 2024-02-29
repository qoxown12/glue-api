package model

type SmbStatus struct {
	Names       string   `json:"names"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	State       string   `json:"state"`
	Hostname    string   `json:"hostname"`
	IpAddress   string   `json:"ip_address"`
	Port        []string `json:"port"`
	ShareFolder string   `json:"folder_name"`
	SharePath   string   `json:"path"`
	Users       Users    `json:"users"`
} //@name SmbStatus
type Users struct {
	Users []string `json:"users"`
}
