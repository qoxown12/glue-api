package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/iscsi"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

// IscsiTargetName 불러 오기
func IscsiName() (ceph_container_name string, hostname string) {
	dat, err := iscsi.IscsiService()
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	var data []string
	for i := 0; i < len(dat); i++ {
		arr_data := dat[i].Placement.Hosts
		data = append(arr_data, data...)
	}

	for i := 0; i < len(data); i++ {
		if data[i] == "gwvm" {
			gwvm_data, err := iscsi.IscsiTargetName("gwvm")
			if err != nil {
				utils.FancyHandleError(err)
				return
			}
			ceph_container_name = gwvm_data
			hostname = "gwvm"
			return
		}
	}
	scvm_data, err := iscsi.IscsiTargetName("scvm")
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	ceph_container_name = scvm_data
	hostname = "scvm"
	return
}

// IscsiServiceCreate godoc
//
//	@Summary		Create of Iscsi Servcie Daemon
//	@Description	Iscsi 서비스 데몬을 생성합니다.
//	@param			hosts 	formData	[]string	true	"Host Name" collectionFormat(multi)
//	@param			service_id	formData	string	true	"ISCSI Service Name"
//	@param			pool 	formData	string	true	"Pool Name"
//	@param			api_port 	formData	int	true	"ISCSI API Port"
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
	if service_count == "" {
		value := model.IscsiServiceCreate{
			Service_Type: "iscsi",
			Service_Id:   service_id,
			Spec: model.Spec{
				Pool:         pool,
				Api_Port:     port,
				Api_User:     api_user,
				Api_Password: api_password},
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
				Pool:         pool,
				Api_Port:     port,
				Api_User:     api_user,
				Api_Password: api_password},
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
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	IscsiTargetList
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/target [get]
func (c *Controller) IscsiTargetList(ctx *gin.Context) {
	dat, hostname := IscsiName()
	dat = strings.Replace(dat, "\n", "", -1)
	data, err := iscsi.IscsiTargetList(dat, hostname)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, data)

}

// IscsiTargetCreate godoc
//
//	@Summary		Create of Iscsi Target
//	@Description	Iscsi 타겟을 생성합니다.
//	@param			iqn_id 	path	string	true	"Iscsi IQN ID" example("iqn.{yyyy-mm}.{naming-authority}:{unique-name}")
//	@param			hostname 	formData	[]string	true	"Gateway Host Name" collectionFormat(multi)
//	@param			ip_address 	formData	[]string	true	"Gateway Host IP Address" collectionFormat(multi)
//	@param			pool_name 	formData	string	true	"Glue Pool Name"
//	@param			disk_name 	formData	string	true	"Iscsi Disk Name"
//	@param			size 	formData	int	false	"Iscsi Disk Image Size(Default GB)"
//	@param			auth    	formData	boolean	true	"Iscsi Authentication" default(false)
//	@param			username 	formData	string	false	"Iscsi Auth User" 	minlength(8) maxlength(64)
//	@param			password 	formData	string	false	"Iscsi Auth Password"  minlength(12) maxlength(16)
//	@param			mutual_username 	formData	string	false	"Iscsi Auth Mutual User" minlength(8) maxlength(64)
//	@param			mutual_password 	formData	string	false	"Iscsi Auth Mutaul Password" minlength(12) maxlength(16)
//	@Tags			IscsiTarget
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/target/{iqn_id} [post]
func (c *Controller) IscsiTargetCreate(ctx *gin.Context) {
	iqn_id := ctx.Param("iqn_id")
	portal, _ := ctx.GetPostFormArray("hostname")
	ip_address, _ := ctx.GetPostFormArray("ip_address")
	disk_name, _ := ctx.GetPostForm("disk_name")
	pool_name, _ := ctx.GetPostForm("pool_name")
	size, _ := ctx.GetPostForm("size")
	auth, _ := ctx.GetPostForm("auth")
	username, _ := ctx.GetPostForm("username")
	password, _ := ctx.GetPostForm("password")
	mutual_username, _ := ctx.GetPostForm("mutual_username")
	mutual_password, _ := ctx.GetPostForm("mutual_password")
	// cmd := exec.Command("sh", "-c", "cat /etc/hosts | sort | grep -w -m 1 'gwvm' | awk '{print $1}'")
	// ip_address, err := cmd.CombinedOutput()
	ceph_container_name, hostname := IscsiName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)

	data, err := iscsi.IscsiTargetCreate(ceph_container_name, hostname, iqn_id, pool_name, disk_name, size, auth, username, password, mutual_username, mutual_password)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	if data == "Success" {
		for i := 0; i < len(portal); i++ {
			gateway, err := iscsi.IscsiGatewayAttach(ceph_container_name, hostname, iqn_id, portal[i], ip_address[i])
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			if i == len(portal)-1 {
				ctx.Header("Access-Control-Allow-Origin", "*")
				ctx.IndentedJSON(http.StatusOK, gateway)
			}
		}
	}

}

