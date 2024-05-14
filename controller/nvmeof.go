package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/glue"
	"Glue-API/utils/nvmeof"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

func (c *Controller) NvmeOption(ctx *gin.Context) {
	SetOptionHeader(ctx)
	ctx.IndentedJSON(http.StatusOK, nil)
}
func NvmeOfServerIPandPort() (server_gateway_ip string, port string, err error) {
	gateway_name, err := nvmeof.NvmeOfGatewayName()
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	if len(gateway_name) == 0 {
		server_gateway_ip = "not"
	} else {
		server_gateway_ip, err = nvmeof.ServerGatewayIp(gateway_name[0].Hostname)
		if err != nil {
			utils.FancyHandleError(err)
			return
		}
	}
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

	value := model.NvmeOfServiceCreate{
		ServiceType: "nvmeof",
		ServiceId:   pool_name,
		Placement: model.NvmeOfPlacement{
			Hosts: hosts,
		},
		Spec: model.NvmeOfSpec{
			Pool:            pool_name,
			TgtCmdExtraArgs: "--cpumask=0xFF",
		},
	}
	yaml_data, err := yaml.Marshal(value)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	nvmeof_conf := "/etc/ceph/nvmeof.conf"
	err = os.WriteFile(nvmeof_conf, yaml_data, 0644)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	} else {
		dat, err := nvmeof.NvmeOfServiceCreate(nvmeof_conf, pool_name)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			if err := os.Remove(nvmeof_conf); err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
		}
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, dat)
	}
}

// NvmeOfImageDownload godoc
//
//	@Summary		Download of NVMe-OF Image
//	@Description	NVMe-OF 이미지를 최신으로 다운받습니다.
//	@param			gateway_ip	formData	string	true	"Glue NVMe-OF Gateway IP"
//	@Tags			NVMe-OF
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{string}	string	"Success"
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nvmeof/image/download [post]
func (c *Controller) NvmeOfImageDownload(ctx *gin.Context) {
	gateway_ip, _ := ctx.GetPostForm("gateway_ip")

	dat, err := nvmeof.NvmeOfCliDownload(gateway_ip)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.IndentedJSON(http.StatusOK, dat)
}

// // NvmeOfTargetCreate godoc
// //
// //	@Summary		Create of NVMe-OF Target
// //	@Description	NVMe-OF 타켓을 생성합니다.
// //	@param			subsystem_nqn_id	formData	string	true	"Glue NVMe-OF Sub System NQN ID"
// //	@param			pool_name 	formData	string	true	"Glue NVMe-OF Use Image Pool Name"
// //	@param			image_name 	formData	string	true	"Glue NVMe-OF Use Image Name"
// //	@param			size 	formData	int	true	"Glue NVMe-OF Image Size(default GB)"
// //	@Tags			NVMe-OF
// //	@Accept			x-www-form-urlencoded
// //	@Produce		json
// //	@Success		200	{string}	string	"Success"
// //	@Failure		400	{object}	httputil.HTTP400BadRequest
// //	@Failure		404	{object}	httputil.HTTP404NotFound
// //	@Failure		500	{object}	httputil.HTTP500InternalServerError
// //	@Router			/api/v1/nvmeof/target [post]
// func (c *Controller) NvmeOfTargetCreate(ctx *gin.Context) {
// 	subsystem_nqn_id, _ := ctx.GetPostForm("subsystem_nqn_id")
// 	pool_name, _ := ctx.GetPostForm("pool_name")
// 	image_name, _ := ctx.GetPostForm("image_name")
// 	size, _ := ctx.GetPostForm("size")

// 	size_int, _ := strconv.Atoi(size)
// 	size_int = size_int * 1024
// 	size = strconv.Itoa(size_int)

