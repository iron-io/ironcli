package mq

import (
	"errors"
	"fmt"

	"github.com/iron-io/iron_go3/mq"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type MqClear struct {
	cli.Command
}

func NewMqClear(settings *common.Settings) *MqClear {
	mqClear := &MqClear{}

	mqClear.Command = cli.Command{
		Name:      "clear",
		Usage:     "clear all messages of queue",
		ArgsUsage: "[QUEUE_NAME]",
		Before: func(c *cli.Context) error {
			if err := common.SetSettings(settings); err != nil {
				return err
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			err := mqClear.Action(c.Args().First(), settings)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return mqClear
}

func (m MqClear) GetCmd() cli.Command {
	return m.Command
}

func (m *MqClear) Action(queueName string, settings *common.Settings) error {
	if queueName == "" {
		return errors.New(`clear requires a queue name`)
	}

	q := mq.ConfigNew(queueName, &settings.Worker)

	if err := q.Clear(); err != nil {
		return fmt.Errorf("create error: %v", err)
	}

	fmt.Println(common.Green(common.LINES, "Queue ", q.Name, " has been successfully cleared"))

	return nil
}
