package mq

import (
	"errors"
	"fmt"

	"github.com/iron-io/iron_go3/mq"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type MqPeek struct {
	number int

	cli.Command
}

func NewMqPeek(settings *common.Settings) *MqPeek {
	mqPeek := &MqPeek{}

	mqPeek.Command = cli.Command{
		Name:      "peek",
		Usage:     "peek at messages in the queue without dequeuing them.",
		ArgsUsage: "[QUEUE_NAME]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:        "number, n",
				Usage:       "",
				Destination: &mqPeek.number,
			},
		},
		Action: func(c *cli.Context) error {
			err := mqPeek.Action(c.Args().First(), settings)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return mqPeek
}

func (m MqPeek) GetCmd() cli.Command {
	return m.Command
}

func (m *MqPeek) Action(queueName string, settings *common.Settings) error {
	if queueName == "" {
		return errors.New(`peek requires one arg`)
	}

	q := mq.ConfigNew(queueName, &settings.Worker)

	msgs, err := q.PeekN(m.number)
	if err != nil {
		return err
	}

	if len(msgs) < 1 {
		return errors.New("Queue is empty.")
	}

	if !common.IsPipedOut() {
		plural := ""
		if m.number > 1 {
			plural = "s"
		}

		fmt.Println(common.LINES, "Message", plural, " successfully peeked")
		fmt.Println()
		fmt.Println("-------- ID ------ | Body")
	}
	common.PrintMessages(msgs)

	return nil
}
