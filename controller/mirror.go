package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/mirror"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/melbahja/goph"
)

// MirrorImageList godoc
//
//	@Summary		Show List of Mirrored Snapshot
//	@Description	미러링중인 이미지의 목록과 상태를 보여줍니다.
//	@param			mirrorPool	path	string	true	"mirrorPool"
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.MirrorList
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror/image/{mirrorPool} [get]
func (c *Controller) MirrorImageList(ctx *gin.Context) {
	pool := ctx.Param("mirrorPool")
	mirrorStatus, err := mirror.GetConfigure()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	if mirrorStatus.Mode != "disabled" {
		dat, err := mirror.ImageList(pool)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(http.StatusOK, dat)
	} else {
		dat := model.MirrorList{}
		ctx.IndentedJSON(http.StatusOK, dat)
	}
}

// MirrorImageInfo godoc
//
//	@Summary		Show Information of Mirrored Snapshot
//	@Description	미러링중인 이미지의 정보를 보여줍니다.
//	@param			mirrorPool	path	string	true	"mirrorPool"
//	@param			imageName	path	string	true	"imageName"
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.ImageMirror
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror/image/{mirrorPool}/{imageName} [get]
func (c *Controller) MirrorImageInfo(ctx *gin.Context) {
	pool := ctx.Param("mirrorPool")
	image := ctx.Param("imageName")
	dat2, err := mirror.ImageList(pool)
	var dat model.MirrorListImages

	for _, mirrorImage := range dat2.Images {
		if mirrorImage.Name == image {
			dat.Name = mirrorImage.Name
			dat.State = mirrorImage.State
			dat.Description = mirrorImage.Description
			dat.PeerSites = mirrorImage.PeerSites
		}
	}
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, dat)
}

// MirrorImageDelete godoc
//
//	@Summary		Delete Mirrored Snapshot
//	@Description	이미지의 미러링을 비활성화 합니다.
//	@param			mirrorPool	path	string	true	"pool"
//	@param			imageName	path	string	true	"imageName"
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	controller.Message
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror/image/{mirrorPool}/{imageName} [delete]
// func (c *Controller) MirrorImageDelete(ctx *gin.Context) {
// 	image := ctx.Param("imageName")
// 	pool := ctx.Param("mirrorPool")
// 	var output string

// 	output, err := mirror.ImageDelete(pool, image)

// 	if err != nil {
// 		utils.FancyHandleError(err)
// 		httputil.NewError(ctx, http.StatusInternalServerError, err)
// 		return
// 	}

// 	output, err = mirror.ImagePreDelete(pool, image)

// 	if err != nil {
// 		if output != "Success" {
// 			utils.FancyHandleError(err)
// 			httputil.NewError(ctx, http.StatusInternalServerError, err)
// 			return
// 		}
// 	}

// 	ctx.IndentedJSON(http.StatusOK, Message{Message: output})
// }

