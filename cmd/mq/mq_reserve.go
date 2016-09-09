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

type MqReserve struct {
	queue_name string
	number     int
	timeout    int
	outputfile string
	file       *os.File

	cli.Command
}

func NewMqReserve(settings *common.Settings) *MqReserve {
	mqReserve := &MqReserve{}

	mqReserve.Command = cli.Command{
		Name:      "reserve",
		Usage:     "reserve meesages by amount of queue",
		ArgsUsage: "[QUEUE_NAME]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:        "number, n",
				Usage:       "",
				Destination: &mqReserve.number,
			},
			cli.IntFlag{
				Name:        "timeout, t",
				Usage:       "0(default) - 3600 max runtime for task in seconds",
				Destination: &mqReserve.timeout,
			},
			cli.StringFlag{
				Name:        "output, o",
				Usage:       "",
				Destination: &mqReserve.outputfile,
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().First() == "" {
				return errors.New(`reserve requires a queue name`)
			}

			if mqReserve.outputfile != "" {
				f, err := os.Create(mqReserve.outputfile)
				if err != nil {
					return err
				}

				mqReserve.file = f
			}

			q := mq.ConfigNew(mqReserve.queue_name, &settings.Worker)
			messages, err := q.GetNWithTimeout(mqReserve.number, mqReserve.timeout)
			if err != nil {
				return err
			}

			// If anything here fails, we still want to print out what was reserved before exiting
			if mqReserve.file != nil {
				b, err := json.Marshal(messages)
				if err != nil {
					common.PrintReservedMessages(messages)

					return err
				}
				_, err = mqReserve.file.Write(b)
				if err != nil {
					common.PrintReservedMessages(messages)

					return err
				}
			}

			if len(messages) < 1 {
				return errors.New("Queue is empty")
			}

			if common.IsPipedOut() {
				common.PrintReservedMessages(messages)
			} else {
				fmt.Println(common.LINES, "Messages successfully reserved")
				fmt.Println("--------- ID ------|------- Reservation ID -------- | Body")
				common.PrintReservedMessages(messages)
			}

			return nil
		},
	}

	return mqReserve
}

func (r MqReserve) GetCmd() cli.Command {
	return r.Command
}
