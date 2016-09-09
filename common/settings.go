package common

import (
	"fmt"

	"github.com/iron-io/iron_go3/config"
)

type Settings struct {
	Product   string
	Env       string
	HUDUrlStr string

	Worker config.Settings
}

func SetSettings(settings *Settings) error {
	oldSettings := *settings

	newSettings := config.ConfigWithEnv("iron_worker", settings.Env)
	settings.Worker = newSettings

	if oldSettings.Worker.Token != "" {
		settings.Worker.Token = oldSettings.Worker.Token
	}

	if oldSettings.Worker.ProjectId != "" {
		settings.Worker.ProjectId = oldSettings.Worker.ProjectId
	}

	if settings.Worker.ProjectId != "" {
		settings.HUDUrlStr = `Check https://hud.iron.io/tq/projects/` + settings.Worker.ProjectId + "/"
	}

	fmt.Println(LINES, `Configuring client`)

	pName, err := projectName(settings)
	if err != nil {
		return err
	}

	fmt.Printf(`%s Project '%s' with id='%s'`, BLANKS, pName, settings.Worker.ProjectId)
	fmt.Println()

	newSettings = config.ConfigWithEnv(settings.Product, settings.Env)
	settings.Worker = newSettings

	return nil
}
