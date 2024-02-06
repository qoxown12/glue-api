package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/nfs"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// NfsClusterList godoc
//
//	@Summary		Show List of Info of Glue NFS Cluster
//	@Description	Glue NFS Cluster의 리스트 및 상세정보를 보여줍니다.
//	@param			cluster_id 	query	string	false	"NFS Cluster Identifier"
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	NfsClusterList
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs [get]
func (c *Controller) NfsClusterList(ctx *gin.Context) {
	cluster_id := ctx.Request.URL.Query().Get("cluster_id")
	dat, err := nfs.NfsClusterList(cluster_id)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
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
	ctx.Header("Access-Control-Allow-Origin", "*")
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
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}
func (c *Controller) NfsClusterDeleteOptions(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}

// NfsExportCreate godoc
//
//	@Summary		Create of Glue NFS Export
//	@Description	Glue NFS Export를 생성합니다.
//	@param			cluster_id 	path	string	true	"NFS Cluster Identifier"
//	@param			access_type formData   string	true    "NFS Access Type" Enums(RW, RO, NONE) default(RW)
//	@param			fs_name     formData   string	true    "FS Name"
//	@param			storage_name formData   string	true    "NFS Storage Name" default(CEPH)
//	@param			path         formData    string true    "Glue FS Path"
//	@param			pseudo     formData   string	true    "NFS Export Path"
//	@param			squash     formData   string	true    "Squash"	Enums(no_root_squash, root_id_squash, all_squash, root_squash) default(no_root_squash)
//	@param			transports     formData   []string	false    "Transports" collectionFormat(multi) default(TCP)
//	@param			security_label     formData   boolean	true    "Security Label" default(false)
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs/export/{cluster_id} [post]
func (c *Controller) NfsExportCreate(ctx *gin.Context) {

	access_type, _ := ctx.GetPostForm("access_type")
	fs_name, _ := ctx.GetPostForm("fs_name")
	storage_name, _ := ctx.GetPostForm("storage_name")
	path, _ := ctx.GetPostForm("path")
	pseudo, _ := ctx.GetPostForm("pseudo")
	squash, _ := ctx.GetPostForm("squash")
	transports, _ := ctx.GetPostFormArray("transports")
	security_label := ctx.GetBool("security_label")

	var protocols = []int{4}
	value := model.NfsExportCreate{
		AccessType: access_type,
		Fsal: model.Fsal{
			Name:   storage_name,
			FsName: fs_name},
		Protocols:     protocols,
		Path:          path,
		Pseudo:        pseudo,
		Squash:        squash,
		SecurityLabel: security_label,
		Transports:    transports}
	//value := []byte("EXPORT {\n \tFSAL {\n\t\tname = \"" + storage_name + "\";\n\t\tfilesystem = \"" + fs_name + "\";\n\t}\n\texport_id = 1;\n\tpath = \"" + path + "\";\n\tpseudo = \"" + pseudo + "\";\n\taccess_type = \"" + access_type + "\";\n\tsquash = \"" + squash + "\";\n\tprotocols = 4;\n\ttrasnports = \"" + transports + "\";\n}")
	json_data, err := json.MarshalIndent(value, "", " ")
	nfs_export_create_conf := "/root/nfs_export_create.conf"
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
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, dat)
	}
	return
}

// NfsExportUpdate godoc
//
//	@Summary		Update of Glue NFS Export
//	@Description	Glue NFS Export를 수정합니다.
//	@param			cluster_id 	path	string	true	"NFS Cluster Identifier"
//	@param			access_type formData   string	true    "NFS Access Type" Enums(RW, RO, NONE) default(RW)
//	@param			fs_name     formData   string	true    "FS Name"
//	@param			storage_name formData   string	true    "NFS Storage Name" default(CEPH)
//	@param			path         formData    string true    "Glue FS Path"
//	@param			pseudo     formData   string	true    "NFS Export Path"
//	@param			squash     formData   string	true    "Squash"	Enums(no_root_squash, root_id_squash, all_squash, root_squash) default(no_root_squash)
//	@param			transports     formData   []string	false    "Transports" collectionFormat(multi) default(TCP)
//	@param			security_label     formData   boolean	true    "Security Label" default(false)
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs/export/{cluster_id} [put]
func (c *Controller) NfsExportUpdate(ctx *gin.Context) {
	cluster_id := ctx.Param("cluster_id")
	access_type, _ := ctx.GetPostForm("access_type")
	fs_name, _ := ctx.GetPostForm("fs_name")
	storage_name, _ := ctx.GetPostForm("storage_name")
	path, _ := ctx.GetPostForm("path")
	pseudo, _ := ctx.GetPostForm("pseudo")
	squash, _ := ctx.GetPostForm("squash")
	transports, _ := ctx.GetPostFormArray("transports")
	security_label := ctx.GetBool("security_label")

	var export_id int
	dat, err := nfs.NfsExportDetailed(cluster_id)

	for i := 0; i < len(dat); i++ {
		if dat[i].Pseudo == pseudo {
			export_id = dat[i].ExportID
		}
	}
	var protocols = []int{4}
	value := model.NfsExportUpdate{
		AccessType: access_type,
		Fsal: model.Fsal{
			Name:   storage_name,
			FsName: fs_name},
		Protocols:     protocols,
		Path:          path,
		Pseudo:        pseudo,
		Squash:        squash,
		SecurityLabel: security_label,
		ExportID:      export_id,
		Transports:    transports}
	json_data, err := json.MarshalIndent(value, "", " ")
	nfs_export_update_conf := "/root/nfs_export_update.conf"
	err = os.WriteFile(nfs_export_update_conf, json_data, 0644)

	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
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
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, dat)
	}
	return
}
func (c *Controller) NfsExportUpdateOptions(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
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
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.IndentedJSON(http.StatusOK, dat)
		}
	}
	// Print the output
}
func (c *Controller) NfsExportDeleteOptions(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}

// NfsExportDetailed godoc
//
//	@Summary		Show Detail of Glue NFS Export
//	@Description	Glue NFS Export 상세 정보를 보여줍니다.
//	@param			cluster_id 	query	string	false	"NFS Cluster Identifier"
//	@Tags			NFS
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	NfsExportDetailed
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nfs/export [get]
func (c *Controller) NfsExportDetailed(ctx *gin.Context) {
	cluster_id := ctx.Request.URL.Query().Get("cluster_id")
	if cluster_id != "" {
		dat, err := nfs.NfsExportDetailed(cluster_id)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(http.StatusOK, dat)
	} else {
		var output model.NfsExportDetailed
		dat2, err := nfs.NfsClusterLs()
		for i := 0; i < len(dat2); i++ {
			dat3, err := nfs.NfsExportDetailed(dat2[i])
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			output = append(dat3, output...)
		}
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(http.StatusOK, output)
	}
	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
}
