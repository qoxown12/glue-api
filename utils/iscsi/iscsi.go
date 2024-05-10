package iscsi

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

func IscsiServiceCreate(iscsi_yaml string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "orch", "apply", "-i", iscsi_yaml)
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
func IscsiService() (dat model.IscsiService, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "orch", "ls", "--service_type", "iscsi", "-f", "json")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	if err = json.Unmarshal(stdout, &dat); err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	return
}
func Ip(hostname string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("sh", "-c", "cat /etc/hosts | grep -v '"+hostname+"-' | grep -w '"+hostname+"' | awk '{print $1}'")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	output = string(stdout)
	return
}
func IscsiNADelete(hostname string, container_id string, iqn_id string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "podman", "exec", "-i", container_id, "gwcli", "/iscsi-targets", "delete", iqn_id)
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
func IscsiHost() (output model.Iscsihosts, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "orch", "ls", "--service-type", "iscsi", "-f", "json")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	if err = json.Unmarshal(stdout, &output); err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	return
}
func ContainerId(hostname string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "podman ps | grep 'tcmu' | awk '{print $1}'")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	output = strings.Split(string(stdout), "\n")[0]
	return
}
