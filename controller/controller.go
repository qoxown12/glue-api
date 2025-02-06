package controller

import (
	"Glue-API/docs"
	"Glue-API/httputil"
	"Glue-API/model"
	"Glue-API/utils"
	"Glue-API/utils/glue"
	"Glue-API/utils/license"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

func LogSetting() {
	logFilePath := "/var/log/glue-api.log"
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		os.Create("/var/log/glue-api.log")
	}
	logFile, _ := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(logFile)
}

// Controller example
type Controller struct {
}

// NewController example
func NewController() *Controller {
	return &Controller{}
}

// Message example
type Message struct {
	Message string `json:"message" example:"message"`
} //@name Message

// Version godoc
//
//	@Summary		Show Versions of API
//	@Description	API 의 버전을 보여줍니다.
//	@Tags			API
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200	{object}	model.Version
//	@Failure		400	{object}	HTTP400BadRequest
//	@Failure		404	{object}	HTTP404NotFound
//	@Failure		500	{object}	HTTP500InternalServerError
//	@Router			/version [get]
func (c *Controller) Version(ctx *gin.Context) {
	dat := model.Version{Version: docs.SwaggerInfo.Version}
	// Print the output
	ctx.IndentedJSON(http.StatusOK, dat)

}
func SetOptionHeader(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "*")
	ctx.Header("Access-Control-Allow-Headers", "*")
	ctx.Header("Access-Control-Max-Age", "3600")
}
func GlueUrl() (output string) {
	dat, err := glue.GlueUrl()
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	url := strings.Split(dat.ActiveName, ".")

	var stdout []byte
	cmd := exec.Command("sh", "-c", "cat /etc/hosts | grep '"+url[0]+"-mngt' | awk '{print $1}'")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		return
	}
	result := strings.Replace(string(stdout), "\n", "", -1)
	settings, _ := utils.ReadConfFile()
	output = string(settings.GlueProtocol+"://") + result + string(":"+settings.GluePort+"/")
	return
}

func GetToken() (output string, err error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	settings, _ := utils.ReadConfFile()
	pw, err := utils.PasswordDecryption(settings.GluePw)

	user_json := model.UserInfo{
		Username: settings.GlueUser,
		Password: pw,
	}
	user, err := json.Marshal(user_json)
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	url := GlueUrl() + "api/auth"
	var jsonStr = bytes.NewBuffer(user)
	request, err := http.NewRequest(http.MethodPost, url, jsonStr)
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("accept", "application/vnd.ceph.api.v1.0+json")

	result, err := client.Do(request)
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		utils.FancyHandleError(err)
		return
	}
	var dat model.Token
	if err = json.Unmarshal(body, &dat); err != nil {
		utils.FancyHandleError(err)
		return
	}
	output = dat.Token
	return
}

func (c *Controller) License(ctx *gin.Context) {
	output, err := license.License()
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, output)
}
