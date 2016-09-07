package common

import "github.com/iron-io/iron_go3/config"

type Settings struct {
	Product string
	Env     string

	Worker config.Settings
}

func SetSettings(settings *Settings) {
	newSettings := config.ConfigWithEnv(settings.Product, settings.Env)
	newSettings.Token = settings.Worker.Token
	newSettings.ProjectId = settings.Worker.ProjectId

	settings.Worker = newSettings
}
