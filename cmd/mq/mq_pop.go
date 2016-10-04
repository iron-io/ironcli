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
		Before: func(c *cli.Context) error {
			if err := common.SetSettings(settings); err != nil {
				return err
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			err := mqPop.Action(c.Args().First(), settings)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return mqPop
}

func (m MqPop) GetCmd() cli.Command {
	return m.Command
}

func (m *MqPop) Action(queueName string, settings *common.Settings) error {
	if queueName == "" {
		return errors.New(`pop requires a queue name`)
	}

	if m.outputfile != "" {
		f, err := os.Create(m.outputfile)
		if err != nil {
			return err
		}

		m.file = f
	}

	q := mq.ConfigNew(queueName, &settings.Worker)

	messages, err := q.PopN(m.number)
	if err != nil {
		return err
	}

	// If anything here fails, we still want to print out what was deleted before exiting
	if m.file != nil {
		b, err := json.Marshal(messages)
		if err != nil {
			common.PrintMessages(messages)

			return err
		}

		_, err = m.file.Write(b)
		if err != nil {
			common.PrintMessages(messages)

			return err
		}
	}

	if common.IsPipedOut() {
		common.PrintMessages(messages)
	} else {
		plural := ""
		if m.number > 1 {
			plural = "s"
		}

		fmt.Println(common.LINES, "Message", plural, " successfully popped off ", q.Name)
		fmt.Println()
		fmt.Println("-------- ID ------ | Body")
		common.PrintMessages(messages)
	}

	return nil
}
