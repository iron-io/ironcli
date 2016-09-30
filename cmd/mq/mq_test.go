package mq

import (
	"testing"

	"github.com/iron-io/ironcli/common"
)

func TestMqCreate(t *testing.T) {
	var settings = &common.Settings{Product: "iron_mq"}

	common.SetSettings(settings)

	mqCreate := NewMqCreate(settings)

	err := mqCreate.Action("testQueue", settings)
	if err != nil {
		t.Error(err)
	}
}

func TestMqClear(t *testing.T) {
	var settings = &common.Settings{Product: "iron_mq"}

	common.SetSettings(settings)

	mqClear := NewMqClear(settings)

	err := mqClear.Action("testQueue", settings)
	if err != nil {
		t.Error(err)
	}
}

func TestMqPush(t *testing.T) {
	var settings = &common.Settings{Product: "iron_mq"}

	common.SetSettings(settings)

	mqPush := NewMqPush(settings)

	err := mqPush.Action("testQueue", []string{"test1", "test2", "test3", "test4"}, settings)
	if err != nil {
		t.Error(err)
	}
}

func TestMqPop(t *testing.T) {
	var settings = &common.Settings{Product: "iron_mq"}

	common.SetSettings(settings)

	mqPop := NewMqPop(settings)

	err := mqPop.Action("testQueue", settings)
	if err != nil {
		t.Error(err)
	}
}

func TestMqReserve(t *testing.T) {
	var settings = &common.Settings{Product: "iron_mq"}

	common.SetSettings(settings)

	mqReserve := NewMqReserve(settings)

	err := mqReserve.Action("testQueue", settings)
	if err != nil {
		t.Error(err)
	}
}

func TestMqPeek(t *testing.T) {
	var settings = &common.Settings{Product: "iron_mq"}

	common.SetSettings(settings)

	mqPeek := NewMqPeek(settings)

	err := mqPeek.Action("testQueue", settings)
	if err != nil {
		t.Error(err)
	}
}

func TestMqList(t *testing.T) {
	var settings = &common.Settings{Product: "iron_mq"}

	common.SetSettings(settings)

	mqList := NewMqList(settings)

	err := mqList.Action(settings)
	if err != nil {
		t.Error(err)
	}
}

func TestMqInfo(t *testing.T) {
	var settings = &common.Settings{Product: "iron_mq"}

	common.SetSettings(settings)

	mqInfo := NewMqInfo(settings)

	err := mqInfo.Action("testQueue", settings)
	if err != nil {
		t.Error(err)
	}
}

func TestMqDelete(t *testing.T) {
	var settings = &common.Settings{Product: "iron_mq"}

	common.SetSettings(settings)

	mqPush := NewMqPush(settings)
	err := mqPush.Action("testQueue", []string{"testForDelete"}, settings)
	if err != nil {
		t.Error(err)
	}

	mqDelete := NewMqDelete(settings)
	err = mqDelete.Action("testQueue", mqPush.ResultIds, settings)
	if err != nil {
		t.Error(err)
	}
}

func TestMqRm(t *testing.T) {
	var settings = &common.Settings{Product: "iron_mq"}

	common.SetSettings(settings)

	mqRm := NewMqRm(settings)
	err := mqRm.Action("testQueue", settings)
	if err != nil {
		t.Error(err)
	}
}
