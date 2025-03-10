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
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
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
	// if clusterConf.Mode == "disabled" {
	// 	err = errors.New("mirroring is disabled")
	// 	return clusterConf, err
	// }
	if len(clusterConf.Peers) > 0 {
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
		clusterConf.ClusterName = peer.SiteName
		clusterConf.Name = peer.ClientName
	}

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
	// if clusterConf.Mode == "disabled" {
	// 	err = errors.New("mirroring is disabled")
	// 	return clusterConf, err
	// }
	if len(clusterConf.Peers) > 0 {
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
	}
	clusterConf.ClusterFileName = tfCluster.Name()
	clusterConf.KeyFileName = tfKey.Name()
	return clusterConf, nil
}

func RbdImage(pool_name string) (pools []string, err error) {
	var stdout []byte
	cmd := exec.Command("rbd", "ls", "-p", pool_name, "--format", "json")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	if err = json.Unmarshal(stdout, &pools); err != nil {
		err_str := strings.ReplaceAll(string(stdout), "\n", "")
		err = errors.New(err_str)
		utils.FancyHandleError(err)
		return
	}
	return
}

func ImageInfo(poolName string, imageName string) (imageInfo model.ImageInfo, err error) {

	var stdoutMirrorPreSetup []byte

	strMirrorPreSetupOutput := exec.Command("rbd", "info", "--image", imageName, "--format", "json", "--pretty-format")
	stdoutMirrorPreSetup, err = strMirrorPreSetupOutput.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdoutMirrorPreSetup))
		utils.FancyHandleError(err)
		return
	}

	if err = json.Unmarshal(stdoutMirrorPreSetup, &imageInfo); err != nil {
		utils.FancyHandleError(err)
		return
	}
	return imageInfo, err
}

func ImageList(pool string) (MirrorList model.MirrorList, err error) {

	var stdout []byte

	cmd := exec.Command("rbd", "mirror", "pool", "status", pool, "--verbose", "--format", "json", "--pretty-format")
	stdout, err = cmd.CombinedOutput()

	if err != nil {
		return
	}

	if err = json.Unmarshal(stdout, &MirrorList); err != nil {
		MirrorList = model.MirrorList{}
	}

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
		} else {
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

func ImagePreDelete(poolName string, imageName string) (output string, err error) {

	var stdoutMirrorPreDelete []byte

	info, err := ImageInfo(poolName, imageName)
	if info.Parent.Image != "" {
		stdoutMirrorPreDeleteOutput := exec.Command("rbd", "mirror", "image", "disable", "--pool", poolName, "--image", info.Parent.Image, "snapshot")
		stdoutMirrorPreDelete, err = stdoutMirrorPreDeleteOutput.CombinedOutput()
		if err != nil {
			if strings.Contains(string(stdoutMirrorPreDelete), "mirroring is enabled on one or more children") {
				output = "Success"
				return
			} else {
				err = errors.New(string(stdoutMirrorPreDelete))
				utils.FancyHandleError(err)
				return
			}
		}
	}

	output = string(stdoutMirrorPreDelete)
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

func ImageDeleteSchedule(poolName string, imageName string) (output string, err error) {

	var stdRemove []byte
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

func ImagePreSetup(poolName string, imageName string) (output string, err error) {

	var stdoutMirrorPreSetupEnable []byte

	info, err := ImageInfo(poolName, imageName)
	if info.Parent.Image != "" {
		stdoutMirrorPreSetupEnableOutput := exec.Command("rbd", "mirror", "image", "enable", "--pool", poolName, "--image", info.Parent.Image, "snapshot")
		stdoutMirrorPreSetupEnable, err = stdoutMirrorPreSetupEnableOutput.CombinedOutput()
		if err != nil {
			err = errors.New(string(stdoutMirrorPreSetupEnable))
			utils.FancyHandleError(err)
			return
		}
	}

	output = string(stdoutMirrorPreSetupEnable)
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

func goCronTask(poolName, hostName, vmName string, imageName []string) (err error) {
	var stdout []byte
	currentTime := time.Now()
	println("start mirror snapshot scheduler --- vm : " + vmName + " --- image : " + strings.Join(imageName, ",") + " --- host : " + hostName + " --- date : " + currentTime.Format("2006-01-02 15:04:05"))
	if hostName != "" {
		println("start domfsfreeze ---")
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostName, "virsh", "domfsfreeze", vmName)
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			println("failed to virsh domfsfreeze")
			println(string(stdout))
		}
	}
	if len(imageName) > 0 {
		for i := 0; i < len(imageName); i++ {
			cmd := exec.Command(poolName, "mirror", "image", "snapshot", poolName+"/"+imageName[i])
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				println("failed to create rbd mirror image snapshot path : " + imageName[i])
				println(string(stdout))
				if hostName != "" {
					exec.Command("ssh", hostName, "virsh", "domfsthaw", vmName)
				}
				break
			}
			host, _ := os.Hostname()
			cmd = exec.Command("rbd", "image-meta", "set", "rbd/MOLD-DR", imageName[i], currentTime.Format("2006-01-02 15:04:05")+","+host)
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				println("failed to update image-meta")
				println(string(stdout))
			}
		}
	}
	if hostName != "" {
		println("start domfsthaw ---")
		cmd := exec.Command("ssh", hostName, "virsh", "domfsthaw", vmName)
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			println("failed to virsh domfsthaw")
			println(string(stdout))
		}
	}
	println("end mirror snapshot scheduler --- vm : " + vmName + " --- image : " + strings.Join(imageName, ",") + " --- host : " + hostName + " --- date : " + currentTime.Format("2006-01-02 15:04:05"))
	return
}

