package user

import (
	"Glue-API/utils"
	"errors"
	"os/exec"
	"strings"
)

func UserCreate(username string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "dashboard", "ac-user-create", username, "-i", "/usr/share/ablestack/keycloak/user.txt", "administrator")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}

func UserDelete(username string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "dashboard", "ac-user-delete", username)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}
