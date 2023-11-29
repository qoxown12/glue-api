package mirror

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/melbahja/goph"
	"os"
	"os/exec"
	"strings"
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
	if gin.IsDebugging() != true {

		cmd := exec.Command("rbd", "mirror", "pool", "info", "--all", "--format", "json", "--pretty-format")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			return clusterConf, err
		}

	} else {
		stdout = []byte("{\n    \"mode\": \"image\",\n    \"site_name\": \"cluster1\",\n    \"peers\": [\n        {\n            \"uuid\": \"c5f36d07-a69e-45b5-a6ae-7e7a5f1266f0\",\n            \"direction\": \"rx-tx\",\n            \"site_name\": \"cluster2\",\n            \"mirror_uuid\": \"d7576286-c14c-47cf-8b87-58c4ec7e08dc\",\n            \"client_name\": \"client.rbd-mirror-peer\",\n            \"key\": \"AQAxDx5lljPfGhAAQhP1voVx5Dogn3f+nzYM8A==\",\n            \"mon_host\": \"[v2:100.100.1.24:3300/0,v1:100.100.1.24:6789/0],[v2:100.100.1.25:3300/0,v1:100.100.1.25:6789/0],[v2:100.100.1.26:3300/0,v1:100.100.1.26:6789/0]\"\n        }\n    ]\n}")
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
	if gin.IsDebugging() != true {

		cmd, err := client.Command("rbd", "mirror", "pool", "info", "--all", "--format", "json", "--pretty-format")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			return clusterConf, err
		}

	} else {
		stdout = []byte("{\n    \"mode\": \"image\",\n    \"site_name\": \"cluster1\",\n    \"peers\": [\n        {\n            \"uuid\": \"c5f36d07-a69e-45b5-a6ae-7e7a5f1266f0\",\n            \"direction\": \"rx-tx\",\n            \"site_name\": \"cluster2\",\n            \"mirror_uuid\": \"d7576286-c14c-47cf-8b87-58c4ec7e08dc\",\n            \"client_name\": \"client.rbd-mirror-peer\",\n            \"key\": \"AQAxDx5lljPfGhAAQhP1voVx5Dogn3f+nzYM8A==\",\n            \"mon_host\": \"[v2:100.100.1.24:3300/0,v1:100.100.1.24:6789/0],[v2:100.100.1.25:3300/0,v1:100.100.1.25:6789/0],[v2:100.100.1.26:3300/0,v1:100.100.1.26:6789/0]\"\n        }\n    ]\n}")
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
	if gin.IsDebugging() != true {

		strRemoteStatus := exec.Command("rbd", "mirror", "snapshot", "schedule", "list", "-R", "--format", "json", "--pretty-format", "-c", mirrorConfig.ClusterName, "-K", mirrorConfig.KeyFileName, "-n", mirrorConfig.Name)
		stdRemote, err = strRemoteStatus.CombinedOutput()
		strLocalStatus := exec.Command("rbd", "mirror", "snapshot", "schedule", "list", "-R", "--format", "json", "--pretty-format")
		stdLocal, err = strLocalStatus.CombinedOutput()

		/*
			   [
			    {
			        "pool": "rbd",
			        "namespace": "",
			        "image": "4f1ff9d5-7cfd-4d5a-97fd-ba3bb6faa17b",
			        "items": [
			            {
			                "interval": "10m",
			                "start_time": ""
			            }
			        ]
			    },
			    {
			        "pool": "rbd",
			        "namespace": "",
			        "image": "5f44786f-ddc8-4f89-b955-5933ecd6ed5e",
			        "items": [
			            {
			                "interval": "10m",
			                "start_time": ""
			            }
			        ]
			    },
			    {
			        "pool": "rbd",
			        "namespace": "",
			        "image": "ac3c34ed-f3a1-403b-8fa8-332286445ebc",
			        "items": [
			            {
			                "interval": "10m",
			                "start_time": ""
			            }
			        ]
			    },
			    {
			        "pool": "rbd",
			        "namespace": "",
			        "image": "ce7f8aba-3171-4c3d-9ecc-6177b1b8fc77",
			        "items": [
			            {
			                "interval": "10m",
			                "start_time": ""
			            }
			        ]
			    },
			    {
			        "pool": "rbd",
			        "namespace": "",
			        "image": "d50cb9bf-b2c7-4bf1-9e21-557d1471d3a9",
			        "items": [
			            {
			                "interval": "10m",
			                "start_time": ""
			            }
			        ]
			    },
			    {
			        "pool": "rbd",
			        "namespace": "",
			        "image": "mirror-test",
			        "items": [
			            {
			                "interval": "10m",
			                "start_time": ""
			            }
			        ]
			    },
			    {
			        "pool": "rbd",
			        "namespace": "",
			        "image": "test2",
			        "items": [
			            {
			                "interval": "10m",
			                "start_time": ""
			            }
			        ]
			    }
			]
		*/

	} else {
		stdRemote = []byte("[\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"4f1ff9d5-7cfd-4d5a-97fd-ba3bb6faa17b\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"5f44786f-ddc8-4f89-b955-5933ecd6ed5e\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"ac3c34ed-f3a1-403b-8fa8-332286445ebc\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"ce7f8aba-3171-4c3d-9ecc-6177b1b8fc77\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"d50cb9bf-b2c7-4bf1-9e21-557d1471d3a9\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"mirror-test\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"test2\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    }\n]")
		stdLocal = []byte("[\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"4f1ff9d5-7cfd-4d5a-97fd-ba3bb6faa17b\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"5f44786f-ddc8-4f89-b955-5933ecd6ed5e\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"ac3c34ed-f3a1-403b-8fa8-332286445ebc\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"ce7f8aba-3171-4c3d-9ecc-6177b1b8fc77\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"d50cb9bf-b2c7-4bf1-9e21-557d1471d3a9\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"mirror-test\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"test2\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    }\n]")
	}
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
	MirrorList.Debug = gin.IsDebugging()
	return MirrorList, err
}

