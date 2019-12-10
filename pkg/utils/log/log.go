package log

import (
	"sync"

	"github.com/labstack/gommon/log"
)

var (
	Logger *log.Logger
	once   sync.Once
)

func InitLog(app string) {
	once.Do(func() {
		Logger = log.New(app)
	})
}
