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