func Status() (mirrorStatus model.MirrorStatus, err error) {
	var tmpdat struct {
		Summary model.MirrorStatus `json:"summary"`
	}
	var stdout []byte
	if gin.IsDebugging() != true {

		cmd := exec.Command("rbd", "mirror", "pool", "status", "--format", "json", "--pretty-format")
		var out strings.Builder
		//cmd.Stderr = &out
		stdout, err = cmd.CombinedOutput()

		if err != nil {
			err = errors.Join(err, errors.New(out.String()))
			utils.FancyHandleError(err)
			return
		}

	} else {

		stdout = []byte("{\n    \"summary\": {\n        \"health\": \"WARNING\",\n        \"daemon_health\": \"OK\",\n        \"image_health\": \"WARNING\",\n        \"states\": {\n            \"unknown\": 14\n        }\n    }\n}")

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

	if gin.IsDebugging() != true {

		strRemoveStatus := exec.Command("rbd", "mirror", "snapshot", "schedule", "rm", "--pool", poolName, "--image", imageName)
		stdRemove, err = strRemoveStatus.CombinedOutput()
	} else {
		stdRemove = []byte("[\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"4f1ff9d5-7cfd-4d5a-97fd-ba3bb6faa17b\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"5f44786f-ddc8-4f89-b955-5933ecd6ed5e\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"ac3c34ed-f3a1-403b-8fa8-332286445ebc\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"ce7f8aba-3171-4c3d-9ecc-6177b1b8fc77\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"d50cb9bf-b2c7-4bf1-9e21-557d1471d3a9\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"mirror-test\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"test2\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    }\n]")
	}
	if err != nil {
		err = errors.New(string(stdRemove))
		utils.FancyHandleError(err)
		return
	}
	if gin.IsDebugging() != true {

		strRemoveStatus := exec.Command("rbd", "mirror", "image", "disable", "--pool", poolName, "--image", imageName)
		stdRemove, err = strRemoveStatus.CombinedOutput()
	} else {
		stdRemove = []byte("[\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"4f1ff9d5-7cfd-4d5a-97fd-ba3bb6faa17b\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"5f44786f-ddc8-4f89-b955-5933ecd6ed5e\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"ac3c34ed-f3a1-403b-8fa8-332286445ebc\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"ce7f8aba-3171-4c3d-9ecc-6177b1b8fc77\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"d50cb9bf-b2c7-4bf1-9e21-557d1471d3a9\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"mirror-test\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"test2\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    }\n]")
	}
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

	if gin.IsDebugging() != true {

		strMirrorEnableOutput := exec.Command("rbd", "mirror", "image", "enable", "--pool", poolName, "--image", imageName, "snapshot")
		stdoutMirrorEnable, err = strMirrorEnableOutput.CombinedOutput()
		if err != nil || string(stdoutMirrorEnable) != "Mirroring enabled\n" {
			err = errors.Join(err, errors.New(string(stdoutMirrorEnable)))
			utils.FancyHandleError(err)
			return
		}

	} else {
		stdoutMirrorEnable = []byte("[\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"4f1ff9d5-7cfd-4d5a-97fd-ba3bb6faa17b\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"5f44786f-ddc8-4f89-b955-5933ecd6ed5e\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"ac3c34ed-f3a1-403b-8fa8-332286445ebc\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"ce7f8aba-3171-4c3d-9ecc-6177b1b8fc77\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"d50cb9bf-b2c7-4bf1-9e21-557d1471d3a9\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"mirror-test\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"test2\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    }\n]")
	}

	output = string(stdoutMirrorEnable)
	return
}

