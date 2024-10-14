package license

import (
	"Glue-API/utils"
	"errors"
	"strings"
)

func License() (output []string, err error) {

	var stdout []byte
	//  For Remote
	settings, _ := utils.ReadConfFile()
	client, err := utils.ConnectSSH(settings.RemoteHostIp, settings.RemoteRootRsaIdPath)
	if err != nil {
		err = err
		utils.FancyHandleError(err)
		return
	}
	//// Defer closing the network connection.
	defer client.Close()
	//// Execute your command.

	// name
	cmd, err := client.Command("cat /root/license_test | grep 'name' | awk '{print $3}'")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	license_info := strings.ReplaceAll(string(stdout), "\n", "")
	output = append(output, string(license_info))

	// date
	cmd, err = client.Command("cat /root/license_test | grep 'date' | awk '{print $3}'")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	license_info = strings.ReplaceAll(string(stdout), "\n", "")
	output = append(output, string(license_info))

	return
}
