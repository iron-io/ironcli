package mq

import (
	"testing"

	"github.com/iron-io/ironcli/common"
)

func TestMqCreate(t *testing.T) {
	var (
		settings = &common.Settings{Product: "iron_mq"}
	)

	common.SetSettings(settings)
	mqCreate := NewMqCreate(settings)

	err := mqCreate.Action("testQueue", settings)
	if err != nil {
		t.Error(err)
	}
}

func TestMqClear(t *testing.T) {
	var (
		settings = &common.Settings{Product: "iron_mq"}
	)

	common.SetSettings(settings)
	mqClear := NewMqClear(settings)

	err := mqClear.Action("testQueue", settings)
	if err != nil {
		t.Error(err)
	}
}

func TestMqPush(t *testing.T) {
	var (
		settings = &common.Settings{Product: "iron_mq"}
	)

	common.SetSettings(settings)
	mqPush := NewMqPush(settings)

	err := mqPush.Action("testQueue", []string{"test1", "test2"}, settings)
	if err != nil {
		t.Error(err)
	}
}

func TestMqPop(t *testing.T) {
	var (
		settings = &common.Settings{Product: "iron_mq"}
	)

	common.SetSettings(settings)
	mqPop := NewMqPop(settings)

	err := mqPop.Action("testQueue", settings)
	if err != nil {
		t.Error(err)
	}
}

func TestMqReserve(t *testing.T) {
	var (
		settings = &common.Settings{Product: "iron_mq"}
	)

	common.SetSettings(settings)
	mqReserve := NewMqReserve(settings)

	err := mqReserve.Action("testQueue", settings)
	if err != nil {
		t.Error(err)
	}
}
