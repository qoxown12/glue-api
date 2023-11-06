package controller

import (
	"Glue-API/docs"
	"Glue-API/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Controller example
type Controller struct {
}

// NewController example
func NewController() *Controller {
	return &Controller{}
}

// Message example
type Message struct {
	Message string `json:"message" example:"message"`
} //@name Message

// Version godoc
//
//	@Summary		Show Versions of API
//	@Description	API 의 버전을 보여줍니다.
//	@Tags			API
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.Version
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/version [get]
func (c *Controller) Version(ctx *gin.Context) {
	dat := model.Version{Version: docs.SwaggerInfo.Version}
	// Print the output
	dat.Debug = gin.IsDebugging()
	ctx.IndentedJSON(http.StatusOK, dat)

}
