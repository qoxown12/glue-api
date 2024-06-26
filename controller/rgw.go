package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/glue"
	"Glue-API/utils/rgw"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (c *Controller) RgwOption(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}

// RgwDaemon godoc
//
//	@Summary		Show List of RADOS Gateway Daemon
//	@Description	RADOS Gateway Daemon 정보를 가져옵니다.
//	@Tags			RGW
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	RgwDaemon
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/rgw [get]
func (c *Controller) RgwDaemon(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	var request *http.Request
	var responseBody []byte
	var err error

	requestUrl := GlueUrl() + "api/rgw/daemon"

	if request, err = http.NewRequest(http.MethodGet, requestUrl, nil); err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
	}
	token, err := GetToken()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	request.Header.Add("accept", "application/vnd.ceph.api.v1.0+json")
	request.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	response, err := client.Do(request)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	defer response.Body.Close()

	responseBody, err = io.ReadAll(response.Body)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	var dat model.RgwDaemon
	if err = json.Unmarshal(responseBody, &dat); err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Print the output
	ctx.IndentedJSON(http.StatusOK, dat)
}

// RgwServiceCreate godoc
//
//	@Summary		Create of RADOS Gateway Service
//	@Description	RADOS Gateway Service를 생성합니다.
//	@Tags			RGW
//	@param			service_name     formData   string	true    "RGW Service Name"
//	@param			realm_name     formData   string	false    "RGW Realm Name"
//	@param			zonegroup_name     formData   string	false    "RGW Zone Group Name"
//	@param			zone_name     formData   string	false    "RGW Zone Name"
//	@param			port     formData   int	false    "Service Port(default: 80)"
//	@param			hosts     formData   []string	true    "Service Placement Host Name" collectionFormat(multi)
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string ""
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/rgw [post]
func (c *Controller) RgwServiceCreate(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	service_name, _ := ctx.GetPostForm("service_name")
	realm_name, _ := ctx.GetPostForm("realm_name")
	zonegroup_name, _ := ctx.GetPostForm("zonegroup_name")
	zone_name, _ := ctx.GetPostForm("zone_name")
	port, _ := ctx.GetPostForm("port")
	hosts, _ := ctx.GetPostFormArray("hosts")

	hosts_str := strings.Join(hosts, ",")
	if port == "" {
		port = "80"
	}
	dat, err := rgw.RgwServiceCreateandUpdate(service_name, realm_name, zonegroup_name, zone_name, hosts_str, port)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		pool, err := glue.PoolReplicatedList("rgw")
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		for i := 0; i < len(pool); i++ {
			_, err := glue.PoolReplicatedSize(pool[i])
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
		}
	}
	ctx.IndentedJSON(http.StatusOK, dat)
}

// RgwServiceUpdate godoc
//
//	@Summary		Update of RADOS Gateway Service
//	@Description	RADOS Gateway Service를 수정합니다.
//	@Tags			RGW
//	@param			service_id     formData   string	true    "RGW Service Name"
//	@param			realm_name     formData   string	false    "RGW Realm Name"
//	@param			zonegroup_name     formData   string	false    "RGW Zone Group Name"
//	@param			zone_name     formData   string	false    "RGW Zone Name"
//	@param			port     formData   int	false    "Service Port(default: 80)"
//	@param			hosts     formData   []string	true    "Service Placement Hosts" collectionFormat(multi)
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string ""
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/rgw [put]
func (c *Controller) RgwServiceUpdate(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	service_id, _ := ctx.GetPostForm("service_id")
	realm_name, _ := ctx.GetPostForm("realm_name")
	zonegroup_name, _ := ctx.GetPostForm("zonegroup_name")
	zone_name, _ := ctx.GetPostForm("zone_name")
	port, _ := ctx.GetPostForm("port")
	hosts, _ := ctx.GetPostFormArray("hosts")

	hosts_str := strings.Join(hosts, ",")

	dat, err := rgw.RgwServiceCreateandUpdate(service_id, realm_name, zonegroup_name, zone_name, hosts_str, port)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, dat)
}