// IscsiTargetDelete godoc
//
//	@Summary		Delete of Iscsi Target
//	@Description	Iscsi 타겟을 삭제합니다.
//	@param			iqn_id 	path	string	true	"Iscsi IQN ID"
//	@param			pool_name 	query	string	true	"Glue Pool Name"
//	@param			disk_name 	query	string	true	"Iscsi Disk Name"
//	@param			image 	query	string	true	"Whether to Delete RBD Image" default(false) Enums(true,false)
//	@Tags			IscsiTarget
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/target/{iqn_id} [delete]
func (c *Controller) IscsiTargetDelete(ctx *gin.Context) {
	iqn_id := ctx.Param("iqn_id")
	pool_name := ctx.Request.URL.Query().Get("pool_name")
	disk_name := ctx.Request.URL.Query().Get("disk_name")
	image := ctx.Request.URL.Query().Get("image")
	ceph_container_name, hostname := IscsiName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)
	data, err := iscsi.IscsiTargetDelete(ceph_container_name, hostname, pool_name, disk_name, iqn_id, image)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, data)

}
func (c *Controller) IscsiTargetDeleteOptions(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}

// IscsiDiskList godoc
//
//	@Summary		Show List of Iscsi Disk
//	@Description	Iscsi 디스크 리스트를 보여줍니다.
//	@Tags			IscsiDisk
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	IscsiDiskList
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/disk [get]
func (c *Controller) IscsiDiskList(ctx *gin.Context) {
	ceph_container_name, hostname := IscsiName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)
	data, err := iscsi.IscsiDiskList(ceph_container_name, hostname)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, data)

}

// IscsiDiskCreate godoc
//
//	@Summary		Create Or Attach of Iscsi Disk
//	@Description	Iscsi 디스크를 생성 또는 부착합니다.
//	@param			pool_name	formData	string	true	"Iscsi Disk Pool Name"
//	@param			disk_name 	formData	string	true	"Iscsi Disk Name"
//	@param			size	formData	int	false	"Iscsi Disk Image Size(Default GB)"
//	@param			iqn_id  formData	string	false	"Iscsi IQN ID"
//	@Tags			IscsiDisk
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/disk [post]
func (c *Controller) IscsiDiskCreate(ctx *gin.Context) {
	disk_name, _ := ctx.GetPostForm("disk_name")
	pool_name, _ := ctx.GetPostForm("pool_name")
	size, _ := ctx.GetPostForm("size")
	iqn_id, _ := ctx.GetPostForm("iqn_id")
	ceph_container_name, hostname := IscsiName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)
	data, err := iscsi.IscsiDiskCreate(ceph_container_name, hostname, pool_name, disk_name, size, iqn_id)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, data)

}