func goCronEventListeners(scheduler gocron.Scheduler, jobID uuid.UUID, beforeIt time.Duration, jobName, imageName, hostName, vmName, poolName string) (host string, clock time.Duration, imageList []string) {
	var afterIt time.Duration
	var exist string
	var interval string
	currentTime := time.Now()

	println("BeforeJobRuns: ", jobID.String(), jobName, currentTime.Format("2006-01-02 15:04:05"))
	mold, _ := utils.ReadMoldFile()
	exist = ""
	if mold.MoldUrl != "moldUrl" {
		drResult := utils.GetDisasterRecoveryClusterList()
		getDisasterRecoveryClusterList := model.GetDisasterRecoveryClusterList{}
		drInfo, _ := json.Marshal(drResult["getdisasterrecoveryclusterlistresponse"])
		json.Unmarshal([]byte(drInfo), &getDisasterRecoveryClusterList)
		if getDisasterRecoveryClusterList.Count == -1 {
			var Ti time.Duration
			interval, _ := ImageMetaGetInterval()
			if strings.Contains(interval, "d") {
				interval = strings.TrimRight(interval, "d\n")
				ti, _ := strconv.Atoi(interval)
				Ti = time.Duration(ti) * 24 * time.Hour
			} else if strings.Contains(interval, "h") {
				interval = strings.TrimRight(interval, "h\n")
				ti, _ := strconv.Atoi(interval)
				Ti = time.Duration(ti) * time.Hour
			} else if strings.Contains(interval, "m") {
				interval = strings.TrimRight(interval, "m\n")
				ti, _ := strconv.Atoi(interval)
				Ti = time.Duration(ti) * time.Minute
			} else {
				Ti = time.Duration(1) * time.Hour
			}
			clock = Ti
			imageList = append(imageList, imageName)
		} else {
			if len(getDisasterRecoveryClusterList.Disasterrecoverycluster) > 0 {
				dr := getDisasterRecoveryClusterList.Disasterrecoverycluster
				for i := 0; i < len(dr); i++ {
					if len(dr[i].Drclustervmmap) > 0 {
						for j := 0; j < len(dr[i].Drclustervmmap); j++ {
							if imageName == dr[i].Drclustervmmap[j].Drclustermirrorvmvolpath {
								exist = "exist"
								interval = dr[i].Details.Mirrorscheduleinterval
								if strings.Contains(interval, "d") {
									interval = strings.TrimRight(interval, "d")
									ti, _ := strconv.Atoi(interval)
									afterIt = time.Duration(ti) * 24 * time.Hour
								} else if strings.Contains(interval, "h") {
									interval = strings.TrimRight(interval, "h")
									ti, _ := strconv.Atoi(interval)
									afterIt = time.Duration(ti) * time.Hour
								} else if strings.Contains(interval, "m") {
									interval = strings.TrimRight(interval, "m")
									ti, _ := strconv.Atoi(interval)
									afterIt = time.Duration(ti) * time.Minute
								} else {
									// 잘못 입력된 경우 1시간으로 설정
									afterIt = time.Duration(1) * time.Hour
								}
								params1 := []utils.MoldParams{
									{"keyword": dr[i].Drclustervmmap[j].Drclustermirrorvmname},
								}
								vmResult := utils.GetListVirtualMachinesMetrics(params1)
								listVirtualMachinesMetrics := model.ListVirtualMachinesMetrics{}
								vmInfo, _ := json.Marshal(vmResult["listvirtualmachinesmetricsresponse"])
								json.Unmarshal([]byte(vmInfo), &listVirtualMachinesMetrics)
								vm := listVirtualMachinesMetrics.Virtualmachine
								for k := 0; k < len(vm); k++ {
									if vm[k].Name == dr[i].Drclustervmmap[j].Drclustermirrorvmname {
										volStatus, _ := ImageStatus("rbd", dr[i].Drclustervmmap[j].Drclustermirrorvmvolpath)
										// 미러링 이미지 상태가 Peer와 정상적으로 ready, sync 인 경우
										if volStatus.Description == "local image is primary" && strings.Contains(volStatus.PeerSites[0].State, "replaying") && strings.Contains(volStatus.PeerSites[0].Description, "idle") {
											params2 := []utils.MoldParams{
												{"virtualmachineid": vm[k].Id},
											}
											volResult := utils.GetListVolumes(params2)
											listVolumes := model.ListVolumes{}
											volInfo, _ := json.Marshal(volResult["listvolumesresponse"])
											json.Unmarshal([]byte(volInfo), &listVolumes)
											vol := listVolumes.Volume
											for v := 0; v < len(vol); v++ {
												imageList = append(imageList, vol[v].Path)
											}
											if vm[k].Hostname != "" {
												hostName = vm[k].Hostname
											} else {
												hostName = ""
											}
											clock = beforeIt
											if beforeIt != afterIt {
												println("updateScheduler : ", jobID.String(), jobName, currentTime.Format("2006-01-02 15:04:05"))
												scheduler.Update(
													uuid.MustParse(imageName),
													gocron.DurationJob(
														afterIt,
													),
													gocron.NewTask(
														func() {
															goCronTask(poolName, hostName, vmName, imageList)
														},
													),
													gocron.WithEventListeners(
														gocron.BeforeJobRuns(
															func(jobID uuid.UUID, jobName string) {
																hostName, clock, imageList = goCronEventListeners(scheduler, jobID, afterIt, jobName, imageName, hostName, vmName, poolName)
																afterIt = clock
															}),
													),
												)
											}
											break
										} else {
											// 미러링 이미지 상태가 Peer와 정상적으로 ready, resync 아닌 경우
											hostName = ""
											imageList = make([]string, 0)
											clock = beforeIt
										}
									}
								}
								break
							}
						}
					}
				}
				if exist != "exist" {
					println("non exist shutdown for scheduler image path : " + imageName)
					hostName = ""
					imageList = make([]string, 0)
					scheduler.Shutdown()
				}
			} else {
				scheduler.Shutdown()
			}
		}
	}
	return hostName, clock, imageList
}

