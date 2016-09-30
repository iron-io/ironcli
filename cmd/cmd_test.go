package cmd

import (
	"testing"

	"github.com/iron-io/ironcli/common"
)

func TestRun(t *testing.T) {
	var (
		settings = &common.Settings{Product: "iron_worker"}
	)

	common.SetSettings(settings)

	run := NewRun(settings)
	run.Zip = "../testdata/test.zip"

	err := run.Action("iron/node:latest", []string{"node", "test.js"}, settings)
	if err != nil {
		t.Error(err)
	}

	worker := common.Worker{Settings: settings.Worker}
	_, err = worker.CodePackageInfo(run.CodeID)
	if err != nil {
		t.Error(err)
	}
}

func TestRegister(t *testing.T) {
	var (
		settings = &common.Settings{Product: "iron_worker"}
	)

	common.SetSettings(settings)

	register := NewRegister(settings)

	err := register.Action("iron/node", []string{}, settings)
	if err != nil {
		t.Error(err)
	}

	worker := common.Worker{Settings: settings.Worker}
	_, err = worker.CodePackageInfo(register.CodeID)
	if err != nil {
		t.Error(err)
	}
}