// 	var dat string
// 	name, err := nvmeof.NvmeOfGatewayName()
// 	if err != nil {
// 		utils.FancyHandleError(err)
// 		httputil.NewError(ctx, http.StatusInternalServerError, err)
// 		return
// 	}
// 	ip, err := nvmeof.HostIp(name[0].Hostname)
// 	if err != nil {
// 		utils.FancyHandleError(err)
// 		httputil.NewError(ctx, http.StatusInternalServerError, err)
// 		return
// 	}
// 	_, err = nvmeof.NvmeOfSubSystemCreate(ip, ip, "5500", subsystem_nqn_id)
// 	if err != nil {
// 		utils.FancyHandleError(err)
// 		httputil.NewError(ctx, http.StatusInternalServerError, err)
// 		return
// 	} else {
// 		for i := 0; i < len(name); i++ {
// 			define_ip, err := nvmeof.HostIp(name[i].Hostname)
// 			if err != nil {
// 				utils.FancyHandleError(err)
// 				httputil.NewError(ctx, http.StatusInternalServerError, err)
// 				return
// 			}
// 			gateway_name := string("client.") + name[i].Daemon_name
// 			dat, err = nvmeof.NvmeOfDefineGateway(ip, define_ip, "5500", subsystem_nqn_id, gateway_name, define_ip)
// 			if err != nil {
// 				utils.FancyHandleError(err)
// 				httputil.NewError(ctx, http.StatusInternalServerError, err)
// 				return
// 			}
// 		}
// 		if dat == "Success" {
// 			_, err = nvmeof.NvmeOfHostAdd(ip, ip, "5500", subsystem_nqn_id)
// 			if err != nil {
// 				utils.FancyHandleError(err)
// 				httputil.NewError(ctx, http.StatusInternalServerError, err)
// 				return
// 			} else {
// 				_, err = glue.CreateImage(image_name, pool_name, size)
// 				if err != nil {
// 					utils.FancyHandleError(err)
// 					httputil.NewError(ctx, http.StatusInternalServerError, err)
// 					return
// 				} else {
// 					dat, err = nvmeof.NvmeOfNameSpaceCreate(ip, ip, "5500", subsystem_nqn_id, pool_name, image_name)
// 					if err != nil {
// 						utils.FancyHandleError(err)
// 						httputil.NewError(ctx, http.StatusInternalServerError, err)
// 						return
// 					}
// 					ctx.Header("Access-Control-Allow-Origin", "*")
// 					ctx.IndentedJSON(http.StatusOK, dat)
// 				}
// 			}
// 		}
// 	}
// }

// NvmeOfTargetCreate godoc
//
//	@Summary		Create of NVMe-OF Target
//	@Description	NVMe-OF 타켓을 생성합니다.
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
//	@Router			/api/v1/nvmeof/target [post]
func (c *Controller) NvmeOfTargetCreate(ctx *gin.Context) {
	gateway_ip, _ := ctx.GetPostForm("gateway_ip")
	subsystem_nqn_id, _ := ctx.GetPostForm("subsystem_nqn_id")
	pool_name, _ := ctx.GetPostForm("pool_name")
	image_name, _ := ctx.GetPostForm("image_name")
	size, _ := ctx.GetPostForm("size")

	size_int, _ := strconv.Atoi(size)
	size_int = size_int * 1024
	size = strconv.Itoa(size_int)

	var dat string
	server_gateway_ip, port, err := NvmeOfServerIPandPort()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	if server_gateway_ip == "not" {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, make([]string, 0))
	} else {
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
		image, _ := glue.ListAndInfoImage(image_name, pool_name)
		if image != nil {
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.IndentedJSON(http.StatusOK, "The Image Name exists. Please Check.")
		} else {
			_, err = nvmeof.NvmeOfSubSystemCreate(server_gateway_ip, server_gateway_ip, port, subsystem_nqn_id)
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
				dat, err = nvmeof.NvmeOfDefineGateway(server_gateway_ip, server_gateway_ip, port, subsystem_nqn_id, gateway_name, gateway_ip)
				if err != nil {
					utils.FancyHandleError(err)
					httputil.NewError(ctx, http.StatusInternalServerError, err)
					return
				}
				if dat == "Success" {
					_, err = nvmeof.NvmeOfHostAdd(server_gateway_ip, server_gateway_ip, port, subsystem_nqn_id)
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
							dat, err = nvmeof.NvmeOfNameSpaceCreate(server_gateway_ip, server_gateway_ip, port, subsystem_nqn_id, pool_name, image_name)
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

	server_gateway_ip, port, err := NvmeOfServerIPandPort()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	if server_gateway_ip == "not" {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, make([]string, 0))
	} else {
		dat, err := nvmeof.NvmeOfSubSystemList(server_gateway_ip, server_gateway_ip, port, subsystem_nqn_id)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, dat)
	}
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

	server_gateway_ip, port, err := NvmeOfServerIPandPort()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	if server_gateway_ip == "not" {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, make([]string, 0))
	} else {
		gat_name, err := nvmeof.NvmeOfGatewayName()
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		hostname, err := nvmeof.Hostname(server_gateway_ip)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		_, err = nvmeof.NvmeOfSubSystemCreate(server_gateway_ip, server_gateway_ip, port, subsystem_nqn_id)
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
			_, err = nvmeof.NvmeOfDefineGateway(server_gateway_ip, server_gateway_ip, port, subsystem_nqn_id, gateway_name, gateway_ip)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			} else {
				dat, err := nvmeof.NvmeOfHostAdd(server_gateway_ip, server_gateway_ip, port, subsystem_nqn_id)
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

	server_gateway_ip, port, err := NvmeOfServerIPandPort()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	if server_gateway_ip == "not" {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, make([]string, 0))
	} else {
		dat, err := nvmeof.NvmeOfSubSystemDelete(server_gateway_ip, server_gateway_ip, port, subsystem_nqn_id)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, dat)
	}
}

