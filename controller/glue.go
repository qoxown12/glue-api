package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/glue"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
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
	dat.Debug = gin.IsDebugging()
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

	if gin.IsDebugging() != true {
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
	} else {
		// Print the output
		versions := []byte("{\n    \"mon\": {\n        \"ceph version Glue-Diplo-4.0.0 (5dd24139a1eada541a3bc16b6941c5dde975e26d) reef (stable)\": 3\n    },\n    \"mgr\": {\n        \"ceph version Glue-Diplo-4.0.0 (5dd24139a1eada541a3bc16b6941c5dde975e26d) reef (stable)\": 2\n    },\n    \"osd\": {\n        \"ceph version Glue-Diplo-4.0.0 (5dd24139a1eada541a3bc16b6941c5dde975e26d) reef (stable)\": 19\n    },\n    \"rbd-mirror\": {\n        \"ceph version Glue-Diplo-4.0.0 (5dd24139a1eada541a3bc16b6941c5dde975e26d) reef (stable)\": 2\n    },\n    \"rgw\": {\n        \"ceph version Glue-Diplo-4.0.0 (5dd24139a1eada541a3bc16b6941c5dde975e26d) reef (stable)\": 1\n    },\n    \"overall\": {\n        \"ceph version Glue-Diplo-4.0.0 (5dd24139a1eada541a3bc16b6941c5dde975e26d) reef (stable)\": 27\n    }\n}")
		if err := json.Unmarshal(versions, &dat); err != nil {
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}
	dat.Debug = gin.IsDebugging()
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
	dat.Debug = gin.IsDebugging()
	dat.Pools = pools
	ctx.IndentedJSON(http.StatusOK, dat)
}

// ListImages godoc
//
//	@Summary		List Images of Pool Glue
//	@Description	Glue 스토리지 풀의 이미지 목록을 보여줍니다.
//	@Tags			Glue
//	@param			pool			path	string	true	"pool"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.GlueVersion
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/glue/pool/{pool} [get]
func (c *Controller) ListImages(ctx *gin.Context) {
	var dat model.SnapshotList
	pool := ctx.Param("pool")
	images, err := glue.ListImage(pool)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	dat.Debug = gin.IsDebugging()
	dat.Images = images
	ctx.IndentedJSON(http.StatusOK, dat)
}