// RgwUserList godoc
//
//	@Summary		List and Info of RADOS Gateway Users
//	@Description	RADOS Gateway User의 리스트 및 정보를 보여줍니다.
//	@param			username     query   string	false    "RGW User Name"
//	@Tags			RGW-User
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	RgwUserInfo
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/rgw/user [get]
func (c *Controller) RgwUserList(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	username := ctx.Request.URL.Query().Get("username")

	if username != "" {
		dat, err := rgw.RgwUserInfo(username)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(http.StatusOK, dat)
	} else {
		var userInfo []model.RgwUserInfoAndStat
		list_dat, err := rgw.RgwUserList()
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		for i := 0; i < len(list_dat); i++ {
			info_dat, err := rgw.RgwUserInfo(list_dat[i])
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			stat_dat, err := rgw.RgwUserStat(list_dat[i])
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			value := model.RgwUserInfoAndStat{
				UserID:              info_dat.UserID,
				DisplayName:         info_dat.DisplayName,
				Email:               info_dat.Email,
				Suspended:           info_dat.Suspended,
				MaxBuckets:          info_dat.MaxBuckets,
				Subusers:            info_dat.Subusers,
				Keys:                info_dat.Keys,
				SwiftKeys:           info_dat.SwiftKeys,
				Caps:                info_dat.Caps,
				OpMask:              info_dat.OpMask,
				System:              info_dat.System,
				DefaultPlacement:    info_dat.DefaultPlacement,
				DefaultStorageClass: info_dat.DefaultStorageClass,
				PlacementTags:       info_dat.PlacementTags,
				BucketQuota:         info_dat.BucketQuota,
				UserQuota:           info_dat.UserQuota,
				TempURLKeys:         info_dat.TempURLKeys,
				Type:                info_dat.Type,
				MfaIds:              info_dat.MfaIds,
				Stats: model.RgwUserStat{
					Stats: stat_dat.Stats,
				},
			}
			userInfo = append(userInfo, value)
		}
		ctx.IndentedJSON(http.StatusOK, userInfo)
	}
}

// RgwUserCreate godoc
//
//	@Summary		Create of RADOS Gateway User
//	@Description	RADOS Gateway User를 생성합니다.
//	@Tags			RGW-User
//	@param			username     formData   string	true    "RGW User ID Name"
//	@param			display_name     formData   string	true    "RGW  User Display Name"
//	@param			email     formData   string	false    "RGW User Email"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/rgw/user [post]
func (c *Controller) RgwUserCreate(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	username, _ := ctx.GetPostForm("username")
	display_name, _ := ctx.GetPostForm("display_name")
	email, _ := ctx.GetPostForm("email")

	dat, err := rgw.RgwUserCreate(username, display_name, email)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, dat)
}

// RgwUserDelete godoc
//
//	@Summary		Delete of RADOS Gateway User
//	@Description	RADOS Gateway User를 삭제합니다.
//	@Tags			RGW-User
//	@param			username     query   string	true    "RGW User ID Name"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/rgw/user [delete]
func (c *Controller) RgwUserDelete(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	username := ctx.Request.URL.Query().Get("username")

	dat, err := rgw.RgwUserDelete(username)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, dat)
}

// RgwUserUpdate godoc
//
//	@Summary		Update of RADOS Gateway User
//	@Description	RADOS Gateway User를 수정합니다.
//	@Tags			RGW-User
//	@param			username     formData   string	true    "RGW User ID Name"
//	@param			display_name     formData   string	false    "RGW User Display Name"
//	@param			email     formData   string	false    "RGW User Email "
//	@param			key_type     formData   string	false    "RGW User S3" Enums(s3)
//	@param			access_key     formData   string	false    "RGW User S3 Access Key"
//	@param			secret_key     formData   string	false    "RGW User S3 Secret Key"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/rgw/user [put]
func (c *Controller) RgwUserUpdate(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	username, _ := ctx.GetPostForm("username")
	display_name, _ := ctx.GetPostForm("display_name")
	email, _ := ctx.GetPostForm("email")
	key_type, _ := ctx.GetPostForm("key_type")
	access_key, _ := ctx.GetPostForm("access_key")
	secret_key, _ := ctx.GetPostForm("secret_key")

	dat, err := rgw.RgwUserUpdate(username, display_name, email, key_type, access_key, secret_key)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, dat)
}

