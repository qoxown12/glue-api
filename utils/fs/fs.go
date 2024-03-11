package fs

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/json"
	"errors"
	"os/exec"
)

func FsStatus() (dat model.FsStatus, err error) {

	var stdout []byte
	cmd := exec.Command("ceph", "fs", "status", "-f", "json")
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
func FsCreate(fs_name string, data_pool_size string, meta_pool_size string) (output string, err error) {
	var stdCreate []byte
	cmd := exec.Command("ceph", "fs", "volume", "create", fs_name, "--placement=label:scvm")
	stdCreate, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdCreate))
		utils.FancyHandleError(err)
		output = "Fail"
		return
	} else {
		cmd := exec.Command("ceph", "osd", "pool", "rename", "cephfs."+fs_name+".data", fs_name+".data")
		stdCreate, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdCreate))
			utils.FancyHandleError(err)
			output = "Fail"
			return
		} else {
			cmd := exec.Command("ceph", "osd", "pool", "rename", "cephfs."+fs_name+".meta", fs_name+".meta")
			stdCreate, err = cmd.CombinedOutput()
			if err != nil {
				err = errors.New(string(stdCreate))
				utils.FancyHandleError(err)
				output = "Fail"
				return
			} else {
				if data_pool_size == "" && meta_pool_size == "" {
					cmd := exec.Command("ceph", "osd", "pool", "set", fs_name+".data", "size", "2")
					stdCreate, err = cmd.CombinedOutput()
					if err != nil {
						err = errors.New(string(stdCreate))
						utils.FancyHandleError(err)
						output = "Fail"
						return
					}
					return
				} else if data_pool_size != "" && meta_pool_size == "" {
					cmd := exec.Command("ceph", "osd", "pool", "set", fs_name+".data", "size", data_pool_size)
					stdCreate, err = cmd.CombinedOutput()
					if err != nil {
						err = errors.New(string(stdCreate))
						utils.FancyHandleError(err)
						output = "Fail"
						return
					}
					return
				} else if data_pool_size == "" && meta_pool_size != "" {
					cmd := exec.Command("ceph", "osd", "pool", "set", fs_name+".meta", "size", meta_pool_size)
					stdCreate, err = cmd.CombinedOutput()
					if err != nil {
						err = errors.New(string(stdCreate))
						utils.FancyHandleError(err)
						output = "Fail"
						return
					}
					return
				} else {
					cmd := exec.Command("ceph", "osd", "pool", "set", fs_name+".data", "size", data_pool_size)
					stdCreate, err = cmd.CombinedOutput()
					if err != nil {
						err = errors.New(string(stdCreate))
						utils.FancyHandleError(err)
						output = "Fail"
						return
					} else {
						cmd := exec.Command("ceph", "osd", "pool", "set", fs_name+".meta", "size", meta_pool_size)
						stdCreate, err = cmd.CombinedOutput()
						if err != nil {
							err = errors.New(string(stdCreate))
							utils.FancyHandleError(err)
							output = "Fail"
							return
						}
					}
					return
				}
			}
		}
	}
}
func FsDelete(fs_name string) (output string, err error) {
	var poolGet []byte
	var poolSet []byte
	var stdDelete []byte
	pool_get_cmd := exec.Command("ceph", "config", "get", "mon", "mon_allow_pool_delete")
	poolGet, err = pool_get_cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(poolGet))
		utils.FancyHandleError(err)
		output = "Fail"
		return
	}
	if string(poolGet) == "true" {
		fs_delete_cmd := exec.Command("ceph", "fs", "volume", "rm", fs_name, "--yes-i-really-mean-it")
		stdDelete, err = fs_delete_cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdDelete))
			utils.FancyHandleError(err)
			output = "Fail"
			return
		} else {
			output = "Success"
		}
		return
	} else {
		pool_set_cmd := exec.Command("ceph", "config", "set", "mon", "mon_allow_pool_delete", "true")
		poolSet, err = pool_set_cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(poolSet))
			utils.FancyHandleError(err)
			output = "Fail"
			return
		} else {
			fs_delete_cmd := exec.Command("ceph", "fs", "volume", "rm", fs_name, "--yes-i-really-mean-it")
			stdDelete, err = fs_delete_cmd.CombinedOutput()
			if err != nil {
				err = errors.New(string(stdDelete))
				utils.FancyHandleError(err)
				output = "Fail"
				return
			} else {
				output = "Success"
			}
			return
		}
	}
}
func FsGetInfo(fs_name string) (dat model.FsGetInfo, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "get", fs_name, "-f", "json")
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
func FsList() (dat model.FsList, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "fs", "ls", "-f", "json")
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
