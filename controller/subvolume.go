package controller

// import (
// 	"Glue-API/httputil"
// 	"Glue-API/model"
// 	"Glue-API/utils"
// 	"Glue-API/utils/fs"
// 	"net/http"
// 	"strconv"
// 	"time"

// 	"github.com/gin-gonic/gin"
// )

// // SubVolumeList godoc
// //
// //	@Summary		Detail Info and List of Glue FS Sub Volumes
// //	@Description	GlueFS의 하위 볼륨에 대한 상세 정보 및 리스트를 보여줍니다.
// //	@param			vol_name 	query	string	true	"Glue FS Sub Volume Name"
// //	@param			group_name 	query	string	false	"Glue FS Sub Volume Group Name"
// //	@Tags			GlueFS-SubVolume
// //	@Accept			x-www-form-urlencoded
// //	@Produce		json
// //	@Success		200	{object}	SubVolumeList
// //	@Failure		400	{object}	httputil.HTTP400BadRequest
// //	@Failure		404	{object}	httputil.HTTP404NotFound
// //	@Failure		500	{object}	httputil.HTTP500InternalServerError
// //	@Router			/api/v1/gluefs/subvolume  [get]
// func (c *Controller) SubVolumeList(ctx *gin.Context) {
// 	ctx.Header("Access-Control-Allow-Origin", "*")

// 	vol_name := ctx.Request.URL.Query().Get("vol_name")
// 	group_name := ctx.Request.URL.Query().Get("group_name")
// 	ls_data, err := fs.SubVolumeLs(vol_name, group_name)
// 	if err != nil {
// 		utils.FancyHandleError(err)
// 		httputil.NewError(ctx, http.StatusInternalServerError, err)
// 		return
// 	}
// 	var value []model.SubVolumeList
// 	for i := 0; i < len(ls_data); i++ {
// 		info_data, err := fs.SubVolumeInfo(vol_name, ls_data[i].Name, group_name)
// 		if err != nil {
// 			utils.FancyHandleError(err)
// 			httputil.NewError(ctx, http.StatusInternalServerError, err)
// 			return
// 		}
// 		snap_data, err := fs.SubVolumeSnapLs(vol_name, ls_data[i].Name, group_name)
// 		if err != nil {
// 			utils.FancyHandleError(err)
// 			httputil.NewError(ctx, http.StatusInternalServerError, err)
// 			return
// 		}
// 		var snap []string
// 		for j := 0; j < len(snap_data); j++ {
// 			snap = append(snap, snap_data[j].Name)
// 		}
// 		value_data := model.SubVolumeList{
// 			Name: ls_data[i].Name,
// 			Info: model.SubVolumeInfo{
// 				Atime:         info_data.Atime,
// 				BytesPcent:    info_data.BytesPcent,
// 				BytesQuota:    info_data.BytesQuota,
// 				BytesUsed:     info_data.BytesUsed,
// 				CreatedAt:     info_data.CreatedAt,
// 				Ctime:         info_data.Ctime,
// 				DataPool:      info_data.DataPool,
// 				Features:      info_data.Features,
// 				Gid:           info_data.Gid,
// 				Mode:          info_data.Mode,
// 				MonAddrs:      info_data.MonAddrs,
// 				Mtime:         info_data.Mtime,
// 				Path:          info_data.Path,
// 				PoolNamespace: info_data.PoolNamespace,
// 				State:         info_data.State,
// 				Type:          info_data.Type,
// 				UID:           info_data.UID,
// 			},
// 			SnapShot: snap,
// 		}
// 		value = append(value, value_data)
// 	}
// 	// Print the output
// 	ctx.IndentedJSON(http.StatusOK, value)
// }

