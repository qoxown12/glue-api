package mirror

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/melbahja/goph"
)

func IsConfigured() (configured bool, err error) {
	config, err := GetConfigure()
	if err != nil || config.Mode == "disabled" {
		return false, err
	} else {
		return true, nil
	}
}

func GetConfigure() (clusterConf model.MirrorConf, err error) {
	var stdout []byte
	//sOut := string(stdout)
	//lines := strings.Split(sOut, "\n")
	tfCluster, err := os.CreateTemp(os.TempDir(), "Glue-Cluster-")
	tfKey, err := os.CreateTemp(os.TempDir(), "Glue-Key-")

	cmd := exec.Command("rbd", "mirror", "pool", "info", "--all", "--format", "json", "--pretty-format")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		return clusterConf, err
	}

	if err = json.Unmarshal(stdout, &clusterConf); err != nil {

		return clusterConf, err
	}
	if clusterConf.Mode == "disabled" {
		err = errors.New("mirroring is disabled")
		return clusterConf, err
	}
	peer := clusterConf.Peers[0]
	strCluster := "[global]\n\tmon host = " + peer.MonHost + "\n"
	// print(strCluster)
	if _, err = tfCluster.WriteString(strCluster); err != nil {
		fmt.Println("Failed to write to temporary file", err)
		return clusterConf, err
	}
	if err = tfCluster.Close(); err != nil {
		return clusterConf, err
	}
	if _, err = tfKey.WriteString(peer.Key); err != nil {
		fmt.Println("Failed to write to temporary file", err)
		return clusterConf, err
	}
	if err = tfKey.Close(); err != nil {
		fmt.Println(err)
		return clusterConf, err
	}

	clusterConf.Name = peer.ClientName
	clusterConf.ClusterFileName = tfCluster.Name()
	clusterConf.KeyFileName = tfKey.Name()
	return clusterConf, nil
}

func GetRemoteConfigure(client *goph.Client) (clusterConf model.MirrorConf, err error) {
	var stdout []byte
	//sOut := string(stdout)
	//lines := strings.Split(sOut, "\n")
	tfCluster, err := os.CreateTemp(os.TempDir(), "Glue-Cluster-")
	tfKey, err := os.CreateTemp(os.TempDir(), "Glue-Key-")

	cmd, err := client.Command("rbd", "mirror", "pool", "info", "--all", "--format", "json", "--pretty-format")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		return clusterConf, err
	}

	if err = json.Unmarshal(stdout, &clusterConf); err != nil {

		return clusterConf, err
	}
	if clusterConf.Mode == "disabled" {
		err = errors.New("mirroring is disabled")
		return clusterConf, err
	}
	peer := clusterConf.Peers[0]
	strCluster := "[global]\n\tmon host = " + peer.MonHost + "\n"
	// print(strCluster)
	if _, err = tfCluster.WriteString(strCluster); err != nil {
		fmt.Println("Failed to write to temporary file", err)
		return clusterConf, err
	}
	if err = tfCluster.Close(); err != nil {
		return clusterConf, err
	}
	if _, err = tfKey.WriteString(peer.Key); err != nil {
		fmt.Println("Failed to write to temporary file", err)
		return clusterConf, err
	}
	if err = tfKey.Close(); err != nil {
		fmt.Println(err)
		return clusterConf, err
	}

	clusterConf.Name = peer.ClientName
	clusterConf.ClusterFileName = tfCluster.Name()
	clusterConf.KeyFileName = tfKey.Name()
	return clusterConf, nil
}

