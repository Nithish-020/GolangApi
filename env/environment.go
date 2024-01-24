package env

import (
	"fcs23pkg/common"
	"fcs23pkg/genpkg"
	"fmt"
	"log"
)

func SetNovoEnvironment() {
	config := genpkg.ReadTomlConfig("../NOVOEnvConfig.toml")

	common.ABHIDomain = fmt.Sprintf("%v", config.(map[string]interface{})["ABHIDomain"])
	common.ABHIAppName = fmt.Sprintf("%v", config.(map[string]interface{})["ABHIAppName"])

	log.Println("common.ABHIAppName", common.ABHIAppName)
	log.Println("common.ABHIDomain", common.ABHIDomain)
}
