package mq

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/iron-io/iron_go3/mq"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type MqPop struct {
	number     int
	outputfile string
	file       *os.File

	cli.Command
}

func NewMqPop(settings *common.Settings) *MqPop {
	mqPop := &MqPop{}

	mqPop.Command = cli.Command{
		Name:      "pop",
		Usage:     "pop messages by amount of queue",
		ArgsUsage: "[QUEUE_NAME]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:        "number, n",
				Usage:       "",
				Destination: &mqPop.number,
			},
			cli.StringFlag{
				Name:        "output, o",
				Usage:       "",
				Destination: &mqPop.outputfile,
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().First() == "" {
				return errors.New(`pop requires a queue name`)
			}

			if mqPop.outputfile != "" {
				f, err := os.Create(mqPop.outputfile)
				if err != nil {
					return err
				}

				mqPop.file = f
			}

			q := mq.ConfigNew(c.Args().First(), &settings.Worker)

			messages, err := q.PopN(mqPop.number)
			if err != nil {
				return err
			}

			// If anything here fails, we still want to print out what was deleted before exiting
			if mqPop.file != nil {
				b, err := json.Marshal(messages)
				if err != nil {
					common.PrintMessages(messages)

					return err
				}

				_, err = mqPop.file.Write(b)
				if err != nil {
					common.PrintMessages(messages)

					return err
				}
			}

			if common.IsPipedOut() {
				common.PrintMessages(messages)
			} else {
				plural := ""
				if mqPop.number > 1 {
					plural = "s"
				}

				fmt.Println(common.LINES, "Message", plural, " successfully popped off ", q.Name)
				fmt.Println()
				fmt.Println("-------- ID ------ | Body")
				common.PrintMessages(messages)
			}

			return nil
		},
	}

	return mqPop
}

func (r MqPop) GetCmd() cli.Command {
	return r.Command
}
