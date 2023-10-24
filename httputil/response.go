package httputil

import (
	"Glue-API/model"
	"github.com/gin-gonic/gin"
)

// NewError example
func NewError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

// HTTPError example
type HTTPError struct {
	model.AbleModel
	Code    int    `json:"code"`
	Message string `json:"message"`
} //@name HTTPError
type HTTP400BadRequest struct {
	HTTPError
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
} //@name HTTP400BadRequest
type HTTP404NotFound struct {
	HTTPError
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" example:"Not Found"`
} //@name HTTP404NotFound
type HTTP500InternalServerError struct {
	HTTPError
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"InternalServerError"`
} //@name HTTP500InternalServerError
