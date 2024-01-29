package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/glue"
	"encoding/json"
	"net/http"
	"os/exec"

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

	ctx.IndentedJSON(http.StatusOK, dat)
}

// ListPools godoc
//
//	@Summary		List Pools of Glue
//	@Description	Glue 의 스토리지 풀 목록을 보여줍니다.
//	@Tags			Glue
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.GlueVersion
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/glue/pool [get]
func (c *Controller) ListPools(ctx *gin.Context) {
	var dat model.GluePools
	pools, err := glue.ListPool()
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	dat.Pools = pools
	ctx.IndentedJSON(http.StatusOK, dat)
}

// ListImages godoc
//
//	@Summary		List RBD Images of Pool Glue
//	@Description	Glue 스토리지 풀의 이미지 목록을 보여줍니다.
//	@Tags			Glue
//	@param			pool_name	path	string	true	"pool_name"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.GlueVersion
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/glue/rbd/{pool_name} [get]
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

// ListImages godoc
//
//	@Summary		List Images of Pool Glue
//	@Description	Glue 스토리지 풀의 이미지 목록을 보여줍니다.
//	@Tags			Glue
//	@param			pool_name	path	string	true	"pool_name"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/glue/pool/{pool_name} [delete]
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
