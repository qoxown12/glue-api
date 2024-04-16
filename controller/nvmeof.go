package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/glue"
	"Glue-API/utils/nvmeof"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (c *Controller) NvmeOption(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}
func NvmeOfServerIPandPort() (ip string, port string, err error) {
	gateway_name, err := nvmeof.NvmeOfGatewayName()
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	gateway_ip, err := nvmeof.HostIp(gateway_name[0].Hostname)
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	ip = gateway_ip
	port = "5500"
	return
}

// NvmeServiceCreate godoc
//
//	@Summary		Create of NVMe-OF Service
//	@Description	NVMe-OF 서비스를 생성합니다.
//	@param			pool_name 	formData	string	true	"Glue NVMe-OF Store Data In Pool Name"
//	@param			hostname	formData	[]string	true	"Glue NVMe-OF Service Placement Hosts" collectionFormat(multi)
//	@Tags			NVMe-OF
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nvmeof [post]
func (c *Controller) NvmeOfServiceCreate(ctx *gin.Context) {
	pool_name, _ := ctx.GetPostForm("pool_name")
	hosts, _ := ctx.GetPostFormArray("hostname")

	hosts_str := strings.Join(hosts, ",")

	dat, err := nvmeof.NvmeOfServiceCreate(pool_name, hosts_str)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// NvmeOfImageDownload godoc
//
//	@Summary		Download of NVMe-OF Image
//	@Description	NVMe-OF 이미지를 최신으로 다운받습니다.
//	@Tags			NVMe-OF
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nvmeof/image/download [post]
func (c *Controller) NvmeOfImageDownload(ctx *gin.Context) {
	dat, err := nvmeof.NvmeOfCliDownload()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// NvmeOfTargetCreate godoc
//
//	@Summary		Create of NVMe-OF Target
//	@Description	NVMe-OF 타켓을 생성합니다.
//	@param			subsystem_nqn_id	formData	string	true	"Glue NVMe-OF Sub System NQN ID"
//	@param			pool_name 	formData	string	true	"Glue NVMe-OF Use Image Pool Name"
//	@param			image_name 	formData	string	true	"Glue NVMe-OF Use Image Name"
//	@param			size 	formData	int	true	"Glue NVMe-OF Image Size(default GB)"
//	@Tags			NVMe-OF
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nvmeof/target [post]
func (c *Controller) NvmeOfTargetCreate(ctx *gin.Context) {
	subsystem_nqn_id, _ := ctx.GetPostForm("subsystem_nqn_id")
	pool_name, _ := ctx.GetPostForm("pool_name")
	image_name, _ := ctx.GetPostForm("image_name")
	size, _ := ctx.GetPostForm("size")

	size_int, _ := strconv.Atoi(size)
	size_int = size_int * 1024
	size = strconv.Itoa(size_int)

	var dat string
	name, err := nvmeof.NvmeOfGatewayName()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	host_nqn, err := nvmeof.NvmeOfHostNqn()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ip, err := nvmeof.HostIp(name[0].Hostname)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	_, err = nvmeof.NvmeOfSubSystemCreate(ip, "5500", subsystem_nqn_id)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		for i := 0; i < len(name); i++ {
			define_ip, err := nvmeof.HostIp(name[i].Hostname)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			gateway_name := string("client.") + name[i].Daemon_name
			dat, err = nvmeof.NvmeOfDefineGateway(define_ip, "5500", subsystem_nqn_id, gateway_name, define_ip)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
		}
		if dat == "Success" {
			_, err = nvmeof.NvmeOfHostAdd(ip, "5500", subsystem_nqn_id, host_nqn)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			} else {
				_, err = glue.CreateImage(image_name, pool_name, size)
				if err != nil {
					utils.FancyHandleError(err)
					httputil.NewError(ctx, http.StatusInternalServerError, err)
					return
				} else {
					dat, err = nvmeof.NvmeOfNameSpaceCreate(ip, "5500", subsystem_nqn_id, pool_name, image_name)
					if err != nil {
						utils.FancyHandleError(err)
						httputil.NewError(ctx, http.StatusInternalServerError, err)
						return
					}
					ctx.Header("Access-Control-Allow-Origin", "*")
					ctx.IndentedJSON(http.StatusOK, dat)
				}
			}
		}
	}
}