func ImageConfig(poolName string, imageName string, interval string, startTime string) (output string, err error) {

	var stdoutScheduleEnable []byte

	if gin.IsDebugging() != true {
		var strScheduleOutput *exec.Cmd
		print(startTime)
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

	} else {
		stdoutScheduleEnable = []byte("[\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"4f1ff9d5-7cfd-4d5a-97fd-ba3bb6faa17b\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"5f44786f-ddc8-4f89-b955-5933ecd6ed5e\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"ac3c34ed-f3a1-403b-8fa8-332286445ebc\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"ce7f8aba-3171-4c3d-9ecc-6177b1b8fc77\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"d50cb9bf-b2c7-4bf1-9e21-557d1471d3a9\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"mirror-test\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"test2\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    }\n]")
	}

	output = string(stdoutScheduleEnable)
	return
}

func ImageStatus(poolName string, imageName string) (imageStatus model.ImageStatus, err error) {

	var stdoutScheduleEnable []byte
	if gin.IsDebugging() != true {
		strScheduleOutput := exec.Command("rbd", "mirror", "image", "status", "--pool", poolName, "--image", imageName, "--format", "json")
		stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()
	} else {
		stdoutScheduleEnable = []byte("{\"name\":\"test\",\"global_id\":\"8a5db4a1-8d16-43df-bc9e-8ea8c8007879\",\"state\":\"up+stopped\",\"description\":\"local image is primary\",\"daemon_service\":{\"service_id\":\"1517082\",\"instance_id\":\"1707685\",\"daemon_id\":\"scvm13.kphwqj\",\"hostname\":\"scvm13\"},\"last_update\":\"2023-11-03 10:11:26\",\"peer_sites\":[{\"site_name\":\"rbd2\",\"mirror_uuids\":\"b69a89e3-ea2b-4ce5-8973-19ac0832281f\",\"state\":\"up+replaying\",\"description\":\"replaying, {\\\"bytes_per_second\\\":0.0,\\\"bytes_per_snapshot\\\":0.0,\\\"last_snapshot_bytes\\\":0,\\\"last_snapshot_sync_seconds\\\":0,\\\"local_snapshot_timestamp\\\":1698973860,\\\"remote_snapshot_timestamp\\\":1698973860,\\\"replay_state\\\":\\\"idle\\\"}\",\"last_update\":\"2023-11-03 10:11:28\"}],\"snapshots\":[{\"id\":6636,\"name\":\".mirror.primary.8a5db4a1-8d16-43df-bc9e-8ea8c8007879.68dfeb05-f9e1-4b4e-894e-5b622821320d\",\"demoted\":false,\"mirror_peer_uuids\":[]},{\"id\":6638,\"name\":\".mirror.primary.8a5db4a1-8d16-43df-bc9e-8ea8c8007879.71bb09d3-651c-4052-846f-c395563df597\",\"demoted\":false,\"mirror_peer_uuids\":[\"231f94df-9c4a-4316-8473-34fbc7f50d6b\"]}]}")
	}

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
	if gin.IsDebugging() != true {
		strScheduleOutput := exec.Command("rbd", "mirror", "image", "promote", "--pool", poolName, "--image", imageName)
		stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()
	} else {
		stdoutScheduleEnable = []byte("{\"name\":\"test\",\"global_id\":\"8a5db4a1-8d16-43df-bc9e-8ea8c8007879\",\"state\":\"up+stopped\",\"description\":\"local image is primary\",\"daemon_service\":{\"service_id\":\"1517082\",\"instance_id\":\"1707685\",\"daemon_id\":\"scvm13.kphwqj\",\"hostname\":\"scvm13\"},\"last_update\":\"2023-11-03 10:11:26\",\"peer_sites\":[{\"site_name\":\"rbd2\",\"mirror_uuids\":\"b69a89e3-ea2b-4ce5-8973-19ac0832281f\",\"state\":\"up+replaying\",\"description\":\"replaying, {\\\"bytes_per_second\\\":0.0,\\\"bytes_per_snapshot\\\":0.0,\\\"last_snapshot_bytes\\\":0,\\\"last_snapshot_sync_seconds\\\":0,\\\"local_snapshot_timestamp\\\":1698973860,\\\"remote_snapshot_timestamp\\\":1698973860,\\\"replay_state\\\":\\\"idle\\\"}\",\"last_update\":\"2023-11-03 10:11:28\"}],\"snapshots\":[{\"id\":6636,\"name\":\".mirror.primary.8a5db4a1-8d16-43df-bc9e-8ea8c8007879.68dfeb05-f9e1-4b4e-894e-5b622821320d\",\"demoted\":false,\"mirror_peer_uuids\":[]},{\"id\":6638,\"name\":\".mirror.primary.8a5db4a1-8d16-43df-bc9e-8ea8c8007879.71bb09d3-651c-4052-846f-c395563df597\",\"demoted\":false,\"mirror_peer_uuids\":[\"231f94df-9c4a-4316-8473-34fbc7f50d6b\"]}]}")
	}

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