// // SubVolumeCreate godoc
// //
// //	@Summary		Create of Glue FS Sub Volume
// //	@Description	GlueFS의 하위 볼륨을 생성합니다.
// //	@param			vol_name 	formData	string	true	"Glue FS Volume Name"
// //	@param			subvol_name 	formData	string	true	"Glue FS Sub Volume Name"
// //	@param			group_name 	formData	string	false	"Glue FS Sub Volume Group Name"
// //	@param			size 	formData	int	true	"Glue FS Sub Volume Size(default GB)"
// //	@param			data_pool_name 	formData	string	true	"Glue FS Sub Volume Data Pool Name"
// //	@param			mode 	formData	int	true	"Glue FS Sub Volume Permissions"
// //	@Tags			GlueFS-SubVolume
// //	@Accept			x-www-form-urlencoded
// //	@Produce		json
// //	@Success		200	{string}	string "Success"
// //	@Failure		400	{object}	httputil.HTTP400BadRequest
// //	@Failure		404	{object}	httputil.HTTP404NotFound
// //	@Failure		500	{object}	httputil.HTTP500InternalServerError
// //	@Router			/api/v1/gluefs/subvolume  [post]
// func (c *Controller) SubVolumeCreate(ctx *gin.Context) {
// 	ctx.Header("Access-Control-Allow-Origin", "*")

// 	vol_name, _ := ctx.GetPostForm("vol_name")
// 	subvol_name, _ := ctx.GetPostForm("subvol_name")
// 	group_name, _ := ctx.GetPostForm("group_name")
// 	size, _ := ctx.GetPostForm("size")
// 	data_pool_name, _ := ctx.GetPostForm("data_pool_name")
// 	mode, _ := ctx.GetPostForm("mode")
// 	size_data, _ := strconv.Atoi(size)
// 	size_int := size_data * 1024 * 1024 * 1024
// 	size_str := strconv.Itoa(size_int)

// 	dat, err := fs.SubVolumeCreate(vol_name, subvol_name, group_name, size_str, data_pool_name, mode)
// 	if err != nil {
// 		utils.FancyHandleError(err)
// 		httputil.NewError(ctx, http.StatusInternalServerError, err)
// 		return
// 	}
// 	// Print the output
// 	ctx.IndentedJSON(http.StatusOK, dat)
// }

// // SubVolumeDelete godoc
// //
// //	@Summary		Delete of Glue FS Volume
// //	@Description	GlueFS 하위 볼륨을 삭제합니다.
// //	@param			vol_name 	query	string	true	"Glue FS Volume Name"
// //	@param			subvol_name 	query	string	true	"Glue FS Sub Volume Name"
// //	@param			group_name 	query	string	false	"Glue FS Volume Group Name"
// //	@Tags			GlueFS-SubVolume
// //	@Accept			x-www-form-urlencoded
// //	@Produce		json
// //	@Success		200	{string}	string "Success"
// //	@Failure		400	{object}	httputil.HTTP400BadRequest
// //	@Failure		404	{object}	httputil.HTTP404NotFound
// //	@Failure		500	{object}	httputil.HTTP500InternalServerError
// //	@Router			/api/v1/gluefs/subvolume [delete]
// func (c *Controller) SubVolumeDelete(ctx *gin.Context) {
// 	ctx.Header("Access-Control-Allow-Origin", "*")
// 	vol_name := ctx.Request.URL.Query().Get("vol_name")
// 	subvol_name := ctx.Request.URL.Query().Get("subvol_name")
// 	group_name := ctx.Request.URL.Query().Get("group_name")

// 	dat, err := fs.SubVolumeDelete(vol_name, subvol_name, group_name)
// 	if err != nil {
// 		utils.FancyHandleError(err)
// 		httputil.NewError(ctx, http.StatusInternalServerError, err)
// 		return
// 	}
// 	// Print the output
// 	ctx.IndentedJSON(http.StatusOK, dat)
// }