// NvmeOfNameSpaceList godoc
//
//	@Summary		Show List of NVMe-OF Sub Systems
//	@Description	NVMe-OF의 Sub System 리스트를 보여줍니다.
//	@param			subsystem_nqn_id	query	string	false	"Glue NVMe-OF Sub System NQN ID"
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
	server_gateway_ip, port, err := NvmeOfServerIPandPort()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	if server_gateway_ip == "not" {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, make([]string, 0))
	} else {
		if subsystem_nqn_id == "" {
			dat, err := nvmeof.NvmeOfSubSystemList(server_gateway_ip, server_gateway_ip, port, subsystem_nqn_id)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			var value []model.NvmeOfNameSpaceList
			for i := 0; i < len(dat.Subsystems); i++ {
				namespace, err := nvmeof.NvmeOfNameSpaceList(server_gateway_ip, server_gateway_ip, port, dat.Subsystems[i].Nqn)
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
			dat, err := nvmeof.NvmeOfNameSpaceList(server_gateway_ip, server_gateway_ip, port, subsystem_nqn_id)
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

// NvmeOfNameSpaceCreate godoc
//
//	@Summary		Create of NVMe-OF NameSpace
//	@Description	NVMe-OF의 네임스페이스를 생성합니다.
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
	subsystem_nqn_id, _ := ctx.GetPostForm("subsystem_nqn_id")
	pool_name, _ := ctx.GetPostForm("pool_name")
	image_name, _ := ctx.GetPostForm("image_name")
	size, _ := ctx.GetPostForm("size")

	size_int, _ := strconv.Atoi(size)
	size_int = size_int * 1024
	size = strconv.Itoa(size_int)

	server_gateway_ip, port, err := NvmeOfServerIPandPort()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	if server_gateway_ip == "not" {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, make([]string, 0))
	} else {
		_, err = glue.CreateImage(image_name, pool_name, size)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else {
			dat, err := nvmeof.NvmeOfNameSpaceCreate(server_gateway_ip, server_gateway_ip, port, subsystem_nqn_id, pool_name, image_name)
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

	server_gateway_ip, port, err := NvmeOfServerIPandPort()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	if server_gateway_ip == "not" {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, make([]string, 0))
	} else {
		if image_del_check == "true" {
			if image_name == "" || pool_name == "" {
				ctx.Header("Access-Control-Allow-Origin", "*")
				ctx.IndentedJSON(http.StatusOK, "Please Check Image Name and Pool Name")
			} else {
				dat, err := nvmeof.NvmeOfNameSpaceDelete(server_gateway_ip, server_gateway_ip, port, subsystem_nqn_id, namespace_uuid)
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
			dat, err := nvmeof.NvmeOfNameSpaceDelete(server_gateway_ip, server_gateway_ip, port, subsystem_nqn_id, namespace_uuid)
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

// NvmeOfTargetList godoc
//
//	@Summary		Show List of NVMe-OF Target
//	@Description	NVMe-OF의 타겟 리스트를 보여줍니다.
//	@param			subsystem_nqn_id	query	string	false	"Glue NVMe-OF Sub System NQN ID"
//	@Tags			NVMe-OF
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	NvmeOfTarget
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/nvmeof/target [get]
func (c *Controller) NvmeOfTargetList(ctx *gin.Context) {
	subsystem_nqn_id := ctx.Request.URL.Query().Get("subsystem_nqn_id")
	server_gateway_ip, port, err := NvmeOfServerIPandPort()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	if server_gateway_ip == "not" {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, make([]string, 0))
	} else {
		container_id, err := nvmeof.Container(server_gateway_ip)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		var list model.NvmeOfTarget
		list, err = nvmeof.NvmeOfTarget(server_gateway_ip, container_id, subsystem_nqn_id)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		for i := 0; i < len(list); i++ {
			con, err := nvmeof.NvmeOfConnection(server_gateway_ip, container_id, list[i].Nqn)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			list[i].Session = len(con)

			image, err := nvmeof.NvmeOfNameSpaceList(server_gateway_ip, server_gateway_ip, port, list[i].Nqn)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
			if len(image.Namespaces) != 0 {
				list[i].Namespaces[0].Block_size = image.Namespaces[0].BlockSize
				list[i].Namespaces[0].Rbd_image_name = image.Namespaces[0].RbdImageName
				list[i].Namespaces[0].Rbd_image_size = image.Namespaces[0].RbdImageSize
				list[i].Namespaces[0].Rbd_pool_name = image.Namespaces[0].RbdPoolName
			}
		}
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.IndentedJSON(http.StatusOK, list)
	}
}
