package nfs

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

func NfsServiceCreate(yaml_file string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "orch", "apply", "-i", yaml_file)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	} else {
		output = "Success"
	}
	return

}
func NfsClusterDelete(cluster_id string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "nfs", "cluster", "rm", cluster_id)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	} else {
		output = "Success"
	}
	return
}
func NfsExportCreateOrUpdate(cluster_id string, json_file string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "nfs", "export", "apply", cluster_id, "-i", json_file)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	} else {
		output = "Success"
	}
	return
}
func NfsExportDelete(cluster_id string, pseudo string) (output string, err error) {
	var stdout []byte

	cmd := exec.Command("ceph", "nfs", "export", "rm", cluster_id, pseudo)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	} else {
		output = "Success"
	}
	return
}
func NfsClusterList(cluster_id string) (dat model.NfsClusterList, err error) {
	var stdout []byte
	if cluster_id == "" {
		cmd := exec.Command("ceph", "nfs", "cluster", "info")
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
	} else {
		cmd := exec.Command("ceph", "nfs", "cluster", "info", cluster_id)
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
	}
	return
}
func NfsExportDetailed(cluster_id string) (dat model.NfsExportDetailed, err error) {
	var stdout []byte

	cmd := exec.Command("ceph", "nfs", "export", "ls", cluster_id, "--detailed")
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
func NfsClusterLs() (dat model.NfsClusterInfoList, err error) {
	var stdout []byte

	cmd := exec.Command("ceph", "nfs", "cluster", "ls")
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
