package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/iscsi"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

// IscsiServiceCreate godoc
//
//	@Summary		Create of Iscsi Servcie Daemon
//	@Description	Iscsi 서비스 데몬을 생성합니다.
//	@param			hosts 	formData	string	true	"Host Name"
//	@param			service_id	formData	string	true	"ISCSI Service Name"
//	@param			pool 	formData	string	true	"Pool Name"
//	@param			api_port 	formData	string	true	"ISCSI API Port"
//	@param			api_user 	formData	string	true	"ISCSI API User"
//	@param			api_password 	formData	string	true	"ISCSI API Password"
//	@Tags			Iscsi
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi [post]
func (c *Controller) IscsiServiceCreate(ctx *gin.Context) {
	service_id, _ := ctx.GetPostForm("service_id")
	hosts, _ := ctx.GetPostForm("hosts")
	pool, _ := ctx.GetPostForm("pool")
	api_port, _ := ctx.GetPostForm("api_port")
	api_user, _ := ctx.GetPostForm("api_user")
	api_password, _ := ctx.GetPostForm("api_password")
	service_type := "iscsi"
	port, _ := strconv.Atoi(api_port)
	value := model.IscsiServiceCreate{
		Service_Type: service_type,
		Service_Id:   service_id,
		Spec: model.Spec{
			Pool:         pool,
			Api_Port:     port,
			Api_User:     api_user,
			Api_Password: api_password},
		Placement: model.Placement{
			Hosts: []string{hosts}},
	}
	yaml_data, err := yaml.Marshal(&value)
	iscsi_yaml := "/etc/ceph/iscsi.yaml"
	err = os.WriteFile(iscsi_yaml, yaml_data, 0644)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	dat, err := iscsi.IscsiServiceCreate(iscsi_yaml)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		if err := os.Remove(iscsi_yaml); err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// IscsiTargetList godoc
//
//	@Summary		Show List of Iscsi Target
//	@Description	Iscsi 타겟 리스트를 가져옵니다.
//	@Tags			IscsiTarget
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	IscsiTargetList
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/target [get]
func (c *Controller) IscsiTargetList(ctx *gin.Context) {
	dat, err := iscsi.IscsiTargetName()
	dat = strings.Replace(dat, "\n", "", -1)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		data, err := iscsi.IscsiTargetList(dat)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, data)
	}
}

// IscsiTargetCreate godoc
//
//	@Summary		Create of Iscsi Target
//	@Description	Iscsi 타겟을 생성합니다.
//	@param			iqn_id 	path	string	true	"Iscsi IQN ID" example("iqn.{yyyy-mm}.{naming-authority}:{unique-name}")
//	@param			hostname 	formData	string	true	"Gateway Host Name"
//	@param			ip_address 	formData	string	true	"Gateway Host IP Address"
//	@param			pool_name 	formData	string	true	"Glue Pool Name"
//	@param			disk_name 	formData	string	true	"Iscsi Disk Name"
//	@param			size 	formData	int	false	"Iscsi Disk Image Size(Default GB)"
//	@Tags			IscsiTarget
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/target/{iqn_id} [post]
func (c *Controller) IscsiTargetCreate(ctx *gin.Context) {
	iqn_id := ctx.Param("iqn_id")
	hostname, _ := ctx.GetPostForm("hostname")
	ip_address, _ := ctx.GetPostForm("ip_address")
	disk_name, _ := ctx.GetPostForm("disk_name")
	pool_name, _ := ctx.GetPostForm("pool_name")
	size, _ := ctx.GetPostForm("size")
	// cmd := exec.Command("sh", "-c", "cat /etc/hosts | sort | grep -w -m 1 'gwvm' | awk '{print $1}'")
	// ip_address, err := cmd.CombinedOutput()
	ceph_container_name, err := iscsi.IscsiTargetName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)

	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		data, err := iscsi.IscsiTargetCreate(ceph_container_name, iqn_id, hostname, ip_address, pool_name, disk_name, size)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, data)
	}
}

