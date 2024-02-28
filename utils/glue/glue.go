package glue

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/json"
	"errors"
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

func ListAndInfoImage(image_name string, pool_name string) (dat model.Images, err error) {
	var stdout []byte
	if image_name == "" && pool_name == "" {
		cmd := exec.Command("rbd", "ls", "--format", "json")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			return
		}

		if err = json.Unmarshal(stdout, &dat); err != nil {
			return
		}
	} else if image_name != "" && pool_name == "" {
		cmd := exec.Command("rbd", "info", image_name, "--format", "json")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			return
		}
		if err = json.Unmarshal(stdout, &dat); err != nil {
			return
		}
	} else if image_name == "" && pool_name != "" {
		cmd := exec.Command("rbd", "ls", "-l", "-p", pool_name, "--format", "json")
		stdout, err = cmd.CombinedOutput()

		if err != nil {
			return
		}

		if err = json.Unmarshal(stdout, &dat); err != nil {
			return
		}
	} else {
		cmd := exec.Command("rbd", "info", pool_name+"/"+image_name, "--format", "json")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			return
		}
		if err = json.Unmarshal(stdout, &dat); err != nil {
			return
		}
	}
	return
}
func CreateImage(image_name string, pool_name string, size string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("rbd", "create", "--size", size, pool_name+"/"+image_name)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		output = "Fail"
		return
	}
	output = "Success"
	return
}
func DeleteImage(image_name string, pool_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("rbd", "rm", pool_name+"/"+image_name)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		output = "Fail"
		return
	}
	output = "Success"
	return
}
func Status() (dat model.GlueStatus, err error) {

	var stdout []byte
	cmd := exec.Command("ceph", "-s", "-f", "json")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}

	if err = json.Unmarshal(stdout, &dat); err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}

	return
}

func PoolDelete(pool_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "osd", "pool", "rm", pool_name, pool_name, "--yes-i-really-really-mean-it")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	if err = json.Unmarshal(stdout, &output); err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}

	return
}
func ServiceLs(service_name string, service_type string) (dat model.ServiceLs, err error) {
	var stdout []byte
	if service_name == "" && service_type == "" {
		cmd := exec.Command("ceph", "orch", "ls", "-f", "json")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			return
		}
		if err = json.Unmarshal(stdout, &dat); err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			return
		}
		return
	} else if service_name == "" && service_type != "" {
		cmd := exec.Command("ceph", "orch", "ls", "--service_type", service_type, "-f", "json")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			return
		}
		if err = json.Unmarshal(stdout, &dat); err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			return
		}
		return
	} else if service_name != "" && service_type == "" {
		cmd := exec.Command("ceph", "orch", "ls", "--service_name", service_name, "-f", "json")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			return
		}
		if err = json.Unmarshal(stdout, &dat); err != nil {
			utils.FancyHandleError(err)
			return
		}
		return
	} else {
		cmd := exec.Command("ceph", "orch", "ls", "--service_type", service_type, "--service_name", service_name, "-f", "json")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			return
		}
		if err = json.Unmarshal(stdout, &dat); err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			return
		}
		return
	}
}

func ServiceControl(control string, service_name string) (output string, err error) {
	var stdout []byte
	if service_name == "smb" {
		cmd := exec.Command("systemctl", control, service_name)
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			return
		}
		output = "Success"
		return
	} else {
		cmd := exec.Command("ceph", "orch", control, service_name)
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			return
		}
		output = string(stdout)
		return
	}
}

func ServiceDelete(service_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "orch", "rm", service_name)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		output = "Fail"
		return
	}
	output = "Success"
	return
}

func HostList() (dat model.HostList, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "orch", "host", "ls", "-f", "json")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	if err = json.Unmarshal(stdout, &dat); err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	return
}
func HostIp() (output []byte, err error) {
	var stdout []byte
	cmd := exec.Command("bash", "-c", "cat /etc/hosts | grep -E '*mngt' | grep -v 'ccvm' | awk '{print $1}'")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = stdout
	return
}
