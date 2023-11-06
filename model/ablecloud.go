package model

// AbleModel
// @description Ablestack Model
type AbleModel struct {
	Debug bool `json:"debug" example:"true" format:"bool"` //Debug info
}

// Version API Version
// @Description API의 버전
type Version struct {
	AbleModel
	Version string `json:"version" example:"1.0" format:"string"`
} //@name Version
