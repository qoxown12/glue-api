package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/smb"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *Controller) SmbOption(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}

// SmbStatus godoc
//
//	@Summary		Show Status of Smb Servcie Daemon
//	@Description	SMB 서비스 데몬 상태를 조회합니다.
//	@Tags			SMB
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	SmbStatus
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/smb [get]
func (c *Controller) SmbStatus(ctx *gin.Context) {

	hosts_data, err := smb.Hosts()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	var hosts []string
	for i := 0; i < len(hosts_data); i++ {
		hosts = append(hosts, hosts_data[i])
		if i == len(hosts_data)-1 {
			hosts = hosts[:len(hosts_data)-1]
		}
	}
	var smb_status []model.SmbStatus
	for i := 0; i < len(hosts); i++ {
		status, _ := smb.SmbStatus(hosts[i])

		smb_status = append(smb_status, status)
		if i == len(hosts)-1 {
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.IndentedJSON(http.StatusOK, smb_status)
		}

	}
}

// SmbCreate godoc
//
//	@Summary		Create of Smb Service
//	@Description	SMB 서비스 전체를 생성합니다.
//	@Tags			SMB
//	@param			hostname     formData   []string	true    "SMB Server Host Name" collectionFormat(multi)
//	@param			username     formData   string	true    "SMB Username"
//	@param			password     formData   string	true    "SMB Password"
//	@param			folder_name     formData   string	true    "SMB Share Folder Name"
//	@param			path    formData   string	true    "SMB Server Actual Shared Path"
//	@param			fs_name     formData   string	true    "Glue File System Name"
//	@param			volume_path    formData   string	true    "Glue File System Volume Path"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/smb [post]
func (c *Controller) SmbCreate(ctx *gin.Context) {
	hostname, _ := ctx.GetPostFormArray("hostname")
	username, _ := ctx.GetPostForm("username")
	password, _ := ctx.GetPostForm("password")
	folder, _ := ctx.GetPostForm("folder_name")
	path, _ := ctx.GetPostForm("path")
	fs_name, _ := ctx.GetPostForm("fs_name")
	volume_path, _ := ctx.GetPostForm("volume_path")

	for i := 0; i < len(hostname); i++ {
		dat, err := smb.SmbCreate(hostname[i], username, password, folder, path, fs_name, volume_path)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		if i == len(hostname)-1 {
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.IndentedJSON(http.StatusOK, dat)
		}
	}
}

// SmbUserCreate godoc
//
//	@Summary		Create User of Smb Service
//	@Description	SMB 서비스 사용자를 생성합니다.
//	@Tags			SMB
//	@param			hostname     formData   []string	true    "SMB Server Host Name" collectionFormat(multi)
//	@param			username     formData   string	true    "SMB Username"
//	@param			password     formData   string	true    "SMB Password"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/smb/user [post]
func (c *Controller) SmbUserCreate(ctx *gin.Context) {
	hostname, _ := ctx.GetPostFormArray("hostname")
	username, _ := ctx.GetPostForm("username")
	password, _ := ctx.GetPostForm("password")

	for i := 0; i < len(hostname); i++ {
		dat, err := smb.SmbUserCreate(hostname[i], username, password)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		if i == len(hostname)-1 {
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.IndentedJSON(http.StatusOK, dat)
		}
	}
}

// SmbUserUpdate godoc
//
//	@Summary		Update User of Smb Service
//	@Description	SMB 서비스 사용자의 패스워드를 변경합니다.
//	@Tags			SMB
//	@param			hostname     formData   []string	true    "SMB Server Host Name" collectionFormat(multi)
//	@param			username     formData   string	true    "SMB Username"
//	@param			password     formData   string	true    "SMB Password"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/smb/user [put]
func (c *Controller) SmbUserUpdate(ctx *gin.Context) {
	hostname, _ := ctx.GetPostFormArray("hostname")
	username, _ := ctx.GetPostForm("username")
	password, _ := ctx.GetPostForm("password")

	for i := 0; i < len(hostname); i++ {
		dat, err := smb.SmbUserUpdate(hostname[i], username, password)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		if i == len(hostname)-1 {
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.IndentedJSON(http.StatusOK, dat)
		}
	}
}

// SmbDelete godoc
//
//	@Summary		Delete of Smb Service
//	@Description	SMB 서비스 전체를 삭제합니다.
//	@Tags			SMB
//	@param			hostname     query   []string	true    "SMB Server Host Name" collectionFormat(multi)
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/smb [delete]
func (c *Controller) SmbDelete(ctx *gin.Context) {
	hostname_arr := ctx.QueryArray("hostname")
	for i := 0; i < len(hostname_arr); i++ {
		dat, err := smb.SmbDelete(hostname_arr[i])
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		if i == len(hostname_arr)-1 {
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.IndentedJSON(http.StatusOK, dat)
		}
	}
}

// SmbUserDelete godoc
//
//	@Summary		Delete User of Smb Service
//	@Description	SMB 서비스 사용자를 삭제합니다.
//	@Tags			SMB
//	@param			hostname     query   []string	true    "SMB Server Host Name" collectionFormat(multi)
//	@param			username     query   string	true    "SMB Username"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/smb/user [delete]
func (c *Controller) SmbUserDelete(ctx *gin.Context) {
	hostname_arr := ctx.QueryArray("hostname")
	username := ctx.Request.URL.Query().Get("username")

	for i := 0; i < len(hostname_arr); i++ {
		dat, err := smb.SmbUserDelete(hostname_arr[i], username)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		if i == len(hostname_arr)-1 {
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.IndentedJSON(http.StatusOK, dat)
		}
	}
}
