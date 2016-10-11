package mq

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/iron-io/iron_go3/mq"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type MqPush struct {
	filename   string
	messages   []string
	queue_name string
	ResultIds  []string

	cli.Command
}

func NewMqPush(settings *common.Settings) *MqPush {
	mqPush := &MqPush{}
	mqPush.Command = cli.Command{
		Name:      "push",
		Usage:     "push messages to queue",
		ArgsUsage: "[QUEUE_NAME] [MESSAGE, ...]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "filename, f",
				Usage:       "",
				Destination: &mqPush.filename,
			},
		},
		Before: func(c *cli.Context) error {
			if err := common.SetSettings(settings); err != nil {
				return err
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			err := mqPush.Action(c.Args().First(), c.Args().Tail(), settings)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return mqPush
}

func (m MqPush) GetCmd() cli.Command {
	return m.Command
}

func (m *MqPush) Execute(queueName string, messages []string) error {
	if queueName == "" {
		return errors.New(`push requires a queue name`)
	}

	m.queue_name = queueName

	if m.filename != "" {
		b, err := ioutil.ReadFile(m.filename)
		if err != nil {
			return err
		}

		messageStruct := struct {
			Messages []string `json:"messages"`
		}{}
		err = json.Unmarshal(b, &messageStruct)
		if err != nil {
			return err
		}

		m.messages = append(m.messages, messageStruct.Messages...)
	}

	if len(messages) < 1 {
		return errors.New(`push requires at least one message`)
	} else {
		m.messages = messages
	}

	return nil
}

func (m *MqPush) Action(queueName string, messages []string, settings *common.Settings) error {
	err := m.Execute(queueName, messages)
	if err != nil {
		return err
	}

	q := mq.ConfigNew(m.queue_name, &settings.Worker)

	ids, err := q.PushStrings(m.messages...)
	if err != nil {
		return err
	}

	if common.IsPipedOut() {
		for _, id := range ids {
			fmt.Println(id)
		}
	} else {
		m.ResultIds = ids

		fmt.Println(common.Green(common.LINES, "Message succesfully pushed!"))
		fmt.Printf("%sMessage IDs:\n", common.BLANKS)
		fmt.Printf("%s", common.BLANKS)

		for _, id := range ids {
			fmt.Printf("%s ", id)
		}

		fmt.Println()
	}

	return nil
}
