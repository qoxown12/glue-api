package nvmeof

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

func NvmeOfServiceCreate(pool_name string, hosts string) (output string, err error) {
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
				cmd := exec.Command("ceph", "orch", "apply", "nvmeof", pool_name, "--placement", hosts)
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
func NvmeOfCliDownload() (output string, err error) {
	var stdout []byte
	cmd := exec.Command("podman", "pull", "quay.io/ceph/nvmeof-cli:latest")
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
func NvmeOfSubSystemCreate(server_gateway_ip string, server_gateway_port string, subsystem_nqn_id string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("podman", "run", "-i", "quay.io/ceph/nvmeof-cli:latest", "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "subsystem", "add", "--subsystem", subsystem_nqn_id)
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
func HostIp(hostname string) (output string, err error) {
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
	cmd := exec.Command("sh", "-c", "cat /etc/hosts | grep -v '-'| grep -w '"+ip_address+"' | awk '{print $2}'")
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
	if err = json.Unmarshal(stdout, &output); err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	return
}
func NvmeOfDefineGateway(server_gateway_ip string, server_gateway_port, subsystem_nqn_id string, gateway_name string, gateway_ip string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("podman", "run", "-i", "quay.io/ceph/nvmeof-cli:latest", "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "listener", "add", "--subsystem", subsystem_nqn_id, "--gateway-name", gateway_name, "--traddr", gateway_ip, "--trsvcid", "4420")
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
func NvmeOfHostNqn() (output string, err error) {
	var stdout []byte
	cmd := exec.Command("cat", "/etc/nvme/hostnqn")
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
func NvmeOfHostAdd(server_gateway_ip string, server_gateway_port string, subsystem_nqn_id string, host_nqn_id string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("podman", "run", "-i", "quay.io/ceph/nvmeof-cli:latest", "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "host", "add", "--subsystem", subsystem_nqn_id, "--host", host_nqn_id)
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
func NvmeOfNameSpaceCreate(server_gateway_ip string, server_gateway_port string, subsystem_nqn_id string, pool_name string, image_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("podman", "run", "-i", "quay.io/ceph/nvmeof-cli:latest", "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "namespace", "add", "--subsystem", subsystem_nqn_id, "--rbd-pool", pool_name, "--rbd-image", image_name)
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
func NvmeOfSubSystemList(server_gateway_ip string, server_gateway_port string, subsystem_nqn_id string) (output model.NvmeOfSubSystemList, err error) {
	var stdout []byte
	if subsystem_nqn_id == "" {
		cmd := exec.Command("podman", "run", "-i", "quay.io/ceph/nvmeof-cli:latest", "--format", "json", "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "subsystem", "list")
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
		cmd := exec.Command("podman", "run", "-i", "quay.io/ceph/nvmeof-cli:latest", "--format", "json", "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "subsystem", "list", "--subsystem", subsystem_nqn_id)
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
func NvmeOfNameSpaceList(server_gateway_ip string, server_gateway_port string, subsystem_nqn_id string) (output model.NvmeOfNameSpaceList, err error) {
	var stdout []byte
	cmd := exec.Command("podman", "run", "-i", "quay.io/ceph/nvmeof-cli:latest", "--format", "json", "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "namespace", "list", "--subsystem", subsystem_nqn_id)
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
func NvmeOfTargetVerify(hostname string, gateway_ip string) (output model.NvmeOfTargetVerify, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "modprobe", "nvme-fabrics")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	} else {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "nvme", "discover", "-t", "tcp", "-a", gateway_ip, "-o", "json")
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
func NvmeOfList(hostname string) (output model.NvmeOfList, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "nvme", "list", "-v", "-o", "json")
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
func NvmeOfPath(hostname string) (output model.NvmeOfPath, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "nvme", "list", "-o", "json")
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
func NvmeOfConnect(hostname string, gateway_ip string, subsystem_nqn_id string, check bool) (output string, err error) {
	var stdout []byte
	if check {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "nvme", "connect-all", "-t", "tcp", "-a", gateway_ip)
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
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "nvme", "connect", "-t", "tcp", "-a", gateway_ip, "-n", subsystem_nqn_id)

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
func NvmeOfDisConnect(hostname string, subsystem_nqn_id string, check bool) (output string, err error) {
	var stdout []byte
	if check {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "nvme", "disconnect-all")
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
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostname, "nvme", "disconnect", "-n", subsystem_nqn_id)
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
func NvmeOfSubSystemDelete(server_gateway_ip string, server_gateway_port string, subsystem_nqn_id string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("podman", "run", "-i", "quay.io/ceph/nvmeof-cli:latest", "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "subsystem", "del", "--subsystem", subsystem_nqn_id)
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
func NvmeOfNameSpaceDelete(server_gateway_ip string, server_gateway_port string, subsystem_nqn_id string, uuid string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("podman", "run", "-i", "quay.io/ceph/nvmeof-cli:latest", "--server-address", server_gateway_ip, "--server-port", server_gateway_port, "namespace", "del", "--subsystem", subsystem_nqn_id, "--uuid", uuid)
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
