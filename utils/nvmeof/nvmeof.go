package nvmeof

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

var nvme_image_version = "quay.io/ceph/nvmeof-cli:1.2.13"

func Container(hostname string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "podman ps | grep 'nvmeof' | awk '{print $1}'")
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
func ServerGatewayIp(hostname string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("sh", "-c", "cat /etc/hosts | grep -v '-'| grep -w '"+hostname+"' | awk '{print $1}'")
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
func Hostname(ip_address string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("sh", "-c", "cat /etc/hosts | grep -w '"+ip_address+"' | awk '{print $2}' | cut -d '-' -f1")
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
func NvmeOfServiceCreate(yaml_file string, pool_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "osd", "pool", "create", pool_name)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	} else {
		cmd := exec.Command("rbd", "pool", "init", pool_name)
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err_str := strings.ReplaceAll(string(stdout), "\n", "")
			err = errors.New(err_str)
			utils.FancyHandleError(err)
			return
		} else {
			cmd := exec.Command("ceph", "osd", "pool", "set", pool_name, "size", "2")
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err_str := strings.ReplaceAll(string(stdout), "\n", "")
				err = errors.New(err_str)
				utils.FancyHandleError(err)
				return
			} else {
				cmd := exec.Command("ceph", "orch", "apply", "-i", yaml_file)
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
}
func NvmeOfCliDownload(hostname string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "podman", "pull", nvme_image_version)
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
func NvmeOfSubSystemCreate(hostname string, server_gateway_ip string, server_gateway_port string, subsystem_nqn_id string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "podman", "run", "-i", nvme_image_version, "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "subsystem", "add", "--subsystem", subsystem_nqn_id)
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
func NvmeOfGatewayName() (output model.NvmeOfGatewayName, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "orch", "ps", "--daemon_type", "nvmeof", "-f", "json")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	if len(stdout) == 4 {
		output = make(model.NvmeOfGatewayName, 0)
	} else {
		if err = json.Unmarshal(stdout, &output); err != nil {
			err_str := strings.ReplaceAll(string(stdout), "\n", "")
			err = errors.New(err_str)
			utils.FancyHandleError(err)
			return
		}
	}
	return
}
func NvmeOfDefineGateway(hostname string, server_gateway_ip string, server_gateway_port, subsystem_nqn_id string, gateway_name string, gateway_ip string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "podman", "run", "-i", nvme_image_version, "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "listener", "add", "--subsystem", subsystem_nqn_id, "--host-name", gateway_name, "--traddr", gateway_ip, "--trsvcid", "4420")
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
func NvmeOfHostAdd(hostname string, server_gateway_ip string, server_gateway_port string, subsystem_nqn_id string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "podman", "run", "-i", nvme_image_version, "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "host", "add", "--subsystem", subsystem_nqn_id, "--host", "'*'")
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
func NvmeOfNameSpaceCreate(hostname string, server_gateway_ip string, server_gateway_port string, subsystem_nqn_id string, pool_name string, image_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "podman", "run", "-i", nvme_image_version, "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "namespace", "add", "--subsystem", subsystem_nqn_id, "--rbd-pool", pool_name, "--rbd-image", image_name)
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
func NvmeOfSubSystemList(hostname string, server_gateway_ip string, server_gateway_port string, subsystem_nqn_id string) (output model.NvmeOfSubSystemList, err error) {
	var stdout []byte
	if subsystem_nqn_id == "" {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "podman", "run", "-i", nvme_image_version, "--format", "json", "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "subsystem", "list")
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
	} else {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "podman", "run", "-i", nvme_image_version, "--format", "json", "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "subsystem", "list", "--subsystem", subsystem_nqn_id)
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
}
func NvmeOfNameSpaceList(hostname string, server_gateway_ip string, server_gateway_port string, subsystem_nqn_id string) (output model.NvmeOfNameSpaceList, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "podman", "run", "-i", nvme_image_version, "--format", "json", "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "namespace", "list", "--subsystem", subsystem_nqn_id)
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
func NvmeOfSubSystemDelete(hostname string, server_gateway_ip string, server_gateway_port string, subsystem_nqn_id string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "podman", "run", "-i", nvme_image_version, "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "subsystem", "del", "--subsystem", subsystem_nqn_id)
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
func NvmeOfNameSpaceDelete(hostname string, server_gateway_ip string, server_gateway_port string, subsystem_nqn_id string, uuid string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "podman", "run", "-i", nvme_image_version, "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "namespace", "del", "--subsystem", subsystem_nqn_id, "--uuid", uuid)
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

func NvmeOfConnection(hostname string, container_id string, subsystem_nqn_id string) (output model.NvmeOfConnection, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "podman", "exec", "-i", container_id, "python3", "/usr/libexec/spdk/scripts/rpc.py", "nvmf_subsystem_get_controllers", subsystem_nqn_id)
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
func NvmeOfTarget(hostname string, container_id string, subsystem_nqn_id string) (output model.NvmeOfTarget, err error) {
	var stdout []byte
	if subsystem_nqn_id == "" {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "podman", "exec", "-i", container_id, "python3", "/usr/libexec/spdk/scripts/rpc.py", "nvmf_get_subsystems")
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
	} else {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "podman", "exec", "-i", container_id, "python3", "/usr/libexec/spdk/scripts/rpc.py", "nvmf_get_subsystems", subsystem_nqn_id)
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
	}
	return
}