// RgwQuota godoc
//
//	@Summary		Setting of RADOS Gateway Quota
//	@Description	RADOS Gateway Quota를 설정 및 활성화합니다.
//	@Tags			RGW
//	@param			username     formData   string	true    "RGW User ID Name"
//	@param			scope     formData   string	true    "RGW Quota Target" Enums(user, bucket)
//	@param			max_objects     formData   int	true    "RGW Quota Max Objects"
//	@param			max_size     formData   string	true    "RGW Quota Max Size(B/K/M/G/T)"
//	@param			state     formData   string	true    "RGW Quota Whether Activated" default(enable) Enums(enable, disable)
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/rgw/quota [post]
func (c *Controller) RgwQuota(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	username, _ := ctx.GetPostForm("username")
	scope, _ := ctx.GetPostForm("scope")
	max_objects, _ := ctx.GetPostForm("max_objects")
	max_size, _ := ctx.GetPostForm("max_size")
	state, _ := ctx.GetPostForm("state")

	dat, err := rgw.RgwQuota(username, scope, max_objects, max_size, state)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, dat)
}

// RgwBucketList godoc
//
//		@Summary		Show List of RADOS Gateway Bucket
//		@Description	RADOS Gateway 버킷을 리스트를 보여줍니다.
//		@Tags			RGW-Bucket
//		@param			bucket_name	query	string	false	"RGW Bucket Name"
//	 	@param          detail 		query 	string 	true 	"RGW Bucket List Detail" Enums(true, false) default(false)
//		@Accept			x-www-form-urlencoded
//		@Produce		json
//		@Success		200	{object}	RgwCommon
//		@Failure		400	{object}	httputil.HTTP400BadRequest
//		@Failure		404	{object}	httputil.HTTP404NotFound
//		@Failure		500	{object}	httputil.HTTP500InternalServerError
//		@Router			/api/v1/rgw/bucket [get]
func (c *Controller) RgwBucketList(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	bucket_name := ctx.Request.URL.Query().Get("bucket_name")
	detail := ctx.Request.URL.Query().Get("detail")
	if detail == "true" {
		dat, err := rgw.RgwBucketDetail(bucket_name)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(http.StatusOK, dat)
	} else {
		dat, err := rgw.RgwBucketList()
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(http.StatusOK, dat)
	}
}

// RgwBucketCreate godoc
//
//	@Summary		Create of RADOS Gateway Bucket
//	@Description	RADOS Gateway 버킷을 생성합니다.
//	@Tags			RGW-Bucket
//	@param			bucket_name	formData	string	true	"RGW Bucket Name"
//	@param			username	formData	string	true	"RGW User Name"
//	@param			lock_enabled 	formData	boolean	true	"RGW Bucket Lock Enabled" Enums(true, false) default(false)
//	@param			lock_mode 	formData	string	false	"RGW Bucket Lock Mode" Enums(compliance, governance)
//	@param			lock_retention_period_days 	formData	int	false	"RGW Bucket Lock Period"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	RgwCommon
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/rgw/bucket [post]
func (c *Controller) RgwBucketCreate(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	bucket_name, _ := ctx.GetPostForm("bucket_name")
	username, _ := ctx.GetPostForm("username")
	lock_enabled, _ := ctx.GetPostForm("lock_enabled")
	lock_mode, _ := ctx.GetPostForm("lock_mode")
	lock_retention_period_days, _ := ctx.GetPostForm("lock_retention_period_days")

	value := model.RgwBucketCreate{
		Bucket:                     bucket_name,
		Uid:                        username,
		Lock_enabled:               lock_enabled,
		Lock_mode:                  lock_mode,
		Lock_retention_period_days: lock_retention_period_days,
	}
	json_data, errs := json.Marshal(value)
	if errs != nil {
		utils.FancyHandleError(errs)
		httputil.NewError(ctx, http.StatusInternalServerError, errs)
		return
	}
	var jsonStr = bytes.NewBuffer(json_data)
	var request *http.Request
	var responseBody []byte
	var err error
	requestUrl := GlueUrl() + "api/rgw/bucket"
	if request, err = http.NewRequest(http.MethodPost, requestUrl, jsonStr); err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
	}
	token, err := GetToken()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	request.Header.Add("accept", "application/vnd.ceph.api.v1.0+json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+token)
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	response, err := client.Do(request)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	defer response.Body.Close()

	responseBody, err = io.ReadAll(response.Body)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	var dat model.RGwCommon
	if err = json.Unmarshal(responseBody, &dat); err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.IndentedJSON(http.StatusOK, "Success")
}

