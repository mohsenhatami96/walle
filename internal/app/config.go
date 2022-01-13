package app

type Config struct {
	GitURL   string `required:"true" envconfig:"DOBBY_GIT_URL"`
	GitToken string `required:"true" envconfig:"DOBBY_GIT_TOKEN"`
}
