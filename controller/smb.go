package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/smb"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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
	status_dat, err := smb.SmbStatus()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	user_dat, err := smb.SmbUserMngt()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	var user []string
	user_data := strings.Split(string(user_dat), "\n")
	for i := 0; i < len(user_data); i++ {
		users := user_data[i]
		user = append(user, users)
		if i == len(user_data)-1 {
			user = user[:len(user_data)-1]
		}
	}
	var names string
	var description string
	var status string
	var state string
	data := strings.FieldsFunc(status_dat, Split)
	for i := 0; i < len(data); i++ {
		if data[i] == "Names" {
			names = data[i+1]
		}
		if data[i] == "Description" {
			description = data[i+1]
		}
		if data[i] == "ActiveState" {
			status = data[i+1]
		}
		if data[i] == "UnitFileState" {
			state = data[i+1]
		}
	}
	dat_hostname, err := smb.Hostname()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	hostname := strings.Replace(dat_hostname, "\n", "", -1)
	dat_ip, err := smb.IpAddress()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ip := strings.Replace(dat_ip, "\n", "", -1)
	dat_port, err := smb.Port()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	port := strings.Split(dat_port, "\n")
	for i := 0; i < len(port); i++ {
		if i == len(port)-1 {
			port = port[:len(port)-1]
		}
	}
	fmt.Print(port)
	value := model.SmbStatus{
		Hostname:    hostname,
		IpAddress:   ip,
		Names:       names,
		Port:        port,
		Description: description,
		Status:      status,
		State:       state,
		Users: model.Users{
			Users: user,
		},
	}
	var ret model.SmbStatus
	json_data, _ := json.Marshal(value)
	json.Unmarshal(json_data, &ret)
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, ret)
}
func Split(r rune) bool {
	return r == '=' || r == '\n'
}

// SmbCreate godoc
//
//	@Summary		Create of Smb Service
//	@Description	SMB 서비스 전체를 생성합니다.
//	@Tags			SMB
//	@param			username     formData   string	true    "SMB Username"
//	@param			password     formData   string	true    "SMB Password"
//	@param			folder_name     formData   string	true    "SMB Share Folder Name"
//	@param			path    formData   string	true    "SMB Server Actual Shared Path"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/smb [post]
func (c *Controller) SmbCreate(ctx *gin.Context) {
	username, _ := ctx.GetPostForm("username")
	password, _ := ctx.GetPostForm("password")
	folder, _ := ctx.GetPostForm("folder_name")
	path, _ := ctx.GetPostForm("path")

	dat, err := smb.SmbCreate(username, password, folder, path)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// SmbUserCreate godoc
//
//	@Summary		Create User of Smb Service
//	@Description	SMB 서비스 사용자를 생성합니다.
//	@Tags			SMB
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
	username, _ := ctx.GetPostForm("username")
	password, _ := ctx.GetPostForm("password")

	dat, err := smb.SmbUserCreate(username, password)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// SmbUserUpdate godoc
//
//	@Summary		Update User of Smb Service
//	@Description	SMB 서비스 사용자의 패스워드를 변경합니다.
//	@Tags			SMB
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
	username, _ := ctx.GetPostForm("username")
	password, _ := ctx.GetPostForm("password")

	dat, err := smb.SmbUserUpdate(username, password)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}
func (c *Controller) SmbUserOptions(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}

// SmbDelete godoc
//
//	@Summary		Delete of Smb Service
//	@Description	SMB 서비스 전체를 삭제합니다.
//	@Tags			SMB
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/smb [delete]
func (c *Controller) SmbDelete(ctx *gin.Context) {

	dat, err := smb.SmbDelete()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}
func (c *Controller) SmbOptions(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}

// SmbUserDelete godoc
//
//	@Summary		Delete User of Smb Service
//	@Description	SMB 서비스 사용자를 삭제합니다.
//	@Tags			SMB
//	@param			username     query   string	true    "SMB Username"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/smb/user [delete]
func (c *Controller) SmbUserDelete(ctx *gin.Context) {
	username := ctx.Request.URL.Query().Get("username")

	dat, err := smb.SmbUserDelete(username)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}
