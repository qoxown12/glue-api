package controller

import (
	// "encoding/json"
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils/gluevm"
	"net/http"

	"github.com/gin-gonic/gin"
	// "os/exec"
)

// ListGlueFs godoc
//
//	@Summary		List Fs of Glue
//	@Description	Glue의 파일 시스템 목록을 보여줍니다.
//	@Tags			Glue
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.GlueVersion
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/gluefs [get]
func (c *Controller) VmState(ctx *gin.Context) {
	var dat = struct {
		model.AbleModel
		Message string
	}{}

	message, err := gluevm.VmState()

	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	dat.Message = message
	ctx.IndentedJSON(http.StatusOK, dat)
}

// FsSetup godoc
//
//	@Summary		Setup Glue File System
//	@Description	Glue의 파일 시스템을 생성합니다.
//	@param			privateKeyFile		formData	file	true	"Remote Cluster PrivateKey"
//	@param			mirrorPool			formData	string	true	"Pool Name for Mirroring"
//	@Tags			Mirror
//	@Accept			multipart/form-data
//	@Produce		json
//	@Success		200	{object}	model.MirrorSetup
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror [post]
func (c *Controller) VmSetup(ctx *gin.Context) {
	var dat = struct {
		model.AbleModel
		Message string
	}{}

	message, err := gluevm.VmStart()

	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	dat.Message = message
	ctx.IndentedJSON(http.StatusOK, dat)
}

// FsSetup godoc
//
//	@Summary		Setup Glue File System
//	@Description	Glue의 파일 시스템을 생성합니다.
//	@param			privateKeyFile		formData	file	true	"Remote Cluster PrivateKey"
//	@param			mirrorPool			formData	string	true	"Pool Name for Mirroring"
//	@Tags			Mirror
//	@Accept			multipart/form-data
//	@Produce		json
//	@Success		200	{object}	model.MirrorSetup
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror [post]
func (c *Controller) VmStart(ctx *gin.Context) {
	var dat = struct {
		model.AbleModel
		Message string
	}{}

	message, err := gluevm.VmStart()

	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	dat.Message = message
	ctx.IndentedJSON(http.StatusOK, dat)
}

// FsSetup godoc
//
//	@Summary		Setup Glue File System
//	@Description	Glue의 파일 시스템을 생성합니다.
//	@param			privateKeyFile		formData	file	true	"Remote Cluster PrivateKey"
//	@param			mirrorPool			formData	string	true	"Pool Name for Mirroring"
//	@Tags			Mirror
//	@Accept			multipart/form-data
//	@Produce		json
//	@Success		200	{object}	model.MirrorSetup
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror [post]
func (c *Controller) VmStop(ctx *gin.Context) {
	var dat = struct {
		model.AbleModel
		Message string
	}{}

	message, err := gluevm.VmStop()

	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	dat.Message = message
	ctx.IndentedJSON(http.StatusOK, dat)
}

// FsSetup godoc
//
//	@Summary		Setup Glue File System
//	@Description	Glue의 파일 시스템을 생성합니다.
//	@param			privateKeyFile		formData	file	true	"Remote Cluster PrivateKey"
//	@param			mirrorPool			formData	string	true	"Pool Name for Mirroring"
//	@Tags			Mirror
//	@Accept			multipart/form-data
//	@Produce		json
//	@Success		200	{object}	model.MirrorSetup
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror [post]
func (c *Controller) VmDelete(ctx *gin.Context) {
	var dat = struct {
		model.AbleModel
		Message string
	}{}

	message, err := gluevm.VmDelete()

	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	dat.Message = message
	ctx.IndentedJSON(http.StatusOK, dat)
}

// FsSetup godoc
//
//	@Summary		Setup Glue File System
//	@Description	Glue의 파일 시스템을 생성합니다.
//	@param			privateKeyFile		formData	file	true	"Remote Cluster PrivateKey"
//	@param			mirrorPool			formData	string	true	"Pool Name for Mirroring"
//	@Tags			Mirror
//	@Accept			multipart/form-data
//	@Produce		json
//	@Success		200	{object}	model.MirrorSetup
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror [post]
func (c *Controller) VmCleanup(ctx *gin.Context) {
	var dat = struct {
		model.AbleModel
		Message string
	}{}

	message, err := gluevm.VmCleanup()

	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	dat.Message = message
	ctx.IndentedJSON(http.StatusOK, dat)
}

// FsSetup godoc
//
//	@Summary		Setup Glue File System
//	@Description	Glue의 파일 시스템을 생성합니다.
//	@param			privateKeyFile		formData	file	true	"Remote Cluster PrivateKey"
//	@param			mirrorPool			formData	string	true	"Pool Name for Mirroring"
//	@Tags			Mirror
//	@Accept			multipart/form-data
//	@Produce		json
//	@Success		200	{object}	model.MirrorSetup
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror [post]
func (c *Controller) VmMigrate(ctx *gin.Context) {
	var dat = struct {
		model.AbleModel
		Message string
	}{}

	target, _ := ctx.GetPostForm("target")

	message, err := gluevm.VmMigrate(target)

	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	dat.Message = message
	ctx.IndentedJSON(http.StatusOK, dat)
}
