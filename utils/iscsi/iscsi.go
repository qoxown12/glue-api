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
func IscsiService() (dat model.IscsiService, err error) {
	var stdout []byte
	cmd := exec.Command("ceph", "orch", "ls", "--service_type", "iscsi", "-f", "json")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		return
	}
	if err = json.Unmarshal(stdout, &dat); err != nil {
		return
	}
	return
}
func IscsiTargetName(hostname string) (output string, err error) {
	var std []byte
	if hostname == "gwvm" {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "ps", "--filter", "name=iscsi.*^?tcmu$", "--format={{.Names}}", "--sort=names")
		std, err = cmd.CombinedOutput()
		output = string(std)
		if err != nil {
			err = errors.New(string(std))
			utils.FancyHandleError(err)
			return
		}
		return
	} else {
		cmd := exec.Command("podman", "ps", "--filter", "name=iscsi.*^?tcmu$", "--format={{.Names}}", "--sort=names")
		std, err = cmd.CombinedOutput()
		output = string(std)
		if err != nil {
			err = errors.New(string(std))
			utils.FancyHandleError(err)
			return
		}
		return
	}
}
func IscsiTargetList(ceph_container_name string, hostname string) (dat model.IscsiTargetList, err error) {
	var std []byte
	if hostname == "gwvm" {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "export", "mode=copy")
		std, err = cmd.CombinedOutput()

		if err != nil {
			return
		}
		if err = json.Unmarshal(std, &dat); err != nil {
			return
		}
		return
	} else {
		cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "export", "mode=copy")
		std, err = cmd.CombinedOutput()

		if err != nil {
			return
		}
		if err = json.Unmarshal(std, &dat); err != nil {
			return
		}
		return
	}
}
func IscsiGatewayAttach(ceph_container_name string, hostname string, iqn_id string, portal string, ip_address string) (output string, err error) {
	var stdout []byte
	if hostname == "gwvm" {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/gateways", "create", "skipchecks=true", portal, ip_address)
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			output = "Fail"
			return
		}
		output = "Success"
		return
	} else {
		cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/gateways", "create", "skipchecks=true", portal, ip_address)
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
func IscsiTargetCreate(ceph_container_name string, hostname string, iqn_id string, pool_name string, disk_name string, size string, auth string, username string, password string, mutual_username string, mutual_password string) (output string, err error) {
	var target []byte
	if hostname == "gwvm" {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets", "create", iqn_id)
		target, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(target))
			utils.FancyHandleError(err)
			output = "Fail"
			return
		} else {
			if size != "" {
				var disk_create []byte
				cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "create", "pool="+pool_name, "image="+disk_name, "size="+size+string("G"))
				disk_create, err = cmd.CombinedOutput()
				if err != nil {
					err = errors.New(string(disk_create))
					utils.FancyHandleError(err)
					output = "Fail"
					return
				} else {
					var disk_attach []byte
					cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "add", pool_name+"/"+disk_name)
					disk_attach, err = cmd.CombinedOutput()
					if err != nil {
						err = errors.New(string(disk_attach))
						utils.FancyHandleError(err)
						output = "Fail"
						return
					} else {
						var hosts_acl []byte
						if auth == "false" {
							cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/hosts", "auth", "disable_acl")
							hosts_acl, err = cmd.CombinedOutput()
							if err != nil {
								err = errors.New(string(hosts_acl))
								utils.FancyHandleError(err)
								output = "Fail"
								return
							} else {
								var target_auth []byte
								if username == "" {
									output = "Success"
									return
								} else {
									if mutual_username == "" {
										cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id, "auth", "username="+username, "password="+password)
										target_auth, err = cmd.CombinedOutput()
										if err != nil {
											err = errors.New(string(target_auth))
											utils.FancyHandleError(err)
											output = "Fail"
											return
										}
										output = "Success"
										return
									} else {
										cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id, "auth", "username="+username, "password="+password, "mutual_username="+mutual_username, "mutual_password="+mutual_password)
										target_auth, err = cmd.CombinedOutput()
										if err != nil {
											err = errors.New(string(target_auth))
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
						} else {
							var target_auth []byte
							if username == "" {
								output = "Success"
								return
							}
							if mutual_username == "" {
								cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id, "auth", "username="+username, "password="+password)
								target_auth, err = cmd.CombinedOutput()
								if err != nil {
									err = errors.New(string(target_auth))
									utils.FancyHandleError(err)
									output = "Fail"
									return
								}
								output = "Success"
								return
							} else {
								cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id, "auth", "username="+username, "password="+password, "mutual_username="+mutual_username, "mutual_password="+mutual_password)
								target_auth, err = cmd.CombinedOutput()
								if err != nil {
									err = errors.New(string(target_auth))
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

				}
			} else {
				var disk_attach []byte
				cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "add", pool_name+"/"+disk_name)
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
	} else {
		cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets", "create", iqn_id)
		target, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(target))
			utils.FancyHandleError(err)
			output = "Fail"
			return
		} else {
			if size != "" {
				var disk_create []byte
				cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "create", "pool="+pool_name, "image="+disk_name, "size="+size+string("G"))
				disk_create, err = cmd.CombinedOutput()
				if err != nil {
					err = errors.New(string(disk_create))
					utils.FancyHandleError(err)
					output = "Fail"
					return
				} else {
					var disk_attach []byte
					cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "add", pool_name+"/"+disk_name)
					disk_attach, err = cmd.CombinedOutput()
					if err != nil {
						err = errors.New(string(disk_attach))
						utils.FancyHandleError(err)
						output = "Fail"
						return
					} else {
						var hosts_acl []byte
						if auth == "false" {
							cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/hosts", "auth", "disable_acl")
							hosts_acl, err = cmd.CombinedOutput()
							if err != nil {
								err = errors.New(string(hosts_acl))
								utils.FancyHandleError(err)
								output = "Fail"
								return
							} else {
								var target_auth []byte
								if username == "" {
									output = "Success"
									return
								} else {
									if mutual_username == "" {
										cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id, "auth", "username="+username, "password="+password)
										target_auth, err = cmd.CombinedOutput()
										if err != nil {
											err = errors.New(string(target_auth))
											utils.FancyHandleError(err)
											output = "Fail"
											return
										}
										output = "Success"
										return
									} else {
										cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id, "auth", "username="+username, "password="+password, "mutual_username="+mutual_username, "mutual_password="+mutual_password)
										target_auth, err = cmd.CombinedOutput()
										if err != nil {
											err = errors.New(string(target_auth))
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
						} else {
							var target_auth []byte
							if username == "" {
								output = "Success"
								return
							}
							if mutual_username == "" {
								cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id, "auth", "username="+username, "password="+password)
								target_auth, err = cmd.CombinedOutput()
								if err != nil {
									err = errors.New(string(target_auth))
									utils.FancyHandleError(err)
									output = "Fail"
									return
								}
								output = "Success"
								return
							} else {
								cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id, "auth", "username="+username, "password="+password, "mutual_username="+mutual_username, "mutual_password="+mutual_password)
								target_auth, err = cmd.CombinedOutput()
								if err != nil {
									err = errors.New(string(target_auth))
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

				}
			} else {
				var disk_attach []byte
				cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "add", pool_name+"/"+disk_name)
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
func IscsiTargetDelete(ceph_container_name string, hostname string, pool_name string, disk_name string, iqn_id string, image string) (output string, err error) {
	var stdout []byte
	if hostname == "gwvm" {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "delete", pool_name+"/"+disk_name)
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			output = "Fail"
			return
		} else {
			cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets", "delete", iqn_id)
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err = errors.New(string(stdout))
				utils.FancyHandleError(err)
				output = "Fail"
				return
			} else {
				if image == "true" {
					cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "delete", pool_name+"/"+disk_name)
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
	} else {
		cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "delete", pool_name+"/"+disk_name)
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			output = "Fail"
			return
		} else {
			cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets", "delete", iqn_id)
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err = errors.New(string(stdout))
				utils.FancyHandleError(err)
				output = "Fail"
				return
			} else {
				if image == "true" {
					cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "delete", pool_name+"/"+disk_name)
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
}

func IscsiDiskList(ceph_container_name string, hostname string) (list model.IscsiDiskList, err error) {
	var stdout []byte
	if hostname == "gwvm" {
		gwcli_cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "export", "mode=copy")
		stdout, err = gwcli_cmd.CombinedOutput()
		if err != nil {
			return
		}
		if err = json.Unmarshal(stdout, &list); err != nil {
			utils.FancyHandleError(err)
			return
		}
		return
	} else {
		gwcli_cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "export", "mode=copy")
		stdout, err = gwcli_cmd.CombinedOutput()
		if err != nil {
			return
		}
		if err = json.Unmarshal(stdout, &list); err != nil {
			utils.FancyHandleError(err)
			return
		}
		return
	}
}
func IscsiDiskCreate(ceph_container_name string, hostname string, pool_name string, disk_name string, size string, iqn_id string) (output string, err error) {
	var stdout []byte
	if hostname == "gwvm" {
		if iqn_id == "" {
			cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "create", "pool="+pool_name, "image="+disk_name, "size="+size+string("G"))
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
			cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "add", pool_name+"/"+disk_name)
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
			cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "create", "pool="+pool_name, "image="+disk_name, "size="+size+string("G"))
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err = errors.New(string(stdout))
				utils.FancyHandleError(err)
				output = "Fail"
				return
			} else {
				cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "add", pool_name+"/"+disk_name)
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
	} else {
		if iqn_id == "" {
			cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "create", "pool="+pool_name, "image="+disk_name, "size="+size+string("G"))
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
			cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "add", pool_name+"/"+disk_name)
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
			cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "create", "pool="+pool_name, "image="+disk_name, "size="+size+string("G"))
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err = errors.New(string(stdout))
				utils.FancyHandleError(err)
				output = "Fail"
				return
			} else {
				cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "add", pool_name+"/"+disk_name)
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

}
func IscsiDiskDelete(ceph_container_name string, hostname string, pool_name string, disk_name string, image string, iqn_id string) (output string, err error) {
	var stdout []byte
	if hostname == "gwvm" {
		if iqn_id == "" {
			if image == "true" {
				cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "delete", pool_name+"/"+disk_name)
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
		} else {
			cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "delete", pool_name+"/"+disk_name)
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err = errors.New(string(stdout))
				utils.FancyHandleError(err)
				output = "Fail"
				return
			}
			if image == "true" {
				cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "delete", pool_name+"/"+disk_name)
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
		return
	} else {
		if iqn_id == "" {
			if image == "true" {
				cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "delete", pool_name+"/"+disk_name)
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
		} else {
			cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets/"+iqn_id+"/disks", "delete", pool_name+"/"+disk_name)
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err = errors.New(string(stdout))
				utils.FancyHandleError(err)
				output = "Fail"
				return
			}
			if image == "true" {
				cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "delete", pool_name+"/"+disk_name)
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
	}
	return
}
func IscsiDiskResize(ceph_container_name string, hostname string, disk_name string, new_size string) (output string, err error) {
	var stdout []byte
	if hostname == "gwvm" {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "resize", disk_name, new_size+string("G"))
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			output = "Fail"
			return
		}
		output = "Success"
		return
	} else {
		cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/disks", "resize", disk_name, new_size+string("G"))
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
func IscsiDiscoveryInfo(ceph_container_name string, hostname string) (dat model.IscsiDiscoveryInfo, err error) {
	var stdout []byte
	if hostname == "gwvm" {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "export", "mode=copy")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			return
		}
		if err = json.Unmarshal(stdout, &dat); err != nil {
			return
		}
		return
	} else {
		cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "export", "mode=copy")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			return
		}
		if err = json.Unmarshal(stdout, &dat); err != nil {
			return
		}
		return
	}
}
func IscsiDiscoveryCreate(ceph_container_name string, hostname string, username string, password string, mutual_username string, mutual_password string) (output string, err error) {
	var stdout []byte
	if hostname == "gwvm" {
		if mutual_username == "" && mutual_password == "" {
			cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets", "discovery_auth", "username="+username, "password="+password)
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err = errors.New(string(stdout))
				utils.FancyHandleError(err)
				output = "Fail"
				return
			}
			output = "Success"
			return
		} else {
			cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets", "discovery_auth", "username="+username, "password="+password, "mutual_username="+mutual_username, "mutual_password="+mutual_password)
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
	} else {
		if mutual_username == "" && mutual_password == "" {
			cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets", "discovery_auth", "username="+username, "password="+password)
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err = errors.New(string(stdout))
				utils.FancyHandleError(err)
				output = "Fail"
				return
			}
			output = "Success"
			return
		} else {
			cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets", "discovery_auth", "username="+username, "password="+password, "mutual_username="+mutual_username, "mutual_password="+mutual_password)
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
}
func IscsiDiscoveryReset(ceph_container_name string, hostname string) (output string, err error) {
	var stdout []byte
	if hostname == "gwvm" {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "gwvm", "podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets", "discovery_auth", "nochap")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdout))
			utils.FancyHandleError(err)
			output = "Fail"
			return
		}
		output = "Success"
		return
	} else {
		cmd := exec.Command("podman", "exec", "-it", ceph_container_name, "gwcli", "/iscsi-targets", "discovery_auth", "nochap")
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
