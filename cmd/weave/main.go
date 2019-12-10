package main

import (
	"flag"
	"fmt"

	"github.com/stretchr/graceful"
	"qinng.io/weave/pkg/database"
	"qinng.io/weave/pkg/router"
	"qinng.io/weave/pkg/utils/log"
)

var (
	port = flag.Int("port", 8080, "server port")
	app  = "weave"
)

func main() {
	flag.Parse()

	address := fmt.Sprintf("0.0.0.0:%d", *port)

	log.InitLog(app)

	if err := database.InitMysql(); err != nil {
		log.Logger.Fatalf("Init database failed: %v", err)
	}

	s, err := router.Routers()
	if err != nil {
		log.Logger.Fatalf("Init router failed: %v", err)
	}

	log.Logger.Infof("Server started on: %s", address)
	graceful.Run(address, 0, s.Server.Handler)
}
