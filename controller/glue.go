package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/glue"
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

// GlueStatus godoc
//
//	@Summary		Show Status of Glue
//	@Description	Glue 의 상태값을 보여줍니다.
//	@Tags			Glue
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.GlueStatus
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/glue [get]
func (c *Controller) GlueStatus(ctx *gin.Context) {

	dat, err := glue.Status()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// GlueVersion godoc
//
//	@Summary		Show Versions of Glue
//	@Description	Glue 의 버전을 보여줍니다.
//	@Tags			Glue
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.GlueVersion
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/glue/version [get]
func (c *Controller) GlueVersion(ctx *gin.Context) {
	var dat model.GlueVersion

	cmd := exec.Command("ceph", "versions")
	stdout, err := cmd.CombinedOutput()

	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	if err := json.Unmarshal(stdout, &dat); err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// ListPools godoc
//
//	@Summary		List Pools of Glue
//	@Description	Glue 의 스토리지 풀 목록을 보여줍니다.
//	@Tags			Pool
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.GluePools
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/pool [get]
func (c *Controller) ListPools(ctx *gin.Context) {
	var dat model.GluePools
	dat, err := glue.ListPool()
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// ListImages godoc
//
//	@Summary		List RBD Images of Pool
//	@Description	Glue 스토리지 풀의 이미지 목록을 보여줍니다.
//	@Tags			Pool
//	@param			pool_name	path	string	true	"Glue Pool Name"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.GluePools
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/pool/{pool_name} [get]
func (c *Controller) ListImages(ctx *gin.Context) {
	var dat model.SnapshotList
	pool_name := ctx.Param("pool_name")
	images, err := glue.ListImage(pool_name)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	dat.Images = images
	ctx.IndentedJSON(http.StatusOK, dat)
}

// InfoImage godoc
//
//	@Summary		Info Images of Pool
//	@Description	Glue 스토리지 풀의 이미지 상세정보를 보여줍니다.
//	@Tags			Pool
//	@param			image_name	path	string	true	"Glue Image Name"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.InfoImage
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/pool/info/{image_name} [get]
func (c *Controller) InfoImage(ctx *gin.Context) {
	image_name := ctx.Param("image_name")
	dat, err := glue.InfoImage(image_name)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// PoolDelete godoc
//
//	@Summary		List Images of Pool
//	@Description	Glue 스토리지 풀을 삭제합니다.
//	@Tags			Pool
//	@param			pool_name	path	string	true	"pool_name"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/pool/{pool_name} [delete]
func (c *Controller) PoolDelete(ctx *gin.Context) {
	pool_name := ctx.Param("pool_name")
	dat, err := glue.PoolDelete(pool_name)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.IndentedJSON(http.StatusOK, dat)
}
func (c *Controller) PoolDeleteOptions(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}

// ServiceLs godoc
//
//	@Summary		Show List or Info of Glue Service
//	@Description	Glue 서비스 목록 또는 정보를 보여줍니다.
//	@Tags			Service
//	@param			service_name	query	string	false	"Glue Service Name"
//	@param			service_type	query	string	false	"Glue Service Type"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.ServiceLs
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/service [get]
func (c *Controller) ServiceLs(ctx *gin.Context) {
	service_name := ctx.Request.URL.Query().Get("service_name")
	service_type := ctx.Request.URL.Query().Get("service_type")
	dat, err := glue.ServiceLs(service_name, service_type)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// ServiceControl godoc
//
//	@Summary		Control of Glue Service
//	@Description	Glue 서비스를 제어합니다.
//	@Tags			Service
//	@param			service_name 	path	string	true	"Glue Service Name"
//	@param			control 	    query	string	true	"Glue Service Control" Enums(start, stop, restart)
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/service/{service_name} [post]
func (c *Controller) ServiceControl(ctx *gin.Context) {
	service_name := ctx.Param("service_name")
	control := ctx.Request.URL.Query().Get("control")
	dat, err := glue.ServiceControl(control, service_name)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// ServiceDelete godoc
//
//	@Summary		Delete of Glue Service
//	@Description	Glue 서비스를 삭제합니다.
//	@Tags			Service
//	@param			service_name 	path	string	true	"Glue Service Name"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/service/{service_name} [delete]
func (c *Controller) ServiceDelete(ctx *gin.Context) {
	service_name := ctx.Param("service_name")
	dat, err := glue.ServiceDelete(service_name)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

func (c *Controller) ServiceDeleteOptions(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}

// HostList godoc
//
//	@Summary		Show List of Glue Hosts
//	@Description	Glue 호스트 리스트를 보여줍니다.
//	@Tags			Hosts
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	HostList
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/glue/hosts [get]
func (c *Controller) HostList(ctx *gin.Context) {
	dat, err := glue.HostList()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ip_address, err := glue.HostIp()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	var str []string
	str_data := strings.Split(string(ip_address), "\n")
	for i := 0; i < len(str_data); i++ {
		strs := str_data[i]
		str = append(str, strs)
		if i == len(str_data)-1 {
			str = str[:len(str_data)-1]
		}
	}
	for i := 0; i < len(dat); i++ {
		dat[i].MainAddr = str[i]
	}
	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}
