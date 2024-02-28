package model

type SmbStatus struct {
	Hostname    string   `json:"hostname"`
	IpAddress   string   `json:"ip_address"`
	Port        []string `json:"port"`
	Names       string   `json:"names"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	State       string   `json:"state"`
	Users       Users    `json:"users"`
} //@name SmbStatus
type Users struct {
	Users []string `json:"users"`
}
