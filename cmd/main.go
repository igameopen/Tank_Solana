package main

import (
	"transferSrv/infra/config"
	"transferSrv/infra/library/log"
	"transferSrv/server"
)

func main() {
	config.LoadCnf("../infra/config/config.toml")
	log.Init()

	server.StartHttpSrv()
}
