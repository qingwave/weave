package main

import (
	"flag"
	"os"

	"github.com/qingwave/weave/pkg/config"
	"github.com/qingwave/weave/pkg/server"
	"github.com/qingwave/weave/pkg/version"

	"github.com/sirupsen/logrus"
)

var (
	printVersion = flag.Bool("v", false, "print version")
	appConfig    = flag.String("config", "config/app.yaml", "application config path")
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

	if *printVersion {
		version.Print()
		os.Exit(0)
	}

	logger := logrus.StandardLogger()
	logger.SetFormatter(&logrus.JSONFormatter{})

	conf, err := config.Parse(*appConfig)
	if err != nil {
		logger.Fatalf("Failed to parse config: %v", err)
	}

	s, err := server.New(conf, logger)
	if err != nil {
		logger.Fatalf("Init server failed: %v", err)
	}

	if err := s.Run(); err != nil {
		logger.Fatalf("Run server failed: %v", err)
	}
}
