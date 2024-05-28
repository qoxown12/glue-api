package rgw

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

func RgwServiceCreateandUpdate(service_name string, realm_name string, zonegroup_name string, zone_name string, hosts string, port string) (output string, err error) {
	var stdout []byte
	if realm_name == "" {
		cmd := exec.Command("ceph", "orch", "apply", "rgw", service_name, "--placement", hosts, "--port", port)
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
		cmd := exec.Command("radosgw-admin", "realm", "create", "--rgw-realm", realm_name)
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err_str := strings.ReplaceAll(string(stdout), "\n", "")
			err = errors.New(err_str)
			utils.FancyHandleError(err)
			return
		} else {
			cmd := exec.Command("radosgw-admin", "zonegroup", "create", "--rgw-zonegroup", zonegroup_name, "--rgw-realm", realm_name, "--master")
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err_str := strings.ReplaceAll(string(stdout), "\n", "")
				err = errors.New(err_str)
				utils.FancyHandleError(err)
				return
			} else {
				cmd := exec.Command("radosgw-admin", "zone", "create", "--rgw-zonegroup", zonegroup_name, "--rgw-zone", zone_name, "--master")
				stdout, err = cmd.CombinedOutput()
				if err != nil {
					err_str := strings.ReplaceAll(string(stdout), "\n", "")
					err = errors.New(err_str)
					utils.FancyHandleError(err)
					return
				} else {
					cmd := exec.Command("ceph", "orch", "apply", "rgw", service_name, "--realm", realm_name, "--zone", zone_name, "--zonegroup", zonegroup_name, "--placement", hosts, "--port", port)
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
func RgwServiceUpdate(yaml_file string) (output string, err error) {
	var stdout []byte
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
func RgwUserList() (output model.RgwUserList, err error) {
	var stdout []byte
	cmd := exec.Command("radosgw-admin", "user", "list")
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
func RgwUserInfo(username string) (output model.RgwUserInfo, err error) {
	var stdout []byte
	cmd := exec.Command("radosgw-admin", "user", "info", "--uid", username)
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
func RgwUserStat(username string) (output model.RgwUserStat, err error) {
	var stdout []byte
	cmd := exec.Command("radosgw-admin", "user", "stats", "--uid", username)
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
func RgwUserCreate(username string, display_name string, email string) (output string, err error) {
	var stdout []byte
	if email != "" {
		cmd := exec.Command("radosgw-admin", "user", "create", "--uid", username, "--display-name", display_name, "--email", email, "--admin")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err_str := strings.ReplaceAll(string(stdout), "\n", "")
			err = errors.New(err_str)
			utils.FancyHandleError(err)
			return
		}
		output = "Success"
	} else {
		cmd := exec.Command("radosgw-admin", "user", "create", "--uid", username, "--display-name", display_name, "--admin")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err_str := strings.ReplaceAll(string(stdout), "\n", "")
			err = errors.New(err_str)
			utils.FancyHandleError(err)
			return
		}
		output = "Success"
	}
	return
}
func RgwUserDelete(username string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("radosgw-admin", "user", "rm", "--uid", username)
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
func RgwUserUpdate(username string, display_name string, email string, key_type string, access_key string, secret_key string) (output string, err error) {
	var stdout []byte
	if display_name == "" {
		if email == "" {
			if key_type != "" {
				cmd := exec.Command("radosgw-admin", "user", "modify", "--uid", username, "--key-type", key_type, "--access-key", access_key, "--secret-key", secret_key)
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
		} else {
			if key_type != "" {
				cmd := exec.Command("radosgw-admin", "user", "modify", "--uid", username, "--email", email, "--key-type", key_type, "--access-key", access_key, "--secret-key", secret_key)
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
				cmd := exec.Command("radosgw-admin", "user", "modify", "--uid", username, "--email", email)
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
	} else {
		if email == "" {
			if key_type != "" {
				cmd := exec.Command("radosgw-admin", "user", "modify", "--uid", username, "--display_name", display_name, "--key-type", key_type, "--access-key", access_key, "--secret-key", secret_key)
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
				cmd := exec.Command("radosgw-admin", "user", "modify", "--uid", username, "--display_name", display_name)
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
		} else {
			if key_type != "" {
				cmd := exec.Command("radosgw-admin", "user", "modify", "--uid", username, "--display_name", display_name, "--email", email, "--key-type", key_type, "--access-key", access_key, "--secret-key", secret_key)
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
				cmd := exec.Command("radosgw-admin", "user", "modify", "--uid", username, "--display_name", display_name, "--email", email)
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
	return
}
func RgwQuota(username string, scope string, max_object string, max_size string, state string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("radosgw-admin", "quota", "set", "--uid", username, "--quota-scope", scope, "--max-objects", max_object, "--max-size", max_size)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	} else {
		if state == "enable" {
			cmd := exec.Command("radosgw-admin", "quota", "enable", "--uid", username, "--quota-scope", scope)
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
			cmd := exec.Command("radosgw-admin", "quota", "disable", "--uid", username, "--quota-scope", scope)
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

func RgwBucketDetail(bucket_name string) (output model.RGwCommon, err error) {
	var stdout []byte
	if bucket_name == "" {
		cmd := exec.Command("radosgw-admin", "bucket", "stats")
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
		cmd := exec.Command("radosgw-admin", "bucket", "stats", "--bucket", bucket_name)
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
func RgwBucketList() (output model.RGwCommon, err error) {
	var stdout []byte
	cmd := exec.Command("radosgw-admin", "bucket", "list")
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
func RgwBucketDelete(bucket_name string) (output string, err error) {
	var stdout []byte
	cmd := exec.Command("radosgw-admin", "bucket", "rm", "--bucket", bucket_name)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	output = strings.ReplaceAll(string(stdout), "\n", "")
	return
}
