package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/iscsi"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

func (c *Controller) IscsiOption(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}

func GlueUrl() (output string) {
	dat, err := iscsi.GlueUrl()
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	url := strings.Split(dat.ActiveName, ".")

	var stdout []byte
	cmd := exec.Command("sh", "-c", "cat /etc/hosts | grep '"+url[0]+"-mngt' | awk '{print $1}'")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		return
	}
	result := strings.Replace(string(stdout), "\n", "", -1)

	output = string("https://") + result + string(":8443/")
	return
}
func GetToken() (output string, err error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	user_json := model.UserInfo{
		Username: "ablecloud",
		Password: "Ablecloud1!",
	}
	user, err := json.Marshal(user_json)
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	url := GlueUrl() + "api/auth"
	var jsonStr = bytes.NewBuffer(user)
	request, err := http.NewRequest(http.MethodPost, url, jsonStr)
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("accept", "application/vnd.ceph.api.v1.0+json")

	result, err := client.Do(request)
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	var dat model.Token
	if err = json.Unmarshal(body, &dat); err != nil {
		utils.FancyHandleError(err)
		return
	}
	output = dat.Token
	return
}

// IscsiServiceCreate godoc
//
//	@Summary		Create of Iscsi Servcie Daemon
//	@Description	Iscsi 서비스 데몬을 생성합니다.
//	@param			hosts	formData	[]string	true	"Host Name" collectionFormat(multi)
//	@param			service_id	formData	string	true	"ISCSI Service Name"
//	@param			pool 	formData	string	true	"Pool Name"
//	@param			api_port 	formData	int	true	"ISCSI API Port" maximum(65535)
//	@param			api_user 	formData	string	true	"ISCSI API User"
//	@param			api_password 	formData	string	true	"ISCSI API Password"
//	@param			count 	formData	int	false	"Iscsi Service Daemon Count"
//	@Tags			Iscsi
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi [post]
func (c *Controller) IscsiServiceCreate(ctx *gin.Context) {
	service_id, _ := ctx.GetPostForm("service_id")
	hosts, _ := ctx.GetPostFormArray("hosts")
	pool, _ := ctx.GetPostForm("pool")
	api_port, _ := ctx.GetPostForm("api_port")
	api_user, _ := ctx.GetPostForm("api_user")
	api_password, _ := ctx.GetPostForm("api_password")
	service_count, _ := ctx.GetPostForm("count")
	port, _ := strconv.Atoi(api_port)
	count, _ := strconv.Atoi(service_count)

	var ip_data []string
	for i := 0; i < len(hosts); i++ {
		dat, err := iscsi.Ip(hosts[i])
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ip := strings.Split(dat, "\n")
		ip_data = append(ip_data, ip[0])
	}
	ip_address := strings.Join(ip_data, ",")
	if service_count == "" {
		value := model.IscsiServiceCreate{
			Service_Type: "iscsi",
			Service_Id:   service_id,
			Spec: model.Spec{
				Pool:          pool,
				Api_Port:      port,
				Api_User:      api_user,
				Api_Password:  api_password,
				TrustedIpList: ip_address},
			Placement: model.Placement{
				Hosts: hosts},
		}
		yaml_data, err := yaml.Marshal(&value)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		iscsi_yaml := "/etc/ceph/iscsi.yaml"
		err = os.WriteFile(iscsi_yaml, yaml_data, 0644)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		dat, err := iscsi.IscsiServiceCreate(iscsi_yaml)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			if err := os.Remove(iscsi_yaml); err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
			}
		}
		// Print the output
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, dat)
	} else {
		value := model.IscsiServiceCreateCount{
			Service_Type: "iscsi",
			Service_Id:   service_id,
			Spec: model.Spec{
				Pool:          pool,
				Api_Port:      port,
				Api_User:      api_user,
				Api_Password:  api_password,
				TrustedIpList: ip_address},
			Placement: model.PlacementCount{
				Count: count,
				Hosts: hosts},
		}
		yaml_data, err := yaml.Marshal(&value)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		iscsi_yaml := "/etc/ceph/iscsi.yaml"
		err = os.WriteFile(iscsi_yaml, yaml_data, 0644)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		dat, err := iscsi.IscsiServiceCreate(iscsi_yaml)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			if err := os.Remove(iscsi_yaml); err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
			}
		}
		// Print the output
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, dat)
	}
}

