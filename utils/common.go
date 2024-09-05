package utils

import (
	"Glue-API/model"
	"encoding/json"
	"log"
	"os"
	"runtime"
)

func HandleError(err error) (b bool) {
	if err != nil {
		// notice that we're using 1, so it will actually log where
		// the error happened, 0 = this function, we don't want that.
		_, filename, line, _ := runtime.Caller(1)
		log.Printf("[error] %s:%d %v", filename, line, err)
		b = true
	}
	return
}

// FancyHandleError this logs the function name as well.
func FancyHandleError(err error) (b bool) {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, filename, line, _ := runtime.Caller(1)

		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)

		b = true
	}
	return
}

// Read the settings file.
func ReadConfFile() (settings model.Settings, err error) {
	content, err := os.ReadFile("./conf.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
		return
	}

	err = json.Unmarshal(content, &settings)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
		return
	}
	return
}

// Read the mold settings file.
func ReadMoldFile() (mold model.Mold, err error) {
	content, err := os.ReadFile("./mold.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
		return
	}

	err = json.Unmarshal(content, &mold)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
		return
	}
	return
}