// RgwBucketUpdate godoc
//
//	@Summary		Update of RADOS Gateway Bucket
//	@Description	RADOS Gateway 버킷을 수정합니다.
//	@Tags			RGW-Bucket
//	@param			bucket_name	formData	string	true	"RGW Bucket Name"
//	@param			bucket_id	formData	string	true	"RGW Bucket ID"
//	@param			username	formData	string	true	"RGW User Name"
//	@param			versioning 	formData	string	false	"RGW Bucket Lock Enabled" Enums(Enabled, Suspended)
//	@param			lock_mode 	formData	string	false	"RGW Bucket Lock Mode(Required value if the lock box is checked)" Enums(compliance, governance)
//	@param			lock_retention_period_days 	formData	int	false	"RGW Bucket Lock Period(Required value if the lock box is checked)"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	RgwCommon
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/rgw/bucket [put]
func (c *Controller) RgwBucketUpdate(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	bucket_name, _ := ctx.GetPostForm("bucket_name")
	bucket_id, _ := ctx.GetPostForm("bucket_id")
	username, _ := ctx.GetPostForm("username")
	versioning, _ := ctx.GetPostForm("versioning")
	lock_mode, _ := ctx.GetPostForm("lock_mode")
	lock_retention_period_days, _ := ctx.GetPostForm("lock_retention_period_days")

	value := model.RgwBucketUpdate{
		Bucket_id:                  bucket_id,
		Uid:                        username,
		Versioning_state:           versioning,
		Lock_mode:                  lock_mode,
		Lock_retention_period_days: lock_retention_period_days,
	}
	json_data, errs := json.Marshal(value)
	if errs != nil {
		utils.FancyHandleError(errs)
		httputil.NewError(ctx, http.StatusInternalServerError, errs)
		return
	}
	var jsonStr = bytes.NewBuffer(json_data)
	var request *http.Request
	var responseBody []byte
	var err error
	requestUrl := GlueUrl() + "api/rgw/bucket/" + bucket_name
	if request, err = http.NewRequest(http.MethodPut, requestUrl, jsonStr); err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
	}
	token, err := GetToken()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	request.Header.Add("accept", "application/vnd.ceph.api.v1.0+json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+token)
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	response, err := client.Do(request)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	defer response.Body.Close()

	responseBody, err = io.ReadAll(response.Body)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	var dat model.RGwCommon
	if err = json.Unmarshal(responseBody, &dat); err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.IndentedJSON(http.StatusOK, "Success")
}

// RgwBucketDelete godoc
//
//	@Summary		Delete of RADOS Gateway Bucket
//	@Description	RADOS Gateway 버킷을 삭제합니다.
//	@Tags			RGW-Bucket
//	@param			bucket_name	query	string	true	"RGW Bucket Name"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	RgwCommon
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/rgw/bucket [delete]
func (c *Controller) RgwBucketDelete(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	bucket_name := ctx.Request.URL.Query().Get("bucket_name")
	dat, err := rgw.RgwBucketDelete(bucket_name)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, dat)
}
