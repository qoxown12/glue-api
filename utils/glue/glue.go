package glue

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/json"
	"os/exec"
)

func ListPool() (pools []string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "osd", "pool", "ls", "--format", "json")
	stdout, err = cmd.CombinedOutput()

	if err != nil {
		return
	}

	if err = json.Unmarshal(stdout, &pools); err != nil {
		return
	}
	return
}

func ListImage(pool_name string) (images []model.Snapshot, err error) {
	var stdout []byte
	cmd := exec.Command("rbd", "ls", "-l", "-p", pool_name, "--format", "json")
	stdout, err = cmd.CombinedOutput()

	if err != nil {
		return
	}

	if err = json.Unmarshal(stdout, &images); err != nil {
		return
	}
	return
}

func Status() (dat model.GlueStatus, err error) {

	var stdout []byte
	cmd := exec.Command("ceph", "-s", "-f", "json")
	stdout, err = cmd.CombinedOutput()

	if err = json.Unmarshal(stdout, &dat); err != nil {
		utils.FancyHandleError(err)
		return
	}

	return
}

func PoolDelete(pool_name string) (output string, err error) {
	var stdDelete []byte
	cmd := exec.Command("ceph", "osd", "pool", "rm", pool_name, pool_name, "--yes-i-really-really-mean-it")
	stdDelete, err = cmd.CombinedOutput()
	if err = json.Unmarshal(stdDelete, &output); err != nil {
		utils.FancyHandleError(err)
		return
	}

	return
}
func ServiceLs(service_name string) (dat model.ServiceLs, err error) {
	var stdout []byte
	if service_name == "" {
		cmd := exec.Command("ceph", "orch", "ls", "-f", "json")
		stdout, err = cmd.CombinedOutput()
		if err = json.Unmarshal(stdout, &dat); err != nil {
			utils.FancyHandleError(err)
			return
		}
	} else {
		cmd := exec.Command("ceph", "orch", "ls", "--service_name", service_name, "-f", "json")
		stdout, err = cmd.CombinedOutput()
		if err = json.Unmarshal(stdout, &dat); err != nil {
			utils.FancyHandleError(err)
			return
		}
	}
	return
}

func ServiceControl(control string, service_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "orch", control, service_name)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		return
	}
	output = string(stdout)
	return
}
