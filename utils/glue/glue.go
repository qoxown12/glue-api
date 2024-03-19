package glue

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

func RbdPool() (pools []string, err error) {
	var stdout []byte
	cmd := exec.Command("sh", "-c", "ceph osd pool ls detail | grep 'rbd' | cut -d \"'\" -f2")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		return
	}
	name := strings.Split(string(stdout), "\n")
	pools = append(pools, name...)
	for i := 0; i < len(name); i++ {
		if i == len(name)-1 {
			pools = pools[:len(name)-1]
		}
	}
	return
}
func RbdImage(pool_name string) (pools []string, err error) {
	var stdout []byte
	cmd := exec.Command("rbd", "ls", "-p", pool_name, "--format", "json")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	if err = json.Unmarshal(stdout, &pools); err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	return
}
func ListPool(pool_name string) (pools []string, err error) {
	var stdout []byte
	if pool_name == "" {
		cmd := exec.Command("ceph", "osd", "pool", "ls", "--format", "json")
		stdout, err = cmd.CombinedOutput()

		if err != nil {
			err_str := strings.ReplaceAll(string(stdout), "\n", "")
			err = errors.New(err_str)
			utils.FancyHandleError(err)
			return
		}

		if err = json.Unmarshal(stdout, &pools); err != nil {
			err_str := strings.ReplaceAll(string(stdout), "\n", "")
			err = errors.New(err_str)
			utils.FancyHandleError(err)
			return
		}
		return
	} else {
		cmd := exec.Command("sh", "-c", "ceph osd pool ls detail | grep \""+pool_name+"\" | cut -d \"'\" -f2")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			return
		}
		name := strings.Split(string(stdout), "\n")
		pools = append(pools, name...)
		for i := 0; i < len(name); i++ {
			if i == len(name)-1 {
				pools = pools[:len(name)-1]
			}
		}
		return
	}
}

func InfoImage(pool_name string) (dat model.Images, err error) {
	var stdout []byte

	cmd := exec.Command("rbd", "ls", "-l", "-p", pool_name, "--format", "json")
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
func ListAndInfoImage(image_name string, pool_name string) (dat model.ImageCommon, err error) {
	var stdout []byte

	if image_name != "" && pool_name == "" {
		cmd := exec.Command("rbd", "info", image_name, "--format", "json")
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
		cmd := exec.Command("rbd", "info", pool_name+"/"+image_name, "--format", "json")
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
func CreateImage(image_name string, pool_name string, size string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("rbd", "create", "--size", size, pool_name+"/"+image_name)
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
func DeleteImage(image_name string, pool_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("rbd", "rm", pool_name+"/"+image_name)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		std := strings.ReplaceAll(string(stdout), "\n", "")
		err_str := strings.ReplaceAll(std, "\r", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
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

func PoolDelete(pool_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "osd", "pool", "rm", pool_name, pool_name, "--yes-i-really-really-mean-it")
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
func ServiceLs(service_name string, service_type string) (dat model.ServiceLs, err error) {
	var stdout []byte
	if service_name == "" && service_type == "" {
		cmd := exec.Command("ceph", "orch", "ls", "-f", "json")
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
	} else if service_name == "" && service_type != "" {
		cmd := exec.Command("ceph", "orch", "ls", "--service_type", service_type, "-f", "json")
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
	} else if service_name != "" && service_type == "" {
		cmd := exec.Command("ceph", "orch", "ls", "--service_name", service_name, "-f", "json")
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
	} else {
		cmd := exec.Command("ceph", "orch", "ls", "--service_type", service_type, "--service_name", service_name, "-f", "json")
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
}

func ServiceControl(control string, service_name string) (output string, err error) {
	var stdout []byte
	if service_name == "smb" {
		cmd := exec.Command("systemctl", control, service_name)
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err_str := strings.ReplaceAll(string(stdout), "\n", "")
			err = errors.New(err_str)
			utils.FancyHandleError(err)
			return
		}
		output = "Success"
		return
	} else {
		cmd := exec.Command("ceph", "orch", control, service_name)
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err_str := strings.ReplaceAll(string(stdout), "\n", "")
			err = errors.New(err_str)
			utils.FancyHandleError(err)
			return
		}
		output = strings.ReplaceAll(string(stdout), "\n", ".")
		return
	}
}

func ServiceDelete(service_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "orch", "rm", service_name)
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

func HostList() (dat model.HostList, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "orch", "host", "ls", "-f", "json")
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
func HostIp() (output []byte, err error) {
	var stdout []byte
	cmd := exec.Command("sh", "-c", "cat /etc/hosts | grep -E '*mngt' | grep -v 'ccvm' | awk '{print $1}'")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	output = stdout
	return
}
