package fs

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

func SubVolumeLs(vol_name string, group_name string) (dat model.SubVolumeAllLs, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "subvolume", "ls", vol_name, group_name)
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
func SubVolumeInfo(vol_name string, subvol_name string, group_name string) (dat model.SubVolumeInfo, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "subvolume", "info", vol_name, subvol_name, group_name)
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
func SubVolumeCreate(vol_name string, subvol_name string, group_name string, size string, data_pool_name string, mode string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "subvolume", "create", vol_name, subvol_name, "--size", size, "--group_name", group_name, "--pool_layout", data_pool_name, "--mode", mode)
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
func SubVolumeDelete(vol_name string, subvol_name string, group_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "subvolume", "rm", vol_name, subvol_name, "--group_name", group_name)
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
func SubVolumeResize(vol_name string, subvol_name string, new_size string, group_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "subvolume", "resize", vol_name, subvol_name, new_size, "--group_name", group_name, "--no_shrink")
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
func SubVolumeSnapLs(vol_name string, subvol_name string, group_name string) (dat model.SubVolumeAllSnapLs, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "subvolume", "snapshot", "ls", vol_name, subvol_name, group_name)
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
func SubVolumeSnapInfo(vol_name string, subvol_name string, snap_name string, group_name string) (dat model.SubVolumeAllSnap, err error) {
	var stdout []byte
	if snap_name == "" {
		cmd := exec.Command("ceph", "fs", "subvolume", "snapshot", "ls", vol_name, subvol_name, group_name)
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
		cmd := exec.Command("ceph", "fs", "subvolume", "snapshot", "info", vol_name, subvol_name, snap_name, group_name)
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
func SubVolumeSnapCreate(vol_name string, subvol_name string, snap_name string, group_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "subvolume", "snapshot", "create", vol_name, subvol_name, snap_name, group_name)
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

func SubVolumeSnapDelete(vol_name string, subvol_name string, snap_name string, group_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "subvolume", "snapshot", "rm", vol_name, subvol_name, snap_name, group_name)
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