// NvmeOfTargetCreate2 godoc
//
//	@Summary		Create of NVMe-OF Target
//	@Description	NVMe-OF 타켓을 생성합니다. 개별적으로 하나 타겟만 만드는 API
//	@param			gateway_ip	formData	string	true	"Glue NVMe-OF Gateway IP"
//	@param			subsystem_nqn_id	formData	string	true	"Glue NVMe-OF Sub System NQN ID"
//	@param			pool_name 	formData	string	true	"Glue NVMe-OF Use Image Pool Name"
//	@param			image_name 	formData	string	true	"Glue NVMe-OF Use Image Name"
//	@param			size 	formData	int	true	"Glue NVMe-OF Image Size(default GB)"
//	@Tags			NVMe-OF
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nvmeof/target2 [post]
func (c *Controller) NvmeOfTargetCreate2(ctx *gin.Context) {
	gateway_ip, _ := ctx.GetPostForm("gateway_ip")
	subsystem_nqn_id, _ := ctx.GetPostForm("subsystem_nqn_id")
	pool_name, _ := ctx.GetPostForm("pool_name")
	image_name, _ := ctx.GetPostForm("image_name")
	size, _ := ctx.GetPostForm("size")

	size_int, _ := strconv.Atoi(size)
	size_int = size_int * 1024
	size = strconv.Itoa(size_int)

	var dat string
	gat_name, err := nvmeof.NvmeOfGatewayName()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	hostname, err := nvmeof.Hostname(gateway_ip)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	host_nqn, err := nvmeof.NvmeOfHostNqn()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	_, err = nvmeof.NvmeOfSubSystemCreate(gateway_ip, "5500", subsystem_nqn_id)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		var gateway_name string
		for i := 0; i < len(gat_name); i++ {
			if strings.Contains(gat_name[i].Daemon_name, hostname) {
				gateway_name = string("client.") + gat_name[i].Daemon_name
			}
		}
		dat, err = nvmeof.NvmeOfDefineGateway(gateway_ip, "5500", subsystem_nqn_id, gateway_name, gateway_ip)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		if dat == "Success" {
			_, err = nvmeof.NvmeOfHostAdd(gateway_ip, "5500", subsystem_nqn_id, host_nqn)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			} else {
				_, err = glue.CreateImage(image_name, pool_name, size)
				if err != nil {
					utils.FancyHandleError(err)
					httputil.NewError(ctx, http.StatusInternalServerError, err)
					return
				} else {
					dat, err = nvmeof.NvmeOfNameSpaceCreate(gateway_ip, "5500", subsystem_nqn_id, pool_name, image_name)
					if err != nil {
						utils.FancyHandleError(err)
						httputil.NewError(ctx, http.StatusInternalServerError, err)
						return
					}
					ctx.Header("Access-Control-Allow-Origin", "*")
					ctx.IndentedJSON(http.StatusOK, dat)
				}
			}
		}
	}
}

