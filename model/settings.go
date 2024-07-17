package model

// import (
// )

type Settings struct {
	ApiPort string `json:"api_port"`
	RemoteHostIp        string `json:"remote_host_ip"`
	RemoteRootRsaIdPath string `json:"remote_root_rsa_id_path"`
	Samba_Security_Type string `json:"samba_security_type"`
	GlueProtocol string `json:"glue_protocol"`
	GluePort string `json:"glue_port"`
	GlueUser string `json:"glue_user"`
	GluePw string `json:"glue_pw"`
}
