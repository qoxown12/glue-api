package glue

import (
	"Glue-API/model"
	"Glue-API/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"os/exec"
	"strings"
)

func ListPool() (pools []string, err error) {
	var stdout []byte
	if gin.IsDebugging() != true {
		cmd := exec.Command("ceph", "osd", "pool", "ls", "--format", "json")
		stdout, err = cmd.CombinedOutput()

		if err != nil {
			return
		}

	} else {
		// Print the output
		stdout = []byte("[\".mgr\",\"rbd\"]")

	}
	if err = json.Unmarshal(stdout, &pools); err != nil {
		return
	}
	return
}

func ListImage(pool string) (images []model.Snapshot, err error) {
	var stdout []byte
	if gin.IsDebugging() != true {
		cmd := exec.Command("rbd", "ls", "-l", "-p", pool, "--format", "json")
		stdout, err = cmd.CombinedOutput()

		if err != nil {
			return
		}

	} else {
		// Print the output
		stdout = []byte("[{\"image\":\"test\",\"id\":\"15de89c519b8\",\"size\":10737418240,\"format\":2},{\"image\":\"test2\",\"id\":\"15dec78c7823\",\"size\":10737418240,\"format\":2},{\"image\":\"test3\",\"id\":\"15de61bec90f\",\"size\":10737418240,\"format\":2}]")

	}
	if err = json.Unmarshal(stdout, &images); err != nil {
		return
	}
	return
}

