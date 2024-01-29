package controller

import (
	"Glue-API/httputil"
	"Glue-API/utils"
	"Glue-API/utils/iscsi"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"sigs.k8s.io/yaml"
)

// IscsiServiceCreate godoc
//
//	@Summary		Create of Iscsi Servcie Daemon
//	@Description	Iscsi 서비스 데몬을 생성합니다.
//	@param			json_file 	body	IscsiServiceCreate  true 	"Iscsi Servcie YAML file"
//	@Tags			ISCSI
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi [post]
func (c *Controller) IscsiServiceCreate(ctx *gin.Context) {
	json_data, _ := io.ReadAll(ctx.Request.Body)
	yaml_data, _ := yaml.JSONToYAML(json_data)
	iscsi_yaml := "/etc/ceph/iscsi.yaml"
	err := os.WriteFile(iscsi_yaml, yaml_data, 0644)

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
	ctx.IndentedJSON(http.StatusOK, dat)
}

// IscsiTargetList godoc
//
//	@Summary		Show List of Iscsi Target
//	@Description	Iscsi 타겟 리스트를 가져옵니다.
//	@Tags			ISCSI
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
		ctx.IndentedJSON(http.StatusOK, data)
	}
}

// IscsiTargetCreate godoc
//
//	@Summary		Create of Iscsi Target
//	@Description	Iscsi 타겟을 생성합니다.
//	@param			iqn_id 	path	string	true	"Iscsi IQN ID" example("iqn.{yyyy-mm}.{naming-authority}:{unique-name}")
//	@param			hostname 	path	string	true	"Gateway Host Name"
//	@Tags			ISCSI
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/target/{iqn_id}/{hostname} [post]
func (c *Controller) IscsiTargetCreate(ctx *gin.Context) {
	iqn_id := ctx.Param("iqn_id")
	hostname := ctx.Param("hostname")
	cmd := exec.Command("sh", "-c", "cat /etc/hosts | sort | grep -w -m 1 'gwvm' | awk '{print $1}'")
	ip_address, err := cmd.CombinedOutput()
	ceph_container_name, err := iscsi.IscsiTargetName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)

	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		data, err := iscsi.IscsiTargetCreate(ceph_container_name, iqn_id, hostname, string(ip_address))
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(http.StatusOK, data)
	}
}

// IscsiTargetDelete godoc
//
//	@Summary		Delete of Iscsi Target
//	@Description	Iscsi 타겟을 삭제합니다.
//	@param			iqn_id 	path	string	true	"Iscsi IQN ID"
//	@Tags			ISCSI
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/target/{iqn_id} [delete]
func (c *Controller) IscsiTargetDelete(ctx *gin.Context) {
	iqn_id := ctx.Param("iqn_id")
	ceph_container_name, err := iscsi.IscsiTargetName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		data, err := iscsi.IscsiTargetDelete(ceph_container_name, iqn_id)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(http.StatusOK, data)
	}
}

// IscsiDiskCreate godoc
//
//	@Summary		Create of Iscsi Disk
//	@Description	Iscsi 디스크를 생성합니다.
//	@param			image_name 	path	string	true	"Iscsi Disk Name"
//	@param			pool_name	formData	string	true	"Iscsi Disk Pool Name"
//	@param			size	formData	string	true	"Iscsi Disk Image Size(Default GB)"
//	@Tags			ISCSI
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/disk/{image_name} [post]
func (c *Controller) IscsiDiskCreate(ctx *gin.Context) {
	image_name := ctx.Param("image_name")
	pool_name, _ := ctx.GetPostForm("pool_name")
	size, _ := ctx.GetPostForm("size")
	ceph_container_name, err := iscsi.IscsiTargetName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		data, err := iscsi.IscsiDiskCreate(ceph_container_name, pool_name, image_name, size+string("G"))
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(http.StatusOK, data)
	}
}
