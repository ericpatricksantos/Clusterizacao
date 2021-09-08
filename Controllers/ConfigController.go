package Controllers

import (
	"main.go/Function/Config"
	Model "main.go/Models"
)

func GetConfig() Model.Configuration {
	return Config.GetConfig()
}