func ImageMirroringSnap(poolName, hostName, vmName string, imageName []string) (output string, err error) {
	var stdout []byte
	currentTime := time.Now()
	println("start mirror snapshot ImageMirroringSnap --- vm : " + vmName + " --- image : " + strings.Join(imageName, ",") + " --- host : " + hostName + " --- date : " + currentTime.Format("2006-01-02 15:04:05"))
	if hostName != "" {
		println("start domfsfreeze ---")
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", hostName, "virsh", "domfsfreeze", vmName)
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			println("failed to virsh domfsfreeze")
			println(string(stdout))
		}
	}
	if len(imageName) > 0 {
		for i := 0; i < len(imageName); i++ {
			cmd := exec.Command(poolName, "mirror", "image", "snapshot", poolName+"/"+imageName[i])
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				println("failed to create rbd mirror image snapshot path : " + imageName[i])
				println(string(stdout))
				exec.Command("ssh", hostName, "virsh", "domfsthaw", vmName)
				break
			}
			host, _ := os.Hostname()
			cmd = exec.Command("rbd", "image-meta", "set", "rbd/MOLD-DR", imageName[i], currentTime.Format("2006-01-02 15:04:05")+","+host)
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				println("failed to update image-meta")
				println(string(stdout))
			}
		}
	}
	if hostName != "" {
		println("start domfsthaw ---")
		cmd := exec.Command("ssh", hostName, "virsh", "domfsthaw", vmName)
		stdout, err = cmd.CombinedOutput()
		if err != nil {
			println("failed to virsh domfsthaw")
			println(string(stdout))
		}
	}
	println("end mirror snapshot ImageMirroringSnap --- vm : " + vmName + " --- image : " + strings.Join(imageName, ",") + " --- host : " + hostName + " --- date : " + currentTime.Format("2006-01-02 15:04:05"))
	output = string(stdout)
	return
}