func ImageList() (MirrorList model.MirrorList, err error) {

	var stdRemote []byte
	var stdLocal []byte
	var Local []model.MirrorImage
	var Remote []model.MirrorImage

	mirrorConfig, err := GetConfigure()
	if err != nil {
		return
	}

	strRemoteStatus := exec.Command("rbd", "mirror", "snapshot", "schedule", "list", "-R", "--format", "json", "--pretty-format", "-c", mirrorConfig.ClusterName, "-K", mirrorConfig.KeyFileName, "-n", mirrorConfig.Name)
	stdRemote, err = strRemoteStatus.CombinedOutput()
	strLocalStatus := exec.Command("rbd", "mirror", "snapshot", "schedule", "list", "-R", "--format", "json", "--pretty-format")
	stdLocal, err = strLocalStatus.CombinedOutput()

	if err != nil {
		return
	}
	if err = json.Unmarshal(stdRemote, &Remote); err != nil {
		Remote = []model.MirrorImage{}
	}
	if err = json.Unmarshal(stdLocal, &Local); err != nil {
		Local = []model.MirrorImage{}
	}
	MirrorList.Local = Local
	MirrorList.Remote = Remote
	return MirrorList, err
}

func Status() (mirrorStatus model.MirrorStatus, err error) {
	var tmpdat struct {
		Summary model.MirrorStatus `json:"summary"`
	}
	var stdout []byte
	cmd := exec.Command("rbd", "mirror", "pool", "status", "--format", "json", "--pretty-format")
	var out strings.Builder
	//cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()

	if err != nil {
		if strings.Contains(string(stdout), "mirroring not enabled on the pool") {
			err = errors.New(string(stdout))
		} else  {
			cmd.Stderr = &out
			err = errors.Join(err, errors.New(out.String()))
		}
		utils.FancyHandleError(err)
		return
	}
	if err = json.Unmarshal(stdout, &tmpdat); err != nil {
		utils.FancyHandleError(err)
		return
	}
	mirrorStatus = tmpdat.Summary
	return
}

func ImageDelete(poolName string, imageName string) (output string, err error) {

	var stdRemove []byte

	strRemoveStatus := exec.Command("rbd", "mirror", "snapshot", "schedule", "rm", "--pool", poolName, "--image", imageName)
	stdRemove, err = strRemoveStatus.CombinedOutput()

	if err != nil {
		err = errors.New(string(stdRemove))
		utils.FancyHandleError(err)
		return
	}

	strRemovestatus := exec.Command("rbd", "mirror", "image", "disable", "--pool", poolName, "--image", imageName)
	stdRemove, err = strRemovestatus.CombinedOutput()

	if err != nil {
		err = errors.New(string(stdRemove))
		utils.FancyHandleError(err)
		return
	}

	output = string(stdRemove)
	return
}

func ImageSetup(poolName string, imageName string) (output string, err error) {

	var stdoutMirrorEnable []byte

	strMirrorEnableOutput := exec.Command("rbd", "mirror", "image", "enable", "--pool", poolName, "--image", imageName, "snapshot")
	stdoutMirrorEnable, err = strMirrorEnableOutput.CombinedOutput()
	if err != nil || string(stdoutMirrorEnable) != "Mirroring enabled\n" {
		err = errors.Join(err, errors.New(string(stdoutMirrorEnable)))
		utils.FancyHandleError(err)
		return
	}

	output = string(stdoutMirrorEnable)
	return
}

func ImageConfig(poolName string, imageName string, interval string, startTime string) (output string, err error) {

	var stdoutScheduleEnable []byte

	var strScheduleOutput *exec.Cmd
	if startTime == "" {
		strScheduleOutput = exec.Command("rbd", "mirror", "snapshot", "schedule", "add", "--pool", poolName, "--image", imageName, interval)
	} else {
		strScheduleOutput = exec.Command("rbd", "mirror", "snapshot", "schedule", "add", "--pool", poolName, "--image", imageName, interval, startTime)
	}
	stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()
	if err != nil {
		err = errors.Join(err, errors.New(string(stdoutScheduleEnable)))
		utils.FancyHandleError(err)
		return
	}

	output = string(stdoutScheduleEnable)
	return
}

func ImageStatus(poolName string, imageName string) (imageStatus model.ImageStatus, err error) {

	var stdoutScheduleEnable []byte

	strScheduleOutput := exec.Command("rbd", "mirror", "image", "status", "--pool", poolName, "--image", imageName, "--format", "json")
	stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()

	if err = json.Unmarshal(stdoutScheduleEnable, &imageStatus); err != nil {
		err = errors.Join(err, errors.New(string(stdoutScheduleEnable)))
		utils.FancyHandleError(err)
		return
	}
	if err != nil {
		err = errors.Join(err, errors.New(string(stdoutScheduleEnable)))
		utils.FancyHandleError(err)
		return
	}
	return
}