// MirrorImageScheduleDelete godoc
//
//	@Summary		Delete Mirrored Snapshot Schedule
//	@Description	이미지의 미러링 스케줄링을 비활성화 합니다.
//	@param			mirrorPool	path	string	true	"pool"
//	@param			imageName	path	string	true	"imageName"
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	controller.Message
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror/image/{mirrorPool}/{imageName} [delete]
func (c *Controller) MirrorImageScheduleDelete(ctx *gin.Context) {
	image := ctx.Param("imageName")
	pool := ctx.Param("mirrorPool")
	var output string

	output, err := mirror.ImageDeleteSchedule(pool, image)

	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	output, err = mirror.ImagePreDelete(pool, image)

	if err != nil {
		if output != "Success" {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	output, err = mirror.ImageMetaRemove(image)
	if err != nil {
		if output != "Success" {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	ctx.IndentedJSON(http.StatusOK, Message{Message: output})
}

// MirrorStatus godoc
//
//	@Summary		Show Status of Mirror
//	@Description	Glue 의 미러링 상태를 보여줍니다.
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.MirrorStatus
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror [get]
func (c *Controller) MirrorStatus(ctx *gin.Context) {
	dat, err := mirror.Status()
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	// Print the output
	ctx.IndentedJSON(http.StatusOK, dat)
}

// MirrorSetup godoc
//
//	@Summary		Setup Mirroring Cluster
//	@Description	Glue 의 미러링 클러스터를 설정합니다.
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.MirrorSetup
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror [POST]
func (c *Controller) MirrorSetup(ctx *gin.Context) {
	var dat model.MirrorSetup

	dat.LocalClusterName, _ = ctx.GetPostForm("localClusterName")
	dat.RemoteClusterName, _ = ctx.GetPostForm("remoteClusterName")
	dat.Host, _ = ctx.GetPostForm("host")
	dat.MirrorPool, _ = ctx.GetPostForm("mirrorPool")
	file, _ := ctx.FormFile("privateKeyFile")
	privkey, err := os.CreateTemp("", "id_rsa-")
	defer privkey.Close()
	defer os.Remove(privkey.Name())
	privkeyname := privkey.Name()

	// Upload the file to specific dst.
	err = ctx.SaveUploadedFile(file, privkeyname)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	EncodedLocalToken, EncodedRemoteToken, err := mirror.ConfigMirror(dat, privkeyname)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	moldUrl, _ := ctx.GetPostForm("moldUrl")
	moldApiKey, _ := ctx.GetPostForm("moldApiKey")
	moldSecretKey, _ := ctx.GetPostForm("moldSecretKey")

	err = mirror.ConfigMold(moldUrl, moldApiKey, moldSecretKey)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	dat.LocalToken = EncodedLocalToken
	dat.RemoteToken = EncodedRemoteToken

	ctx.IndentedJSON(http.StatusOK, dat)
}

// MirrorUpdate godoc
//
//		@Summary		Put Mirroring Cluster
//		@Description	Glue 의 미러링 클러스터의 설정을 변경합니다.
//	 	@param			interval		formData 	string	true	"Mirroring Schedule Interval"
//		@param			moldUrl			formData	string	true	"Mold API request URL"
//		@param			moldApiKey		formData	string	true	"Mold Admin Api Key"
//		@param			moldSecretKey	formData	string	true	"Mold Admin Secret Key"
//		@Tags			Mirror
//		@Accept			x-www-form-urlencoded
//		@Produce		json
//		@Success		200	{object}	model.Mold
//		@Failure		400	{object}	httputil.HTTP400BadRequest
//		@Failure		404	{object}	httputil.HTTP404NotFound
//		@Failure		500	{object}	httputil.HTTP500InternalServerError
//		@Router			/api/v1/mirror [put]
func (c *Controller) MirrorUpdate(ctx *gin.Context) {

	var mold = model.Mold{}

	interval, _ := ctx.GetPostForm("interval")
	moldUrl, _ := ctx.GetPostForm("moldUrl")
	moldApiKey, _ := ctx.GetPostForm("moldApiKey")
	moldSecretKey, _ := ctx.GetPostForm("moldSecretKey")

	err := mirror.ConfigMold(moldUrl, moldApiKey, moldSecretKey)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	rbd_image, err := mirror.RbdImage("rbd")
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	for i := 0; i < len(rbd_image); i++ {
		if rbd_image[i] == "MOLD-DR" {
			err := mirror.ImageMetaUpdate(interval)
			if err != nil {
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
		}
	}
	ctx.IndentedJSON(http.StatusOK, mold)
}

// MirrorDelete godoc
//
//	@Summary		Delete Mirroring Cluster
//	@Description	Glue 의 미러링 클러스터를 제거합니다.
//	@param			host			formData	string	true	"Remote Cluster Host Address"
//	@param			privateKeyFile	formData	file	true	"Remote Cluster PrivateKey"
//	@param			mirrorPool		formData	string	true	"Pool Name for Mirroring"
//	@Tags			Mirror
//	@Accept			multipart/form-data
//	@Produce		json
//	@Success		200	{object}	model.MirrorSetup
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror [delete]
func (c *Controller) MirrorDelete(ctx *gin.Context) {
	var dat model.MirrorSetup
	var stdout []byte

	var out strings.Builder
	var output string
	dat.Host, _ = ctx.GetPostForm("host")
	dat.MirrorPool, _ = ctx.GetPostForm("mirrorPool")
	file, _ := ctx.FormFile("privateKeyFile")

	privkey, err := os.CreateTemp("", "id_rsa-")
	defer privkey.Close()
	defer os.Remove(privkey.Name())
	privkeyname := privkey.Name()

	// Upload the file to specific dst.
	err = ctx.SaveUploadedFile(file, privkeyname)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Get Mirroring Images
	mirrorStatus, err := mirror.GetConfigure()
	if mirrorStatus.Mode != "disabled" {
		MirroredImage, err := mirror.ImageList(dat.MirrorPool)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		for _, image := range MirroredImage.Images {
			_, errt := mirror.ImageDeleteSchedule(dat.MirrorPool, image.Name)
			if errt != nil {
				err = errors.Join(err, errt)
			}
			output, errt = mirror.ImagePreDelete(dat.MirrorPool, image.Name)
			if errt != nil {
				if output != "Success" {
					err = errors.Join(err, errt)
				}
			}
		}
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}
	//remote local peer

	if len(mirrorStatus.Peers) > 0 {
		peerUUID := mirrorStatus.Peers[0].Uuid
		cmd := exec.Command("rbd", "mirror", "pool", "peer", "remove", "--pool", dat.MirrorPool, peerUUID)
		stdout, err = cmd.CombinedOutput()
		println("out: " + string(stdout))
		println("err: " + out.String())
		// if err != nil || (out.String() != "" && out.String() != "rbd: mirroring is already configured for image mode") {
		if err != nil {
			cmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		cmd = exec.Command("ceph", "auth", "del", "client.rbd-mirror-peer")
		stdout, err = cmd.CombinedOutput()
		println("out: " + string(stdout))
		println("err: " + out.String())
		// if err != nil || (out.String() != "" && out.String() != "rbd: mirroring is already configured for image mode") {
		if err != nil {
			cmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// Mirror Disable
	if mirrorStatus.Mode != "disabled" {
		cmd := exec.Command("rbd", "mirror", "pool", "disable")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			cmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// Mirror Daemon Destroy
	cmd := exec.Command("ceph", "orch", "rm", "rbd-mirror")
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		cmd.Stderr = &out
		if !strings.Contains(out.String(), "Invalid service 'rbd-mirror'.") {
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// DR Mirror Image Destroy
	cmd = exec.Command("rbd", "rm", "rbd/MOLD-DR")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		cmd.Stderr = &out
		if !strings.Contains(string(stdout), "No such file or directory") {
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// DR mold conf reset
	err = mirror.ConfigMold("moldUrl", "moldApiKey", "moldSecretKey")
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	//remote local peer
	out.Reset()
	client, err := utils.ConnectSSH(dat.Host, privkeyname)
	if err != nil {
		err = errors.Join(err, errors.New("failed to connect ssh."))
		utils.FancyHandleError(err)
		return
	}
	defer func(client *goph.Client) {
		err := client.Close()
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}(client)
	remoteMirrorStatus, err := mirror.GetRemoteConfigure(client)

	if len(remoteMirrorStatus.Peers) > 0 {
		peerUUID := remoteMirrorStatus.Peers[0].Uuid
		sshcmd, err := client.Command("rbd", "mirror", "pool", "peer", "remove", "--pool", dat.MirrorPool, peerUUID)
		if err != nil {
			sshcmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			return
		}

		stdout, err = sshcmd.CombinedOutput()
		if err != nil {
			sshcmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		sshcmd, err = client.Command("ceph", "auth", "del", "client.rbd-mirror-peer")
		if err != nil {
			sshcmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			return
		}

		stdout, err = sshcmd.CombinedOutput()
		if err != nil {
			sshcmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

	}

	// Mirror Disable
	if remoteMirrorStatus.Mode != "disabled" {
		sshcmd, err := client.Command("rbd", "mirror", "pool", "disable")
		if err != nil {
			sshcmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			return
		}
		stdout, err = sshcmd.CombinedOutput()
		if err != nil {
			sshcmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// Mirror Daemon Destroy
	sshcmd, err := client.Command("ceph", "orch", "rm", "rbd-mirror")
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		sshcmd.Stderr = &out
		if !strings.Contains(out.String(), "Invalid service 'rbd-mirror'.") {
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// DR Mirror Image Destroy
	sshcmd, err = client.Command("rbd", "ls", "-p", "rbd", "--format", "json")
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	var pools []string
	stdout, err = sshcmd.CombinedOutput()
	if err = json.Unmarshal(stdout, &pools); err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	for i := 0; i < len(pools); i++ {
		if pools[i] == "MOLD-DR" {
			sshcmd, err = client.Command("rbd", "rm", "rbd/MOLD-DR")
			if err != nil {
				sshcmd.Stderr = &out
				err = errors.Join(err, errors.New(out.String()))
				utils.FancyHandleError(err)
				return
			}
			stdout, err = sshcmd.CombinedOutput()
			if err != nil {
				sshcmd.Stderr = &out
				err = errors.Join(err, errors.New(out.String()))
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
		}
	}

	// Secondary DR mold conf reset 추가 필요
	// err = mirror.ConfigMold("moldUrl", "moldApiKey", "moldSecretKey")
	// if err != nil {
	// 	utils.FancyHandleError(err)
	// 	httputil.NewError(ctx, http.StatusInternalServerError, err)
	// 	return
	// }

	ctx.IndentedJSON(http.StatusOK, dat)
}

// MirrorImageSetup godoc
//
//	@Summary		Setup Image Mirroring
//	@Description	Glue 의 이미지에 미러링을 설정합니다.
//	@param			mirrorPool	path		string	true	"Pool Name for Mirroring"
//	@param			imageName	path		string	true	"Image Name for Mirroring"
//	@param			interval	formData	string	true	"Interval of image snapshot"
//	@param			startTime	formData	string	false	"StartTime of image snapshot"
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.ImageMirror
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror/image/{mirrorPool}/{imageName} [post]
// func (c *Controller) MirrorImageSetup(ctx *gin.Context) {
// 	//var dat model.MirrorSetup
// 	var dat = struct {
// 		Message string
// 	}{}

// 	//mirrorPool, _ := ctx.GetPostForm("mirrorPool")
// 	mirrorPool := ctx.Param("mirrorPool")
// 	imageName := ctx.Param("imageName")
// 	//imageName, _ := ctx.GetPostForm("imageName")
// 	interval, _ := ctx.GetPostForm("interval")
// 	startTime, _ := ctx.GetPostForm("startTime")
// 	print(startTime)

// 	message, err := mirror.ImagePreSetup(mirrorPool, imageName)
// 	if err != nil {
// 		utils.FancyHandleError(err)
// 		httputil.NewError(ctx, http.StatusInternalServerError, err)
// 		return
// 	}

// 	message, err = mirror.ImageSetup(mirrorPool, imageName)
// 	if err != nil {
// 		utils.FancyHandleError(err)
// 		httputil.NewError(ctx, http.StatusInternalServerError, err)
// 		return
// 	}

// 	message, err = mirror.ImageConfig(mirrorPool, imageName, interval, startTime)
// 	if err != nil {
// 		utils.FancyHandleError(err)
// 		httputil.NewError(ctx, http.StatusInternalServerError, err)
// 		return
// 	}
// 	dat.Message = message
// 	ctx.IndentedJSON(http.StatusOK, dat)

// }

// MirrorImageScheduleSetup godoc
//
//		@Summary		Setup Image Mirroring Schedule
//		@Description	Glue 의 이미지에 미러링 스케줄을 설정합니다.
//		@param			mirrorPool	path		string	true	"Pool Name for Mirroring"
//		@param			imageName	path		string	true	"Image Name for Mirroring"
//	 	@param          hostName    path    	string  true    "Host Name"
//	 	@param          vmName      path    	string  true    "VM Name"
//		@param			volType		formData	string	true	"Volume Type"
//		@Tags			Mirror
//		@Accept			x-www-form-urlencoded
//		@Produce		json
//		@Success		200	{object}	model.ImageMirror
//		@Failure		400	{object}	httputil.HTTP400BadRequest
//		@Failure		404	{object}	httputil.HTTP404NotFound
//		@Failure		500	{object}	httputil.HTTP500InternalServerError
//		@Router			/api/v1/mirror/image/{mirrorPool}/{imageName}/{hostName}/{vmName} [post]
func (c *Controller) MirrorImageScheduleSetup(ctx *gin.Context) {
	//var dat model.MirrorSetup
	var dat = struct {
		Message string
	}{}

	mirrorPool := ctx.Param("mirrorPool")
	imageName := ctx.Param("imageName")
	hostName := ctx.Param("hostName")
	vmName := ctx.Param("vmName")
	volType, _ := ctx.GetPostForm("volType")

	message, err := mirror.ImagePreSetup(mirrorPool, imageName)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	message, err = mirror.ImageSetup(mirrorPool, imageName)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	if volType == "ROOT" {
		interval, err := mirror.ImageMetaGetInterval()
		if err != nil {
			mirror.ImageDeleteSchedule(mirrorPool, imageName)
			mirror.ImagePreDelete(mirrorPool, imageName)
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		_, err = mirror.ImageConfigSchedule(mirrorPool, imageName, hostName, vmName, interval)
		if err != nil {
			mirror.ImageDeleteSchedule(mirrorPool, imageName)
			mirror.ImagePreDelete(mirrorPool, imageName)
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	dat.Message = message
	ctx.IndentedJSON(http.StatusOK, dat)

}

// MirrorImageSnap godoc
//
//		@Summary		Take Image Mirroring Snapshot or Setup Image Mirroring Snapshot Schedule
//		@Description	Glue 의 이미지에 미러링 스냅샷을 생성하거나 스케줄을 설정합니다.
//		@param			mirrorPool	path		string	true	"Pool Name for Mirroring"
//		@param			vmName   	path		string	true	"VM Name for Mirroring"
//	 	@param          hostName    formData    string  false   "Host Name for Mirroring VMe"
//	 	@param          imageName   formData    string  false   "Image Name for Mirroring (Schedule)"
//	 	@param          imageList   formData    string  false   "Image List for Mirroring (Manual)"
//		@Tags			Mirror
//		@Accept			x-www-form-urlencoded
//		@Produce		json
//		@Success		200	{object}	model.ImageMirror
//		@Failure		400	{object}	httputil.HTTP400BadRequest
//		@Failure		404	{object}	httputil.HTTP404NotFound
//		@Failure		500	{object}	httputil.HTTP500InternalServerError
//		@Router			/api/v1/mirror/image/snapshot/{mirrorPool}/{vmName} [post]
func (c *Controller) MirrorImageSnap(ctx *gin.Context) {

	var dat = struct {
		Message string
	}{}

	mirrorPool := ctx.Param("mirrorPool")
	vmName := ctx.Param("vmName")
	hostName, _ := ctx.GetPostForm("hostName")
	imageName, _ := ctx.GetPostForm("imageName")
	imageList, _ := ctx.GetPostForm("imageList")

	// 수동 스냅샷 생성
	if imageList != "" {
		volList := strings.Split(imageList, ",")
		message, _ := mirror.ImageMirroringSnap(mirrorPool, hostName, vmName, volList)
		dat.Message = message
	}

	// 스냅샷 스케줄 설정
	if imageName != "" {
		interval, err := mirror.ImageMetaGetInterval()
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		message, err := mirror.ImageConfigSchedule(mirrorPool, imageName, hostName, vmName, interval)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		dat.Message = message
	}

	ctx.IndentedJSON(http.StatusOK, dat)
}

// MirrorImageUpdate godoc
//
//	@Summary		Put Image Mirroring
//	@Description	Glue 의 이미지에 미러링의 설정을 변경합니다.
//	@param			mirrorPool	path		string	true	"Pool Name for Mirroring"
//	@param			imageName	path		string	true	"Image Name for Mirroring"
//	@param			interval	formData	string	true	"Interval of image snapshot"
//	@param			startTime	formData	string	false	"Starttime of image snapshot"
//	@param			imageRegion	formData	string	false	"Current Image region"
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.ImageMirror
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror/image/{mirrorPool}/{imageName} [put]
// func (c *Controller) MirrorImageUpdate(ctx *gin.Context) {
// 	//var dat model.MirrorSetup
// 	var dat = struct {
// 		Message  string
// 		schedule []model.MirrorImageItem
// 	}{}

// 	mirrorPool := ctx.Param("mirrorPool")
// 	imageName := ctx.Param("imageName")
// 	interval, _ := ctx.GetPostForm("interval")
// 	startTime, _ := ctx.GetPostForm("startTime")
// 	imageRegion, _ := ctx.GetPostForm("imageRegion")

// 	MirroredImage, err := mirror.ImageList()
// 	if err != nil {
// 		utils.FancyHandleError(err)
// 		httputil.NewError(ctx, http.StatusInternalServerError, err)
// 		return
// 	}
// 	for _, image := range MirroredImage.Local {
// 		if image.Image == imageName {
// 			if len(image.Items) > 0 {
// 				dat.schedule = image.Items
// 			}
// 		}
// 	}
// 	for _, image := range MirroredImage.Remote {
// 		if image.Image == imageName {
// 			if len(image.Items) > 0 {
// 				dat.schedule = image.Items
// 			}
// 		}
// 	}

// 	if imageRegion == "remote" {
// 		message, err := mirror.ImageRemoteUpdate(mirrorPool, imageName, interval, startTime)
// 		if err != nil {
// 			utils.FancyHandleError(err)
// 			httputil.NewError(ctx, http.StatusInternalServerError, err)
// 			return
// 		}
// 		dat.Message = message
// 	} else {
// 		message, err := mirror.ImageUpdate(mirrorPool, imageName, interval, startTime, dat.schedule)
// 		if err != nil {
// 			utils.FancyHandleError(err)
// 			httputil.NewError(ctx, http.StatusInternalServerError, err)
// 			return
// 		}
// 		dat.Message = message
// 	}

// 	ctx.IndentedJSON(http.StatusOK, dat)
// }

// MirrorImageParentInfo godoc
//
//	@Summary		Show Mirroring Image Parent Info
//	@Description	Glue 의 이미지에 미러링 정보를 확인합니다.
//	@param			mirrorPool	path	string	true	"Pool Name for Mirroring"
//	@param			imageName	path	string	true	"Image Name for Mirroring"
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	ImageInfo
//	@Failure		400	{object}	HTTP400BadRequest
//	@Failure		404	{object}	HTTP404NotFound
//	@Failure		500	{object}	HTTP500InternalServerError
//	@Router			/api/v1/mirror/image/info/{mirrorPool}/{imageName} [get]
func (c *Controller) MirrorImageParentInfo(ctx *gin.Context) {

	var dat model.ImageInfo

	mirrorPool := ctx.Param("mirrorPool")
	imageName := ctx.Param("imageName")

	dat, err := mirror.ImageInfo(mirrorPool, imageName)
	if err != nil {
		utils.FancyHandleError(err)
		print(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, dat)
}

// MirrorImageStatus godoc
//
//	@Summary		Show Mirroring Image Status
//	@Description	Glue 의 이미지에 미러링상태를 확인합니다.
//	@param			mirrorPool	path	string	true	"Pool Name for Mirroring"
//	@param			imageName	path	string	true	"Image Name for Mirroring"
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	ImageStatus
//	@Failure		400	{object}	HTTP400BadRequest
//	@Failure		404	{object}	HTTP404NotFound
//	@Failure		500	{object}	HTTP500InternalServerError
//	@Router			/api/v1/mirror/image/status/{mirrorPool}/{imageName} [get]
func (c *Controller) MirrorImageStatus(ctx *gin.Context) {

	var dat model.ImageStatus

	mirrorPool := ctx.Param("mirrorPool")
	imageName := ctx.Param("imageName")

	dat, err := mirror.ImageStatus(mirrorPool, imageName)
	if err != nil {
		utils.FancyHandleError(err)
		print(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	print(dat.Description)
	ctx.IndentedJSON(http.StatusOK, dat)
}

// MirrorImagePromote godoc
//
//	@Summary		Promote Image Mirroring
//	@Description	Glue 의 이미지를 Promote 합니다.
//	@param			mirrorPool	path	string	true	"Pool Name for Mirroring"
//	@param			imageName	path	string	true	"Image Name for Mirroring"
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.ImageStatus
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror/image/promote/{mirrorPool}/{imageName} [post]
func (c *Controller) MirrorImagePromote(ctx *gin.Context) {

	var dat = struct {
		Message string
	}{}

	mirrorPool := ctx.Param("mirrorPool")
	imageName := ctx.Param("imageName")

	message, err := mirror.ImagePromote(mirrorPool, imageName)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	dat.Message = message
	ctx.IndentedJSON(http.StatusOK, dat)
}

// MirrorImagePromotePeer godoc
//
//	@Summary		Peer Promote Image Mirroring
//	@Description	Peer Glue 의 이미지를 Promote 합니다.
//	@param			mirrorPool	path	string	true	"Pool Name for Mirroring"
//	@param			imageName	path	string	true	"Image Name for Mirroring"
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.ImageStatus
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror/image/promote/peer/{mirrorPool}/{imageName} [post]
func (c *Controller) MirrorImagePromotePeer(ctx *gin.Context) {

	var dat = struct {
		Message string
	}{}

	mirrorPool := ctx.Param("mirrorPool")
	imageName := ctx.Param("imageName")

	message, err := mirror.RemoteImagePromote(mirrorPool, imageName)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	dat.Message = message
	ctx.IndentedJSON(http.StatusOK, dat)
}

// MirrorImageDemote godoc
//
//	@Summary		Demote Image Mirroring
//	@Description	Glue 의 이미지를 Demote 합니다.
//	@param			mirrorPool	path	string	true	"Pool Name for Mirroring"
//	@param			imageName	path	string	true	"Image Name for Mirroring"
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.ImageStatus
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror/image/demote/{mirrorPool}/{imageName} [delete]
func (c *Controller) MirrorImageDemote(ctx *gin.Context) {

	var dat = struct {
		Message string
	}{}

	mirrorPool := ctx.Param("mirrorPool")
	imageName := ctx.Param("imageName")

	message, err := mirror.ImageDemote(mirrorPool, imageName)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	dat.Message = message
	ctx.IndentedJSON(http.StatusOK, dat)
}

// MirrorImageDemotePeer godoc
//
//	@Summary		Peer Demote Image Mirroring
//	@Description	Peer Glue 의 이미지를 Demote 합니다.
//	@param			mirrorPool	path	string	true	"Pool Name for Mirroring"
//	@param			imageName	path	string	true	"Image Name for Mirroring"
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.ImageStatus
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror/image/demote/{mirrorPool}/{imageName} [delete]
func (c *Controller) MirrorImageDemotePeer(ctx *gin.Context) {

	var dat = struct {
		Message string
	}{}

	mirrorPool := ctx.Param("mirrorPool")
	imageName := ctx.Param("imageName")

	message, err := mirror.RemoteImageDemote(mirrorPool, imageName)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	dat.Message = message
	ctx.IndentedJSON(http.StatusOK, dat)
}

// MirrorImageResync godoc
//
//	@Summary		Resync Image Mirroring
//	@Description	Glue 의 이미지를 resync 합니다.
//	@param			mirrorPool	path	string	true	"Pool Name for Mirroring"
//	@param			imageName	path	string	true	"Image Name for Mirroring"
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.ImageStatus
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror/image/resync/{mirrorPool}/{imageName} [put]
func (c *Controller) MirrorImageResync(ctx *gin.Context) {

	var dat = struct {
		Message string
	}{}

	mirrorPool := ctx.Param("mirrorPool")
	imageName := ctx.Param("imageName")

	message, err := mirror.ImageResync(mirrorPool, imageName)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	dat.Message = message
	ctx.IndentedJSON(http.StatusOK, dat)
}

// MirrorImageResyncPeer godoc
//
//	@Summary		Peer Resync Image Mirroring
//	@Description	Peer Glue 의 이미지를 resync 합니다.
//	@param			mirrorPool	path	string	true	"Pool Name for Mirroring"
//	@param			imageName	path	string	true	"Image Name for Mirroring"
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.ImageStatus
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror/image/resync/peer/{mirrorPool}/{imageName} [put]
func (c *Controller) MirrorImageResyncPeer(ctx *gin.Context) {

	var dat = struct {
		Message string
	}{}

	mirrorPool := ctx.Param("mirrorPool")
	imageName := ctx.Param("imageName")

	message, err := mirror.RemoteImageResync(mirrorPool, imageName)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	dat.Message = message
	ctx.IndentedJSON(http.StatusOK, dat)
}

// MirrorPoolEnable godoc
//
//	@Summary		Enable Mirroring
//	@Description	Glue 의 미러링 클러스터를 활성화합니다.
//	@param			host			formData	string	true	"Remote Cluster Host Address"
//	@param			privateKeyFile	formData	file	true	"Remote Cluster PrivateKey"
//	@param			mirrorPool		formData	string	true	"Pool Name for Mirroring"
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.ImageStatus
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror/pool/{mirrorPool} [POST]
func (c *Controller) MirrorPoolEnable(ctx *gin.Context) {

	var dat model.MirrorSetup

	dat.LocalClusterName, _ = ctx.GetPostForm("localClusterName")
	dat.RemoteClusterName, _ = ctx.GetPostForm("remoteClusterName")
	dat.Host, _ = ctx.GetPostForm("host")
	dat.MirrorPool, _ = ctx.GetPostForm("mirrorPool")
	file, _ := ctx.FormFile("privateKeyFile")
	privkey, err := os.CreateTemp("", "id_rsa-")
	defer privkey.Close()
	defer os.Remove(privkey.Name())
	privkeyname := privkey.Name()

	// Upload the file to specific dst.
	err = ctx.SaveUploadedFile(file, privkeyname)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	EncodedLocalToken, EncodedRemoteToken, err := mirror.EnableMirror(dat, privkeyname)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	dat.LocalToken = EncodedLocalToken
	dat.RemoteToken = EncodedRemoteToken

	ctx.IndentedJSON(http.StatusOK, dat)

}

// MirrorPoolDisable godoc
//
//	@Summary		Disable Mirroring
//	@Description	Glue 의 미러링 클러스터를 비활성화합니다.
//	@param			host			formData	string	true	"Remote Cluster Host Address"
//	@param			privateKeyFile	formData	file	true	"Remote Cluster PrivateKey"
//	@param			mirrorPool		formData	string	true	"Pool Name for Mirroring"
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.ImageStatus
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror/pool/{mirrorPool} [delete]
func (c *Controller) MirrorPoolDisable(ctx *gin.Context) {

	var dat model.MirrorSetup
	var EncodedLocalToken string
	var EncodedRemoteToken string
	var stdout []byte

	var out strings.Builder
	var output string
	dat.Host, _ = ctx.GetPostForm("host")
	dat.MirrorPool, _ = ctx.GetPostForm("mirrorPool")
	file, _ := ctx.FormFile("privateKeyFile")

	privkey, err := os.CreateTemp("", "id_rsa-")
	defer privkey.Close()
	defer os.Remove(privkey.Name())
	privkeyname := privkey.Name()

	// Upload the file to specific dst.
	err = ctx.SaveUploadedFile(file, privkeyname)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Get Mirroring Images
	MirroredImage, err := mirror.ImageList(dat.MirrorPool)
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	for _, image := range MirroredImage.Images {
		_, errt := mirror.ImageDeleteSchedule(dat.MirrorPool, image.Name)
		if errt != nil {
			err = errors.Join(err, errt)
		}
		output, errt = mirror.ImagePreDelete(dat.MirrorPool, image.Name)
		if errt != nil {
			if output != "Success" {
				err = errors.Join(err, errt)
			}
		}
	}
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	//remote local peer
	mirrorStatus, err := mirror.GetConfigure()

	if len(mirrorStatus.Peers) > 0 {
		peerUUID := mirrorStatus.Peers[0].Uuid
		cmd := exec.Command("rbd", "mirror", "pool", "peer", "remove", "--pool", dat.MirrorPool, peerUUID)
		stdout, err = cmd.CombinedOutput()
		println("out: " + string(stdout))
		println("err: " + out.String())
		// if err != nil || (out.String() != "" && out.String() != "rbd: mirroring is already configured for image mode") {
		if err != nil {
			cmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// Mirror Disable
	if mirrorStatus.Mode != "disabled" {
		cmd := exec.Command("rbd", "mirror", "pool", "disable")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			cmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	//remote local peer
	out.Reset()
	client, err := utils.ConnectSSH(dat.Host, privkeyname)
	defer func(client *goph.Client) {
		err := client.Close()
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}(client)
	remoteMirrorStatus, err := mirror.GetRemoteConfigure(client)

	// Mirror Peer Remove
	if len(remoteMirrorStatus.Peers) > 0 {
		peerUUID := remoteMirrorStatus.Peers[0].Uuid
		sshcmd, err := client.Command("rbd", "mirror", "pool", "peer", "remove", "--pool", dat.MirrorPool, peerUUID)
		if err != nil {
			sshcmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			return
		}
		stdout, err = sshcmd.CombinedOutput()
		if err != nil {
			sshcmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// Mirror Disable
	if remoteMirrorStatus.Mode != "disabled" {
		sshcmd, err := client.Command("rbd", "mirror", "pool", "disable")
		if err != nil {
			sshcmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			return
		}
		stdout, err = sshcmd.CombinedOutput()
		if err != nil {
			sshcmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	dat.LocalToken = EncodedLocalToken
	dat.RemoteToken = EncodedRemoteToken
	ctx.IndentedJSON(http.StatusOK, dat)
}

// MirrorDeleteGarbage godoc
//
//	@Summary		Delete Mirroring Cluster Garbage
//	@Description	Glue 의 미러링 클러스터 가비지를 제거합니다.
//	@Tags			Mirror
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	controller.Message
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror/garbage [delete]
func (c *Controller) MirrorDeleteGarbage(ctx *gin.Context) {

	var stdout []byte
	var out strings.Builder

	mirrorPool := ctx.Param("mirrorPool")
	mirrorStatus, err := mirror.GetConfigure()

	// Mirror Peer Remove
	if len(mirrorStatus.Peers) > 0 {
		peerUUID := mirrorStatus.Peers[0].Uuid
		cmd := exec.Command("rbd", "mirror", "pool", "peer", "remove", "--pool", mirrorPool, peerUUID)
		stdout, err = cmd.CombinedOutput()
		println("out: " + string(stdout))
		if err != nil {
			cmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		cmd = exec.Command("ceph", "auth", "del", "client.rbd-mirror-peer")
		stdout, err = cmd.CombinedOutput()
		println("out: " + string(stdout))
		println("err: " + out.String())
		// if err != nil || (out.String() != "" && out.String() != "rbd: mirroring is already configured for image mode") {
		if err != nil {
			cmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// Mirror Disable
	if mirrorStatus.Mode != "disabled" {
		cmd := exec.Command("rbd", "mirror", "pool", "disable")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			cmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// Mirror Daemon Destroy
	cmd := exec.Command("ceph", "orch", "rm", "rbd-mirror")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		cmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	// DR Mirror Image Destroy
	cmd = exec.Command("rbd", "rm", "rbd/MOLD-DR")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		cmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	// DR mold conf reset
	err = mirror.ConfigMold("moldUrl", "moldApiKey", "moldSecretKey")
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, string(stdout))
}
