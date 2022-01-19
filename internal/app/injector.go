package app

import (
	"github.com/mohsenhatami96/walle/pkg/cloner"
	"github.com/mohsenhatami96/walle/services"
)

type Injector struct {
	Cloner services.Cloner
}

func (injector *Injector) Inject(config Config) {
	injector.Cloner = cloner.New(
		config.GitURL,
		config.GitToken,
		config.GitUsername,
		config.SSHAuth,
		config.SSHPrivateKeyPath,
	)
}
