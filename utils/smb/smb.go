package smb

import (
	"Glue-API/utils"
	"errors"
	"os/exec"
)

func SmbStatus() (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "systemctl", "show", "--no-pager", "smb")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = string(stdout)
	return
}
func SmbUserMngt() (output []byte, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "/usr/local/samba/bin/pdbedit -L | grep -v 'root' | cut -d ':' -f1 ")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = stdout
	return
}
func SmbCreate(username string, password string, folder string, path string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "sh", "/usr/local/samba/sbin/Samba-Execute.sh", "create", "--username", username, "--password", password, "--folder", folder, "--path", path)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}
func SmbUserCreate(username string, password string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "sh", "/usr/local/samba/sbin/Samba-Execute.sh", "user_create", "--username", username, "--password", password)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}
func SmbUserUpdate(username string, password string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "sh", "/usr/local/samba/sbin/Samba-Execute.sh", "update", "--username", username, "--password", password)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}
func SmbUserDelete(username string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "sh", "/usr/local/samba/sbin/Samba-Execute.sh", "user_delete", "--username", username)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}
func SmbDelete() (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "sh", "/usr/local/samba/sbin/Samba-Execute.sh", "delete")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}
func Hostname() (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "hostname -I | awk '{print $1}'")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = string(stdout)
	return
}
func IpAddress() (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "hostname")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = string(stdout)
	return
}
func Port() (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "netstat -ltnp | grep  smb | grep -v tcp6 | awk '{print $4}' | cut -d ':' -f2")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = string(stdout)
	return
}
func SharePath() (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "cat /usr/local/samba/etc/smb.conf | grep path | awk '{print $3}'")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = string(stdout)
	return
}
func ShareFolder() (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "grep -F '[' /usr/local/samba/etc/smb.conf | grep -v 'global' | tr -d '[]'")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = string(stdout)
	return
}
