package mq

import (
	"errors"
	"fmt"

	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/iron_go3/mq"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type MqPeek struct {
	number int

	cli.Command
}

func NewMqPeek(settings *config.Settings) *MqPeek {
	mqPeek := &MqPeek{}

	mqPeek.Command = cli.Command{
		Name:      "peek",
		Usage:     "peek messages by amount of queue",
		ArgsUsage: "[QUEUE_NAME] [args]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:        "number, n",
				Usage:       "number usage",
				Destination: &mqPeek.number,
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().First() == "" {
				return errors.New(`peek requires one arg`)
			}

			q := mq.ConfigNew(c.Args().First(), settings)

			msgs, err := q.PeekN(mqPeek.number)
			if err != nil {
				return err
			}

			if len(msgs) < 1 {
				return errors.New("Queue is empty.")
			}

			if !common.IsPipedOut() {
				plural := ""
				if mqPeek.number > 1 {
					plural = "s"
				}

				fmt.Println(common.LINES, "Message", plural, " successfully peeked")
				fmt.Println()
				fmt.Println("-------- ID ------ | Body")
			}
			common.PrintMessages(msgs)

			return nil
		},
	}

	return mqPeek
}

func (r MqPeek) GetCmd() cli.Command {
	return r.Command
}
