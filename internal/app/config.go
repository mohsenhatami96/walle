package app

type Config struct {
	GitURL            string `required:"true" envconfig:"WAllE_GITLAB_URL"`
	GitToken          string `required:"true" envconfig:"WALLE_GITLAB_TOKEN"`
	GitUsername       string `envconfig:"WALLE_GITLAB_USERNAME" default:"root"`
	SSHAuth           bool   `envconfig:"WALLE_SSH_AUTH" default:"false"`
	SSHPrivateKeyPath string `envconfig:"WALLE_SSH_PKEY" default:"~/.ssh/id_rsa"`
}