func ImageConfigSchedule(poolName, imageName, hostName, vmName, interval string) (output string, err error) {

	var beforeIt, clock time.Duration
	var imageList []string

	if strings.Contains(interval, "d") {
		interval = strings.TrimRight(interval, "d\n")
		ti, _ := strconv.Atoi(interval)
		beforeIt = time.Duration(ti) * 24 * time.Hour
	} else if strings.Contains(interval, "h") {
		interval = strings.TrimRight(interval, "h\n")
		ti, _ := strconv.Atoi(interval)
		beforeIt = time.Duration(ti) * time.Hour
	} else if strings.Contains(interval, "m") {
		interval = strings.TrimRight(interval, "m\n")
		ti, _ := strconv.Atoi(interval)
		beforeIt = time.Duration(ti) * time.Minute
	} else {
		err = errors.Join(err, errors.New("The interval must include d, h, and m, and the scheduler setup failed because it was set incorrectly."))
		utils.FancyHandleError(err)
		return
	}

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		err = errors.Join(err, errors.New("failed to create mirror image snapshot scheduler."))
		utils.FancyHandleError(err)
		return
	}

	j, err := scheduler.NewJob(
		gocron.DurationJob(
			beforeIt,
		),
		gocron.NewTask(
			func() {
				goCronTask(poolName, hostName, vmName, imageList)
			},
		),
		gocron.WithIdentifier(uuid.MustParse(imageName)),
		gocron.WithName(vmName),
		gocron.WithEventListeners(
			gocron.BeforeJobRuns(
				func(jobID uuid.UUID, jobName string) {
					hostName, clock, imageList = goCronEventListeners(scheduler, jobID, beforeIt, jobName, imageName, hostName, vmName, poolName)
					beforeIt = clock
				}),
		),
	)
	if err != nil {
		err = errors.Join(err, errors.New("failed to create mirror image snapshot scheduler job."))
		utils.FancyHandleError(err)
		return
	}
	scheduler.Start()

	output = "Success"
	println(j.ID().ID())
	return
}

