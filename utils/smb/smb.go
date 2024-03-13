package smb

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

func SmbStatus(hostname string, name string) (dat model.SmbStatus, err error) {
	var stdout []byte

	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "sh", "/usr/local/samba/sbin/Samba-Execute.sh", "select")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		dat = model.SmbStatus{
			Hostname:  name,
			IpAddress: hostname,
			Status:    "Error",
			State:     string(stdout)}
		return
	}
	if err = json.Unmarshal(stdout, &dat); err != nil {
		return
	}
	return

}
func SmbCreate(hostname string, username string, password string, folder string, path string, fs_name string, volume_path string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "sh", "/usr/local/samba/sbin/Samba-Execute.sh", "create", "--username", username, "--password", password, "--folder", folder, "--path", path, "--fs_name", fs_name, "--volume_path", volume_path)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}
func SmbUserCreate(hostname string, username string, password string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "sh", "/usr/local/samba/sbin/Samba-Execute.sh", "user_create", "--username", username, "--password", password)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}
func SmbUserUpdate(hostname string, username string, password string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "sh", "/usr/local/samba/sbin/Samba-Execute.sh", "update", "--username", username, "--password", password)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}
func SmbUserDelete(hostname string, username string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "sh", "/usr/local/samba/sbin/Samba-Execute.sh", "user_delete", "--username", username)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}
func SmbDelete(hostname string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "sh", "/usr/local/samba/sbin/Samba-Execute.sh", "delete")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
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
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	host_data := strings.Split(string(stdout), "\n")
	output = host_data
	return
}
