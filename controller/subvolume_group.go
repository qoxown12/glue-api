package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/fs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (c *Controller) SubVolumeGroupOption(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}

// SubVolumeGroupList godoc
//
//	@Summary		Detail Info and List of Glue FS Volume Groups
//	@Description	GlueFS볼륨의 그룹에 대한 상세 정보 및 리스트를 보여줍니다.
//	@param			vol_name 	query	string	true	"Glue FS Volume Name"
//	@Tags			GlueFS-SubVolume-Group
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	SubVolumeGroupList
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/gluefs/subvolume/group  [get]
func (c *Controller) SubVolumeGroupList(ctx *gin.Context) {
	vol_name := ctx.Request.URL.Query().Get("vol_name")
	ls_data, err := fs.SubVolumeGroupLs(vol_name)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	var value []model.SubVolumeGroupList
	for i := 0; i < len(ls_data); i++ {
		info_data, err := fs.SubVolumeGroupInfo(vol_name, ls_data[i].Name)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		path_data, err := fs.SubVolumeGroupGetPath(vol_name, ls_data[i].Name)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		snap_data, err := fs.SubVolumeGroupSnapLs(vol_name, ls_data[i].Name)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		var snap []string
		for j := 0; j < len(snap_data); j++ {
			snap = append(snap, snap_data[j].Name)
		}
		value_data := model.SubVolumeGroupList{
			Name: ls_data[i].Name,
			Info: model.SubVolumeGroupInfo{
				Atime:      info_data.Atime,
				BytesPcent: info_data.BytesPcent,
				BytesQuota: info_data.BytesQuota,
				BytesUsed:  info_data.BytesUsed,
				CreatedAt:  info_data.CreatedAt,
				Ctime:      info_data.Ctime,
				DataPool:   info_data.DataPool,
				Gid:        info_data.Gid,
				Mode:       info_data.Mode,
				MonAddrs:   info_data.MonAddrs,
				Mtime:      info_data.Mtime,
				UID:        info_data.UID,
			},
			Path:     path_data,
			Snapshot: snap,
		}
		value = append(value, value_data)
	}
	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, value)
}

// SubVolumeGroupCreate godoc
//
//	@Summary		Create of Glue FS Volume Group
//	@Description	GlueFS 볼륨의 그룹을 생성합니다.
//	@param			vol_name 	formData	string	true	"Glue FS Volume Name"
//	@param			group_name 	formData	string	true	"Glue FS Volume Group Name"
//	@param			size 	formData	int	true	"Glue FS Volume Group Size(default GB)"
//	@param			data_pool_name 	formData	string	true	"Glue FS Volume Group Data Pool Name"
//	@param			mode 	formData	int	true	"Glue FS Volume Group Permissions"
//	@Tags			GlueFS-SubVolume-Group
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/gluefs/subvolume/group  [post]
func (c *Controller) SubVolumeGroupCreate(ctx *gin.Context) {
	vol_name, _ := ctx.GetPostForm("vol_name")
	group_name, _ := ctx.GetPostForm("group_name")
	size, _ := ctx.GetPostForm("size")
	data_pool_name, _ := ctx.GetPostForm("data_pool_name")
	mode, _ := ctx.GetPostForm("mode")
	size_data, _ := strconv.Atoi(size)
	size_int := size_data * 1024 * 1024 * 1024
	size_str := strconv.Itoa(size_int)

	dat, err := fs.SubVolumeGroupCreate(vol_name, group_name, size_str, data_pool_name, mode)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// SubVolumeGroupDelete godoc
//
//	@Summary		Delete of Glue FS Volume Group
//	@Description	GlueFS 볼륨의 그룹을 삭제합니다.
//	@param			vol_name 	query	string	true	"Glue FS Volume Name"
//	@param			group_name 	query	string	true	"Glue FS Volume Group Name"
//	@param			path 	query	string	true	"Glue FS Volume Group Path"
//	@Tags			GlueFS-SubVolume-Group
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/gluefs/subvolume/group  [delete]
func (c *Controller) SubVolumeGroupDelete(ctx *gin.Context) {
	vol_name := ctx.Request.URL.Query().Get("vol_name")
	group_name := ctx.Request.URL.Query().Get("group_name")
	path := ctx.Request.URL.Query().Get("path")

	dat, err := fs.SubVolumeGroupDelete(vol_name, group_name, path)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// SubVolumeGroupResize godoc
//
//	@Summary		Update Size of Glue FS Volume Group
//	@Description	GlueFS 볼륨의 그룹의 할당된 사이즈를 수정합니다.
//	@param			vol_name 	formData	string	true	"Glue FS Volume Name"
//	@param			group_name 	formData	string	true	"Glue FS Volume Group Name"
//	@param			new_size 	formData	string	true	"Glue FS Volume Group New Size(default GB)"
//	@Tags			GlueFS-SubVolume-Group
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/gluefs/subvolume/group [put]
func (c *Controller) SubVolumeGroupResize(ctx *gin.Context) {
	vol_name, _ := ctx.GetPostForm("vol_name")
	group_name, _ := ctx.GetPostForm("group_name")
	new_size, _ := ctx.GetPostForm("new_size")
	new_size_data, _ := strconv.Atoi(new_size)
	size_int := new_size_data * 1024 * 1024 * 1024
	new_size_str := strconv.Itoa(size_int)

	dat, err := fs.SubVolumeGroupResize(vol_name, group_name, new_size_str)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// // SubVolumeGroupSnapDelete godoc
// //
// //	@Summary		Delete of Glue FS Volume Group Snapshot
// //	@Description	GlueFS 볼륨의 그룹의 스냅샷을 삭제합니다.
// //	@param			vol_name 	query	string	true	"Glue FS Volume Name"
// //	@param			group_name 	query	string	true	"Glue FS Volume Group Name"
// //	@param			snap_name 	query	string	true	"Glue FS Volume Group SnapShot Name"
// //	@Tags			GlueFS-SubVolume-Group
// //	@Accept			x-www-form-urlencoded
// //	@Produce		json
// //	@Success		200	{string}	string "Success"
// //	@Failure		400	{object}	httputil.HTTP400BadRequest
// //	@Failure		404	{object}	httputil.HTTP404NotFound
// //	@Failure		500	{object}	httputil.HTTP500InternalServerError
// //	@Router			/api/v1/gluefs/subvolume/group/snapshot  [delete]
// func (c *Controller) SubVolumeGroupSnapDelete(ctx *gin.Context) {
// 	vol_name := ctx.Request.URL.Query().Get("vol_name")
// 	group_name := ctx.Request.URL.Query().Get("group_name")
// 	snap_name := ctx.Request.URL.Query().Get("snap_name")

// 	dat, err := fs.SubVolumeGroupSnapDelete(vol_name, group_name, snap_name)
// 	if err != nil {
// 		utils.FancyHandleError(err)
// 		httputil.NewError(ctx, http.StatusInternalServerError, err)
// 		return
// 	}
// 	// Print the output
// 	ctx.Header("Access-Control-Allow-Origin", "*")
// 	ctx.IndentedJSON(http.StatusOK, dat)
// }
