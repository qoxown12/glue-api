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
	cmd := exec.Command("ceph", "config", "get", "mon", "mon_allow_pool_delete")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	if string(stdout) == "true" {
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
	} else {
		cmd := exec.Command("ceph", "config", "set", "mon", "mon_allow_pool_delete", "true")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err_str := strings.ReplaceAll(string(stdout), "\n", "")
			err = errors.New(err_str)
			utils.FancyHandleError(err)
			return
		} else {
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
	}
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
		if strings.Contains(string(stdout), "No services reported") {
			dat = make([]string, 0)
			return
		} else {
			if err = json.Unmarshal(stdout, &dat); err != nil {
				err_str := strings.ReplaceAll(string(stdout), "\n", "")
				err = errors.New(err_str)
				utils.FancyHandleError(err)
				return
			}
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
		if strings.Contains(string(stdout), "No services reported") {
			dat = make([]string, 0)
			return
		} else {
			if err = json.Unmarshal(stdout, &dat); err != nil {
				err_str := strings.ReplaceAll(string(stdout), "\n", "")
				err = errors.New(err_str)
				utils.FancyHandleError(err)
				return
			}
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
		if strings.Contains(string(stdout), "No services reported") {
			dat = make([]string, 0)
			return
		} else {
			if err = json.Unmarshal(stdout, &dat); err != nil {
				err_str := strings.ReplaceAll(string(stdout), "\n", "")
				err = errors.New(err_str)
				utils.FancyHandleError(err)
				return
			}
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
	cmd := exec.Command("sh", "-c", "cat /etc/hosts | grep -E '*mngt' | grep -v 'ccvm' | awk '{print $1, $2}'")
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
func RgwPool() (output []string, err error) {
	var stdout []byte
	cmd := exec.Command("sh", "-c", "ceph osd pool ls | grep 'rgw' | sort")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	var str []string
	str_data := strings.Split(string(stdout), "\n")
	for i := 0; i < len(str_data); i++ {
		strs := str_data[i]
		str = append(str, strs)
		if i == len(str_data)-1 {
			str = str[:len(str_data)-1]
		}
	}
	output = str
	return
}
func PoolReplicatedList(pool_type string) (output []string, err error) {
	var stdout []byte
	cmd := exec.Command("sh", "-c", "ceph osd pool ls | grep '"+pool_type+"'")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	str := strings.Split(string(stdout), "\n")
	for i := 0; i < len(str); i++ {
		output = append(output, str[i])
		if i == len(str)-1 {
			output = output[:len(str)-1]
		}
	}
	return
}
func PoolReplicatedSize(pool_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "osd", "pool", "set", pool_name, "size", "2")
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
func ServiceReDeploy(service_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "orch", "redeploy", service_name)
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
func GlueUrl() (dat model.GlueUrl, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "mgr", "stat")
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