func ImageUpdate(poolName string, imageName string, interval string, startTime string, schedule []model.MirrorImageItem) (output string, err error) {

	var stdoutScheduleEnable []byte
	var strScheduleOutput *exec.Cmd

	for _, scd := range schedule {
		if scd.StartTime == "" {
			strScheduleOutput = exec.Command("rbd", "mirror", "snapshot", "schedule", "rm", "--pool", poolName, "--image", imageName, scd.Interval)
		} else {
			strScheduleOutput = exec.Command("rbd", "mirror", "snapshot", "schedule", "rm", "--pool", poolName, "--image", imageName, scd.Interval, scd.StartTime)
		}
		stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()
		if err != nil {
			err = errors.Join(err, errors.New(string(stdoutScheduleEnable)))
			utils.FancyHandleError(err)
			return
		}
	}

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

func ImageRemoteUpdate(poolName string, imageName string, interval string, startTime string) (output string, err error) {

	var stdoutScheduleEnable []byte
	var strScheduleOutput *exec.Cmd

	mirrorConfig, err := GetConfigure()
	if err != nil {
		return
	}

	if startTime == "" {
		strScheduleOutput = exec.Command("rbd", "-c", mirrorConfig.ClusterFileName, "--cluster", mirrorConfig.ClusterName, "--name", mirrorConfig.Peers[0].ClientName, "--keyfile", mirrorConfig.KeyFileName, "mirror", "snapshot", "schedule", "add", "--pool", poolName, "--image", imageName, interval)
	} else {
		strScheduleOutput = exec.Command("rbd", "-c", mirrorConfig.ClusterFileName, "--cluster", mirrorConfig.ClusterName, "--name", mirrorConfig.Peers[0].ClientName, "--keyfile", mirrorConfig.KeyFileName, "mirror", "snapshot", "schedule", "add", "--pool", poolName, "--image", imageName, interval, startTime)
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
	if err != nil {
		err = errors.Join(err, errors.New(string(stdoutScheduleEnable)))
		utils.FancyHandleError(err)
		return
	}

	if err = json.Unmarshal(stdoutScheduleEnable, &imageStatus); err != nil {
		utils.FancyHandleError(err)
		return
	}

	return imageStatus, err
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
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		utils.FancyHandleError(err)
		return
	}

	// Mirror Daemon Deploy
	cmd = exec.Command("ceph", "orch", "apply", "rbd-mirror")
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		utils.FancyHandleError(err)
		return
	}

	// Mirror Bootstrap
	cmd = exec.Command("rbd", "mirror", "pool", "peer", "bootstrap", "create", "--site-name", dat.LocalClusterName, "-p", dat.MirrorPool)
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	DecodedLocalToken, err := base64.StdEncoding.DecodeString(string(stdout))
	if err != nil {
		utils.FancyHandleError(err)
		return
	}

	if err = json.Unmarshal(DecodedLocalToken, &LocalToken); err != nil {
		utils.FancyHandleError(err)
		return
	}

	cmd = exec.Command("ceph", "auth", "caps", "client."+LocalToken.ClientId, "mgr", "profile rbd", "mon", "profile rbd-mirror-peer", "osd", "profile rbd")
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		utils.FancyHandleError(err)
		return
	}

	cmd = exec.Command("ceph", "auth", "get-key", "client."+LocalToken.ClientId, "--format", "json")
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	if err = json.Unmarshal(stdout, &LocalKey); err != nil {
		utils.FancyHandleError(err)
		return
	}

	// Generate Token
	LocalToken.Key = LocalKey.Key
	JsonLocalKey, err := json.Marshal(LocalToken)
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	EncodedLocalToken = base64.StdEncoding.EncodeToString(JsonLocalKey)
	localTokenFile, err := os.CreateTemp("", "localtoken")
	if err != nil {
		utils.FancyHandleError(err)
		return
	}

	// defer localTokenFile.Close()
	// defer os.Remove(localTokenFile.Name())
	localTokenFile.WriteString(EncodedLocalToken)

	// For Remote
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
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	// println("out: " + string(stdout))
	// println("err: " + out.String())
	if err != nil {
		utils.FancyHandleError(err)
		return
	}

	// Mirror Daemon Deploy
	sshcmd, err = client.Command("ceph", "orch", "apply", "rbd-mirror")
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
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
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	//println("out: " + string(stdout))
	//println("err: " + out.String())
	if err != nil {
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

	sshcmd, err = client.Command("ceph", "auth", "caps", "client."+RemoteToken.ClientId, "mgr", "'profile rbd'", "mon", "'profile rbd-mirror-peer'", "osd", "'profile rbd'")
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}

	sshcmd, err = client.Command("ceph", "auth", "get-key", "client."+RemoteToken.ClientId, "--format", "json")
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	if err = json.Unmarshal(stdout, &RemoteKey); err != nil {
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
	remoteTokenFile, err := os.CreateTemp("", "remotetoken")
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	remoteTokenFile.WriteString(EncodedRemoteToken)

	// token import
	sshcmd, err = client.Command("echo", EncodedLocalToken, ">", remoteTokenFileName)
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}

	sshcmd, err = client.Command("rbd", "mirror", "pool", "info", "--pool", dat.MirrorPool, "--format", "json")
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		sshcmd.Stderr = &out
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
				sshcmd.Stderr = &out
				err = errors.Join(err, errors.New(out.String()))
				utils.FancyHandleError(err)
				return
			}
			// sshcmd.Stderr = &out
			stdout, err = sshcmd.CombinedOutput()
			if err != nil {
				sshcmd.Stderr = &out
				err = errors.Join(err, errors.New(out.String()))
				utils.FancyHandleError(err)
				return
			}
		}
	}

	sshcmd, err = client.Command("rbd", "mirror", "pool", "peer", "bootstrap", "import", "--pool", dat.MirrorPool, "--token-path", remoteTokenFileName)
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}

	sshcmd, err = client.Command("rbd", "create", "--size", "1", "rbd/MOLD-DR")
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}

	sshcmd, err = client.Command("rbd", "image-meta", "set", "rbd/MOLD-DR", "interval", "1h")
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		sshcmd.Stderr = &out
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

	cmd = exec.Command("rbd", "mirror", "pool", "info", "--pool", dat.MirrorPool, "--format", "json")
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		cmd.Stderr = &out
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
			// cmd.Stderr = &out
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				cmd.Stderr = &out
				err = errors.Join(err, errors.New(out.String()))
				utils.FancyHandleError(err)
				return
			}
		}
	}

	cmd = exec.Command("rbd", "mirror", "pool", "peer", "bootstrap", "import", "--pool", dat.MirrorPool, "--token-path", remoteTokenFile.Name())
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		cmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}

	cmd = exec.Command("rbd", "create", "--size", "1", "rbd/MOLD-DR")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		cmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}

	cmd = exec.Command("rbd", "image-meta", "set", "rbd/MOLD-DR", "interval", "1h")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		cmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	return
}

