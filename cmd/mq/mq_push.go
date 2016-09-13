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
		Action: func(c *cli.Context) error {
			err := mqPush.Execute(c.Args().First(), c.Args().Tail())
			if err != nil {
				return err
			}

			q := mq.ConfigNew(mqPush.queue_name, &settings.Worker)

			ids, err := q.PushStrings(mqPush.messages...)
			if err != nil {
				return err
			}

			if common.IsPipedOut() {
				for _, id := range ids {
					fmt.Println(id)
				}
			} else {
				fmt.Println(common.Green(common.LINES, "Message succesfully pushed!"))
				fmt.Printf("%sMessage IDs:\n", common.BLANKS)
				fmt.Printf("%s", common.BLANKS)

				for _, id := range ids {
					fmt.Printf("%s ", id)
				}

				fmt.Println()
			}

			return nil
		},
	}

	return mqPush
}

func (r MqPush) GetCmd() cli.Command {
	return r.Command
}

func (r *MqPush) Execute(queueName string, messages []string) error {
	if queueName == "" {
		return errors.New(`push requires a queue name`)
	}

	r.queue_name = queueName

	if r.filename != "" {
		b, err := ioutil.ReadFile(r.filename)
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

		r.messages = append(r.messages, messageStruct.Messages...)
	}

	if len(messages) < 1 {
		return errors.New(`push requires at least one message`)
	} else {
		r.messages = messages
	}

	return nil
}
