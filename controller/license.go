package controller

import (
	"Glue-API/httputil"
	"Glue-API/utils"
	"Glue-API/utils/license"
	"net/http"

	"github.com/gin-gonic/gin"
)

// License godoc
//
//	@Summary		Show License
//	@Description	라이센스를 조회합니다.
//	@Tags			License
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	LicenseList
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/license [get]
func (c *Controller) License(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	license_data, err := license.License()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, license_data)
}