// NvmeOfTargetVerify godoc
//
//	@Summary		Show List of Verify the NVMe-OF Target is Reachable
//	@Description	NVMe-OF의 타겟에 연결 할 수 있는지에 대한 리스트를 보여줍니다.
//	@Tags			NVMe-OF
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	NvmeOfTargetVerify
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nvmeof/target [get]
func (c *Controller) NvmeOfTargetVerify(ctx *gin.Context) {

	ip, _, err := NvmeOfServerIPandPort()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	gateway, err := nvmeof.NvmeOfGatewayName()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	dat, err := nvmeof.NvmeOfTargetVerify(gateway[0].Hostname, ip)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// NvmeOfList godoc
//
//	@Summary		Show Detail and List of Connected NVMe-OF Infomation
//	@Description	연결된 NVMe-OF의 상세정보 및 리스트를 보여줍니다.
//	@param			detail 	query	string	true	"Glue NVMe-OF List detail" Enums(true, false) default(true)
//	@param			hostname	query	string	false	"Glue NVMe-OF Run Host Name"
//	@Tags			NVMe-OF
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	NvmeOfList
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nvmeof/list [get]
func (c *Controller) NvmeOfList(ctx *gin.Context) {
	hostname := ctx.Request.URL.Query().Get("hostname")
	detail := ctx.Request.URL.Query().Get("detail")

	gateway, err := nvmeof.NvmeOfGatewayName()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	if detail == "true" {
		if hostname == "" {
			var value []model.HostNvmeOfList
			for i := 0; i < len(gateway); i++ {
				dat, err := nvmeof.NvmeOfList(gateway[i].Hostname)
				if err != nil {
					utils.FancyHandleError(err)
					httputil.NewError(ctx, http.StatusInternalServerError, err)
					return
				}
				host_list := model.HostNvmeOfList{
					Hostname: gateway[i].Hostname,
					Detail: model.NvmeOfList{
						Devices: dat.Devices,
					},
				}
				value = append(value, host_list)
			}
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.IndentedJSON(http.StatusOK, value)
		} else {
			dat, err := nvmeof.NvmeOfList(hostname)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.IndentedJSON(http.StatusOK, dat)
		}
	} else {
		if hostname == "" {
			var value []model.NvmeOfPath
			for i := 0; i < len(gateway); i++ {
				dat, err := nvmeof.NvmeOfPath(gateway[i].Hostname)
				if err != nil {
					utils.FancyHandleError(err)
					httputil.NewError(ctx, http.StatusInternalServerError, err)
					return
				}
				path := model.NvmeOfPath{
					Hostname: gateway[i].Hostname,
					Devices:  dat.Devices,
				}
				value = append(value, path)
			}
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.IndentedJSON(http.StatusOK, value)
		} else {
			dat, err := nvmeof.NvmeOfPath(hostname)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.IndentedJSON(http.StatusOK, dat)
		}
	}
}

// NvmeOfConnect godoc
//
//	@Summary		Connect of NVMe-OF Sub System
//	@Description	NVMe-OF 하위 시스템에 연결합니다.
//	@param			hostname	formData	string	true	"Glue NVMe-OF Run Host Name"
//	@param			subsystem_nqn_id	formData	string	false	"Glue NVMe-OF Sub System NQN ID"
//	@param			full_connection 	formData	bool	true	"Glue NVMe-OF Full Connection Check" Enums(true, false) default(false)
//	@Tags			NVMe-OF
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nvmeof/connect [post]
func (c *Controller) NvmeOfConnect(ctx *gin.Context) {
	hostname, _ := ctx.GetPostForm("hostname")
	subsystem_nqn_id, _ := ctx.GetPostForm("subsystem_nqn_id")
	check := ctx.GetBool("full_connection")
	ip, _, err := NvmeOfServerIPandPort()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	if !check && subsystem_nqn_id == "" {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusBadRequest, "Please Sub System Nqn ID")
	} else {
		dat, err := nvmeof.NvmeOfConnect(hostname, ip, subsystem_nqn_id, check)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, dat)
	}
}

// NvmeOfDisConnect godoc
//
//	@Summary		DisConnect of NVMe-OF Sub System
//	@Description	NVMe-OF 하위 시스템에 연결을 해제합니다.
//	@param			hostname	formData	string	true	"Glue NVMe-OF Run Host Name"
//	@param			subsystem_nqn_id	formData	string	false	"Glue NVMe-OF Sub System NQN ID"
//	@param			full_DisConnection 	formData	bool	true	"Glue NVMe-OF Full DisConnection Check" Enums(true, false) default(false)
//	@Tags			NVMe-OF
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string "Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nvmeof/disconnect [post]
func (c *Controller) NvmeOfDisConnect(ctx *gin.Context) {
	hostname, _ := ctx.GetPostForm("hostname")
	subsystem_nqn_id, _ := ctx.GetPostForm("subsystem_nqn_id")
	check := ctx.GetBool("full_connection")
	if !check && subsystem_nqn_id == "" {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusBadRequest, "Please Sub System Nqn ID")
	} else {
		dat, err := nvmeof.NvmeOfDisConnect(hostname, subsystem_nqn_id, check)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, dat)
	}
}

