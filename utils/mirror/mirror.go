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
		err = errors.Join(err, errors.New(out.String()))
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
	cmd := exec.Command("rbd", "mirror", "pool", "enable", "--site-name", dat.LocalClusterName, "-p", dat.MirrorPool, "image")
	cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil || (out.String() != "" && out.String() != "rbd: mirroring is already configured for image mode") {
		utils.FancyHandleError(err)
		return
	}

	// Mirror Daemon Deploy
	cmd = exec.Command("ceph", "orch", "apply", "rbd-mirror")
	cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		utils.FancyHandleError(err)
		return
	}

	// Mirror Bootstrap
	cmd = exec.Command("rbd", "mirror", "pool", "peer", "bootstrap", "create", "--site-name", dat.LocalClusterName, "-p", dat.MirrorPool)
	cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	DecodedLocalToken, err := base64.StdEncoding.DecodeString(string(stdout))
	if err != nil || out.String() != "" {
		utils.FancyHandleError(err)
		return
	}

	if err = json.Unmarshal(DecodedLocalToken, &LocalToken); err != nil {
		utils.FancyHandleError(err)
		return
	}
	//println("ceph", "auth", "caps", "client."+LocalToken.ClientId, "mgr", "'profile rbd'", "mon", "'profile rbd-mirror-peer'", "osd", "'profile rbd'")
	cmd = exec.Command("ceph", "auth", "caps", "client."+LocalToken.ClientId, "mgr", "profile rbd", "mon", "profile rbd-mirror-peer", "osd", "profile rbd")
	cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()

	if err != nil {
		utils.FancyHandleError(err)
		return
	}

	cmd = exec.Command("ceph", "auth", "get-key", "client."+LocalToken.ClientId, "--format", "json")
	cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()

	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	if err = json.Unmarshal(stdout, &LocalKey); err != nil {
		utils.FancyHandleError(err)
		return
	}

	//Generate Token
	LocalToken.Key = LocalKey.Key
	JsonLocalKey, err := json.Marshal(LocalToken)
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	EncodedLocalToken = base64.StdEncoding.EncodeToString(JsonLocalKey)
	localTokenFile, err := os.CreateTemp("", "localtoken")

	defer localTokenFile.Close()
	defer os.Remove(localTokenFile.Name())
	localTokenFile.WriteString(EncodedLocalToken)

	//  For Remote
	client, err := utils.ConnectSSH(dat.Host, privkeyname)
	utils.FancyHandleError(err)
	if err != nil {
		utils.FancyHandleError(err)
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
		return
	}
	sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	//println("out: " + string(stdout))
	//println("err: " + out.String())
	if err != nil {
		utils.FancyHandleError(err)
		return
	} else if out.String() != "" && out.String() == "rbd: mirroring is already configured for image mode" {
		err = errors.New(out.String())
		utils.FancyHandleError(err)
		return
	}

	// Mirror Daemon Deploy
	sshcmd, err = client.Command("ceph", "orch", "apply", "rbd-mirror")
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		utils.FancyHandleError(err)
		return
	}

	// Mirror Bootstrap
	sshcmd, err = client.Command("rbd", "mirror", "pool", "peer", "bootstrap", "create", "--site-name", dat.RemoteClusterName, "-p", dat.MirrorPool)
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	//println("out: " + string(stdout))
	//println("err: " + out.String())
	if err != nil || out.String() != "" {
		utils.FancyHandleError(err)
		return
	}
	DecodedRemoteoken, err := base64.StdEncoding.DecodeString(string(stdout))
	if err != nil {
		utils.FancyHandleError(err)
		return
	}

	if err = json.Unmarshal(DecodedRemoteoken, &RemoteToken); err != nil {
		utils.FancyHandleError(err)
		return
	}
	//println("ceph", "auth", "caps", "client."+RemoteToken.ClientId, "mgr", "'profile rbd'", "mon", "'profile rbd-mirror-peer'", "osd", "'profile rbd'")
	sshcmd, err = client.Command("ceph", "auth", "caps", "client."+RemoteToken.ClientId, "mgr", "'profile rbd'", "mon", "'profile rbd-mirror-peer'", "osd", "'profile rbd'")
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}

	sshcmd, err = client.Command("ceph", "auth", "get-key", "client."+RemoteToken.ClientId, "--format", "json")
	if err != nil {
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	if err = json.Unmarshal(stdout, &LocalKey); err != nil {
		utils.FancyHandleError(err)
		return
	}

	//Generate Token
	RemoteToken.Key = RemoteKey.Key
	JsonRemoteKey, err := json.Marshal(RemoteToken)
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	EncodedRemoteToken = base64.StdEncoding.EncodeToString(JsonRemoteKey)

	// token import

	sshcmd, err = client.Command("echo", EncodedLocalToken, ">", remoteTokenFileName)
	if err != nil {
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}

	sshcmd, err = client.Command("rbd", "mirror", "pool", "info", "--pool", dat.MirrorPool, "--format", "json")
	if err != nil {
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	var remoteMirrorInfo model.MirrorInfo
	if err = json.Unmarshal(stdout, &remoteMirrorInfo); err != nil {
		utils.FancyHandleError(err)
		return
	}

	if len(remoteMirrorInfo.Peers) != 0 {
		for _, peer := range remoteMirrorInfo.Peers {
			sshcmd, err = client.Command("rbd", "mirror", "pool", "peer", "remove", "--pool", dat.MirrorPool, peer.Uuid)
			if err != nil {
				err = errors.Join(err, errors.New(out.String()))
				utils.FancyHandleError(err)
				return
			}
			sshcmd.Stderr = &out
			stdout, err = sshcmd.CombinedOutput()
			if err != nil {
				err = errors.Join(err, errors.New(out.String()))
				utils.FancyHandleError(err)
				return
			}
		}
	}

	sshcmd, err = client.Command("rbd", "mirror", "pool", "peer", "bootstrap", "import", "--pool", dat.MirrorPool, "--token-path", remoteTokenFileName)
	if err != nil {
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}

	out.Reset()
	println(EncodedRemoteToken)
	cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	println("out: " + string(stdout))
	println("err:" + string(out.String()))
	if err != nil {
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}

	cmd = exec.Command("rbd", "mirror", "pool", "info", "--pool", dat.MirrorPool, "--format", "json")
	cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}

	var localMirrorInfo model.MirrorInfo
	if err = json.Unmarshal(stdout, &localMirrorInfo); err != nil {
		utils.FancyHandleError(err)
		return
	}

	if len(localMirrorInfo.Peers) != 0 {
		for _, peer := range localMirrorInfo.Peers {
			cmd = exec.Command("rbd", "mirror", "pool", "peer", "remove", "--pool", dat.MirrorPool, peer.Uuid)
			cmd.Stderr = &out
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				err = errors.Join(err, errors.New(out.String()))
				utils.FancyHandleError(err)
				return
			}
		}
	}

	cmd = exec.Command("cat", localTokenFile.Name())
	cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	println(string(stdout))
	if err != nil {
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return

	}
	cmd = exec.Command("rbd", "mirror", "pool", "peer", "bootstrap", "import", "--pool", dat.MirrorPool, "--token-path", localTokenFile.Name())
	cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
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
