package app

import (
	"github.com/mohsenhatami96/dobby/pkg/cloner"
	"github.com/mohsenhatami96/dobby/services"
)

type Injector struct {
	Cloner services.Cloner
}

func (injector *Injector) Inject(config Config) {
	injector.Cloner = cloner.New(config.GitURL, config.GitToken)
}
