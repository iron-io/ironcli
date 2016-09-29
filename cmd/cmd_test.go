package cmd

import (
	"testing"

	"github.com/iron-io/ironcli/common"
)

func TestRegister(t *testing.T) {
	var (
		settings = &common.Settings{Product: "iron_worker"}
	)

	common.SetSettings(settings)

	register := NewRegister(settings)

	err := register.Action("test", []string{}, settings)
	if err != nil {
		t.Error(err)
	}
}

func TestRun(t *testing.T) {
	var (
		settings = &common.Settings{Product: "iron_worker"}
	)

	common.SetSettings(settings)

	run := NewRun(settings)

	err := run.Action("test", []string{}, settings)
	if err != nil {
		t.Error(err)
	}
}
