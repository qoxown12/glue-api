package controller

import (
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/mirror"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// MirrorImageList godoc
//
//	@Summary		Show List of Mirrored Image
//	@Description	미러링중인 이미지의 목록과 상태를 보여줍니다.
//	@Tags			Mirror
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.MirrorList
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror/image [get]
func (c *Controller) MirrorImageList(ctx *gin.Context) {
	dat, err := mirror.ImageList()
	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, dat)
}

// MirrorImageDelete godoc
//
//	@Summary		Delete Mirrored Image
//	@Description	이미지의 미러링을 비활성화 합니다.
//	@param			imageName			path		string				true	"imageName"
//	@param			pool		path		string				true	"pool"
//	@Tags			Mirror
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	controller.Message
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror/image/{pool}/{imagename} [delete]
func (c *Controller) MirrorImageDelete(ctx *gin.Context) {
	image := ctx.Param("imageName")
	pool := ctx.Param("pool")
	var output string

	output, err := mirror.ImageDelete(pool, image)

	if err != nil {
		utils.FancyHandleError(err)
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, Message{Message: output})
}

// MirrorStatus godoc
//
//	@Summary		Show Status of Mirror
//	@Description	Glue 의 미러링 상태를 보여줍니다.
//	@Tags			Mirror
//	@Accept			json
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
	dat.Debug = gin.IsDebugging()
	ctx.IndentedJSON(http.StatusOK, dat)
}

