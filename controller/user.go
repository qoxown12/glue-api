package controller

import (
	"Glue-API/httputil"
	"Glue-API/utils"
	"Glue-API/utils/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserCreate godoc
//
//	@Summary		Create User of Service
//	@Description	서비스 사용자를 생성합니다.
//	@param			username 	path	string	true	"Username"
//	@Tags			USER
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/user/{username} [POST]
func (c *Controller) UserCreate(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	username := ctx.Param("username")

	dat, err := user.UserCreate(username)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.IndentedJSON(http.StatusOK, dat)
}

// UserDelete godoc
//
//	@Summary		Delete User of Service
//	@Description	서비스 사용자를 삭제합니다.
//	@Tags			USER
//	@param			username     path   string	true    "Username"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/user/{username} [DELETE]
func (c *Controller) UserDelete(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	username := ctx.Param("username")

	dat, err := user.UserDelete(username)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, dat)
}