// IscsiServiceUpdate godoc
//
//	@Summary		Update of Iscsi Servcie Daemon
//	@Description	Iscsi 서비스 데몬을 수정합니다.
//	@param			hosts	formData	[]string	true	"Host Name" collectionFormat(multi)
//	@param			service_id	formData	string	true	"ISCSI Service Name"
//	@param			pool 	formData	string	true	"Pool Name"
//	@param			api_port 	formData	int	true	"ISCSI API Port" maximum(65535)
//	@param			api_user 	formData	string	true	"ISCSI API User"
//	@param			api_password 	formData	string	true	"ISCSI API Password"
//	@param			count 	formData	int	false	"Iscsi Service Daemon Count"
//	@Tags			Iscsi
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi [put]
func (c *Controller) IscsiServiceUpdate(ctx *gin.Context) {
	service_id, _ := ctx.GetPostForm("service_id")
	hosts, _ := ctx.GetPostFormArray("hosts")
	pool, _ := ctx.GetPostForm("pool")
	api_port, _ := ctx.GetPostForm("api_port")
	api_user, _ := ctx.GetPostForm("api_user")
	api_password, _ := ctx.GetPostForm("api_password")
	service_count, _ := ctx.GetPostForm("count")
	port, _ := strconv.Atoi(api_port)
	count, _ := strconv.Atoi(service_count)

	var ip_data []string
	for i := 0; i < len(hosts); i++ {
		dat, err := iscsi.Ip(hosts[i])
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ip := strings.Split(dat, "\n")
		ip_data = append(ip_data, ip[0])
	}
	ip_address := strings.Join(ip_data, ",")
	if service_count == "" {
		value := model.IscsiServiceCreate{
			Service_Type: "iscsi",
			Service_Id:   service_id,
			Spec: model.Spec{
				Pool:          pool,
				Api_Port:      port,
				Api_User:      api_user,
				Api_Password:  api_password,
				TrustedIpList: ip_address},
			Placement: model.Placement{
				Hosts: hosts},
		}
		yaml_data, err := yaml.Marshal(&value)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		iscsi_yaml := "/etc/ceph/iscsi.yaml"
		err = os.WriteFile(iscsi_yaml, yaml_data, 0644)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		dat, err := iscsi.IscsiServiceCreate(iscsi_yaml)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			if err := os.Remove(iscsi_yaml); err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
			}
		}
		// Print the output
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, dat)
	} else {
		value := model.IscsiServiceCreateCount{
			Service_Type: "iscsi",
			Service_Id:   service_id,
			Spec: model.Spec{
				Pool:          pool,
				Api_Port:      port,
				Api_User:      api_user,
				Api_Password:  api_password,
				TrustedIpList: ip_address},
			Placement: model.PlacementCount{
				Count: count,
				Hosts: hosts},
		}
		yaml_data, err := yaml.Marshal(&value)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		iscsi_yaml := "/etc/ceph/iscsi.yaml"
		err = os.WriteFile(iscsi_yaml, yaml_data, 0644)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		dat, err := iscsi.IscsiServiceCreate(iscsi_yaml)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			if err := os.Remove(iscsi_yaml); err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
			}
		}
		// Print the output
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, dat)
	}
}

// IscsiTargetList godoc
//
//	@Summary		Show List of Iscsi Target
//	@Description	Iscsi 타겟 리스트를 가져옵니다.
//	@Tags			IscsiTarget
//	@param			iqn_id	query	string	false	"Iscsi Target IQN Name"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	IscsiCommon
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/target [get]
func (c *Controller) IscsiTargetList(ctx *gin.Context) {
	var request *http.Request
	var responseBody []byte
	var err error
	var requestUrl string
	iqn_id := ctx.Request.URL.Query().Get("iqn_id")
	if iqn_id == "" {
		requestUrl = GlueUrl() + "api/iscsi/target"
	} else {
		requestUrl = GlueUrl() + "api/iscsi/target/" + iqn_id
	}

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
	var dat model.IscsiCommon
	if err = json.Unmarshal(responseBody, &dat); err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)

}

