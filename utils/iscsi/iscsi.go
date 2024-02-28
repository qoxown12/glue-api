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
func IscsiService() (dat model.IscsiService, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "orch", "ls", "--service_type", "iscsi", "-f", "json")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		return
	}
	if err = json.Unmarshal(stdout, &dat); err != nil {
		return
	}
	return
}

func GlueUrl() (dat model.GlueUrl, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "mgr", "stat")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		return
	}
	if err = json.Unmarshal(stdout, &dat); err != nil {
		return
	}
	return
}
