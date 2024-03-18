package fs

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

func SubVolumeGroupCreate(vol_name string, group_name string, size string, data_pool_name string, mode string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "subvolumegroup", "create", vol_name, group_name, size, data_pool_name, "--mode", mode)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}
func SubVolumeGroupInfo(vol_name string, group_name string) (dat model.SubVolumeGroupInfo, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "subvolumegroup", "info", vol_name, group_name)
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
}
func SubVolumeGroupLs(vol_name string) (dat model.SubVolumeAllLs, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "subvolumegroup", "ls", vol_name)
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
}
func SubVolumeGroupGetPath(vol_name string, group_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "subvolumegroup", "getpath", vol_name, group_name)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	result := strings.Split(string(stdout), "\n")
	output = result[0]
	return
}
func SubVolumeGroupDelete(vol_name string, group_name string, path string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("mount", "-t", "ceph", "admin@."+vol_name+"="+path, "/gluefs/not")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	} else {
		cmd := exec.Command("sh", "-c", "rm -rf /gluefs/not/*")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			return
		} else {
			cmd := exec.Command("umount", "-l", "-f", "/gluefs/not")
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err = errors.New(string(stdout))
				utils.FancyHandleError(err)
				return
			} else {
				cmd := exec.Command("ceph", "fs", "subvolumegroup", "rm", vol_name, group_name)
				stdout, err = cmd.CombinedOutput()
				if err != nil {
					err = errors.New(string(stdout))
					utils.FancyHandleError(err)
					return
				}
				output = "Success"
				return
			}
		}
	}
}
func SubVolumeGroupResize(vol_name string, group_name string, new_size string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "subvolumegroup", "resize", vol_name, group_name, new_size, "--no_shrink")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}
func SubVolumeGroupSnapDelete(vol_name string, group_name string, snap_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "subvolumegroup", "snapshot", "rm", vol_name, group_name, snap_name, "--force")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdout))
		utils.FancyHandleError(err)
		return
	}
	output = "Success"
	return
}
func SubVolumeGroupSnapLs(vol_name string, group_name string) (dat model.SubVolumeAllSnapLs, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "subvolumegroup", "snapshot", "ls", vol_name, group_name)
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
}
