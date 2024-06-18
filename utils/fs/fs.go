package fs

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

func FsStatus() (dat model.FsStatus, err error) {

	var stdout []byte
	cmd := exec.Command("ceph", "fs", "status", "-f", "json")
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
func CephHost() (dat model.CephHost, err error) {
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
func FsCreate(fs_name string, hosts string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "volume", "create", fs_name, "--placement", hosts)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	} else {
		cmd := exec.Command("ceph", "osd", "pool", "rename", "cephfs."+fs_name+".data", fs_name+".data")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err_str := strings.ReplaceAll(string(stdout), "\n", "")
			err = errors.New(err_str)
			utils.FancyHandleError(err)
			return
		} else {
			cmd := exec.Command("ceph", "osd", "pool", "rename", "cephfs."+fs_name+".meta", fs_name+".meta")
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err_str := strings.ReplaceAll(string(stdout), "\n", "")
				err = errors.New(err_str)
				utils.FancyHandleError(err)
				return
			} else {
				cmd := exec.Command("ceph", "osd", "pool", "set", fs_name+".data", "size", "2")
				stdout, err = cmd.CombinedOutput()
				if err != nil {
					err_str := strings.ReplaceAll(string(stdout), "\n", "")
					err = errors.New(err_str)
					utils.FancyHandleError(err)
					return
				} else {
					cmd := exec.Command("ceph", "osd", "pool", "set", fs_name+".meta", "size", "2")
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
}
func FsDelete(fs_name string) (output string, err error) {
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
		cmd := exec.Command("ceph", "fs", "volume", "rm", fs_name, "--yes-i-really-mean-it")
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
			cmd := exec.Command("ceph", "fs", "volume", "rm", fs_name, "--yes-i-really-mean-it")
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
func FsGetInfo(fs_name string) (dat model.FsGetInfo, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "get", fs_name, "-f", "json")
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
func FsList() (dat model.FsList, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "ls", "-f", "json")
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
func FsUpdate(old_name string, new_name string, hosts string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "rename", old_name, new_name, "--yes-i-really-mean-it")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	} else {
		cmd := exec.Command("ceph", "osd", "pool", "rename", old_name+".data", new_name+".data")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err_str := strings.ReplaceAll(string(stdout), "\n", "")
			err = errors.New(err_str)
			utils.FancyHandleError(err)
			return
		} else {
			cmd := exec.Command("ceph", "osd", "pool", "rename", old_name+".meta", new_name+".meta")
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err_str := strings.ReplaceAll(string(stdout), "\n", "")
				err = errors.New(err_str)
				utils.FancyHandleError(err)
				return
			} else {
				if hosts != "" {
					cmd := exec.Command("ceph", "orch", "apply", "mds", new_name, hosts)
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
				output = "Success"
				return
			}
		}
	}
}