// IscsiTargetDelete godoc
//
//	@Summary		Delete of Iscsi Target
//	@Description	Iscsi 타겟을 삭제합니다.
//	@param			iqn_id 	path	string	true	"Iscsi IQN ID"
//	@param			pool_name 	query	string	true	"Glue Pool Name"
//	@param			disk_name 	query	string	true	"Iscsi Disk Name"
//	@param			image 	query	string	true	"Whether to Delete RBD Image" default(false) Enums(true,false)
//	@Tags			IscsiTarget
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/target/{iqn_id} [delete]
func (c *Controller) IscsiTargetDelete(ctx *gin.Context) {
	iqn_id := ctx.Param("iqn_id")
	pool_name := ctx.Request.URL.Query().Get("pool_name")
	disk_name := ctx.Request.URL.Query().Get("disk_name")
	image := ctx.Request.URL.Query().Get("image")
	ceph_container_name, err := iscsi.IscsiTargetName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		data, err := iscsi.IscsiTargetDelete(ceph_container_name, pool_name, disk_name, iqn_id, image)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, data)
	}
}
func (c *Controller) IscsiTargetDeleteOptions(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}

// IscsiDiskList godoc
//
//	@Summary		Show List of Iscsi Disk
//	@Description	Iscsi 디스크 리스트를 보여줍니다.
//	@Tags			IscsiDisk
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	IscsiDiskList
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/disk [get]
func (c *Controller) IscsiDiskList(ctx *gin.Context) {
	ceph_container_name, err := iscsi.IscsiTargetName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		data, err := iscsi.IscsiDiskList(ceph_container_name)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, data)
	}
}

// IscsiDiskCreate godoc
//
//	@Summary		Create Or Attach of Iscsi Disk
//	@Description	Iscsi 디스크를 생성 또는 부착합니다.
//	@param			pool_name	formData	string	true	"Iscsi Disk Pool Name"
//	@param			disk_name 	formData	string	true	"Iscsi Disk Name"
//	@param			size	formData	int	false	"Iscsi Disk Image Size(Default GB)"
//	@param			iqn_id  formData	string	false	"Iscsi IQN ID"
//	@Tags			IscsiDisk
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/disk [post]
func (c *Controller) IscsiDiskCreate(ctx *gin.Context) {
	disk_name, _ := ctx.GetPostForm("disk_name")
	pool_name, _ := ctx.GetPostForm("pool_name")
	size, _ := ctx.GetPostForm("size")
	iqn_id, _ := ctx.GetPostForm("iqn_id")
	ceph_container_name, err := iscsi.IscsiTargetName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		data, err := iscsi.IscsiDiskCreate(ceph_container_name, pool_name, disk_name, size, iqn_id)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, data)
	}
}

// IscsiDiskDelete godoc
//
//	@Summary		Delete of Iscsi Disk
//	@Description	Iscsi 디스크를 삭제합니다.
//	@param			pool_name 	query	string	true	"Iscsi Disk Pool Name"
//	@param			disk_name 	query	string	true	"Iscsi Disk Name"
//	@param			iqn_id  query	string	false	"Iscsi IQN ID"
//	@param			image 	query	string	true	"Whether to Delete RBD Image" default(false) Enums(true,false)
//	@Tags			IscsiDisk
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest``
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/disk [delete]
func (c *Controller) IscsiDiskDelete(ctx *gin.Context) {
	disk_name := ctx.Request.URL.Query().Get("disk_name")
	pool_name := ctx.Request.URL.Query().Get("pool_name")
	image := ctx.Request.URL.Query().Get("image")
	iqn_id := ctx.Request.URL.Query().Get("iqn_id")
	ceph_container_name, err := iscsi.IscsiTargetName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		data, err := iscsi.IscsiDiskDelete(ceph_container_name, pool_name, disk_name, image, iqn_id)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, data)
	}
}

// IscsiDiskResize godoc
//
//	@Summary		Change Size of Iscsi Disk
//	@Description	Iscsi 디스크 용량을 변경합니다.
//	@param			disk_name 	query	string	true	"Iscsi Disk Name"
//	@param			new_size 	query	int	true	"Iscsi Disk New Size(default GB)"
//	@Tags			IscsiDisk
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/disk [put]
func (c *Controller) IscsiDiskResize(ctx *gin.Context) {
	disk_name := ctx.Request.URL.Query().Get("disk_name")
	new_size := ctx.Request.URL.Query().Get("new_size")
	ceph_container_name, err := iscsi.IscsiTargetName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		data, err := iscsi.IscsiDiskResize(ceph_container_name, disk_name, new_size)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, data)
	}
}
func (c *Controller) IscsiDiskOptions(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}
