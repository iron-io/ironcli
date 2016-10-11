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

type MqDelete struct {
	filequeue_name string
	queue_name     string
	ids            []string

	cli.Command
}

func NewMqDelete(settings *common.Settings) *MqDelete {
	mqDelete := &MqDelete{}

	mqDelete.Command = cli.Command{
		Name:      "delete",
		Usage:     "delete messages by ID of queue",
		ArgsUsage: "[QUEUE_NAME] [MSG_ID, ...]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "filequeue, i",
				Usage:       "",
				Destination: &mqDelete.filequeue_name,
			},
		},
		Before: func(c *cli.Context) error {
			if err := common.SetSettings(settings); err != nil {
				return err
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			err := mqDelete.Action(c.Args().First(), c.Args().Tail(), settings)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return mqDelete
}

func (r MqDelete) GetCmd() cli.Command {
	return r.Command
}

func (r *MqDelete) Execute(queueName string, ids []string) error {
	if queueName == "" {
		return errors.New(`delete requires a queue name`)
	}

	r.queue_name = queueName

	// Read and parse piped info
	if common.IsPipedIn() {
		ids, err := common.ReadIds()
		if err != nil {
			return err
		}

		r.ids = append(r.ids, ids...)
	}

	if r.filequeue_name != "" {
		b, err := ioutil.ReadFile(r.filequeue_name)
		if err != nil {
			return err
		}

		// Use the message struct so its compatible with output files from reserve
		var msgs []mq.Message
		err = json.Unmarshal(b, &msgs)
		if err != nil {
			return err
		}
		for _, msg := range msgs {
			r.ids = append(r.ids, msg.Id)
		}
	}

	if len(ids) >= 1 {
		r.ids = append(r.ids, ids...)
	}

	if len(r.ids) < 1 {
		return errors.New("delete requires at least one message id")
	}

	return nil
}

func (m *MqDelete) Action(queueName string, ids []string, settings *common.Settings) error {
	err := m.Execute(queueName, ids)
	if err != nil {
		return err
	}

	q := mq.ConfigNew(m.queue_name, &settings.Worker)

	err = q.DeleteMessages(m.ids)
	if err != nil {
		return err
	}

	plural := ""
	if len(m.ids) > 1 {
		plural = "s"
	}

	fmt.Println(common.BLANKS, "Done deleting message", plural)

	return nil
}