// // SubVolumeResize godoc
// //
// //	@Summary		Update Size of Glue FS Volume
// //	@Description	GlueFS 볼륨의 할당된 사이즈를 수정합니다.
// //	@param			vol_name 	formData	string	true	"Glue FS Volume Name"
// //	@param			subvol_name 	formData	string	true	"Glue FS Sub Volume Name"
// //	@param			group_name 	formData	string	false	"Glue FS Volume Group Name"
// //	@param			new_size 	formData	string	true	"Glue FS Sub Volume New Size(default GB)"
// //	@Tags			GlueFS-SubVolume
// //	@Accept			x-www-form-urlencoded
// //	@Produce		json
// //	@Success		200	{string}	string "Success"
// //	@Failure		400	{object}	httputil.HTTP400BadRequest
// //	@Failure		404	{object}	httputil.HTTP404NotFound
// //	@Failure		500	{object}	httputil.HTTP500InternalServerError
// //	@Router			/api/v1/gluefs/subvolume [put]
// func (c *Controller) SubVolumeResize(ctx *gin.Context) {
// 	ctx.Header("Access-Control-Allow-Origin", "*")

// 	vol_name, _ := ctx.GetPostForm("vol_name")
// 	subvol_name, _ := ctx.GetPostForm("subvol_name")
// 	group_name, _ := ctx.GetPostForm("group_name")
// 	new_size, _ := ctx.GetPostForm("new_size")
// 	new_size_data, _ := strconv.Atoi(new_size)
// 	size_int := new_size_data * 1024 * 1024 * 1024
// 	new_size_str := strconv.Itoa(size_int)

// 	dat, err := fs.SubVolumeResize(vol_name, subvol_name, new_size_str, group_name)
// 	if err != nil {
// 		utils.FancyHandleError(err)
// 		httputil.NewError(ctx, http.StatusInternalServerError, err)
// 		return
// 	}
// 	// Print the output
// 	ctx.IndentedJSON(http.StatusOK, dat)
// }

// // SubVolumeSnapList godoc
// //
// //	@Summary		Show List or Info of Glue FS Sub Volume Snapshot
// //	@Description	GlueFS의 하위 볼륨 스냅샷의 리스트 및 상세 정보를 보여줍니다.
// //	@param			vol_name 	query	string	true	"Glue FS Sub Volume Name"
// //	@param			subvol_name 	query	string	true	"Glue FS Volume Sub Volume Name"
// //	@param			group_name 	query	string	false	"Glue FS Volume Group Name"
// //	@param			snap_name 	query	string	false	"Glue FS Volume SnapShot Name"
// //	@Tags			GlueFS-SubVolume-Snapshot
// //	@Accept			x-www-form-urlencoded
// //	@Produce		json
// //	@Success		200	{object}	SubVolumeAllSnap
// //	@Failure		400	{object}	httputil.HTTP400BadRequest
// //	@Failure		404	{object}	httputil.HTTP404NotFound
// //	@Failure		500	{object}	httputil.HTTP500InternalServerError
// //	@Router			/api/v1/gluefs/subvolume/snapshot  [get]
// func (c *Controller) SubVolumeSnapList(ctx *gin.Context) {
// 	ctx.Header("Access-Control-Allow-Origin", "*")
// 	vol_name := ctx.Request.URL.Query().Get("vol_name")
// 	subvol_name := ctx.Request.URL.Query().Get("subvol_name")
// 	group_name := ctx.Request.URL.Query().Get("group_name")
// 	snap_name := ctx.Request.URL.Query().Get("snap_name")

// 	if snap_name == "" {
// 		dat, err := fs.SubVolumeSnapLs(vol_name, subvol_name, group_name)
// 		if err != nil {
// 			utils.FancyHandleError(err)
// 			httputil.NewError(ctx, http.StatusInternalServerError, err)
// 			return
// 		}
// 		// Print the output
// 		ctx.IndentedJSON(http.StatusOK, dat)
// 	} else {
// 		dat, err := fs.SubVolumeSnapInfo(vol_name, subvol_name, snap_name, group_name)
// 		if err != nil {
// 			utils.FancyHandleError(err)
// 			httputil.NewError(ctx, http.StatusInternalServerError, err)
// 			return
// 		}
// 		// Print the output
// 		ctx.IndentedJSON(http.StatusOK, dat)
// 	}
// }

