package main

import (
	"Neo4jPlayground/cmd/web"
	"Neo4jPlayground/internal/config"
	"Neo4jPlayground/internal/handlers"
)

func main() {
	var conf config.Conf
	conf.GetConf()
	handlers.InitDriver(&conf)
	web.Routes(&conf)
	defer handlers.CloseDriver()
}
