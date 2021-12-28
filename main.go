package main

import (
	"flag"

	"weave/pkg/config"
	"weave/pkg/server"

	"github.com/sirupsen/logrus"
)

// @title           Weave Server API
// @version         2.0
// @description     This is a weave server api doc.

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
func main() {
	flag.Parse()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	conf, err := config.Parse()
	if err != nil {
		logger.Fatalf("Failed to parse config: %v", err)
	}

	s, err := server.New(conf, logger)
	if err != nil {
		logger.Fatalf("Init server failed: %v", err)
	}

	s.Run()
}
