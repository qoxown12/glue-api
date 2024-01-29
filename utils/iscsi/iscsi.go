package iscsi

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/json"
	"errors"
	"os/exec"
)

func IscsiServiceCreate(iscsi_yaml string) (output string, err error) {
	var stdCreate []byte
	cluster_create_cmd := exec.Command("ceph", "orch", "apply", "-i", iscsi_yaml)
	stdCreate, err = cluster_create_cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdCreate))
		utils.FancyHandleError(err)
		output = "Fail"
		return
	} else {
		output = "Success"
	}
	return
}

func IscsiTargetName() (output string, err error) {
	var std []byte
	cmd := exec.Command("ssh", "gwvm", "podman", "ps", "--filter", "name=iscsi.*^?tcmu$", "--format={{.Names}}", "--sort=names")
	std, err = cmd.CombinedOutput()
	output = string(std)
	if err != nil {
		err = errors.New(string(std))
		utils.FancyHandleError(err)
		return
	}
	return
}
func IscsiTargetList(ceph_container_name string) (dat model.IscsiTargetList, err error) {
	var std []byte
	cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "export", "mode=copy")
	std, err = cmd.CombinedOutput()

	if err != nil {
		return
	}
	if err = json.Unmarshal(std, &dat); err != nil {
		return
	}
	return
}
func IscsiTargetCreate(ceph_container_name string, iqn_id string, hostname string, ip_address string) (output string, err error) {
	var target []byte
	cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets", "create", iqn_id)
	target, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(target))
		utils.FancyHandleError(err)
		output = "Fail"
		return
	} else {
		var gateway []byte
		cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/gateways", "create", hostname, ip_address)
		gateway, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(gateway))
			utils.FancyHandleError(err)
			output = "Fail"
			return
		}
		output = "Success"
	}

	return
}
func IscsiTargetDelete(ceph_container_name string, iqn_id string) (output string, err error) {
	var stdDelete []byte
	cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets", "delete", iqn_id)
	stdDelete, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdDelete))
		utils.FancyHandleError(err)
		output = "Fail"
		return
	}
	output = "Success"
	return
}
