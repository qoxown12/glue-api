package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/glue"
	"Glue-API/utils/nfs"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

func (c *Controller) NfsOption(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}

// NfsClusterList godoc
//
//	@Summary		Show List of Info of Glue NFS Cluster
//	@Description	Glue NFS Cluster의 리스트 및 상세정보를 보여줍니다.
//	@param			cluster_id 	query	string	false	"NFS Cluster Identifier"
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	NfsClusterList
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs [get]
func (c *Controller) NfsClusterList(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	cluster_id := ctx.Request.URL.Query().Get("cluster_id")
	dat, err := nfs.NfsClusterList(cluster_id)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.IndentedJSON(http.StatusOK, dat)
}

// NfsClusterCreate godoc
//
//	@Summary		Create of Glue NFS Cluster
//	@Description	Glue NFS Cluster를 생성합니다.
//	@param			cluster_id 	path	string	true	"NFS Cluster Identifier"
//	@param			port 	path	string	true	"Cluster Port"
//	@param			hosts 		formData	[]string	true	"Cluster Daemon Hostname" collectionFormat(multi)
//	@param			service_count 	formData	int		false	"Cluster Daemon Service Count"
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs/{cluster_id}/{port} [post]
func (c *Controller) NfsClusterCreate(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	cluster_id := ctx.Param("cluster_id")
	hosts, _ := ctx.GetPostFormArray("hosts")
	service_count, _ := ctx.GetPostForm("service_count")
	port_swag := ctx.Param("port")
	port, _ := strconv.Atoi(port_swag)
	count, _ := strconv.Atoi(service_count)
	if service_count == "" {
		value := model.NfsClusterCreate{
			ServiceType: "nfs",
			ServiceID:   cluster_id,
			Placement: model.NfsPlacement{
				Hosts: hosts,
			},
			Spec: model.NfsSpec{
				Port: port,
			},
		}
		yaml_data, err := yaml.Marshal(&value)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		nfs_yaml := "/etc/ceph/nfs.yaml"
		err = os.WriteFile(nfs_yaml, yaml_data, 0644)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		dat, err := nfs.NfsServiceCreate(nfs_yaml)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			if err := os.Remove(nfs_yaml); err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
			}
			time.Sleep(3 * time.Second)
			pool, err := glue.PoolReplicatedList("nfs")
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			for i := 0; i < len(pool); i++ {
				_, err := glue.PoolReplicatedSize(pool[i])
				if err != nil {
					utils.FancyHandleError(err)
					httputil.NewError(ctx, http.StatusInternalServerError, err)
					return
				}
			}
		}
		ctx.IndentedJSON(http.StatusOK, dat)
	} else {
		value := model.NfsClusterCreateCount{
			ServiceType: "nfs",
			ServiceID:   cluster_id,
			Placement: model.NfsPlacementCount{
				Count: count,
				Hosts: hosts,
			},
			Spec: model.NfsSpec{
				Port: port,
			}}
		yaml_data, err := yaml.Marshal(&value)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		nfs_yaml := "/etc/ceph/nfs.yaml"
		err = os.WriteFile(nfs_yaml, yaml_data, 0644)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		dat, err := nfs.NfsServiceCreate(nfs_yaml)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			if err := os.Remove(nfs_yaml); err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
			}
			time.Sleep(3 * time.Second)
			pool, err := glue.PoolReplicatedList("nfs")
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			for i := 0; i < len(pool); i++ {
				_, err := glue.PoolReplicatedSize(pool[i])
				if err != nil {
					utils.FancyHandleError(err)
					httputil.NewError(ctx, http.StatusInternalServerError, err)
					return
				}
			}
		}
		ctx.IndentedJSON(http.StatusOK, dat)
	}
}

