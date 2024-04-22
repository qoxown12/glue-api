package gluevm

import (

	// "Glue-API/utils"
	"Glue-API/utils"
	"errors"
	"os/exec"

	"github.com/gin-gonic/gin"
	// "strings"
)

func VmState(hypervisorType string) (output string, err error) {
	var stdoutVmState []byte

	if gin.IsDebugging() == true {
		if hypervisorType == "cell" {
			strVmStateOutput := exec.Command("python3", "/usr/share/cockpit/ablestack/python/gwvm/gwvm_status_check.py", "check")

			stdoutVmState, err = strVmStateOutput.CombinedOutput()
			if err != nil {
				err = errors.Join(err, errors.New(string(stdoutVmState)))
				utils.FancyHandleError(err)
				return
			}
		} else {
			output = "This hypervisor type is not supported."
			return
		}
	} else {
		stdoutVmState = []byte("{\"code\": 200, \"val\": {\"role\": \"Started\", \"started\": \"100.100.1.3\", \"ip\": \"10.10.1.20\", \"mac\": \"00:24:81:7f:30:1e\", \"nictype\": \"bridge\", \"nicbridge\": \"bridge0\", \"Id\": \"20\", \"Name\": \"gwvm\", \"UUID\": \"637d8fe6-4797-4fa6-be94-4823eb61ddd8\", \"OS Type\": \"hvm\", \"State\": \"running\", \"CPU(s)\": \"4\", \"CPU time\": \"196.6s\", \"Max memory\": \"8388608 KiB\", \"Used memory\": \"8388608 KiB\", \"Persistent\": \"no\", \"Autostart\": \"disable\", \"Managed save\": \"no\", \"Security model\": \"none\", \"Security DOI\": \"0\", \"prefix\": \"24\", \"gw\": \"10.10.0.1\", \"disk_cap\": \"83G\", \"disk_alloc\": \"13G\", \"disk_phy\": \"70G\", \"disk_usage_rate\": \"16%\"}, \"name\": \"check\", \"type\": \"dict\"}")
	}

	output = string(stdoutVmState)
	return
}

func VmDetail(hypervisorType string) (output string, err error) {

	var stdoutVmStart []byte

	if gin.IsDebugging() == true {
		if hypervisorType == "cell" {
			//  For Remote
			settings, _ := utils.ReadConfFile()
			client, err_val := utils.ConnectSSH(settings.RemoteHostIp, settings.RemoteRootRsaIdPath)
			if err_val != nil {
				err = err_val
				utils.FancyHandleError(err_val)
				return
			}
			//// Defer closing the network connection.
			defer client.Close()
			//// Execute your command.

			sshcmd, err_val := client.Command("python3", "/usr/share/cockpit/ablestack/python/pcs/main.py", "status", "--resource", "gateway_res")
			if err_val != nil {
				err = err_val
				utils.FancyHandleError(err_val)
				return

			}
			stdoutVmStart, err = sshcmd.CombinedOutput()
		} else {
			output = "This hypervisor type is not supported."
			return
		}
	} else {
		stdoutVmStart = []byte("{\"debug\": false, \"Message\": \"{\n    \"code\": 200,\n    \"val\": {\n        \"clustered_host\": [\n            \"100.100.2.1\",\n            \"100.100.2.2\",\n            \"100.100.2.3\"\n        ],\n        \"nodes\": [\n            {\n                \"host\": \"100.100.2.1\",\n                \"online\": \"true\",\n                \"resources_running\": \"1\",\n                \"standby\": \"false\",\n                \"standby_onfail\": \"false\",\n                \"maintenance\": \"false\",\n                \"pending\": \"false\",\n                \"unclean\": \"false\",\n                \"shutdown\": \"false\",\n                \"expected_up\": \"true\",\n                \"is_dc\": \"true\",\n                \"type\": \"member\"\n            },\n            {\n                \"host\": \"100.100.2.2\",\n                \"online\": \"true\",\n                \"resources_running\": \"0\",\n                \"standby\": \"false\",\n                \"standby_onfail\": \"false\",\n                \"maintenance\": \"false\",\n                \"pending\": \"false\",\n                \"unclean\": \"false\",\n                \"shutdown\": \"false\",\n                \"expected_up\": \"true\",\n                \"is_dc\": \"false\",\n                \"type\": \"member\"\n            },\n            {\n                \"host\": \"100.100.2.3\",\n                \"online\": \"true\",\n                \"resources_running\": \"1\",\n                \"standby\": \"false\",\n                \"standby_onfail\": \"false\",\n                \"maintenance\": \"false\",\n                \"pending\": \"false\",\n                \"unclean\": \"false\",\n                \"shutdown\": \"false\",\n                \"expected_up\": \"true\",\n                \"is_dc\": \"false\",\n                \"type\": \"member\"\n            }\n        ],\n        \"started\": \"100.100.2.3\",\n        \"role\": \"Started\",\n        \"active\": \"true\",\n        \"blocked\": \"false\",\n        \"failed\": \"false\"\n    },\n    \"name\": \"statusResource\",\n    \"type\": \"dict\"\n}\n\"}")
	}
	output = string(stdoutVmStart)
	return
}

