package main

import (
    "encoding/json"
    "fmt"
    "gopkg.in/ini.v1"
    "os"
    "flag"
    "os/exec"
)

func sectionList() (string){
	// config 파일 읽기
    cfg, err := ini.Load("/etc/samba/smb.conf")
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

func arrayList() (string){
	// config 파일 읽기
    cfg, err := ini.Load("/etc/samba/smb.conf")
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
                    if stdout != nil {
                        jsonData[sectionName]["mount_yn"] = "true"
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

func confAdd(sectionStr string, pathStr string, vfs_objectsStr string, cacheStr string) (){
    // smb.conf 파일 읽기
    cfg, err := ini.Load("/etc/samba/smb.conf")
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
    err = cfg.SaveTo("/etc/samba/smb.conf")
    if err != nil {
        fmt.Printf("Fail to save file: %v", err)
        os.Exit(1)
    }
    // fmt.Println("conf 파일이 성공적으로 수정되었습니다.")
}

func confDelete(sectionStr string) (){
        // config 파일 읽기
        cfg, err := ini.Load("/etc/samba/smb.conf")
        if err != nil {
            fmt.Printf("Fail to read file: %v", err)
            os.Exit(1)
        }
    
        // 섹션 삭제
        cfg.DeleteSection(sectionStr)
    
        // config 파일에 변경 사항 저장
        err = cfg.SaveTo("/etc/samba/smb.conf")
        if err != nil {
            fmt.Printf("Fail to save file: %v", err)
            os.Exit(1)
        }
        // fmt.Println("conf 파일이 성공적으로 수정되었습니다.")
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