func Status() (dat model.GlueStatus, err error) {

	if gin.IsDebugging() != true {
		var stdout []byte
		var out strings.Builder
		cmd := exec.Command("ceph", "-s", "-f", "json")
		cmd.Stderr = &out
		stdout, err = cmd.CombinedOutput()

		if err = json.Unmarshal(stdout, &dat); err != nil {
			utils.FancyHandleError(err)
			return
		}
	} else {

		stdout := []byte("{\n  \"fsid\":\"9980ffe8-4bc1-11ee-9b1f-002481004170\",\n  \"health\":{\n    \"status\":\"HEALTH_WARN\",\n    \"checks\":{\n      \"RECENT_MGR_MODULE_CRASH\":{\n        \"severity\":\"HEALTH_WARN\",\n        \"summary\":{\n          \"message\":\"4 mgr modules have recently crashed\",\n          \"count\":4\n        },\n        \"muted\":false\n      }\n    },\n    \"mutes\":[\n\n    ]\n  },\n  \"election_epoch\":148,\n  \"quorum\":[\n    0,\n    1,\n    2\n  ],\n  \"quorum_names\":[\n    \"scvm1\",\n    \"scvm3\",\n    \"scvm2\"\n  ],\n  \"quorum_age\":1320385,\n  \"monmap\":{\n    \"epoch\":9,\n    \"min_mon_release_name\":\"reef\",\n    \"num_mons\":3\n  },\n  \"osdmap\":{\n    \"epoch\":13906,\n    \"num_osds\":19,\n    \"num_up_osds\":19,\n    \"osd_up_since\":1694672928,\n    \"num_in_osds\":19,\n    \"osd_in_since\":1693900905,\n    \"num_remapped_pgs\":0\n  },\n  \"pgmap\":{\n    \"pgs_by_state\":[\n      {\n        \"state_name\":\"active+clean\",\n        \"count\":801\n      }\n    ],\n    \"num_pgs\":801,\n    \"num_pools\":8,\n    \"num_objects\":255687,\n    \"data_bytes\":1010750055765,\n    \"bytes_used\":1945430351872,\n    \"bytes_avail\":16298248544256,\n    \"bytes_total\":18243678896128,\n    \"read_bytes_sec\":5370,\n    \"write_bytes_sec\":3247913,\n    \"read_op_per_sec\":75,\n    \"write_op_per_sec\":99\n  },\n  \"fsmap\":{\n    \"epoch\":1,\n    \"by_rank\":[\n\n    ],\n    \"up:standby\":0\n  },\n  \"mgrmap\":{\n    \"available\":true,\n    \"num_standbys\":1,\n    \"modules\":[\n      \"cephadm\",\n      \"dashboard\",\n      \"iostat\",\n      \"nfs\",\n      \"prometheus\",\n      \"restful\"\n    ],\n    \"services\":{\n      \"dashboard\":\"https://10.10.1.13:8443/\",\n      \"prometheus\":\"http://100.100.1.13:9283/\"\n    }\n  },\n  \"servicemap\":{\n    \"epoch\":14150,\n    \"modified\":\"2023-10-11T08:41:52.944377+0000\",\n    \"services\":{\n      \"rbd-mirror\":{\n        \"daemons\":{\n          \"summary\":\"\",\n          \"20610527\":{\n            \"start_epoch\":13092,\n            \"start_stamp\":\"2023-10-05T02:04:38.425678+0000\",\n            \"gid\":20610527,\n            \"addr\":\"100.100.1.12:0/3953042359\",\n            \"metadata\":{\n              \"arch\":\"x86_64\",\n              \"ceph_release\":\"reef\",\n              \"ceph_version\":\"ceph version Glue-Diplo-4.0.0 (5dd24139a1eada541a3bc16b6941c5dde975e26d) reef (stable)\",\n              \"ceph_version_short\":\"Glue-Diplo-4.0.0\",\n              \"container_hostname\":\"scvm2\",\n              \"container_image\":\"localhost:5000/glue/daemon@sha256:87d1dba17511fc6dd0f89cead2aac496ac427226f475b2251f72f5a933268e59\",\n              \"cpu\":\"Intel(R) Xeon(R) Silver 4210 CPU @ 2.20GHz\",\n              \"distro\":\"rocky\",\n              \"distro_description\":\"Rocky Linux 9.2 (Blue Onyx)\",\n              \"distro_version\":\"9.2\",\n              \"hostname\":\"scvm2\",\n              \"id\":\"scvm2.zpdohp\",\n              \"instance_id\":\"20610527\",\n              \"kernel_description\":\"#1 SMP PREEMPT_DYNAMIC Wed Aug 16 10:08:14 KST 2023\",\n              \"kernel_version\":\"5.14.0-284.25.2.ablecloud.el9.x86_64\",\n              \"mem_swap_kb\":\"16777212\",\n              \"mem_total_kb\":\"32600252\",\n              \"os\":\"Linux\"\n            },\n            \"task_status\":{\n\n            }\n          },\n          \"21129580\":{\n            \"start_epoch\":13093,\n            \"start_stamp\":\"2023-10-05T02:04:50.518201+0000\",\n            \"gid\":21129580,\n            \"addr\":\"100.100.1.13:0/2840805437\",\n            \"metadata\":{\n              \"arch\":\"x86_64\",\n              \"ceph_release\":\"reef\",\n              \"ceph_version\":\"ceph version Glue-Diplo-4.0.0 (5dd24139a1eada541a3bc16b6941c5dde975e26d) reef (stable)\",\n              \"ceph_version_short\":\"Glue-Diplo-4.0.0\",\n              \"container_hostname\":\"scvm3\",\n              \"container_image\":\"localhost:5000/glue/daemon@sha256:87d1dba17511fc6dd0f89cead2aac496ac427226f475b2251f72f5a933268e59\",\n              \"cpu\":\"Intel(R) Xeon(R) Silver 4210 CPU @ 2.20GHz\",\n              \"distro\":\"rocky\",\n              \"distro_description\":\"Rocky Linux 9.2 (Blue Onyx)\",\n              \"distro_version\":\"9.2\",\n              \"hostname\":\"scvm3\",\n              \"id\":\"scvm3.yfdixv\",\n              \"instance_id\":\"21129580\",\n              \"kernel_description\":\"#1 SMP PREEMPT_DYNAMIC Wed Aug 16 10:08:14 KST 2023\",\n              \"kernel_version\":\"5.14.0-284.25.2.ablecloud.el9.x86_64\",\n              \"mem_swap_kb\":\"16777212\",\n              \"mem_total_kb\":\"32600256\",\n              \"os\":\"Linux\"\n            },\n            \"task_status\":{\n\n            }\n          }\n        }\n      },\n      \"rgw\":{\n        \"daemons\":{\n          \"summary\":\"\",\n          \"12167529\":{\n            \"start_epoch\":3806,\n            \"start_stamp\":\"2023-09-15T07:04:49.891558+0000\",\n            \"gid\":12167529,\n            \"addr\":\"100.100.1.11:0/2444620655\",\n            \"metadata\":{\n              \"arch\":\"x86_64\",\n              \"ceph_release\":\"reef\",\n              \"ceph_version\":\"ceph version Glue-Diplo-4.0.0 (5dd24139a1eada541a3bc16b6941c5dde975e26d) reef (stable)\",\n              \"ceph_version_short\":\"Glue-Diplo-4.0.0\",\n              \"container_hostname\":\"scvm1\",\n              \"container_image\":\"localhost:5000/glue/daemon@sha256:87d1dba17511fc6dd0f89cead2aac496ac427226f475b2251f72f5a933268e59\",\n              \"cpu\":\"Intel(R) Xeon(R) Silver 4210 CPU @ 2.20GHz\",\n              \"distro\":\"rocky\",\n              \"distro_description\":\"Rocky Linux 9.2 (Blue Onyx)\",\n              \"distro_version\":\"9.2\",\n              \"frontend_config#0\":\"beast port=80\",\n              \"frontend_type#0\":\"beast\",\n              \"hostname\":\"scvm1\",\n              \"id\":\"glue.scvm1.lzjtbp\",\n              \"kernel_description\":\"#1 SMP PREEMPT_DYNAMIC Wed Aug 16 10:08:14 KST 2023\",\n              \"kernel_version\":\"5.14.0-284.25.2.ablecloud.el9.x86_64\",\n              \"mem_swap_kb\":\"16777212\",\n              \"mem_total_kb\":\"32600252\",\n              \"num_handles\":\"1\",\n              \"os\":\"Linux\",\n              \"pid\":\"2\",\n              \"realm_id\":\"\",\n              \"realm_name\":\"\",\n              \"zone_id\":\"e4de42e3-5687-4aec-a230-0679a51b2e5e\",\n              \"zone_name\":\"default\",\n              \"zonegroup_id\":\"325fa56b-993b-46ab-8f98-857d4735a189\",\n              \"zonegroup_name\":\"default\"\n            },\n            \"task_status\":{\n\n            }\n          }\n        }\n      }\n    }\n  },\n  \"progress_events\":{\n\n  }\n}")

		if err = json.Unmarshal(stdout, &dat); err != nil {
			utils.FancyHandleError(err)
			return
		}
	}
	return
}