func VmSetup(hypervisorType string, gwvmMngtNicParen string, gwvmMngtNicIp string, gwvmStorageNicParent string, gwvmStorageNicIp string) (output string, err error) {

	var stdoutVmSetup []byte

	if gin.IsDebugging() == true {
		if hypervisorType == "cell" {
			strVmSetupOutput := exec.Command("python3", "/usr/share/cockpit/ablestack/python/gwvm/gwvm_create.py", "create", "-mnb", gwvmMngtNicParen, "-mi", gwvmMngtNicIp, "-snb", gwvmStorageNicParent, "-si", gwvmStorageNicIp)

			stdoutVmSetup, err = strVmSetupOutput.CombinedOutput()
			if err != nil {
				err = errors.Join(err, errors.New(string(stdoutVmSetup)))
				utils.FancyHandleError(err)
				return
			}
		} else {
			output = "This hypervisor type is not supported."
			return
		}
	} else {
		stdoutVmSetup = []byte("{\"code\": 200, \"val\": \"Gateway VM Create Success\", \"name\": \"create\", \"type\": \"str\"}")
	}

	output = string(stdoutVmSetup)
	return
}

func VmStart(hypervisorType string) (output string, err error) {

	var stdoutVmStart []byte

	if gin.IsDebugging() == true {
		if hypervisorType == "cell" {
			//  For Remote
			settings, _ := utils.ReadConfFile()
			client, err_val := utils.ConnectSSH(settings.RemoteHostIp, settings.RemoteRootRsaIdPath)
			if err_val != nil {
				err = err_val
				utils.FancyHandleError(err_val)
				return
			}
			//// Defer closing the network connection.
			defer client.Close()
			//// Execute your command.

			sshcmd, err_val := client.Command("python3", "/usr/share/cockpit/ablestack/python/pcs/main.py", "enable", "--resource", "gateway_res")
			if err_val != nil {
				err = err_val
				utils.FancyHandleError(err_val)
				return

			}
			stdoutVmStart, err = sshcmd.CombinedOutput()
		} else {
			output = "This hypervisor type is not supported."
			return
		}
	} else {
		stdoutVmStart = []byte("{ \"code\": 200, \"val\": \"enable\", \"name\": \"enableResource\", \"type\": \"str\"}")
	}
	output = string(stdoutVmStart)
	return
}

func VmStop(hypervisorType string) (output string, err error) {

	var stdoutVmStop []byte

	if gin.IsDebugging() == true {
		if hypervisorType == "cell" {
			//  For Remote
			settings, _ := utils.ReadConfFile()
			client, err_val := utils.ConnectSSH(settings.RemoteHostIp, settings.RemoteRootRsaIdPath)
			if err_val != nil {
				err = err_val
				utils.FancyHandleError(err_val)
				return
			}
			//// Defer closing the network connection.
			defer client.Close()
			//// Execute your command.

			sshcmd, err_val := client.Command("python3", "/usr/share/cockpit/ablestack/python/pcs/main.py", "disable", "--resource", "gateway_res")
			if err_val != nil {
				err = err_val
				utils.FancyHandleError(err_val)
				return

			}
			stdoutVmStop, err = sshcmd.CombinedOutput()
		} else {
			output = "This hypervisor type is not supported."
			return
		}
	} else {
		stdoutVmStop = []byte("{ \"code\": 200, \"val\": \"disable\", \"name\": \"disableResource\", \"type\": \"str\"}")
	}

	output = string(stdoutVmStop)
	return
}