// NfsClusterUpdate godoc
//
//	@Summary		Update of Glue NFS Cluster
//	@Description	Glue NFS Cluster를 수정합니다.
//	@param			cluster_id 	path	string	true	"NFS Cluster Identifier"
//	@param			port 	path	string	true	"Cluster Port"
//	@param			hosts 		formData	[]string	true	"Cluster Daemon Hostname" collectionFormat(multi)
//	@param			service_count 	formData	int		false	"Cluster Daemon Service Count"
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs/{cluster_id}/{port} [put]
func (c *Controller) NfsClusterUpdate(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	cluster_id := ctx.Param("cluster_id")
	hosts, _ := ctx.GetPostFormArray("hosts")
	service_count, _ := ctx.GetPostForm("service_count")
	port_swag := ctx.Param("port")
	port, _ := strconv.Atoi(port_swag)
	count, _ := strconv.Atoi(service_count)
	if service_count == "" {
		value := model.NfsClusterCreate{
			ServiceType: "nfs",
			ServiceID:   cluster_id,
			Placement: model.NfsPlacement{
				Hosts: hosts,
			},
			Spec: model.NfsSpec{
				Port: port,
			},
		}
		yaml_data, err := yaml.Marshal(&value)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		nfs_yaml := "/etc/ceph/nfs.yaml"
		err = os.WriteFile(nfs_yaml, yaml_data, 0644)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		_, err = nfs.NfsServiceCreate(nfs_yaml)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			if err := os.Remove(nfs_yaml); err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
			}
			dat, err := glue.ServiceReDeploy("nfs." + cluster_id)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			ctx.IndentedJSON(http.StatusOK, dat)
		}
	} else {
		value := model.NfsClusterCreateCount{
			ServiceType: "nfs",
			ServiceID:   cluster_id,
			Placement: model.NfsPlacementCount{
				Count: count,
				Hosts: hosts,
			},
			Spec: model.NfsSpec{
				Port: port,
			}}
		yaml_data, err := yaml.Marshal(&value)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		nfs_yaml := "/etc/ceph/nfs.yaml"
		err = os.WriteFile(nfs_yaml, yaml_data, 0644)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		_, err = nfs.NfsServiceCreate(nfs_yaml)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			if err := os.Remove(nfs_yaml); err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
			}
			dat, err := glue.ServiceReDeploy("nfs." + cluster_id)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			ctx.IndentedJSON(http.StatusOK, dat)
		}
	}

}

// NfsClusterDelete godoc
//
//	@Summary		Delete of Glue NFS Cluster
//	@Description	Glue NFS Cluster를 삭제합니다.
//	@param			cluster_id 	path	string	true	"NFS Cluster Identifier"
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs/{cluster_id} [delete]
func (c *Controller) NfsClusterDelete(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	cluster_id := ctx.Param("cluster_id")
	dat, err := nfs.NfsClusterDelete(cluster_id)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.IndentedJSON(http.StatusOK, dat)
}

