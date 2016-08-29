package mq

import (
	"fmt"

	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/iron_go3/mq"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type MqCreate struct {
	cli.Command
}

func NewMqCreate(settings *config.Settings) *MqCreate {
	mqCreate := &MqCreate{
		Command: cli.Command{
			Name:      "create",
			Usage:     "create queue",
			ArgsUsage: "[QUEUE_NAME]",
			Action: func(c *cli.Context) error {
				fmt.Printf("%sCreating queue \"%s\"\n", common.BLANKS, c.Args().First())

				q := mq.ConfigNew(c.Args().First(), settings)
				_, err := q.PushStrings("")
				if err != nil {
					return err
				}

				err = q.Clear()
				if err != nil {
					return err
				}

				fmt.Println(common.LINES, "Queue ", q.Name, " has been successfully created.")
				common.PrintQueueHudURL(common.BLANKS, q)

				return nil
			},
		},
	}

	return mqCreate
}

func (r MqCreate) GetCmd() cli.Command {
	return r.Command
}
