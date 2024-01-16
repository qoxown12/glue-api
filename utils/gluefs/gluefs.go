package gluefs

import (
	"Glue-API/model"
	// "Glue-API/utils"
	"encoding/json"
	"os/exec"

	"github.com/gin-gonic/gin"
	// "strings"
)

func ListFs() (glue_fs_list []model.GlueFs, err error) {
	var stdout []byte
	if gin.IsDebugging() == true {
		cmd := exec.Command("ceph", "fs", "ls", "--format", "json")
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			return
		}
	} else {
		// Print the output
		stdout = []byte("[{\"name\":\"fs\",\"metadata_pool\":\"gluefs.fs.meta\",\"metadata_pool_id\":4,\"data_pool_ids\":[5],\"data_pools\":[\"gluefs.fs.data\"]}]")

	}
	if err = json.Unmarshal(stdout, &glue_fs_list); err != nil {
		return
	}
	return
}

// func FsSetup(poolName string, imageName string) (output string, err error) {
func FsSetup(fsName string) {
	// var stdoutMirrorEnable []byte

	// if gin.IsDebugging() == true {

	// 	strMirrorEnableOutput := exec.Command("rbd", "mirror", "image", "enable", "--pool", poolName, "--image", imageName, "snapshot")
	// 	stdoutMirrorEnable, err = strMirrorEnableOutput.CombinedOutput()
	// 	if err != nil || string(stdoutMirrorEnable) != "Mirroring enabled\n" {
	// 		err = errors.Join(err, errors.New(string(stdoutMirrorEnable)))
	// 		utils.FancyHandleError(err)
	// 		return
	// 	}

	// } else {
	// 	stdoutMirrorEnable = []byte("[\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"4f1ff9d5-7cfd-4d5a-97fd-ba3bb6faa17b\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"5f44786f-ddc8-4f89-b955-5933ecd6ed5e\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"ac3c34ed-f3a1-403b-8fa8-332286445ebc\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"ce7f8aba-3171-4c3d-9ecc-6177b1b8fc77\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"d50cb9bf-b2c7-4bf1-9e21-557d1471d3a9\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"mirror-test\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    },\n    {\n        \"pool\": \"rbd\",\n        \"namespace\": \"\",\n        \"image\": \"test2\",\n        \"items\": [\n            {\n                \"interval\": \"10m\",\n                \"start_time\": \"\"\n            }\n        ]\n    }\n]")
	// }

	// output = string(stdoutMirrorEnable)
	// return
}