// IscsiTargetDelete godoc
//
//	@Summary		Delete of Iscsi Target
//	@Description	Iscsi 타겟을 삭제합니다.
//	@Tags			IscsiTarget
//	@param			iqn_id	query	string	true	"Iscsi Target IQN Name"
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/target [delete]
func (c *Controller) IscsiTargetDelete(ctx *gin.Context) {
	var request *http.Request
	var responseBody []byte
	var err error

	iqn_id := ctx.Request.URL.Query().Get("iqn_id")

	requestUrl := GlueUrl() + "api/iscsi/target/" + iqn_id

	if request, err = http.NewRequest(http.MethodDelete, requestUrl, nil); err != nil {
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
	var dat model.IscsiCommon
	if response.StatusCode == 204 {
		dat = "Success"
	} else {
		if err = json.Unmarshal(responseBody, &dat); err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}
	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// IscsiTargetCreate godoc
//
//	@Summary		Create of Iscsi Target
//	@Description	Iscsi 타겟을 생성합니다.
//	@Tags			IscsiTarget
//	@param			iqn_id	formData	string	true	"Iscsi Target IQN Name"
//	@param			hostname 	formData	[]string	true	"Gateway Host Name" collectionFormat(multi)
//	@param			ip_address 	formData	[]string	true	"Gateway Host IP Address" collectionFormat(multi)
//	@param			pool_name 	formData	[]string	false	"Glue Pool Name" collectionFormat(multi)
//	@param			image_name 	formData	[]string	false	"Glue Image Name" collectionFormat(multi)
//	@param			acl_enabled    	formData	boolean	true	"Iscsi Authentication" default(false)
//	@param			username 	formData	string	false	"Iscsi Auth User" 	minlength(8) maxlength(64)
//	@param			password 	formData	string	false	"Iscsi Auth Password"  minlength(12) maxlength(16)
//	@param			mutual_username 	formData	string	false	"Iscsi Auth Mutual User" minlength(8) maxlength(64)
//	@param			mutual_password 	formData	string	false	"Iscsi Auth Mutaul Password" minlength(12) maxlength(16)
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/target [post]
func (c *Controller) IscsiTargetCreate(ctx *gin.Context) {
	iqn_id, _ := ctx.GetPostForm("iqn_id")
	hostname, _ := ctx.GetPostFormArray("hostname")
	ip_address, _ := ctx.GetPostFormArray("ip_address")
	image_name, _ := ctx.GetPostFormArray("image_name")
	pool_name, _ := ctx.GetPostFormArray("pool_name")
	acl_enabled, _ := ctx.GetPostForm("acl_enabled")
	username, _ := ctx.GetPostForm("username")
	password, _ := ctx.GetPostForm("password")
	mutual_username, _ := ctx.GetPostForm("mutual_username")
	mutual_password, _ := ctx.GetPostForm("mutual_password")

	var portal model.Portals
	portals := make([]model.Portals, 0)
	for i := 0; i < len(hostname); i++ {
		portal = model.Portals{
			Host: hostname[i],
			Ip:   ip_address[i],
		}
		portals = append(portals, portal)
	}
	var disk model.Disks
	disks := make([]model.Disks, 0)
	for i := 0; i < len(image_name); i++ {
		disk = model.Disks{
			Pool:      pool_name[i],
			Image:     image_name[i],
			Backstore: "user:rbd",
			Lun:       i,
		}
		disks = append(disks, disk)
	}
	acl, errs := strconv.ParseBool(acl_enabled)
	if errs != nil {
		utils.FancyHandleError(errs)
		httputil.NewError(ctx, http.StatusInternalServerError, errs)
		return
	}
	value := model.IscsiTargetCreate{
		Target_Iqn:  iqn_id,
		Portals:     portals,
		Disks:       disks,
		Acl_Enabled: acl,
		Auth: model.Auth{
			User:            username,
			Password:        password,
			Mutual_User:     mutual_username,
			Mutual_Password: mutual_password,
		},
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
	requestUrl := GlueUrl() + "api/iscsi/target"
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
	var dat model.IscsiCommon
	if response.StatusCode == 201 {
		dat = "Success"
	} else {
		if err = json.Unmarshal(responseBody, &dat); err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}
	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// IscsiTargetUpdate godoc
//
//	@Summary		Update of Iscsi Target
//	@Description	Iscsi 타겟을 수정합니다.
//	@Tags			IscsiTarget
//	@param			iqn_id	formData	string	true	"Iscsi Target Old IQN Name"
//	@param			new_iqn_id	formData	string	true	"Iscsi Target New IQN Name"
//	@param			hostname 	formData	[]string	true	"Gateway Host Name" collectionFormat(multi)
//	@param			ip_address 	formData	[]string	true	"Gateway Host IP Address" collectionFormat(multi)
//	@param			pool_name 	formData	[]string	false	"Glue Pool Name" collectionFormat(multi)
//	@param			image_name 	formData	[]string	false	"Glue Image Name" collectionFormat(multi)
//	@param			acl_enabled    	formData	boolean	true	"Iscsi Authentication" default(false)
//	@param			username 	formData	string	false	"Iscsi Auth User" 	minlength(8) maxlength(64)
//	@param			password 	formData	string	false	"Iscsi Auth Password"  minlength(12) maxlength(16)
//	@param			mutual_username 	formData	string	false	"Iscsi Auth Mutual User" minlength(8) maxlength(64)
//	@param			mutual_password 	formData	string	false	"Iscsi Auth Mutaul Password" minlength(12) maxlength(16)
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	IscsiCommon
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/target [put]
func (c *Controller) IscsiTargetUpdate(ctx *gin.Context) {
	iqn_id, _ := ctx.GetPostForm("iqn_id")
	new_iqn_id, _ := ctx.GetPostForm("new_iqn_id")
	hostname, _ := ctx.GetPostFormArray("hostname")
	ip_address, _ := ctx.GetPostFormArray("ip_address")
	image_name, _ := ctx.GetPostFormArray("image_name")
	pool_name, _ := ctx.GetPostFormArray("pool_name")
	acl_enabled, _ := ctx.GetPostForm("acl_enabled")
	username, _ := ctx.GetPostForm("username")
	password, _ := ctx.GetPostForm("password")
	mutual_username, _ := ctx.GetPostForm("mutual_username")
	mutual_password, _ := ctx.GetPostForm("mutual_password")

	var portal model.Portals
	portals := make([]model.Portals, 0)
	for i := 0; i < len(hostname); i++ {
		portal = model.Portals{
			Host: hostname[i],
			Ip:   ip_address[i],
		}
		portals = append(portals, portal)
	}
	var disk model.Disks
	disks := make([]model.Disks, 0)
	for i := 0; i < len(image_name); i++ {
		disk = model.Disks{
			Pool:      pool_name[i],
			Image:     image_name[i],
			Backstore: "user:rbd",
			Lun:       i,
		}
		disks = append(disks, disk)
	}
	acl, errs := strconv.ParseBool(acl_enabled)
	if errs != nil {
		utils.FancyHandleError(errs)
		httputil.NewError(ctx, http.StatusInternalServerError, errs)
		return
	}
	value := model.IscsiTargetUpdate{
		New_Target_Iqn: new_iqn_id,
		Portals:        portals,
		Disks:          disks,
		Acl_Enabled:    acl,
		Auth: model.Auth{
			User:            username,
			Password:        password,
			Mutual_User:     mutual_username,
			Mutual_Password: mutual_password,
		},
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
	requestUrl := GlueUrl() + "api/iscsi/target/" + iqn_id
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
	var dat model.IscsiCommon
	if err = json.Unmarshal(responseBody, &dat); err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// IscsiGetDiscoveryAuth godoc
//
//	@Summary		Show of Iscsi Discovery Auth Details
//	@Description	Iscsi 계정 정보를 가져옵니다.
//	@Tags			Iscsi
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	Auth
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/discovery [get]
func (c *Controller) IscsiGetDiscoveryAuth(ctx *gin.Context) {
	var request *http.Request
	var responseBody []byte
	var err error

	requestUrl := GlueUrl() + "api/iscsi/discoveryauth"
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
	var dat model.IscsiCommon
	if err = json.Unmarshal(responseBody, &dat); err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// IscsiUpdateDiscoveryAuth godoc
//
//	@Summary		Update of Iscsi Discovery Auth Details
//	@Description	Iscsi 계정 정보를 수정합니다.
//	@Tags			Iscsi
//	@param			user	formData	string	false	"Iscsi Discovery Authorization Username" minlength(8) maxlength(64)
//	@param			password	formData	string	false	"Iscsi Discovery Authorization Password" minlength(12) maxlength(16)
//	@param			mutual_user	formData	string	false	"Iscsi Discovery Authorization Mutual Username" minlength(8) maxlength(64)
//	@param			mutual_password	formData	string	false	"Iscsi Discovery Authorization Mutual Password" minlength(12) maxlength(16)
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	Auth
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/discovery [put]
func (c *Controller) IscsiUpdateDiscoveryAuth(ctx *gin.Context) {
	user, _ := ctx.GetPostForm("user")
	password, _ := ctx.GetPostForm("password")
	mutual_user, _ := ctx.GetPostForm("mutual_user")
	mutual_password, _ := ctx.GetPostForm("mutual_password")

	var request *http.Request
	var responseBody []byte
	var err error

	requestUrl := GlueUrl() + "api/iscsi/discoveryauth?user=%20&password=%20&mutual_user=%20&mutual_password=%20"

	value := model.Auth{
		User:            user,
		Password:        password,
		Mutual_User:     mutual_user,
		Mutual_Password: mutual_password,
	}
	json_data, _ := json.Marshal(value)
	var jsonStr = bytes.NewBuffer(json_data)
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
	request.Header.Add("Content-Type", "application/json")
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
	var dat model.IscsiCommon
	if err = json.Unmarshal(responseBody, &dat); err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}
