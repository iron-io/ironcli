package mq

import (
	"fmt"

	"github.com/iron-io/iron_go3/mq"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type MqCreate struct {
	cli.Command
}

func NewMqCreate(settings *common.Settings) *MqCreate {
	mqCreate := &MqCreate{}

	mqCreate.Command = cli.Command{
		Name:      "create",
		Usage:     "create queue",
		ArgsUsage: "[QUEUE_NAME]",
		Before: func(c *cli.Context) error {
			if err := common.SetSettings(settings); err != nil {
				return err
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			err := mqCreate.Action(c.Args().First(), settings)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return mqCreate
}

func (m MqCreate) GetCmd() cli.Command {
	return m.Command
}

func (m *MqCreate) Action(queueName string, settings *common.Settings) error {
	fmt.Printf("%sCreating queue \"%s\"\n", common.BLANKS, queueName)

	q := mq.ConfigNew(queueName, &settings.Worker)
	_, err := q.PushStrings("")
	if err != nil {
		return err
	}

	err = q.Clear()
	if err != nil {
		return err
	}

	fmt.Println(common.LINES, "Queue", q.Name, "has been successfully created.")
	common.PrintQueueHudURL(common.BLANKS, q)

	return nil
}
