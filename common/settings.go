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

func settingValues(settings *Settings) {
	oldSettings := *settings

	newSettings := config.ConfigWithEnv(settings.Product, settings.Env)
	settings.Worker = newSettings

	if oldSettings.Worker.Token != "" {
		settings.Worker.Token = oldSettings.Worker.Token
	}

	if oldSettings.Worker.ProjectId != "" {
		settings.Worker.ProjectId = oldSettings.Worker.ProjectId
	}

	if settings.Worker.ProjectId != "" {
		settings.HUDUrlStr = `Check https://hud-e.iron.io/worker/projects/` + settings.Worker.ProjectId + "/"
	}
}

func setMq(settings *Settings) error {
	settingValues(settings)

	if !IsPipedOut() {
		fmt.Printf("%sConfiguring client\n", LINES)

		pName, err := MqProjectName(settings.Worker)
		if err != nil {
			return err
		}

		if pName == "" {
			fmt.Printf("%sCould not find project name.", BLANKS)
		} else {
			fmt.Printf(`%sProject '%s' with id='%s'`, BLANKS, pName, settings.Worker.ProjectId)
		}

		fmt.Println()
	}

	return nil
}

func setProject(settings *Settings) error {
	settingValues(settings)

	fmt.Println(LINES, `Configuring client`)

	pName, err := projectName(settings)
	if err != nil {
		return err
	}

	fmt.Printf(`%s Project '%s' with id='%s'`, BLANKS, pName, settings.Worker.ProjectId)
	fmt.Println()

	return nil
}

func SetSettings(settings *Settings) error {
	var err error

	switch settings.Product {
	case "iron_worker":
		err = setProject(settings)
	case "iron_mq":
		err = setMq(settings)
	}

	return err
}
