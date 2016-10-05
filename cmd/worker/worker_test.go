package worker

import (
	"strings"
	"testing"
	"time"

	"github.com/iron-io/ironcli/common"
)

func TestWorkerUpload(t *testing.T) {
	var settings = &common.Settings{Product: "iron_worker"}

	common.SetSettings(settings)

	workerUpload := NewWorkerUpload(settings)
	workerUpload.Zip = "../../testdata/test.zip"

	err := workerUpload.Action("iron/node:latest", []string{"node", "test.js"}, settings)
	if err != nil {
		t.Error(err)
	}
}

func TestWorkerSchedule(t *testing.T) {
	var settings = &common.Settings{Product: "iron_worker"}

	common.SetSettings(settings)

	workerSchedule := NewWorkerSchedule(settings)

	err := workerSchedule.Action("iron/node", settings)
	if err != nil {
		t.Error(err)
	}
}

func TestWorkerQueue(t *testing.T) {
	var settings = &common.Settings{Product: "iron_worker"}

	common.SetSettings(settings)

	workerQueue := NewWorkerQueue(settings)

	err := workerQueue.Action("iron/node", settings)
	if err != nil {
		t.Error(err)
	}
}

func TestWorkerStatus(t *testing.T) {
	var settings = &common.Settings{Product: "iron_worker"}

	common.SetSettings(settings)

	workerQueue := NewWorkerQueue(settings)

	err := workerQueue.Action("iron/node", settings)
	if err != nil {
		t.Error(err)
	}

	startTime := time.Now()

	// Wait a new task to getting a complete status
	for {
		workerStatus := NewWorkerStatus(settings)

		err = workerStatus.Action(workerQueue.TaskID, settings)
		if err != nil {
			t.Error(err)
		}

		if workerStatus.Status != "complete" && time.Now().Sub(startTime).Seconds() > 60 ||
			(workerStatus.Status == "error" || workerStatus.Status == "cancelled") {
			t.Error("Status should be complete")
			break
		} else if workerStatus.Status == "complete" {
			break
		}

		time.Sleep(2 * time.Second)
	}
}

func TestWorkerLog(t *testing.T) {
	var settings = &common.Settings{Product: "iron_worker"}

	common.SetSettings(settings)

	workerQueue := NewWorkerQueue(settings)

	err := workerQueue.Action("iron/node", settings)
	if err != nil {
		t.Error(err)
	}

	startTime := time.Now()

	// Wait a new task to getting a log
	for {
		workerLog := NewWorkerLog(settings)
		workerLog.Action(workerQueue.TaskID, settings)

		if !strings.Contains(workerLog.Log, "test") && time.Now().Sub(startTime).Seconds() > 60 {
			t.Error("Log has another output from script")
			break
		} else if strings.Contains(workerLog.Log, "test") {
			break
		}

		time.Sleep(2 * time.Second)
	}
}