// NvmeOfSubSystemList godoc
//
//	@Summary		Show List of NVMe-OF Sub System
//	@Description	NVMe-OF의 Sub System 리스트를 보여줍니다.
//	@param			subsystem_nqn_id	query	string	false	"Glue NVMe-OF Sub System NQN ID"
//	@Tags			NVMe-OF-SubSystem
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	NvmeOfSubSystemList
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nvmeof/subsystem [get]
func (c *Controller) NvmeOfSubSystemList(ctx *gin.Context) {
	subsystem_nqn_id := ctx.Request.URL.Query().Get("subsystem_nqn_id")

	ip, port, err := NvmeOfServerIPandPort()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	dat, err := nvmeof.NvmeOfSubSystemList(ip, port, subsystem_nqn_id)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// NvmeOfSubSystemCreate godoc
//
//	@Summary		Create of NVMe-OF Sub System
//	@Description	NVMe-OF의 Sub System을 생성합니다.
//	@param			gateway_ip	formData	string	true	"Glue NVMe-OF Gateway IP"
//	@param			subsystem_nqn_id	formData	string	true	"Glue NVMe-OF Sub System NQN ID"
//	@Tags			NVMe-OF-SubSystem
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nvmeof/subsystem [post]
func (c *Controller) NvmeOfSubSystemCreate(ctx *gin.Context) {
	gateway_ip, _ := ctx.GetPostForm("gateway_ip")
	subsystem_nqn_id, _ := ctx.GetPostForm("subsystem_nqn_id")

	gat_name, err := nvmeof.NvmeOfGatewayName()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	hostname, err := nvmeof.Hostname(gateway_ip)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	host_nqn, err := nvmeof.NvmeOfHostNqn()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	_, err = nvmeof.NvmeOfSubSystemCreate(gateway_ip, "5500", subsystem_nqn_id)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		var gateway_name string
		for i := 0; i < len(gat_name); i++ {
			if strings.Contains(gat_name[i].Daemon_name, hostname) {
				gateway_name = string("client.") + gat_name[i].Daemon_name
			}
		}
		_, err = nvmeof.NvmeOfDefineGateway(gateway_ip, "5500", subsystem_nqn_id, gateway_name, gateway_ip)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			dat, err := nvmeof.NvmeOfHostAdd(gateway_ip, "5500", subsystem_nqn_id, host_nqn)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.IndentedJSON(http.StatusOK, dat)
		}
	}
}

// NvmeOfSubSystemDelete godoc
//
//	@Summary		Delete of NVMe-OF Sub System
//	@Description	NVMe-OF Sub System을 삭제합니다.
//	@param			subsystem_nqn_id	query	string	true	"Glue NVMe-OF Sub System NQN ID"
//	@Tags			NVMe-OF-SubSystem
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nvmeof/subsystem [delete]
func (c *Controller) NvmeOfSubSystemDelete(ctx *gin.Context) {
	subsystem_nqn_id := ctx.Request.URL.Query().Get("subsystem_nqn_id")

	ip, port, err := NvmeOfServerIPandPort()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	dat, err := nvmeof.NvmeOfSubSystemDelete(ip, port, subsystem_nqn_id)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// NvmeOfNameSpaceList godoc
//
//	@Summary		Show List of NVMe-OF Sub Systems
//	@Description	NVMe-OF의 Sub System 리스트를 보여줍니다.
//	@param			subsystem_nqn_id	query	string	true	"Glue NVMe-OF Sub System NQN ID"
//	@Tags			NVMe-OF-NameSpace
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	NvmeOfNameSpaceList
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nvmeof/namespace [get]
func (c *Controller) NvmeOfNameSpaceList(ctx *gin.Context) {
	subsystem_nqn_id := ctx.Request.URL.Query().Get("subsystem_nqn_id")
	ip, port, err := NvmeOfServerIPandPort()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	if subsystem_nqn_id == "" {
		dat, err := nvmeof.NvmeOfSubSystemList(ip, port, subsystem_nqn_id)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		var value []model.NvmeOfNameSpaceList
		for i := 0; i < len(dat.Subsystems); i++ {
			namespace, err := nvmeof.NvmeOfNameSpaceList(ip, port, dat.Subsystems[i].Nqn)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			value = append(value, namespace)
		}
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, value)
	} else {
		dat, err := nvmeof.NvmeOfNameSpaceList(ip, port, subsystem_nqn_id)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, dat)
	}

}

