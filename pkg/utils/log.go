package utils

import (
	"github.com/labstack/gommon/log"
)

var Logger *log.Logger

func init() {
	Logger = log.New("weave")
	Logger.Info("Init log weave")
}