// NfsExportCreate godoc
//
//	@Summary		Create of Glue NFS Export
//	@Description	Glue NFS Export를 생성합니다.
//	@param			cluster_id 	path	string	true	"NFS Cluster Identifier"
//	@param			access_type formData   string	true    "NFS Access Type" Enums(RW, RO, NONE) default(RW)
//	@param			fs_name     formData   string	false    "FS Name(if stroage_name is CEPH Required value)"
//	@param			storage_name formData   string	true    "NFS Storage Name" Enums(CEPH, RGW)
//	@param			path         formData    string true    "Glue FS Path or Glue RGW Bucket Name"
//	@param			pseudo     formData   string	true    "NFS Export Path"
//	@param			squash     formData   string	true    "Squash"	Enums(no_root_squash, root_id_squash, all_squash, root_squash) default(no_root_squash)
//	@param			transports     formData   []string	false    "Transports" collectionFormat(multi) default(TCP)
//	@param			security_label     formData   boolean	false    "Security Label" default(false)
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs/export/{cluster_id} [post]
func (c *Controller) NfsExportCreate(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	access_type, _ := ctx.GetPostForm("access_type")
	fs_name, _ := ctx.GetPostForm("fs_name")
	storage_name, _ := ctx.GetPostForm("storage_name")
	path, _ := ctx.GetPostForm("path")
	pseudo, _ := ctx.GetPostForm("pseudo")
	squash, _ := ctx.GetPostForm("squash")
	transports, _ := ctx.GetPostFormArray("transports")
	security_label := ctx.GetBool("security_label")

	var protocols = []int{4}
	var value model.NfsExportAll
	if storage_name == "CEPH" {
		value = model.NfsExportCreate{
			AccessType: access_type,
			Fsal: model.NfsFsal{
				Name:   storage_name,
				FsName: fs_name},
			Protocols:     protocols,
			Path:          path,
			Pseudo:        pseudo,
			Squash:        squash,
			SecurityLabel: security_label,
			Transports:    transports}
	} else {
		value = model.NfsExportRgwCreate{
			AccessType: access_type,
			Fsal: model.RgwFsal{
				Name: storage_name,
			},
			Protocols:  protocols,
			Path:       path,
			Pseudo:     pseudo,
			Squash:     squash,
			Transports: transports}
	}
	json_data, err := json.MarshalIndent(value, "", " ")
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	nfs_export_create_conf := "/root/nfs_export_create.conf"
	err = os.WriteFile(nfs_export_create_conf, json_data, 0644)

	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		cluster_id := ctx.Param("cluster_id")
		dat, err := nfs.NfsExportCreateOrUpdate(cluster_id, nfs_export_create_conf)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			if err := os.Remove(nfs_export_create_conf); err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
		}
		ctx.IndentedJSON(http.StatusOK, dat)
	}
}

// NfsExportUpdate godoc
//
//	@Summary		Update of Glue NFS Export
//	@Description	Glue NFS Export를 수정합니다.
//	@param			cluster_id 	path	string	true	"NFS Cluster Identifier"
//	@param			export_id 	formData	int	true	"NFS Export ID"
//	@param			access_type formData   string	true    "NFS Access Type" Enums(RW, RO, NONE) default(RW)
//	@param			fs_name     formData   string	false    "FS Name(if stroage_name is CEPH Required value)"
//	@param			storage_name formData   string	true    "NFS Storage Name" Enums(CEPH, RGW)
//	@param			path         formData    string true    "Glue FS Path or Glue RGW Bucket Name"
//	@param			pseudo     formData   string	true    "NFS Export Path"
//	@param			squash     formData   string	true    "Squash"	Enums(no_root_squash, root_id_squash, all_squash, root_squash) default(no_root_squash)
//	@param			transports     formData   []string	false    "Transports" collectionFormat(multi) default(TCP)
//	@param			security_label     formData   boolean	false    "Security Label" default(false)
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs/export/{cluster_id} [put]
func (c *Controller) NfsExportUpdate(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	cluster_id := ctx.Param("cluster_id")
	export_id_data, _ := ctx.GetPostForm("export_id")
	access_type, _ := ctx.GetPostForm("access_type")
	fs_name, _ := ctx.GetPostForm("fs_name")
	storage_name, _ := ctx.GetPostForm("storage_name")
	path, _ := ctx.GetPostForm("path")
	pseudo, _ := ctx.GetPostForm("pseudo")
	squash, _ := ctx.GetPostForm("squash")
	transports, _ := ctx.GetPostFormArray("transports")
	security_label := ctx.GetBool("security_label")
	export_id, _ := strconv.Atoi(export_id_data)

	var protocols = []int{4}
	var value model.NfsExportAll
	if storage_name == "CEPH" {
		value = model.NfsExportUpdate{
			AccessType: access_type,
			Fsal: model.NfsFsal{
				Name:   storage_name,
				FsName: fs_name},
			Protocols:     protocols,
			Path:          path,
			Pseudo:        pseudo,
			Squash:        squash,
			SecurityLabel: security_label,
			ExportID:      export_id,
			Transports:    transports}
	} else {
		value = model.NfsExportRgwUpdate{
			AccessType: access_type,
			Fsal: model.RgwFsal{
				Name: storage_name,
			},
			Protocols:  protocols,
			Path:       path,
			Pseudo:     pseudo,
			Squash:     squash,
			ExportID:   export_id,
			Transports: transports}
	}

	json_data, err := json.MarshalIndent(value, "", " ")
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	nfs_export_update_conf := "/root/nfs_export_update.conf"
	err = os.WriteFile(nfs_export_update_conf, json_data, 0644)

	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		dat, err := nfs.NfsExportCreateOrUpdate(cluster_id, nfs_export_update_conf)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			if err := os.Remove(nfs_export_update_conf); err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
		}
		ctx.IndentedJSON(http.StatusOK, dat)
	}
}

