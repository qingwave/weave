package main

import (
	"flag"
	"fmt"

	"qinng.io/weave/pkg/server"
	"qinng.io/weave/pkg/utils"
)

var (
	port = flag.Int("port", 8080, "server port")
)

func main() {
	flag.Parse()

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	s := server.NewServer(address)
	utils.Logger.Infof("start server in %s", address)
	s.Run()
}