func ConfigMold(moldUrl, moldApiKey, moldSecretKey string) (err error) {

	mold := model.Mold{MoldUrl: moldUrl, MoldApiKey: moldApiKey, MoldSecretKey: moldSecretKey}

	jsonFile, _ := json.MarshalIndent(mold, "", " ")
	os.WriteFile("./mold.json", jsonFile, 0644)

	var stdout []byte

	cmd := exec.Command("sh", "-c", "cat /etc/hosts | grep -E 'scvm.*-mngt' | grep -v $(hostname) | awk '{print $1}'")
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	var str []string
	ipAddress := strings.Split(string(stdout), "\n")
	for i := 0; i < len(ipAddress); i++ {
		strs := ipAddress[i]
		str = append(str, strs)
		if i == len(ipAddress)-1 {
			str = str[:len(ipAddress)-1]
		}
	}
	for j := 0; j < len(str); j++ {
		cmd := exec.Command("sh", "-c", "scp -o StrictHostKeyChecking=no /usr/local/glue-api/mold.json "+str[j]+":/usr/local/glue-api/mold.json")
		// cmd.Stderr = &out
		stdout, err = cmd.CombinedOutput()
		println(string(stdout))
		if err != nil {
			utils.FancyHandleError(err)
			return
		}
	}
	return
}

func ImageMetaUpdate(interval string) (err error) {

	var stdout []byte
	cmd := exec.Command("rbd", "image-meta", "set", "rbd/MOLD-DR", "interval", interval)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.Join(err, errors.New(string(stdout)))
		utils.FancyHandleError(err)
		return
	}
	return
}

func ImageMetaRemove(imageName string) (output string, err error) {

	var stdout []byte
	cmd := exec.Command("rbd", "image-meta", "remove", "rbd/MOLD-DR", imageName)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(stdout), "no existing metadata key") {
			output = "Success"
			return
		} else {
			err = errors.Join(err, errors.New(string(stdout)))
			utils.FancyHandleError(err)
			return
		}
	}
	output = string(stdout)
	return
}

func ImageMetaGetTime(imageName string) (output string, err error) {

	var stdout []byte
	cmd := exec.Command("rbd", "image-meta", "get", "rbd/MOLD-DR", imageName)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.Join(err, errors.New(string(stdout)))
		utils.FancyHandleError(err)
		return
	}
	output = string(stdout)
	return
}

func ImageMetaGetInterval() (output string, err error) {

	var stdout []byte
	cmd := exec.Command("rbd", "image-meta", "get", "rbd/MOLD-DR", "interval")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.Join(err, errors.New(string(stdout)))
		utils.FancyHandleError(err)
		return
	}
	output = string(stdout)
	return
}

func ImagePromote(poolName string, imageName string) (output string, err error) {

	var stdoutScheduleEnable []byte
	strScheduleOutput := exec.Command("rbd", "mirror", "image", "promote", "--pool", poolName, "--image", imageName, "--force")
	stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()
	if strings.Contains(string(stdoutScheduleEnable), "unrecognised option") {
		strScheduleOutput = exec.Command("rbd", "mirror", "image", "promote", "--pool", poolName, "--image", imageName)
		stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()
	}
	if !strings.Contains(string(stdoutScheduleEnable), "Image promoted") {
		err = errors.Join(err, errors.New(string(stdoutScheduleEnable)))
	}
	if err != nil {
		err = errors.Join(err, errors.New(string(stdoutScheduleEnable)))
		utils.FancyHandleError(err)
		return
	}

	output = string(stdoutScheduleEnable)
	return
}

