package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	// "Glue-API/utils"
	"Glue-API/utils/gluefs"
	// "encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
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
func (c *Controller) ListFs(ctx *gin.Context) {
	var dat model.GlueFsList
	glue_fs_list, err := gluefs.ListFs()
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	dat.Debug = gin.IsDebugging()
	dat.GlueFsList = glue_fs_list
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
func (c *Controller) FsSetup(ctx *gin.Context) {
	// var dat model.FsSetup

	fsName := ctx.Query("fsName")
	print("1111:")
	print(fsName)
	print(":::")
	
	id2 := ctx.Param("id")
	print("2222:")
	print(id2)
	print(":::")

	fsName,aaa := ctx.GetPostForm("fsName")
	print("3333:")
	print(fsName)
	print(aaa)
	print(":::")

	// // Upload the file to specific dst.
	// err = ctx.SaveUploadedFile(file, privkeyname)
	// if err != nil {
	// 	httputil.NewError(ctx, http.StatusInternalServerError, err)
	// 	return
	// }

	// if gin.IsDebugging() == true {
	// 	EncodedLocalToken, EncodedRemoteToken, err := mirror.ConfigMirror(dat, privkeyname)
	// 	if err != nil {
	// 		utils.FancyHandleError(err)
	// 		httputil.NewError(ctx, http.StatusInternalServerError, err)
	// 		return
	// 	}

	// 	dat.LocalToken = EncodedLocalToken
	// 	dat.RemoteToken = EncodedRemoteToken

	// } else {
	// 	println("...")
	// }
	// dat.Debug = gin.IsDebugging()
	// ctx.IndentedJSON(http.StatusOK, dat)
}