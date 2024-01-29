package controller

import (
	"Glue-API/httputil"
	"Glue-API/utils"
	"Glue-API/utils/nfs"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// NfsClusterLs godoc
//
//	@Summary		Show List of Glue NFS Cluster
//	@Description	Glue NFS Cluster의 리스트를 보여줍니다.
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	NfsClusterLs
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs [get]
func (c *Controller) NfsClusterLs(ctx *gin.Context) {
	dat, err := nfs.NfsClusterLs()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.IndentedJSON(http.StatusOK, dat)
}

// NfsClusterInfo godoc
//
//	@Summary		Show Info of Glue NFS Cluster
//	@Description	Glue NFS Cluster의 상세 정보를 보여줍니다.
//	@param			cluster_id 	path	string	true	"NFS Cluster Identifier"
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	NfsClusterInfo
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs/{cluster_id} [get]
func (c *Controller) NfsClusterInfo(ctx *gin.Context) {

	cluster_id := ctx.Param("cluster_id")

	dat, err := nfs.NfsClusterInfo(cluster_id)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.IndentedJSON(http.StatusOK, dat)
}

// NfsClusterCreate godoc
//
//	@Summary		Create of Glue NFS Cluster
//	@Description	Glue NFS Cluster를 생성합니다.
//	@param			cluster_id 	path	string	true	"NFS Cluster Identifier"
//	@param			port 	path	string	true	"Cluster Port"
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs/{cluster_id}/{port} [post]
func (c *Controller) NfsClusterCreate(ctx *gin.Context) {
	cluster_id := ctx.Param("cluster_id")
	port := ctx.Param("port")
	dat, err := nfs.NfsClusterCreate(cluster_id, port)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.IndentedJSON(http.StatusOK, dat)
}

// NfsClusterDelete godoc
//
//	@Summary		Delete of Glue NFS Cluster
//	@Description	Glue NFS Cluster를 삭제합니다.
//	@param			cluster_id 	path	string	true	"NFS Cluster Identifier"
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs/{cluster_id} [delete]
func (c *Controller) NfsClusterDelete(ctx *gin.Context) {
	cluster_id := ctx.Param("cluster_id")
	dat, err := nfs.NfsClusterDelete(cluster_id)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.IndentedJSON(http.StatusOK, dat)
}

// NfsExportCreate godoc
//
//	@Summary		Create of Glue NFS Export
//	@Description	Glue NFS Export를 생성합니다.
//	@param			cluster_id 	path	string	true	"NFS Cluster Identifier"
//	@param			json_file 	body	NfsExportCreate true 	"NFS Export JSON file"
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs/export/{cluster_id} [post]
func (c *Controller) NfsExportCreate(ctx *gin.Context) {
	json_data, err := io.ReadAll(ctx.Request.Body)
	nfs_export_create_conf := "/usr/share/cockpit/ablestack/tools/properties/nfs_export_create.conf"
	err = os.WriteFile(nfs_export_create_conf, json_data, 0644)

	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		cluster_id := ctx.Param("cluster_id")
		dat, err := nfs.NfsExportCreateOrUpdate(cluster_id, nfs_export_create_conf)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			if err := os.Remove(nfs_export_create_conf); err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
			}
		}
		ctx.IndentedJSON(http.StatusOK, dat)
	}
	return
}

// NfsExportUpdate godoc
//
//	@Summary		Update of Glue NFS Export
//	@Description	Glue NFS Export를 수정합니다.
//	@param			cluster_id 	path	string	true	"NFS Cluster Identifier"
//	@param			json_file 	body	NfsExportUpdate true 	"NFS Export JSON file"
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs/export/{cluster_id} [put]
func (c *Controller) NfsExportUpdate(ctx *gin.Context) {
	json_data, err := io.ReadAll(ctx.Request.Body)
	nfs_export_update_conf := "/usr/share/cockpit/ablestack/tools/properties/nfs_export_update.conf"
	err = os.WriteFile(nfs_export_update_conf, json_data, 0644)

	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		cluster_id := ctx.Param("cluster_id")

		dat, err := nfs.NfsExportCreateOrUpdate(cluster_id, nfs_export_update_conf)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			if err := os.Remove(nfs_export_update_conf); err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
			}
		}
		ctx.IndentedJSON(http.StatusOK, dat)
	}
	return
}

// NfsExportDelete godoc
//
//	@Summary		Delete of Glue NFS Export
//	@Description	Glue NFS Export를 삭제합니다.
//	@param			cluster_id 	path	string	true	"NFS Cluster Identifier"
//	@param			export_id 	path	string	true	"NFS Export ID"
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs/export/{cluster_id}/{export_id} [delete]
func (c *Controller) NfsExportDelete(ctx *gin.Context) {
	cluster_id := ctx.Param("cluster_id")
	export_id, err := strconv.Atoi(ctx.Param("export_id"))
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	detail, err := nfs.NfsExportDetailed(cluster_id)

	for i := 0; i < len(detail); i++ {
		if detail[i].ExportID == export_id {
			dat, err := nfs.NfsExportDelete(cluster_id, detail[i].Pseudo)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			ctx.IndentedJSON(http.StatusOK, dat)
		}
	}
	// Print the output
}

// NfsExportDetailed godoc
//
//	@Summary		Show Detail of Glue NFS Export
//	@Description	Glue NFS Export 상세 정보를 보여줍니다.
//	@param			cluster_id 	path	string	true	"NFS Cluster Identifier"
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	NfsExportDetailed
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs/export/{cluster_id} [get]
func (c *Controller) NfsExportDetailed(ctx *gin.Context) {
	cluster_id := ctx.Param("cluster_id")
	dat, err := nfs.NfsExportDetailed(cluster_id)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.IndentedJSON(http.StatusOK, dat)
}