// IscsiDiskDelete godoc
//
//	@Summary		Delete of Iscsi Disk
//	@Description	Iscsi 디스크를 삭제합니다.
//	@param			pool_name 	query	string	true	"Iscsi Disk Pool Name"
//	@param			disk_name 	query	string	true	"Iscsi Disk Name"
//	@param			iqn_id  query	string	false	"Iscsi IQN ID"
//	@param			image 	query	string	true	"Whether to Delete RBD Image" default(false) Enums(true,false)
//	@Tags			IscsiDisk
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest``
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/disk [delete]
func (c *Controller) IscsiDiskDelete(ctx *gin.Context) {
	disk_name := ctx.Request.URL.Query().Get("disk_name")
	pool_name := ctx.Request.URL.Query().Get("pool_name")
	image := ctx.Request.URL.Query().Get("image")
	iqn_id := ctx.Request.URL.Query().Get("iqn_id")
	ceph_container_name, hostname := IscsiName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)

	data, err := iscsi.IscsiDiskDelete(ceph_container_name, hostname, pool_name, disk_name, image, iqn_id)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, data)

}

// IscsiDiskResize godoc
//
//	@Summary		Change Size of Iscsi Disk
//	@Description	Iscsi 디스크 용량을 변경합니다.
//	@param			disk_name 	query	string	true	"Iscsi Disk Name"
//	@param			new_size 	query	int	true	"Iscsi Disk New Size(default GB)"
//	@Tags			IscsiDisk
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/disk [put]
func (c *Controller) IscsiDiskResize(ctx *gin.Context) {
	disk_name := ctx.Request.URL.Query().Get("disk_name")
	new_size := ctx.Request.URL.Query().Get("new_size")
	ceph_container_name, hostname := IscsiName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)

	data, err := iscsi.IscsiDiskResize(ceph_container_name, hostname, disk_name, new_size)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, data)

}
func (c *Controller) IscsiDiskOptions(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}

// IscsiDiscoveryCreate godoc
//
//	@Summary		Create User of Iscsi Discovery Auth
//	@Description	Iscsi Discovery Auth 계정을 생성합니다.
//	@param			username 	formData	string	true	"Iscsi Discovery Authentication User"  minlength(8) maxlength(64)
//	@param			password 	formData	string	true	"Iscsi Discovery Authentication Password" minlength(12) maxlength(16)
//	@param			mutual_username 	formData	string	false	"Iscsi Discovery Authentication Mutual User" minlength(8) maxlength(64)
//	@param			mutual_password 	formData	string	false	"Iscsi Discovery Authentication Mutual Password" minlength(12) maxlength(16)
//	@Tags			IscsiDiscovery
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/discovery [post]
func (c *Controller) IscsiDiscoveryCreate(ctx *gin.Context) {
	username, _ := ctx.GetPostForm("username")
	password, _ := ctx.GetPostForm("password")
	mutual_username, _ := ctx.GetPostForm("mutual_username")
	mutual_password, _ := ctx.GetPostForm("mutual_password")

	ceph_container_name, hostname := IscsiName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)
	dat, err := iscsi.IscsiDiscoveryCreate(ceph_container_name, hostname, username, password, mutual_username, mutual_password)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// IscsiDiscoveryInfo godoc
//
//	@Summary		Info User of Iscsi Discovery Auth
//	@Description	Iscsi Discovery Auth 계정 정보를 보여줍니다.
//	@Tags			IscsiDiscovery
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	IscsiDiscoveryInfo
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/discovery [get]
func (c *Controller) IscsiDiscoveryInfo(ctx *gin.Context) {

	ceph_container_name, hostname := IscsiName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)
	dat, err := iscsi.IscsiDiscoveryInfo(ceph_container_name, hostname)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// IscsiDiscoveryReset godoc
//
//	@Summary		Reset User of Iscsi Discovery Auth
//	@Description	Iscsi Discovery Auth 계정 정보를 초기화 합니다.
//	@Tags			IscsiDiscovery
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/iscsi/discovery [delete]
func (c *Controller) IscsiDiscoveryReset(ctx *gin.Context) {
	ceph_container_name, hostname := IscsiName()
	ceph_container_name = strings.Replace(ceph_container_name, "\n", "", -1)
	dat, err := iscsi.IscsiDiscoveryReset(ceph_container_name, hostname)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}
func (c *Controller) IscsiDiscoveryOptions(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}