func ImageDemote(poolName string, imageName string) (output string, err error) {

	var stdoutScheduleEnable []byte

	strScheduleOutput := exec.Command("rbd", "mirror", "image", "demote", "--pool", poolName, "--image", imageName, "--force")
	stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()
	if strings.Contains(string(stdoutScheduleEnable), "unrecognised option") {
		strScheduleOutput = exec.Command("rbd", "mirror", "image", "demote", "--pool", poolName, "--image", imageName)
		stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()
	}
	if !strings.Contains(string(stdoutScheduleEnable), "Image demoted") {
		err = errors.Join(err, errors.New(string(stdoutScheduleEnable)))
	}
	if err != nil {
		utils.FancyHandleError(err)
	}

	output = string(stdoutScheduleEnable)
	return
}

func RemoteImagePromote(poolName string, imageName string) (output string, err error) {

	var stdoutScheduleEnable []byte
	conf, err := GetConfigure()
	strScheduleOutput := exec.Command("rbd", "-c", conf.ClusterFileName, "--cluster", conf.ClusterName, "--name", conf.Peers[0].ClientName, "--keyfile", conf.KeyFileName, "mirror", "image", "promote", "--pool", poolName, "--image", imageName, "--force")
	stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()
	if strings.Contains(string(stdoutScheduleEnable), "unrecognised option") {
		strScheduleOutput = exec.Command("rbd", "-c", conf.ClusterFileName, "--cluster", conf.ClusterName, "--name", conf.Peers[0].ClientName, "--keyfile", conf.KeyFileName, "mirror", "image", "promote", "--pool", poolName, "--image", imageName)
		stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()
	}
	if !strings.Contains(string(stdoutScheduleEnable), "Image promoted") {
		err = errors.Join(err, errors.New(string(stdoutScheduleEnable)))
	}
	if err != nil {
		utils.FancyHandleError(err)
	}

	output = string(stdoutScheduleEnable)
	return
}

func RemoteImageDemote(poolName string, imageName string) (output string, err error) {

	var stdoutScheduleEnable []byte
	conf, err := GetConfigure()
	strScheduleOutput := exec.Command("rbd", "-c", conf.ClusterFileName, "--cluster", conf.ClusterName, "--name", conf.Peers[0].ClientName, "--keyfile", conf.KeyFileName, "mirror", "image", "demote", "--pool", poolName, "--image", imageName, "--force")
	stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()
	if strings.Contains(string(stdoutScheduleEnable), "unrecognised option") {
		strScheduleOutput = exec.Command("rbd", "-c", conf.ClusterFileName, "--cluster", conf.ClusterName, "--name", conf.Peers[0].ClientName, "--keyfile", conf.KeyFileName, "mirror", "image", "demote", "--pool", poolName, "--image", imageName)
		stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()
	}
	if !strings.Contains(string(stdoutScheduleEnable), "Image demoted") {
		err = errors.Join(err, errors.New(string(stdoutScheduleEnable)))
	}
	if err != nil {
		utils.FancyHandleError(err)
	}

	output = string(stdoutScheduleEnable)
	return
}

func RemoteImageResync(poolName string, imageName string) (output string, err error) {

	var stdoutScheduleEnable []byte
	conf, err := GetConfigure()
	strScheduleOutput := exec.Command("rbd", "-c", conf.ClusterFileName, "--cluster", conf.ClusterName, "--name", conf.Peers[0].ClientName, "--keyfile", conf.KeyFileName, "mirror", "image", "resync", "--pool", poolName, "--image", imageName)
	stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()

	if !strings.Contains(string(stdoutScheduleEnable), "Flagged image") {
		err = errors.Join(err, errors.New(string(stdoutScheduleEnable)))
	}
	if err != nil {
		utils.FancyHandleError(err)
	}

	output = string(stdoutScheduleEnable)
	return
}

func ImageResync(poolName string, imageName string) (output string, err error) {

	var stdoutScheduleEnable []byte

	strScheduleOutput := exec.Command("rbd", "mirror", "image", "resync", "--pool", poolName, "--image", imageName)
	stdoutScheduleEnable, err = strScheduleOutput.CombinedOutput()

	if !strings.Contains(string(stdoutScheduleEnable), "Flagged image") {
		err = errors.Join(err, errors.New(string(stdoutScheduleEnable)))
	}
	if err != nil {
		utils.FancyHandleError(err)
	}

	output = string(stdoutScheduleEnable)
	return
}

