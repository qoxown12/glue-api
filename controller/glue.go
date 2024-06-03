package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/glue"
	"encoding/json"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (c *Controller) GlueOption(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}

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
	ctx.Header("Access-Control-Allow-Origin", "*")

	dat, err := glue.Status()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
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
	ctx.Header("Access-Control-Allow-Origin", "*")

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
	ctx.IndentedJSON(http.StatusOK, dat)
}

// ListPools godoc
//
//	@Summary		List Pools of Glue
//	@Description	Glue 의 스토리지 풀 목록을 보여줍니다.
//	@Tags			Pool
//	@param			pool_type	query	string	false	"pool_type"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.GluePools
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/pool [get]
func (c *Controller) ListPools(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	pool_type := ctx.Request.URL.Query().Get("pool_type")
	dat, err := glue.ListPool(pool_type)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, dat)
}

// PoolDelete godoc
//
//	@Summary		Delete of Pool
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
	ctx.Header("Access-Control-Allow-Origin", "*")

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

// ListAndInfoImage godoc
//
//	@Summary		Show List or Info Images of Pool
//	@Description	Glue 스토리지 풀의 이미지 목록을 보여줍니다.
//	@Tags			Image
//	@param			pool_name	query	string	false	"Glue Pool Name"
//	@param			image_name	query	string	false	"Glue Image Name"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	GluePools
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/image [get]
func (c *Controller) ListAndInfoImage(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	pool_name := ctx.Request.URL.Query().Get("pool_name")
	image_name := ctx.Request.URL.Query().Get("image_name")

	if image_name == "" && pool_name == "" {
		rbd_pool_dat, err := glue.RbdPool()
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		var pools []string
		for i := 0; i < len(rbd_pool_dat); i++ {
			rbd_image_dat, err := glue.RbdImage(rbd_pool_dat[i])
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			for j := 0; j < len(rbd_image_dat); j++ {
				name := rbd_pool_dat[i] + string("/") + rbd_image_dat[j]
				pools = append(pools, name)
			}
		}
		ctx.IndentedJSON(http.StatusOK, pools)
	} else if image_name == "" && pool_name != "" {
		dat, err := glue.InfoImage(pool_name)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(http.StatusOK, dat)
	} else {
		dat, err := glue.ListAndInfoImage(image_name, pool_name)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(http.StatusOK, dat)
	}
}

// CreateImage godoc
//
//	@Summary		Create Images of Pool
//	@Description	Glue 스토리지 풀의 이미지를 생성합니다.
//	@Tags			Image
//	@param			image_name	formData	string	true	"Glue Image Name"
//	@param			pool_name	formData	string	true	"Glue Pool Name"
//	@param			size	formData	int 	true	"Image Size(default:GB)"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/image [post]
func (c *Controller) CreateImage(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	image_name, _ := ctx.GetPostForm("image_name")
	pool_name, _ := ctx.GetPostForm("pool_name")
	size, _ := ctx.GetPostForm("size")
	size_int, err := strconv.Atoi(size)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	size_int = size_int * 1024
	size_st := strconv.Itoa(size_int)
	dat, err := glue.CreateImage(image_name, pool_name, size_st)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, dat)
}

// DeleteImage godoc
//
//	@Summary		Delete Images of Pool
//	@Description	Glue 스토리지 풀의 이미지를 삭제합니다.
//	@Tags			Image
//	@param			image_name	query	string	true	"Glue Image Name"
//	@param			pool_name	query	string	true	"Glue Pool Name"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/image [delete]
func (c *Controller) DeleteImage(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	image_name := ctx.Request.URL.Query().Get("image_name")
	pool_name := ctx.Request.URL.Query().Get("pool_name")
	dat, err := glue.DeleteImage(image_name, pool_name)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, dat)
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
	ctx.Header("Access-Control-Allow-Origin", "*")

	service_name := ctx.Request.URL.Query().Get("service_name")
	service_type := ctx.Request.URL.Query().Get("service_type")
	dat, err := glue.ServiceLs(service_name, service_type)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Print the output
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
	ctx.Header("Access-Control-Allow-Origin", "*")

	service_name := ctx.Param("service_name")
	control := ctx.Request.URL.Query().Get("control")
	dat, err := glue.ServiceControl(control, service_name)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
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
	ctx.Header("Access-Control-Allow-Origin", "*")

	service_name := ctx.Param("service_name")
	// if strings.Contains(service_name, "rgw") {
	// 	rgw_dat, err := glue.RgwPool()
	// 	if err != nil {
	// 		utils.FancyHandleError(err)
	// 		httputil.NewError(ctx, http.StatusInternalServerError, err)
	// 		return
	// 	}
	// 	for i := 0; i < len(rgw_dat); i++ {
	// 		a, err := glue.PoolDelete(rgw_dat[i])
	// 		if err != nil {
	// 			utils.FancyHandleError(err)
	// 			httputil.NewError(ctx, http.StatusInternalServerError, err)
	// 			return
	// 		}
	// 		if i == len(rgw_dat)-1 {
	// 			ctx.IndentedJSON(http.StatusOK, "Success")
	// 		}
	// 	}
	// }
	dat, err := glue.ServiceDelete(service_name)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.IndentedJSON(http.StatusOK, dat)
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
	ctx.Header("Access-Control-Allow-Origin", "*")

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
		dat[i].Ip_Address = str[i]
	}
	// Print the output
	ctx.IndentedJSON(http.StatusOK, dat)
}
