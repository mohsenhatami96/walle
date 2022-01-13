package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/mohsenhatami96/dobby/internal/app"
)

func main() {
	var config app.Config
	err := envconfig.Process("dobby", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	injector := app.Injector{}
	injector.Inject(config)
	injector.Cloner.CloneAll()
}