// NfsExportDelete godoc
//
//	@Summary		Delete of Glue NFS Export
//	@Description	Glue NFS Export를 삭제합니다.
//	@param			cluster_id 	path	string	true	"NFS Cluster Identifier"
//	@param			export_id 	path	string	true	"NFS Export ID"
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs/export/{cluster_id}/{export_id} [delete]
func (c *Controller) NfsExportDelete(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	cluster_id := ctx.Param("cluster_id")
	export_id, err := strconv.Atoi(ctx.Param("export_id"))
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	detail, err := nfs.NfsExportDetailed(cluster_id)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	for i := 0; i < len(detail); i++ {
		if detail[i].ExportID == export_id {
			dat, err := nfs.NfsExportDelete(cluster_id, detail[i].Pseudo)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return

			}
			ctx.IndentedJSON(http.StatusOK, dat)
		}
	}
	// Print the output
}

// NfsExportDetailed godoc
//
//	@Summary		Show Detail of Glue NFS Export
//	@Description	Glue NFS Export 상세 정보를 보여줍니다.
//	@param			cluster_id 	query	string	false	"NFS Cluster Identifier"
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	NfsExportDetailed
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs/export [get]
func (c *Controller) NfsExportDetailed(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	cluster_id := ctx.Request.URL.Query().Get("cluster_id")
	if cluster_id != "" {
		dat, err := nfs.NfsExportDetailed(cluster_id)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(http.StatusOK, dat)
	} else {
		var output model.NfsExportDetailed
		dat2, err := nfs.NfsClusterLs()
		for i := 0; i < len(dat2); i++ {
			dat3, err := nfs.NfsExportDetailed(dat2[i])
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			output = append(dat3, output...)
		}
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(http.StatusOK, output)
	}
}

// IngressCreate godoc
//
//	@Summary		Create of Glue Ingress Service
//	@Description	Glue Ingress Service를 생성합니다.
//	@param			service_id 	formData	string	true	"NFS or RGW Ingress Service Name"
//	@param			hosts     formData   []string	true    "NFS or RGW Ingress Host Name" collectionFormat(multi)
//	@param			backend_service formData   string	true    "NFS or RGW Cluster Type"
//	@param			virtual_ip     formData   string	true    "NFS or RGW Ingress Virtual Ip"
//	@param			frontend_port     formData   int	true    "NFS or RGW Ingress Access Port" maximum(65535)
//	@param			monitor_port     formData   int	true    "NFS or RGW Ingress HA Proxy for Load Balancer Port" maximum(65535)
//	@param			virtual_interface_networks     formData   []string	false    "NFS or RGW Ingress Vitual IP of CIDR Networks" collectionFormat(multi)
//	@Tags			Ingress
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/ingress [post]
func (c *Controller) IngressCreate(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	service_id, _ := ctx.GetPostForm("service_id")
	hosts, _ := ctx.GetPostFormArray("hosts")
	backend_service, _ := ctx.GetPostForm("backend_service")
	virtual_ip, _ := ctx.GetPostForm("virtual_ip")
	frontend_port_data, _ := ctx.GetPostForm("frontend_port")
	monitor_port_data, _ := ctx.GetPostForm("monitor_port")
	virtual_interface_networks, _ := ctx.GetPostFormArray("virtual_interface_networks")
	frontend_port, _ := strconv.Atoi(frontend_port_data)
	monitor_port, _ := strconv.Atoi(monitor_port_data)

	value := model.Ingress{
		ServiceType: "ingress",
		ServiceID:   service_id,
		Placement: model.NfsPlacement{
			Hosts: hosts,
		},
		Spec: model.IngressSpec{
			BackendService:           backend_service,
			VirtualIp:                virtual_ip,
			FrontendPort:             frontend_port,
			MonitorPort:              monitor_port,
			VirtualInterfaceNetworks: virtual_interface_networks,
			UseKeepalivedMulticast:   false,
		},
	}
	yaml_data, err := yaml.Marshal(value)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ingress_conf := "/root/ingress.conf"
	err = os.WriteFile(ingress_conf, yaml_data, 0644)

	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		dat, err := nfs.NfsServiceCreate(ingress_conf)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			if err := os.Remove(ingress_conf); err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
		}
		ctx.IndentedJSON(http.StatusOK, dat)
	}
}