func ConfigMirror(dat model.MirrorSetup, privkeyname string) (EncodedLocalToken string, EncodedRemoteToken string, err error) {

	var LocalToken model.MirrorToken
	var RemoteToken model.MirrorToken
	var out strings.Builder
	var LocalKey model.AuthKey
	var RemoteKey model.AuthKey
	var stdout []byte

	remoteTokenFileName := "/tmp/remoteToken"

	// Mirror Enable
	println("ConfigMirror :: Mirror Enable")
	cmd := exec.Command("rbd", "mirror", "pool", "enable", "--site-name", dat.LocalClusterName, "-p", dat.MirrorPool, "image")
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		println("ConfigMirror :: Mirror Enable error : " + string(stdout))
		utils.FancyHandleError(err)
		return
	}

	// Mirror Daemon Deploy
	println("ConfigMirror :: Mirror Daemon Deploy")
	cmd = exec.Command("ceph", "orch", "apply", "rbd-mirror")
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		println("ConfigMirror :: Mirror Daemon Deploy error" + string(stdout))
		utils.FancyHandleError(err)
		return
	}

	// Mirror Bootstrap
	println("ConfigMirror :: Mirror Bootstrap")
	cmd = exec.Command("rbd", "mirror", "pool", "peer", "bootstrap", "create", "--site-name", dat.LocalClusterName, "-p", dat.MirrorPool)
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	DecodedLocalToken, err := base64.StdEncoding.DecodeString(string(stdout))
	if err != nil {
		println("ConfigMirror :: Mirror Bootstrap error" + string(stdout))
		utils.FancyHandleError(err)
		return
	}

	if err = json.Unmarshal(DecodedLocalToken, &LocalToken); err != nil {
		println("ConfigMirror :: Mirror Bootstrap error json unmarshal")
		utils.FancyHandleError(err)
		return
	}

	println("ConfigMirror :: Mirror auth caps rbd-mirror-peer")
	cmd = exec.Command("ceph", "auth", "caps", "client."+LocalToken.ClientId, "mgr", "profile rbd", "mon", "profile rbd-mirror-peer", "osd", "profile rbd")
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		println("ConfigMirror :: Mirror auth caps rbd-mirror-peer error" + string(stdout))
		utils.FancyHandleError(err)
		return
	}

	println("ConfigMirror :: Mirror auth get-key rbd-mirror-peer")
	cmd = exec.Command("ceph", "auth", "get-key", "client."+LocalToken.ClientId, "--format", "json")
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		println("ConfigMirror :: Mirror auth get-key rbd-mirror-peer error" + string(stdout))
		utils.FancyHandleError(err)
		return
	}
	if err = json.Unmarshal(stdout, &LocalKey); err != nil {
		println("ConfigMirror :: Mirror auth get-key rbd-mirror-peer error json unmarshal")
		utils.FancyHandleError(err)
		return
	}

	// Generate Token
	println("ConfigMirror :: Mirror Generate Token")
	LocalToken.Key = LocalKey.Key
	JsonLocalKey, err := json.Marshal(LocalToken)
	if err != nil {
		println("ConfigMirror :: Mirror Generate Token err json marshal")
		utils.FancyHandleError(err)
		return
	}
	EncodedLocalToken = base64.StdEncoding.EncodeToString(JsonLocalKey)
	localTokenFile, err := os.CreateTemp("", "localtoken")
	if err != nil {
		println("ConfigMirror :: Mirror Generate Token err os create temp")
		utils.FancyHandleError(err)
		return
	}

	// defer localTokenFile.Close()
	// defer os.Remove(localTokenFile.Name())
	localTokenFile.WriteString(EncodedLocalToken)

	// For Remote
	println("ConfigMirror :: Mirror For Remote Connect")
	client, err := utils.ConnectSSH(dat.Host, privkeyname)
	utils.FancyHandleError(err)
	if err != nil {
		println("ConfigMirror :: Mirror For Remote Connect err")
		utils.FancyHandleError(err)
		return
	}
	//// Defer closing the network connection.
	defer client.Close()
	//
	//// Execute your command.

	// Mirror Enable
	out.Reset()
	println("ConfigMirror :: Mirror For Remote Mirror Enable")
	sshcmd, err := client.Command("rbd", "mirror", "pool", "enable", "--site-name", dat.RemoteClusterName, "-p", dat.MirrorPool, "image")
	if err != nil {
		println("ConfigMirror :: Mirror For Remote Mirror Enable err ConnectSSH")
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	// println("out: " + string(stdout))
	// println("err: " + out.String())
	if err != nil {
		println("ConfigMirror :: Mirror For Remote Mirror Enable err" + string(stdout))
		utils.FancyHandleError(err)
		return
	}

	// Mirror Daemon Deploy
	println("ConfigMirror :: Mirror For Remote Mirror Daemon Deploy")
	sshcmd, err = client.Command("ceph", "orch", "apply", "rbd-mirror")
	if err != nil {
		println("ConfigMirror :: Mirror For Remote Mirror Daemon Deploy err ConnectSSH")
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		println("ConfigMirror :: Mirror For Remote Mirror Daemon Deploy err" + string(stdout))
		utils.FancyHandleError(err)
		return
	}

	// Mirror Bootstrap
	println("ConfigMirror :: Mirror For Remote Mirror Bootstrap")
	sshcmd, err = client.Command("rbd", "mirror", "pool", "peer", "bootstrap", "create", "--site-name", dat.RemoteClusterName, "-p", dat.MirrorPool)
	if err != nil {
		println("ConfigMirror :: Mirror For Remote Mirror Bootstrap err ConnectSSH")
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	//println("out: " + string(stdout))
	//println("err: " + out.String())
	if err != nil {
		println("ConfigMirror :: Mirror For Remote Mirror Bootstrap err" + string(stdout))
		utils.FancyHandleError(err)
		return
	}

	DecodedRemoteoken, err := base64.StdEncoding.DecodeString(string(stdout))
	if err != nil {
		println("ConfigMirror :: Mirror For Remote Mirror Bootstrap err decode remote token")
		utils.FancyHandleError(err)
		return
	}

	if err = json.Unmarshal(DecodedRemoteoken, &RemoteToken); err != nil {
		println("ConfigMirror :: Mirror For Remote Mirror Bootstrap err json unmarshal remote token")
		utils.FancyHandleError(err)
		return
	}

	println("ConfigMirror :: Mirror For Remote Mirror auth caps rbd-mirror-peer")
	sshcmd, err = client.Command("ceph", "auth", "caps", "client."+RemoteToken.ClientId, "mgr", "'profile rbd'", "mon", "'profile rbd-mirror-peer'", "osd", "'profile rbd'")
	if err != nil {
		println("ConfigMirror :: Mirror For Remote Mirror auth caps rbd-mirror-peer err connectSSH")
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		sshcmd.Stderr = &out
		println("ConfigMirror :: Mirror For Remote Mirror auth caps rbd-mirror-peer err" + out.String())
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}

	println("ConfigMirror :: Mirror For Remote Mirror auth get-key rbd-mirror-peer")
	sshcmd, err = client.Command("ceph", "auth", "get-key", "client."+RemoteToken.ClientId, "--format", "json")
	if err != nil {
		println("ConfigMirror :: Mirror For Remote Mirror auth get-key rbd-mirror-peer err connectSSH")
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		sshcmd.Stderr = &out
		println("ConfigMirror :: Mirror For Remote Mirror auth get-key rbd-mirror-peer err" + out.String())
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	if err = json.Unmarshal(stdout, &LocalKey); err != nil {
		println("ConfigMirror :: Mirror For Remote Mirror auth get-key rbd-mirror-peer err json unmarshal")
		utils.FancyHandleError(err)
		return
	}

	//Generate Token
	println("ConfigMirror :: Mirror For Remote Mirror Generate Token")
	RemoteToken.Key = RemoteKey.Key
	JsonRemoteKey, err := json.Marshal(RemoteToken)
	if err != nil {
		println("ConfigMirror :: Mirror For Remote Mirror Generate Token err json marshal")
		utils.FancyHandleError(err)
		return
	}
	EncodedRemoteToken = base64.StdEncoding.EncodeToString(JsonRemoteKey)
	remoteTokenFile, err := os.CreateTemp("", "remotetoken")
    if err != nil {
        println("ConfigMirror :: Mirror For Remote Mirror Generate Token err os create temp remote token")
        utils.FancyHandleError(err)
        return
    }
	remoteTokenFile.WriteString(EncodedRemoteToken)

	// token import
	println("ConfigMirror :: Mirror For Remote Mirror token import")
	sshcmd, err = client.Command("echo", EncodedLocalToken, ">", remoteTokenFileName)
	if err != nil {
		sshcmd.Stderr = &out
		println("ConfigMirror :: Mirror For Remote Mirror token import err connectSSH")
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		sshcmd.Stderr = &out
		println("ConfigMirror :: Mirror For Remote Mirror token import err" + out.String())
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}

	println("ConfigMirror :: Mirror For Remote Mirror pool info")
	sshcmd, err = client.Command("rbd", "mirror", "pool", "info", "--pool", dat.MirrorPool, "--format", "json")
	if err != nil {
		println("ConfigMirror :: Mirror For Remote Mirror pool info err connectSSH")
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		sshcmd.Stderr = &out
		println("ConfigMirror :: Mirror For Remote Mirror pool info err" + out.String())
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	var remoteMirrorInfo model.MirrorInfo
	if err = json.Unmarshal(stdout, &remoteMirrorInfo); err != nil {
		println("ConfigMirror :: Mirror For Remote Mirror pool info err json unmarshal")
		utils.FancyHandleError(err)
		return
	}

	if len(remoteMirrorInfo.Peers) != 0 {
		for _, peer := range remoteMirrorInfo.Peers {
			println("ConfigMirror :: Mirror For Remote Mirror pool peer remove")
			sshcmd, err = client.Command("rbd", "mirror", "pool", "peer", "remove", "--pool", dat.MirrorPool, peer.Uuid)
			if err != nil {
				println("ConfigMirror :: Mirror For Remote Mirror pool peer remove err connectSSH")
				sshcmd.Stderr = &out
				err = errors.Join(err, errors.New(out.String()))
				utils.FancyHandleError(err)
				return
			}
			// sshcmd.Stderr = &out
			stdout, err = sshcmd.CombinedOutput()
			if err != nil {
				sshcmd.Stderr = &out
				println("ConfigMirror :: Mirror For Remote Mirror pool peer remove err" + out.String())
				err = errors.Join(err, errors.New(out.String()))
				utils.FancyHandleError(err)
				return
			}
		}
	}

	println("ConfigMirror :: Mirror For Remote Mirror pool peer bootstrap import")
	sshcmd, err = client.Command("rbd", "mirror", "pool", "peer", "bootstrap", "import", "--pool", dat.MirrorPool, "--token-path", remoteTokenFileName)
	if err != nil {
		println("ConfigMirror :: Mirror For Remote Mirror pool peer bootstrap import err connectSSH")
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		sshcmd.Stderr = &out
		println("ConfigMirror :: Mirror For Remote Mirror pool peer bootstrap import err" + out.String())
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}

	out.Reset()
	// println(EncodedRemoteToken)
	// cmd.Stderr = &out
	// stdout, err = cmd.CombinedOutput()
	// println("out: " + string(stdout))
	// println("err:" + string(out.String()))
	// if err != nil {
	// 	println("out.Reset() err")
	// 	cmd.Stderr = &out
	// 	err = errors.Join(err, errors.New(out.String()))
	// 	utils.FancyHandleError(err)
	// 	return
	// }

	println("ConfigMirror :: Mirror pool info")
	cmd = exec.Command("rbd", "mirror", "pool", "info", "--pool", dat.MirrorPool, "--format", "json")
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		cmd.Stderr = &out
		println("ConfigMirror :: Mirror pool info err" + out.String())
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}

	var localMirrorInfo model.MirrorInfo
	if err = json.Unmarshal(stdout, &localMirrorInfo); err != nil {
		println("ConfigMirror :: Mirror pool info err json unmarshal")
		utils.FancyHandleError(err)
		return
	}

	if len(localMirrorInfo.Peers) != 0 {
		for _, peer := range localMirrorInfo.Peers {
			println("ConfigMirror :: Mirror pool peer remove")
			cmd = exec.Command("rbd", "mirror", "pool", "peer", "remove", "--pool", dat.MirrorPool, peer.Uuid)
			// cmd.Stderr = &out
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				cmd.Stderr = &out
				println("ConfigMirror :: Mirror pool peer remove err" + out.String())
				err = errors.Join(err, errors.New(out.String()))
				utils.FancyHandleError(err)
				return
			}
		}
	}

	println("ConfigMirror :: Mirror pool peer bootstrap import")
	cmd = exec.Command("rbd", "mirror", "pool", "peer", "bootstrap", "import", "--pool", dat.MirrorPool, "--token-path", remoteTokenFile.Name())
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		cmd.Stderr = &out
		println("ConfigMirror :: Mirror pool peer bootstrap import err" + out.String())
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	return
}

func ImagePromote(poolName string, imageName string) (imageStatus model.ImageStatus, err error) {

	var stdoutScheduleEnable []byte
	strScheduleOutput := exec.Command("rbd", "mirror", "image", "promote", "--pool", poolName, "--image", imageName)
	stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()

	if !strings.Contains(string(stdoutScheduleEnable), "Image promoted") {
		err = errors.Join(err, errors.New(string(stdoutScheduleEnable)))
	}
	if err != nil {
		err = errors.Join(err, errors.New(string(stdoutScheduleEnable)))
		utils.FancyHandleError(err)
		return
	}
	return
}

func ImageDemote(poolName string, imageName string) (imageStatus model.ImageStatus, err error) {

	var stdoutScheduleEnable []byte

	strScheduleOutput := exec.Command("rbd", "mirror", "image", "demote", "--pool", poolName, "--image", imageName)
	stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()
	if !strings.Contains(string(stdoutScheduleEnable), "Image demoted") {
		err = errors.Join(err, errors.New(string(stdoutScheduleEnable)))
	}
	if err != nil {
		utils.FancyHandleError(err)
	}
	return
}

func RemoteImagePromote(poolName string, imageName string) (imageStatus model.ImageStatus, err error) {

	var stdoutScheduleEnable []byte
	conf, err := GetConfigure()
	strScheduleOutput := exec.Command("rbd", "-c", conf.ClusterFileName, "--cluster", conf.ClusterName, "--name", conf.Peers[0].ClientName, "--keyfile", conf.KeyFileName, "mirror", "image", "promote", "--pool", poolName, "--image", imageName)
	stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()
	if !strings.Contains(string(stdoutScheduleEnable), "Image promoted") {
		err = errors.Join(err, errors.New(string(stdoutScheduleEnable)))
	}
	if err != nil {
		utils.FancyHandleError(err)
	}
	return
}

func RemoteImageDemote(poolName string, imageName string) (imageStatus model.ImageStatus, err error) {

	var stdoutScheduleEnable []byte
	conf, err := GetConfigure()
	strScheduleOutput := exec.Command("rbd", "-c", conf.ClusterFileName, "--cluster", conf.ClusterName, "--name", conf.Peers[0].ClientName, "--keyfile", conf.KeyFileName, "mirror", "image", "demote", "--pool", poolName, "--image", imageName)
	stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()
	if !strings.Contains(string(stdoutScheduleEnable), "Image demoted") {
		err = errors.Join(err, errors.New(string(stdoutScheduleEnable)))
	}
	if err != nil {
		utils.FancyHandleError(err)
	}
	return
}
