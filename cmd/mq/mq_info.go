package mq

import (
	"errors"
	"fmt"

	"github.com/iron-io/iron_go3/mq"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type MqInfo struct {
	subscriberList bool

	cli.Command
}

func NewMqInfo(settings *common.Settings) *MqInfo {
	mqInfo := &MqInfo{}

	mqInfo.Command = cli.Command{
		Name:      "info",
		Usage:     "get info about queue",
		ArgsUsage: "[QUEUE_NAME]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "subscriber-list",
				Usage:       "",
				Destination: &mqInfo.subscriberList,
			},
		},
		Before: func(c *cli.Context) error {
			if err := common.SetSettings(settings); err != nil {
				return err
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			err := mqInfo.Action(c.Args().First(), settings)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return mqInfo
}

func (m MqInfo) GetCmd() cli.Command {
	return m.Command
}

func (m *MqInfo) Action(queueName string, settings *common.Settings) error {
	if queueName == "" {
		return errors.New(`info requires a queue name`)
	}

	q := mq.ConfigNew(queueName, &settings.Worker)

	info, err := q.Info()
	if err != nil {
		return err
	}

	fmt.Printf("%sName: %s\n", common.BLANKS, info.Name)
	fmt.Printf("%sCurrent Size: %d\n", common.BLANKS, info.Size)
	fmt.Printf("%sTotal messages: %d\n", common.BLANKS, info.TotalMessages)
	fmt.Printf("%sMessage expiration: %d\n", common.BLANKS, info.MessageExpiration)
	fmt.Printf("%sMessage timeout: %d\n", common.BLANKS, info.MessageTimeout)

	if info.Push != nil {
		fmt.Printf("%sType: %s\n", common.BLANKS, info.Type)
		fmt.Printf("%sSubscribers: %d\n", common.BLANKS, len(info.Push.Subscribers))
		fmt.Printf("%sRetries: %d\n", common.BLANKS, info.Push.Retries)
		fmt.Printf("%sRetries delay: %d\n", common.BLANKS, info.Push.RetriesDelay)

		if m.subscriberList {
			fmt.Printf("%sSubscriber list\n", common.LINES)
			common.PrintSubscribers(info)
			fmt.Println()
		}
	}
	common.PrintQueueHudURL(common.BLANKS, q)

	return nil
}