// // SubVolumeGroupSnapCreate godoc
// //
// //	@Summary		Create of Glue FS Sub Volume Group Snapshot
// //	@Description	GlueFS의 하위 볼륨의 그룹의 스냅샷을 생성합니다.
// //	@param			vol_name 	formData	string	true	"Glue FS Sub Volume Name"
// //	@param			subvol_name 	formData	string	true	"Glue FS Volume Sub Volume Name"
// //	@param			group_name 	formData	string	false	"Glue FS Volume Group Name"
// //	@Tags			GlueFS-SubVolume-Snapshot
// //	@Accept			x-www-form-urlencoded
// //	@Produce		json
// //	@Success		200	{string}	string "Success"
// //	@Failure		400	{object}	httputil.HTTP400BadRequest
// //	@Failure		404	{object}	httputil.HTTP404NotFound
// //	@Failure		500	{object}	httputil.HTTP500InternalServerError
// //	@Router			/api/v1/gluefs/subvolume/snapshot  [post]
// func (c *Controller) SubVolumeSnapCreate(ctx *gin.Context) {
// 	ctx.Header("Access-Control-Allow-Origin", "*")

// 	vol_name, _ := ctx.GetPostForm("vol_name")
// 	subvol_name, _ := ctx.GetPostForm("subvol_name")
// 	group_name, _ := ctx.GetPostForm("group_name")
// 	snap_name := time.Now().Format(time.RFC3339)

// 	dat, err := fs.SubVolumeSnapCreate(vol_name, subvol_name, snap_name, group_name)

// 	if err != nil {
// 		utils.FancyHandleError(err)
// 		httputil.NewError(ctx, http.StatusInternalServerError, err)
// 		return
// 	}
// 	// Print the output
// 	ctx.IndentedJSON(http.StatusOK, dat)
// }

// // SubVolumeSnapDelete godoc
// //
// //	@Summary		Delete of Glue FS Volume Snapshot
// //	@Description	GlueFS 볼륨의 스냅샷을 삭제합니다.
// //	@param			vol_name 	query	string	true	"Glue FS Volume Name"
// //	@param			subvol_name 	query	string	true	"Glue FS Volume Sub Volume Name"
// //	@param			group_name 	query	string	false	"Glue FS Volume Group Name"
// //	@param			snap_name 	query	string	true	"Glue FS Volume Group SnapShot Name"
// //	@Tags			GlueFS-SubVolume-Snapshot
// //	@Accept			x-www-form-urlencoded
// //	@Produce		json
// //	@Success		200	{string}	string "Success"
// //	@Failure		400	{object}	httputil.HTTP400BadRequest
// //	@Failure		404	{object}	httputil.HTTP404NotFound
// //	@Failure		500	{object}	httputil.HTTP500InternalServerError
// //	@Router			/api/v1/gluefs/subvolume/snapshot  [delete]
// func (c *Controller) SubVolumeSnapDelete(ctx *gin.Context) {
// 	ctx.Header("Access-Control-Allow-Origin", "*")
// 	vol_name := ctx.Request.URL.Query().Get("vol_name")
// 	subvol_name := ctx.Request.URL.Query().Get("subvol_name")
// 	group_name := ctx.Request.URL.Query().Get("group_name")
// 	snap_name := ctx.Request.URL.Query().Get("snap_name")

// 	dat, err := fs.SubVolumeSnapDelete(vol_name, subvol_name, snap_name, group_name)
// 	if err != nil {
// 		utils.FancyHandleError(err)
// 		httputil.NewError(ctx, http.StatusInternalServerError, err)
// 		return
// 	}
// 	// Print the output
// 	ctx.IndentedJSON(http.StatusOK, dat)
// }