func EnableMirror(dat model.MirrorSetup, privkeyname string) (EncodedLocalToken string, EncodedRemoteToken string, err error) {

	var LocalToken model.MirrorToken
	var RemoteToken model.MirrorToken
	var out strings.Builder
	var LocalKey model.AuthKey
	var RemoteKey model.AuthKey
	var stdout []byte

	remoteTokenFileName := "/tmp/remoteToken"

	// Mirror Enable
	cmd := exec.Command("rbd", "mirror", "pool", "enable", "--site-name", dat.LocalClusterName, "-p", dat.MirrorPool, "image")
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		utils.FancyHandleError(err)
		return
	}

	// Mirror Bootstrap
	cmd = exec.Command("rbd", "mirror", "pool", "peer", "bootstrap", "create", "--site-name", dat.LocalClusterName, "-p", dat.MirrorPool)
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	DecodedLocalToken, err := base64.StdEncoding.DecodeString(string(stdout))
	if err != nil {
		utils.FancyHandleError(err)
		return
	}

	if err = json.Unmarshal(DecodedLocalToken, &LocalToken); err != nil {
		utils.FancyHandleError(err)
		return
	}

	cmd = exec.Command("ceph", "auth", "get-key", "client."+LocalToken.ClientId, "--format", "json")
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	if err = json.Unmarshal(stdout, &LocalKey); err != nil {
		utils.FancyHandleError(err)
		return
	}

	// Generate Token
	LocalToken.Key = LocalKey.Key
	JsonLocalKey, err := json.Marshal(LocalToken)
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	EncodedLocalToken = base64.StdEncoding.EncodeToString(JsonLocalKey)
	localTokenFile, err := os.CreateTemp("", "localtoken")
	if err != nil {
		utils.FancyHandleError(err)
		return
	}

	// defer localTokenFile.Close()
	// defer os.Remove(localTokenFile.Name())
	localTokenFile.WriteString(EncodedLocalToken)

	// For Remote
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
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	// println("out: " + string(stdout))
	// println("err: " + out.String())
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
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	//println("out: " + string(stdout))
	//println("err: " + out.String())
	if err != nil {
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

	sshcmd, err = client.Command("ceph", "auth", "get-key", "client."+RemoteToken.ClientId, "--format", "json")
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	if err = json.Unmarshal(stdout, &RemoteKey); err != nil {
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
	remoteTokenFile, err := os.CreateTemp("", "remotetoken")
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	remoteTokenFile.WriteString(EncodedRemoteToken)

	// token import
	sshcmd, err = client.Command("echo", EncodedLocalToken, ">", remoteTokenFileName)
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}

	sshcmd, err = client.Command("rbd", "mirror", "pool", "info", "--pool", dat.MirrorPool, "--format", "json")
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		sshcmd.Stderr = &out
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
				sshcmd.Stderr = &out
				err = errors.Join(err, errors.New(out.String()))
				utils.FancyHandleError(err)
				return
			}
			// sshcmd.Stderr = &out
			stdout, err = sshcmd.CombinedOutput()
			if err != nil {
				sshcmd.Stderr = &out
				err = errors.Join(err, errors.New(out.String()))
				utils.FancyHandleError(err)
				return
			}
		}
	}

	sshcmd, err = client.Command("rbd", "mirror", "pool", "peer", "bootstrap", "import", "--pool", dat.MirrorPool, "--token-path", remoteTokenFileName)
	if err != nil {
		sshcmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	// sshcmd.Stderr = &out
	stdout, err = sshcmd.CombinedOutput()
	if err != nil {
		sshcmd.Stderr = &out
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

	cmd = exec.Command("rbd", "mirror", "pool", "info", "--pool", dat.MirrorPool, "--format", "json")
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		cmd.Stderr = &out
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
			// cmd.Stderr = &out
			stdout, err = cmd.CombinedOutput()
			if err != nil {
				cmd.Stderr = &out
				err = errors.Join(err, errors.New(out.String()))
				utils.FancyHandleError(err)
				return
			}
		}
	}

	cmd = exec.Command("rbd", "mirror", "pool", "peer", "bootstrap", "import", "--pool", dat.MirrorPool, "--token-path", remoteTokenFile.Name())
	// cmd.Stderr = &out
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		cmd.Stderr = &out
		err = errors.Join(err, errors.New(out.String()))
		utils.FancyHandleError(err)
		return
	}
	return
}
