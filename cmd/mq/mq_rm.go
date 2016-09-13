package mq

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/iron-io/iron_go3/mq"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type MqRm struct {
	cli.Command
}

func NewMqRm(settings *common.Settings) *MqRm {
	mqRm := &MqRm{
		Command: cli.Command{
			Name:      "rm",
			Usage:     "rm queues by name",
			ArgsUsage: "[QUEUE_NAME]",
			Action: func(c *cli.Context) error {
				if c.Args().First() == "" {
					return errors.New(`rm requires a queue name`)
				}

				var queues []mq.Queue

				if common.IsPipedIn() {
					scanner := bufio.NewScanner(os.Stdin)
					for scanner.Scan() {
						name := scanner.Text()
						queues = append(queues, mq.ConfigNew(name, &settings.Worker))
					}
					if err := scanner.Err(); err != nil {
						fmt.Fprintln(os.Stderr, err)
					}
				} else {
					queues = append(queues, mq.ConfigNew(c.Args().First(), &settings.Worker))
				}

				for _, q := range queues {
					err := q.Delete()
					if err != nil {
						fmt.Println(common.Red("Error deleting queue ", q.Name, ": ", err))
					} else {
						fmt.Println(common.Green(common.LINES, q.Name, " has been sucessfully deleted."))
					}
				}

				return nil
			},
		},
	}

	return mqRm
}

func (r MqRm) GetCmd() cli.Command {
	return r.Command
}
