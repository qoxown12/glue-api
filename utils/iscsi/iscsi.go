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
func IscsiTargetCreate(ceph_container_name string, iqn_id string, hostname string, ip_address string, pool_name string, disk_name string, size string) (output string, err error) {
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
		cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/gateways", "create", "skipchecks=true", hostname, ip_address)
		gateway, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(gateway))
			utils.FancyHandleError(err)
			output = "Fail"
			return
		} else {
			if size != "" {
				var disk_create []byte
				cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "create", "pool="+pool_name, "image="+disk_name, "size="+size+string("G"))
				disk_create, err = cmd.CombinedOutput()
				if err != nil {
					err = errors.New(string(disk_create))
					utils.FancyHandleError(err)
					output = "Fail"
					return
				} else {
					var disk_attach []byte
					cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "add", pool_name+"/"+disk_name)
					disk_attach, err = cmd.CombinedOutput()
					if err != nil {
						err = errors.New(string(disk_attach))
						utils.FancyHandleError(err)
						output = "Fail"
						return
					}
					output = "Success"
					return
				}
			} else {
				var disk_attach []byte
				cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "add", pool_name+"/"+disk_name)
				disk_attach, err = cmd.CombinedOutput()
				if err != nil {
					err = errors.New(string(disk_attach))
					utils.FancyHandleError(err)
					output = "Fail"
					return
				}
				output = "Success"
				return
			}
		}

	}
}
func IscsiTargetDelete(ceph_container_name string, pool_name string, disk_name string, iqn_id string, image string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "delete", pool_name+"/"+disk_name)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		output = "Fail"
		return
	} else {
		cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets", "delete", iqn_id)
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			output = "Fail"
			return
		} else {
			if image == "true" {
				cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "delete", pool_name+"/"+disk_name)
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
		}
		output = "Success"
		return
	}
}

func IscsiDiskList(ceph_container_name string) (list model.IscsiDiskList, err error) {
	var stdout []byte
	gwcli_cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "export", "mode=copy")
	stdout, err = gwcli_cmd.CombinedOutput()
	if err = json.Unmarshal(stdout, &list); err != nil {
		utils.FancyHandleError(err)
		return
	}
	return
}
func IscsiDiskCreate(ceph_container_name string, pool_name string, disk_name string, size string, iqn_id string) (output string, err error) {
	var stdout []byte
	if iqn_id == "" {
		cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "create", "pool="+pool_name, "image="+disk_name, "size="+size+string("G"))
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			output = "Fail"
			return
		} else {
			output = "Success"
			return
		}
	} else if size == "" {
		cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "add", pool_name+"/"+disk_name)
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			output = "Fail"
			return
		} else {
			output = "Success"
			return
		}
	} else {
		cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "create", "pool="+pool_name, "image="+disk_name, "size="+size+string("G"))
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			output = "Fail"
			return
		} else {
			cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "add", pool_name+"/"+disk_name)
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err = errors.New(string(stdout))
				utils.FancyHandleError(err)
				output = "Fail"
				return
			} else {
				output = "Success"
				return
			}
		}
	}

}
func IscsiDiskDelete(ceph_container_name string, pool_name string, disk_name string, image string, iqn_id string) (output string, err error) {
	var stdout []byte
	if iqn_id == "" {
		if image == "true" {
			cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "delete", pool_name+"/"+disk_name)
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err = errors.New(string(stdout))
				utils.FancyHandleError(err)
				output = "Fail"
				return
			}
			return
		}
	} else {
		cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "delete", pool_name+"/"+disk_name)
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			output = "Fail"
			return
		} else {
			if image == "true" {
				cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "delete", pool_name+"/"+disk_name)
				stdout, err = cmd.CombinedOutput()
				if err != nil {
					err = errors.New(string(stdout))
					utils.FancyHandleError(err)
					output = "Fail"
					return
				}
				return
			}
		}
		return
	}
	return
}
func IscsiDiskResize(ceph_container_name string, disk_name string, new_size string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ssh", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "resize", disk_name, new_size+string("G"))
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
