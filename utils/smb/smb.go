package smb

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

var Samba_Execute_sh = "/usr/local/glue-api/shell/Samba-Execute.sh"

func SmbStatus(hostname string, name string) (dat model.SmbStatus, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "sh", Samba_Execute_sh, "select")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(stdout), "kex_exchange_identification") {
			cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "sh", Samba_Execute_sh, "select")
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				if strings.Contains(string(stdout), "kex_exchange_identification") {
					cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "sh", Samba_Execute_sh, "select")
					stdout, err = cmd.CombinedOutput()
					if err != nil {
						dat = model.SmbNormalStatus{
							Hostname:  name,
							IpAddress: hostname,
							Status:    "Warn",
							State:     "Please refresh"}

					}
					if err = json.Unmarshal(stdout, &dat); err != nil {
						return
					}
				}
			} else {
				if err = json.Unmarshal(stdout, &dat); err != nil {
					return
				}
			}
			return
		} else {
			dat = model.SmbNormalStatus{
				Hostname:  name,
				IpAddress: hostname,
				Status:    "Error",
				State:     string(stdout)}
		}
		return
	}
	if err = json.Unmarshal(stdout, &dat); err != nil {
		return
	}
	return
}
func SmbCreate(hostname string, sec_type string, username string, password string, folder string, path string, fs_name string, volume_path string, realm string, dns string) (output string, err error) {
	var stdout []byte
	if sec_type == "normal" {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "sh", Samba_Execute_sh, "delete")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string("(") + hostname + string(") ") + string(stdout))
			utils.FancyHandleError(err)
			return
		} else {
			cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "sh", Samba_Execute_sh, "create", sec_type, "--username", username, "--password", password, "--folder", folder, "--path", path, "--fs_name", fs_name, "--volume_path", volume_path)
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err = errors.New(string("(") + hostname + string(") ") + string(stdout))
				utils.FancyHandleError(err)
				return
			}
			output = "Success"
		}
	} else {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "sh", Samba_Execute_sh, "delete")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string("(") + hostname + string(") ") + string(stdout))
			utils.FancyHandleError(err)
			return
		} else {
			cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "sh", Samba_Execute_sh, "create", sec_type, "--username", username, "--password", password, "--folder", folder, "--path", path, "--fs_name", fs_name, "--volume_path", volume_path, "--realm", realm, "--dns", dns)
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err = errors.New(string("(") + hostname + string(") ") + string(stdout))
				utils.FancyHandleError(err)
				return
			}
			output = "Success"
		}
	}
	return
}
func SmbUserCreate(hostname string, username string, password string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "sh", Samba_Execute_sh, "user_create", "--username", username, "--password", password)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string("(") + hostname + string(") ") + string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}
func SmbUserUpdate(hostname string, username string, password string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "sh", Samba_Execute_sh, "update", "--username", username, "--password", password)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string("(") + hostname + string(") ") + string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}
func SmbUserDelete(hostname string, username string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "sh", Samba_Execute_sh, "user_delete", "--username", username)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string("(") + hostname + string(") ") + string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}
func SmbDelete(hostname string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "sh", Samba_Execute_sh, "delete")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string("(") + hostname + string(") ") + string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}
func Hosts() (output []string, err error) {
	var stdout []byte
	cmd := exec.Command("sh", "-c", "cat /etc/hosts | grep -v 'ccvm' | grep 'mngt' | awk '{print $1}'")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	host_data := strings.Split(string(stdout), "\n")
	output = host_data
	return
}
