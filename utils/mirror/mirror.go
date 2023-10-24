package mirror

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
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
		stdout, err = cmd.Output()
		/*
			   {
			    "mode": "image",
			    "site_name": "cluster1",
			    "peers": [
			        {
			            "uuid": "c5f36d07-a69e-45b5-a6ae-7e7a5f1266f0",
			            "direction": "rx-tx",
			            "site_name": "cluster2",
			            "mirror_uuid": "d7576286-c14c-47cf-8b87-58c4ec7e08dc",
			            "client_name": "client.rbd-mirror-peer",
			            "key": "AQAxDx5lljPfGhAAQhP1voVx5Dogn3f+nzYM8A==",
			            "mon_host": "[v2:100.100.1.24:3300/0,v1:100.100.1.24:6789/0],[v2:100.100.1.25:3300/0,v1:100.100.1.25:6789/0],[v2:100.100.1.26:3300/0,v1:100.100.1.26:6789/0]"
			        }
			    ]
			}
		*/
		if err != nil {
			return clusterConf, err
		}

	} else {
		stdout = []byte("{\n    \"mode\": \"image\",\n    \"site_name\": \"cluster1\",\n    \"peers\": [\n        {\n            \"uuid\": \"c5f36d07-a69e-45b5-a6ae-7e7a5f1266f0\",\n            \"direction\": \"rx-tx\",\n            \"site_name\": \"cluster2\",\n            \"mirror_uuid\": \"d7576286-c14c-47cf-8b87-58c4ec7e08dc\",\n            \"client_name\": \"client.rbd-mirror-peer\",\n            \"key\": \"AQAxDx5lljPfGhAAQhP1voVx5Dogn3f+nzYM8A==\",\n            \"mon_host\": \"[v2:100.100.1.24:3300/0,v1:100.100.1.24:6789/0],[v2:100.100.1.25:3300/0,v1:100.100.1.25:6789/0],[v2:100.100.1.26:3300/0,v1:100.100.1.26:6789/0]\"\n        }\n    ]\n}")
	}

	if err := json.Unmarshal(stdout, &clusterConf); err != nil {

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
	if err := tfCluster.Close(); err != nil {
		return clusterConf, err
	}
	if _, err = tfKey.WriteString(peer.Key); err != nil {
		fmt.Println("Failed to write to temporary file", err)
		return clusterConf, err
	}
	if err := tfKey.Close(); err != nil {
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
		stdRemote, err = strRemoteStatus.Output()
		strLocalStatus := exec.Command("rbd", "mirror", "snapshot", "schedule", "list", "-R", "--format", "json", "--pretty-format")
		stdLocal, err = strLocalStatus.Output()

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
	if err := json.Unmarshal(stdRemote, &Remote); err != nil {
		Remote = []model.MirrorImage{}
	}
	if err := json.Unmarshal(stdLocal, &Local); err != nil {
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
		cmd.Stderr = &out
		stdout, err = cmd.Output()

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

		strRemoveStatus := exec.Command("rbd", "mirror", "image", "disable", "--pool", poolName, "--image", imageName)
		stdRemove, err = strRemoveStatus.Output()

	} else {
		stdRemove = []byte("[\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"4f1ff9d5-7cfd-4d5a-97fd-ba3bb6faa17b\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"5f44786f-ddc8-4f89-b955-5933ecd6ed5e\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"ac3c34ed-f3a1-403b-8fa8-332286445ebc\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"ce7f8aba-3171-4c3d-9ecc-6177b1b8fc77\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"d50cb9bf-b2c7-4bf1-9e21-557d1471d3a9\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"mirror-test\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"test2\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    }\n]")
	}
	if err != nil {
		err = errors.Join(err, errors.New(string(stdRemove)))
		utils.FancyHandleError(err)
		return
	}
	output = string(stdRemove)
	return
}