func VmDelete(hypervisorType string) (output string, err error) {

	var stdoutVmDelete []byte

	if gin.IsDebugging() == true {
		if hypervisorType == "cell" {
			//  For Remote
			settings, _ := utils.ReadConfFile()
			client, err_val := utils.ConnectSSH(settings.RemoteHostIp, settings.RemoteRootRsaIdPath)
			if err_val != nil {
				err = err_val
				utils.FancyHandleError(err_val)
				return
			}
			//// Defer closing the network connection.
			defer client.Close()
			//// Execute your command.

			sshcmd, err_val := client.Command("python3", "/usr/share/cockpit/ablestack/python/pcs/main.py", "remove", "--resource", "gateway_res")
			if err_val != nil {
				err = err_val
				utils.FancyHandleError(err_val)
				return
			}

			stdoutVmDelete, err = sshcmd.CombinedOutput()
		} else {
			output = "This hypervisor type is not supported."
			return
		}
	} else {
		stdoutVmDelete = []byte("{ \"code\": 200, \"val\": \"remove\", \"name\": \"removeResource\", \"type\": \"str\"}")
	}

	output = string(stdoutVmDelete)
	return
}

func VmCleanup(hypervisorType string) (output string, err error) {

	var stdoutVmCleanup []byte

	if gin.IsDebugging() == true {
		if hypervisorType == "cell" {
			//  For Remote
			settings, _ := utils.ReadConfFile()
			client, err_val := utils.ConnectSSH(settings.RemoteHostIp, settings.RemoteRootRsaIdPath)
			if err_val != nil {
				err = err_val
				utils.FancyHandleError(err_val)
				return
			}
			//// Defer closing the network connection.
			defer client.Close()
			//// Execute your command.

			sshcmd, err_val := client.Command("python3", "/usr/share/cockpit/ablestack/python/pcs/main.py", "cleanup", "--resource", "gateway_res")
			if err_val != nil {
				err = err_val
				utils.FancyHandleError(err_val)
				return

			}
			stdoutVmCleanup, err = sshcmd.CombinedOutput()
		} else {
			output = "This hypervisor type is not supported."
			return
		}
	} else {
		stdoutVmCleanup = []byte("{ \"code\": 200, \"val\": \"cleanup\", \"name\": \"cleanupResource\", \"type\": \"str\"}")
	}

	output = string(stdoutVmCleanup)
	return
}

func VmMigrate(hypervisorType string, target string) (output string, err error) {

	var stdoutVmMigrate []byte

	if gin.IsDebugging() == true {
		if hypervisorType == "cell" {
			//  For Remote
			settings, _ := utils.ReadConfFile()
			client, err_val := utils.ConnectSSH(settings.RemoteHostIp, settings.RemoteRootRsaIdPath)
			if err_val != nil {
				err = err_val
				utils.FancyHandleError(err_val)
				return
			}
			//// Defer closing the network connection.
			defer client.Close()
			//// Execute your command.

			sshcmd, err_val := client.Command("python3", "/usr/share/cockpit/ablestack/python/pcs/main.py", "move", "--resource", "gateway_res", "--target", target)
			if err_val != nil {
				err = err_val
				utils.FancyHandleError(err_val)
				return

			}
			stdoutVmMigrate, err = sshcmd.CombinedOutput()
		} else {
			output = "This hypervisor type is not supported."
			return
		}
	} else {
		stdoutVmMigrate = []byte("{\"code\": 200, \"val\": { \"current host\": \"100.100.5.3\", \"target host\": \"100.100.5.1\" }, \"name\": \"moveResource\", \"type\": \"dict\"}")
	}

	output = string(stdoutVmMigrate)
	return
}
