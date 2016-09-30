package worker

import (
	"strings"
	"testing"
	"time"

	"github.com/iron-io/ironcli/common"
)

func TestWorkerUpload(t *testing.T) {
	var (
		settings = &common.Settings{Product: "iron_worker"}
	)

	common.SetSettings(settings)

	workerUpload := NewWorkerUpload(settings)
	workerUpload.Zip = "../../testdata/test.zip"

	err := workerUpload.Action("iron/node:latest", []string{"node", "test.js"}, settings)
	if err != nil {
		t.Error(err)
	}
}

func TestWorkerSchedule(t *testing.T) {
	var (
		settings = &common.Settings{Product: "iron_worker"}
	)

	common.SetSettings(settings)

	workerSchedule := NewWorkerSchedule(settings)

	err := workerSchedule.Action("iron/node", settings)
	if err != nil {
		t.Error(err)
	}
}

func TestWorkerQueue(t *testing.T) {
	var (
		settings = &common.Settings{Product: "iron_worker"}
	)

	common.SetSettings(settings)

	workerQueue := NewWorkerQueue(settings)

	err := workerQueue.Action("iron/node", settings)
	if err != nil {
		t.Error(err)
	}
}

func TestWorkerStatus(t *testing.T) {
	var (
		settings = &common.Settings{Product: "iron_worker"}
	)

	common.SetSettings(settings)

	workerQueue := NewWorkerQueue(settings)

	err := workerQueue.Action("iron/node", settings)
	if err != nil {
		t.Error(err)
	}

	// Wait a new task to getting a complete status
	time.Sleep(5 * time.Second)

	workerStatus := NewWorkerStatus(settings)

	err = workerStatus.Action(workerQueue.TaskID, settings)
	if err != nil {
		t.Error(err)
	}

	if workerStatus.Status != "complete" {
		t.Error("Status should be complete")
	}
}

func TestWorkerLog(t *testing.T) {
	var (
		settings = &common.Settings{Product: "iron_worker"}
	)

	common.SetSettings(settings)

	workerQueue := NewWorkerQueue(settings)

	err := workerQueue.Action("iron/node", settings)
	if err != nil {
		t.Error(err)
	}

	// Wait a new task to getting a log
	time.Sleep(5 * time.Second)

	workerLog := NewWorkerLog(settings)

	err = workerLog.Action(workerQueue.TaskID, settings)
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(workerLog.Log, "test") {
		t.Error("Log has another output from script")
	}
}