// IngressUpdate godoc
//
//	@Summary		Update of Glue Ingress Service
//	@Description	Glue Ingress Service를 수정합니다.
//	@param			service_id 	formData	string	true	"NFS or RGW Ingress Service Name"
//	@param			hosts     formData   []string	true    "NFS or RGW Ingress Host Name" collectionFormat(multi)
//	@param			backend_service formData   string	true    "NFS or RGW Cluster Type"
//	@param			virtual_ip     formData   string	true    "NFS or RGW Ingress Virtual Ip"
//	@param			frontend_port     formData   int	true    "NFS or RGW Ingress Access Port" maximum(65535)
//	@param			monitor_port     formData   int	true    "NFS or RGW Ingress HA Proxy for Load Balancer Port" maximum(65535)
//	@param			virtual_interface_networks     formData   []string	false    "NFS or RGW Ingress Vitual IP of CIDR Networks" collectionFormat(multi)
//	@Tags			Ingress
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/ingress [put]
func (c *Controller) IngressUpdate(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	service_id, _ := ctx.GetPostForm("service_id")
	hosts, _ := ctx.GetPostFormArray("hosts")
	backend_service, _ := ctx.GetPostForm("backend_service")
	virtual_ip, _ := ctx.GetPostForm("virtual_ip")
	frontend_port_data, _ := ctx.GetPostForm("frontend_port")
	monitor_port_data, _ := ctx.GetPostForm("monitor_port")
	virtual_interface_networks, _ := ctx.GetPostFormArray("virtual_interface_networks")
	frontend_port, _ := strconv.Atoi(frontend_port_data)
	monitor_port, _ := strconv.Atoi(monitor_port_data)

	value := model.Ingress{
		ServiceType: "ingress",
		ServiceID:   service_id,
		Placement: model.NfsPlacement{
			Hosts: hosts,
		},
		Spec: model.IngressSpec{
			BackendService:           backend_service,
			VirtualIp:                virtual_ip,
			FrontendPort:             frontend_port,
			MonitorPort:              monitor_port,
			VirtualInterfaceNetworks: virtual_interface_networks,
			UseKeepalivedMulticast:   false,
		},
	}
	yaml_data, err := yaml.Marshal(value)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ingress_conf := "/etc/ceph/ingress.conf"
	err = os.WriteFile(ingress_conf, yaml_data, 0644)

	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		_, err := nfs.NfsServiceCreate(ingress_conf)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			dat, err := glue.ServiceReDeploy("ingress." + service_id)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			} else {
				if err := os.Remove(ingress_conf); err != nil {
					utils.FancyHandleError(err)
					httputil.NewError(ctx, http.StatusInternalServerError, err)
					return
				}
			}
			ctx.IndentedJSON(http.StatusOK, dat)
		}
	}
}
