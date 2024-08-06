package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/ini.v1"
)

var conf_path string = "/etc/samba/smb.conf"

func sectionList() string {
	// config 파일 읽기
	cfg, err := ini.Load(conf_path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// config 파일을 JSON 구조로 변환
	// jsonData := make(map[string]map[string]string)
	var sectionArray []string

	for _, section := range cfg.Sections() {
		sectionName := section.Name()
		if sectionName == ini.DEFAULT_SECTION {
			sectionName = "default"
		}
		if sectionName != "default" && sectionName != "global" {
			sectionArray = append(sectionArray, string(sectionName))
		}
	}

	// 리스트를 JSON 형식으로 변환
	jsonData, err := json.Marshal(sectionArray)
	if err != nil {
		fmt.Printf("Fail to marshal JSON: %v", err)
		os.Exit(1)
	}

	// JSON 출력
	jsonString := string(jsonData)
	fmt.Println(jsonString)
	return string(jsonString)
}

func arrayList() string {
	// config 파일 읽기
	cfg, err := ini.Load(conf_path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// config 파일을 JSON 구조로 변환
	jsonData := make(map[string]map[string]string)

	for _, section := range cfg.Sections() {
		sectionName := section.Name()
		if sectionName == ini.DEFAULT_SECTION {
			sectionName = "default"
		}
		if sectionName != "default" && sectionName != "global" {
			jsonData[sectionName] = make(map[string]string)
			for _, key := range section.Keys() {
				jsonData[sectionName][key.Name()] = key.Value()
				if key.Name() == "path" {
					var stdout []byte
					cmd := exec.Command("findmnt", key.Value())
					stdout, _ = cmd.CombinedOutput()
					if string(stdout) != "" {
						jsonData[sectionName]["mount_yn"] = "true"
					} else {
						jsonData[sectionName]["mount_yn"] = "false"
					}
				}
			}
		}
	}

	// JSON으로 변환
	jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		fmt.Printf("Fail to marshal JSON: %v", err)
		os.Exit(1)
	}

	// JSON 출력
	fmt.Println(string(jsonBytes))
	return string(jsonBytes)
}

func confAdd(sectionStr string, pathStr string, vfs_objectsStr string, cacheStr string) {
	//global 세션 하위 주석처리
	confAddAnnotation()

	// smb.conf 파일 읽기
	cfg, err := ini.Load(conf_path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	section := sectionStr
	comment := "Share Directories"
	path := pathStr
	writable := "yes"
	public := "yes"
	create_mask := "0777"
	directory_mask := "0777"
	vfs_objects := vfs_objectsStr
	cache := cacheStr

	cfg.Section(section)
	cfg.Section(section).Key("comment").SetValue(comment)
	cfg.Section(section).Key("path").SetValue(path)
	cfg.Section(section).Key("writable").SetValue(writable)
	cfg.Section(section).Key("public").SetValue(public)
	cfg.Section(section).Key("create mask").SetValue(create_mask)
	cfg.Section(section).Key("directory mask").SetValue(directory_mask)

	// ads 일 경우 활성화
	if vfs_objects == "true" {
		cfg.Section(section).Key("vfs objects").SetValue("fake_compression")
	}

	if cache == "true" {
		cfg.Section(section).Key("csc policy").SetValue("programs")
	} else {
		cfg.Section(section).DeleteKey("csc policy")
	}

	// INI 파일에 변경 사항 저장
	err = cfg.SaveTo(conf_path)
	if err != nil {
		fmt.Printf("Fail to save file: %v", err)
		os.Exit(1)
	}

	//global 세션 하위 주석해제
	confRemoveAnnotation()
	// fmt.Println("conf 파일이 성공적으로 수정되었습니다.")
}

func confDelete(sectionStr string) {
	//global 세션 하위 주석처리
	confAddAnnotation()

	// config 파일 읽기
	cfg, err := ini.Load(conf_path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// 섹션 삭제
	cfg.DeleteSection(sectionStr)

	// config 파일에 변경 사항 저장
	err = cfg.SaveTo(conf_path)
	if err != nil {
		fmt.Printf("Fail to save file: %v", err)
		os.Exit(1)
	}
	//global 세션 하위 주석해제
	confRemoveAnnotation()
	// fmt.Println("conf 파일이 성공적으로 수정되었습니다.")
}

// smb.conf [global] 하위 주석처리
func confAddAnnotation() {
	// 파일을 열거나 생성합니다.
	file, err := os.Open(conf_path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 파일을 읽어들입니다.
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// 파일 내용을 처리합니다.
	var result []string
	insideGlobal := false
	for _, line := range lines {
		if strings.HasPrefix(line, "[global]") {
			insideGlobal = true
		} else if strings.HasPrefix(line, "[") && insideGlobal {
			insideGlobal = false
		}
		if insideGlobal {
			result = append(result, "# "+line) // 주석 처리합니다.
		} else {
			result = append(result, line)
		}
	}

	// 수정된 내용을 파일에 씁니다.
	err = os.WriteFile(conf_path, []byte(strings.Join(result, "\n")), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
	}
}

// smb.conf [global] 하위 주석해제
func confRemoveAnnotation() {
	// 파일을 열거나 생성합니다.
	file, err := os.Open(conf_path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 파일을 읽어들입니다.
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// 파일 내용을 처리합니다.
	var result []string
	insideGlobal := false
	for _, line := range lines {
		if strings.HasPrefix(line, "[global]") || strings.HasPrefix(line, "#[global]") || strings.HasPrefix(line, "# [global]") {
			insideGlobal = true
			if strings.HasPrefix(line, "[global]") {
				result = append(result, line)
			} else if strings.HasPrefix(line, "#[global]") {
				result = append(result, strings.TrimPrefix(line, "#"))
			} else if strings.HasPrefix(line, "# [global]") {
				result = append(result, strings.TrimPrefix(line, "# "))
			}
		} else if strings.HasPrefix(line, "[") && insideGlobal {
			insideGlobal = false
			result = append(result, line)
		} else if insideGlobal {
			if strings.HasPrefix(line, "# ") {
				// 주석을 제거하고 원래의 내용을 복원합니다.
				result = append(result, strings.TrimPrefix(line, "# "))
			} else {
				result = append(result, line)
			}
		} else {
			result = append(result, line)
		}
	}

	// 수정된 내용을 파일에 씁니다.
	err = os.WriteFile(conf_path, []byte(strings.Join(result, "\n")), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
	}
}

func main() {
	action := flag.String("a", "", "action")
	section := flag.String("s", "", "section value")
	path := flag.String("p", "", "path value")
	vfs_objects := flag.String("f", "fasle", "vfs objects")
	cache := flag.String("c", "false", "section value")
	flag.Parse()

	if *action == "" {
		fmt.Printf("Please enter action.")
		os.Exit(1)
	} else if *action == "sectionList" {
		sectionList()
	} else if *action == "arrayList" {
		arrayList()
	} else if *action == "confAdd" {
		confAdd(*section, *path, *vfs_objects, *cache)
	} else if *action == "confDelete" {
		confDelete(*section)
	} else {
		fmt.Printf("This action does not exist.")
		os.Exit(1)
	}
}