// MirrorSetup godoc
//
//	@Summary		Setup Mirroring Cluster
//	@Description	Glue 의 미러링 클러스터를 설정합니다..
//	@param			localClusterName	formData	string	true	"Local Cluster Name"
//	@param			remoteClusterName		formData	string	true	"Remote Cluster Name"
//	@param			host		formData	string	true	"Remote Cluster Host Address"
//	@param			privateKeyFile	formData	file	true	"Remote Cluster PrivateKey"
//	@param			mirrorPool		formData	string	true	"Pool Name for Mirroring"
//	@Tags			Mirror
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.MirrorSetup
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror [post]
func (c *Controller) MirrorSetup(ctx *gin.Context) {
	var dat model.MirrorSetup
	var LocalToken model.MirrorToken
	var RemoteToken model.MirrorToken
	var EncodedLocalToken string
	var EncodedRemoteToken string

	var out strings.Builder
	dat.LocalClusterName, _ = ctx.GetPostForm("localClusterName")
	dat.RemoteClusterName, _ = ctx.GetPostForm("remoteClusterName")
	dat.Host, _ = ctx.GetPostForm("host")
	dat.MirrorPool, _ = ctx.GetPostForm("mirrorPool")
	file, _ := ctx.FormFile("privateKeyFile")
	log.Println(file.Filename)
	privkey, err := os.CreateTemp("", "id_rsa-")
	defer privkey.Close()
	defer os.Remove(privkey.Name())
	privkeyname := privkey.Name()

	var LocalKey model.AuthKey
	var RemoteKey model.AuthKey
	localTokenFileName := "/tmp/localToken"
	remoteTokenFileName := "/tmp/remoteToken"

	// Upload the file to specific dst.
	err = ctx.SaveUploadedFile(file, privkeyname)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	var stdout []byte

	if gin.IsDebugging() != true {

		// Mirror Enable
		cmd := exec.Command("rbd", "mirror", "pool", "enable", "--site-name", dat.LocalClusterName, "-p", dat.MirrorPool, "image")
		cmd.Stderr = &out
		stdout, err = cmd.Output()
		println("out: " + string(stdout))
		println("err: " + out.String())
		if err != nil || (out.String() != "" && out.String() != "rbd: mirroring is already configured for image mode") {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		// Mirror Daemon Deploy
		cmd = exec.Command("ceph", "orch", "apply", "rbd-mirror")
		cmd.Stderr = &out
		stdout, err = cmd.Output()
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		// Mirror Bootstrap
		cmd = exec.Command("rbd", "mirror", "pool", "peer", "bootstrap", "create", "--site-name", dat.LocalClusterName, "-p", dat.MirrorPool)
		cmd.Stderr = &out
		stdout, err = cmd.Output()
		println("out: " + string(stdout))
		println("err: " + out.String())
		DecodedLocalToken, err := base64.StdEncoding.DecodeString(string(stdout))
		println(string(DecodedLocalToken))
		if err != nil || out.String() != "" {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		if err := json.Unmarshal(DecodedLocalToken, &LocalToken); err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		println("ceph", "auth", "caps", "client."+LocalToken.ClientId, "mgr", "'profile rbd'", "mon", "'profile rbd-mirror-peer'", "osd", "'profile rbd'")
		cmd = exec.Command("ceph", "auth", "caps", "client."+LocalToken.ClientId, "mgr", "profile rbd", "mon", "profile rbd-mirror-peer", "osd", "profile rbd")
		cmd.Stderr = &out
		stdout, err = cmd.Output()

		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		cmd = exec.Command("ceph", "auth", "get-key", "client."+LocalToken.ClientId, "--format", "json")
		cmd.Stderr = &out
		stdout, err = cmd.Output()

		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		if err := json.Unmarshal(stdout, &LocalKey); err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		//Generate Token
		LocalToken.Key = LocalKey.Key
		JsonLocalKey, err := json.Marshal(LocalToken)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		EncodedLocalToken = base64.StdEncoding.EncodeToString(JsonLocalKey)
		localTokenFile, err := os.CreateTemp("", "localtoken")
		defer localTokenFile.Close()
		defer os.Remove(localTokenFile.Name())
		localTokenFile.WriteString(EncodedLocalToken)

		//  For Remote
		client, err := ConnectSSH(dat.Host, privkeyname)
		utils.FancyHandleError(err)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		//// Defer closing the network connection.
		defer client.Close()
		//
		//// Execute your command.

		// Mirror Enable
		out.Reset()
		sshcmd, err := client.Command("rbd", "mirror", "pool", "enable", "--site-name", dat.RemoteClusterName, "-p", dat.MirrorPool, "image")
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		sshcmd.Stderr = &out
		stdout, err = sshcmd.Output()
		println("out: " + string(stdout))
		println("err: " + out.String())
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		} else if out.String() != "" && out.String() == "rbd: mirroring is already configured for image mode" {
			err = errors.New(out.String())
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		// Mirror Daemon Deploy
		sshcmd, err = client.Command("ceph", "orch", "apply", "rbd-mirror")
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		sshcmd.Stderr = &out
		stdout, err = sshcmd.Output()
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		// Mirror Bootstrap
		sshcmd, err = client.Command("rbd", "mirror", "pool", "peer", "bootstrap", "create", "--site-name", dat.RemoteClusterName, "-p", dat.MirrorPool)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		sshcmd.Stderr = &out
		stdout, err = sshcmd.Output()
		println("out: " + string(stdout))
		println("err: " + out.String())
		if err != nil || out.String() != "" {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		DecodedRemoteoken, err := base64.StdEncoding.DecodeString(string(stdout))
		println(string(DecodedRemoteoken))
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		if err := json.Unmarshal(DecodedRemoteoken, &RemoteToken); err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		println("ceph", "auth", "caps", "client."+RemoteToken.ClientId, "mgr", "'profile rbd'", "mon", "'profile rbd-mirror-peer'", "osd", "'profile rbd'")
		sshcmd, err = client.Command("ceph", "auth", "caps", "client."+RemoteToken.ClientId, "mgr", "'profile rbd'", "mon", "'profile rbd-mirror-peer'", "osd", "'profile rbd'")
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		sshcmd.Stderr = &out
		stdout, err = sshcmd.Output()
		if err != nil {
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		sshcmd, err = client.Command("ceph", "auth", "get-key", "client."+RemoteToken.ClientId, "--format", "json")
		if err != nil {
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		sshcmd.Stderr = &out
		stdout, err = sshcmd.Output()
		if err != nil {
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		if err := json.Unmarshal(stdout, &LocalKey); err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		//Generate Token
		RemoteToken.Key = RemoteKey.Key
		JsonRemoteKey, err := json.Marshal(RemoteToken)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		EncodedRemoteToken = base64.StdEncoding.EncodeToString(JsonRemoteKey)

		// token import

		sshcmd, err = client.Command("echo", EncodedLocalToken, ">", remoteTokenFileName)
		if err != nil {
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		sshcmd.Stderr = &out
		stdout, err = sshcmd.Output()
		if err != nil {
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		sshcmd, err = client.Command("rbd", "mirror", "pool", "peer", "bootstrap", "import", "--pool", dat.MirrorPool, "--token-path", remoteTokenFileName)
		if err != nil {
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		sshcmd.Stderr = &out
		stdout, err = sshcmd.Output()
		if err != nil {
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		out.Reset()
		cmd = exec.Command("echo", EncodedRemoteToken, ">", localTokenFileName)
		cmd.Stderr = &out
		stdout, err = cmd.Output()
		if err != nil {
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		cmd = exec.Command("rbd", "mirror", "pool", "peer", "bootstrap", "import", "--pool", dat.MirrorPool, "--token-path", localTokenFileName)
		if err != nil {
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		cmd.Stderr = &out
		stdout, err = cmd.Output()
		if err != nil {
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return

		}

	} else {

		stdout = []byte("{\n    \"summary\": {\n        \"health\": \"WARNING\",\n        \"daemon_health\": \"OK\",\n        \"image_health\": \"WARNING\",\n        \"states\": {\n            \"unknown\": 14\n        }\n    }\n}")
		client, err := ConnectSSH(dat.Host, privkeyname)
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		//// Defer closing the network connection.
		defer client.Close()
		//
		//// Execute your command.

		// Mirror Enable
		sshcmd, err := client.Command("ls", "-al")
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		sshcmd.Stderr = &out
		stdout, err = sshcmd.Output()
		println("out: " + string(stdout))
		println("err: " + out.String())
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}
	//if err := json.Unmarshal(stdout, &dat); err != nil {
	//	httputil.NewError(ctx, http.StatusInternalServerError, err)
	//	return
	//}
	println(EncodedLocalToken)
	// Print the output
	dat.Debug = gin.IsDebugging()
	dat.LocalToken = EncodedLocalToken
	dat.RemoteToken = EncodedRemoteToken
	ctx.IndentedJSON(http.StatusOK, dat)
}

// MirrorDelete godoc
//
//	@Summary		Delete Mirroring Cluster
//	@Description	Glue 의 미러링 클러스터를 제거합니다.
//	@param			host		formData	string	true	"Remote Cluster Host Address"
//	@param			privateKeyFile	formData	file	true	"Remote Cluster PrivateKey"
//	@param			mirrorPool		formData	string	true	"Pool Name for Mirroring"
//	@Tags			Mirror
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.MirrorSetup
//	@Failure		400	{object}	httputil.HTTP400BadRequest
//	@Failure		404	{object}	httputil.HTTP404NotFound
//	@Failure		500	{object}	httputil.HTTP500InternalServerError
//	@Router			/api/v1/mirror [delete]
func (c *Controller) MirrorDelete(ctx *gin.Context) {
	var dat model.MirrorSetup
	var EncodedLocalToken string
	var EncodedRemoteToken string
	var stdout []byte

	var out strings.Builder
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
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	if gin.IsDebugging() != true {

		// Get Mirroring Images
		MirroredImage, err := mirror.ImageList()
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		for _, image := range MirroredImage.Local {
			_, errt := mirror.ImageDelete(image.Pool, image.Image)
			if errt != nil {
				err = errors.Join(err, errt)
			}
		}
		if err != nil {
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		mirrorStatus, err := mirror.GetConfigure()

		if len(mirrorStatus.Peers) > 0 {
			peerUUID := mirrorStatus.Peers[0].Uuid
			cmd := exec.Command("rbd", "mirror", "pool", "peer", "remove", "--pool", dat.MirrorPool, peerUUID)
			cmd.Stderr = &out
			stdout, err = cmd.Output()
			println("out: " + string(stdout))
			println("err: " + out.String())
			if err != nil || (out.String() != "" && out.String() != "rbd: mirroring is already configured for image mode") {
				err = errors.Join(err, errors.New(out.String()))
				utils.FancyHandleError(err)
				httputil.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
		}
		// Mirror peer remove
		/*
			rbd mirror pool info rbd --all
			Mode: image
			Site Name: cluster1

			Peer Sites:

			UUID: 8afad32c-7e5e-4744-b391-cee7dd5ef2d9
			Name: cluster2
			Mirror UUID: c1b81c8b-8973-4b3c-ad3c-63e1420f60ba
			Direction: rx-tx
			Client: client.rbd-mirror-peer
			Mon Host: [v2:100.100.14.21:3300/0,v1:100.100.14.21:6789/0],[v2:100.100.14.22:3300/0,v1:100.100.14.22:6789/0],[v2:100.100.14.23:3300/0,v1:100.100.14.23:6789/0]
			Key: AQDspDBlfMiOAhAAdgJqKs0DOJ4wSg/3UuIo3A==
			[root@scvm11 ~]# rbd help mirror pool peer remove
			usage: rbd mirror pool peer remove [--pool <pool>]
			                                   <pool-name> <uuid>

			Remove a mirroring peer from a pool.

			Positional arguments
			  <pool-name>          pool name
			  <uuid>               peer uuid

			Optional arguments
			  -p [ --pool ] arg    pool name

		*/

		// Mirror Disable
		cmd := exec.Command("rbd", "mirror", "pool", "disable")
		cmd.Stderr = &out
		stdout, err = cmd.Output()
		println("out: " + string(stdout))
		println("err: " + out.String())
		if err != nil || (out.String() != "" && out.String() != "rbd: mirroring is already configured for image mode") {
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		// Mirror Daemon Destroy
		cmd = exec.Command("ceph", "orch", "rm", "rbd-mirror")
		cmd.Stderr = &out
		stdout, err = cmd.Output()
		if err != nil {
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// Print the output
	dat.Debug = gin.IsDebugging()
	dat.LocalToken = EncodedLocalToken
	dat.RemoteToken = EncodedRemoteToken
	ctx.IndentedJSON(http.StatusOK, dat)
}