// NvmeOfNameSpaceCreate godoc
//
//	@Summary		Create of NVMe-OF NameSpace
//	@Description	NVMe-OF의 네임스페이스를 생성합니다.
//	@param			gateway_ip	formData	string	true	"Glue NVMe-OF Gateway IP"
//	@param			subsystem_nqn_id	formData	string	true	"Glue NVMe-OF Sub System NQN ID"
//	@param			pool_name	formData	string	true	"Glue Pool Name"
//	@param			image_name	formData	string	true	"Glue Image Name"
//	@param			size 	formData	int	true	"Glue NVMe-OF Image Size(default GB)"
//	@Tags			NVMe-OF-NameSpace
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nvmeof/namespace [post]
func (c *Controller) NvmeOfNameSpaceCreate(ctx *gin.Context) {
	gateway_ip, _ := ctx.GetPostForm("gateway_ip")
	subsystem_nqn_id, _ := ctx.GetPostForm("subsystem_nqn_id")
	pool_name, _ := ctx.GetPostForm("pool_name")
	image_name, _ := ctx.GetPostForm("image_name")
	size, _ := ctx.GetPostForm("size")

	size_int, _ := strconv.Atoi(size)
	size_int = size_int * 1024
	size = strconv.Itoa(size_int)

	_, port, err := NvmeOfServerIPandPort()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	_, err = glue.CreateImage(image_name, pool_name, size)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		dat, err := nvmeof.NvmeOfNameSpaceCreate(gateway_ip, port, subsystem_nqn_id, pool_name, image_name)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, dat)
	}
}

// NvmeOfNameSpaceDelete godoc
//
//	@Summary		Delete of NVMe-OF NameSpace
//	@Description	NVMe-OF의 네임스페이스를 삭제합니다.
//	@param			subsystem_nqn_id	query	string	true	"Glue NVMe-OF Sub System NQN ID"
//	@param			namespace_uuid	query	string	true	"Glue NVMe-OF NameSpace UUID"
//	@param			image_del_check query	bool	true	"Glue NVMe-OF Image Delete Check" Enums(true, false) default(false)
//	@param			pool_name	query	string	false	"Glue Pool Name"
//	@param			image_name	query	string	false	"Glue Image Name"
//	@Tags			NVMe-OF-NameSpace
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nvmeof/namespace [delete]
func (c *Controller) NvmeOfNameSpaceDelete(ctx *gin.Context) {
	subsystem_nqn_id := ctx.Request.URL.Query().Get("subsystem_nqn_id")
	namespace_uuid := ctx.Request.URL.Query().Get("namespace_uuid")
	image_del_check := ctx.Request.URL.Query().Get("image_del_check")
	image_name := ctx.Request.URL.Query().Get("image_name")
	pool_name := ctx.Request.URL.Query().Get("pool_name")

	ip, port, err := NvmeOfServerIPandPort()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	if image_del_check == "true" {
		if image_name == "" || pool_name == "" {
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.IndentedJSON(http.StatusBadRequest, "Please Check Image Name and Pool Name")
		} else {
			dat, err := nvmeof.NvmeOfNameSpaceDelete(ip, port, subsystem_nqn_id, namespace_uuid)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			} else {
				_, err = glue.DeleteImage(image_name, pool_name)
				if err != nil {
					utils.FancyHandleError(err)
					httputil.NewError(ctx, http.StatusInternalServerError, err)
					return
				}
			}
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.IndentedJSON(http.StatusOK, dat)
		}
	} else {
		dat, err := nvmeof.NvmeOfNameSpaceDelete(ip, port, subsystem_nqn_id, namespace_uuid)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, dat)
	}
}
