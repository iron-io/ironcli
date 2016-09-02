package mq

import (
	"fmt"

	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/iron_go3/mq"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type MqList struct {
	page    string
	perPage int
	filter  string

	cli.Command
}

func NewMqList(settings *config.Settings) *MqList {
	mqList := &MqList{}
	mqList.Command = cli.Command{
		Name:      "list",
		Usage:     "list of queues",
		ArgsUsage: "[args]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "page",
				Usage:       "",
				Destination: &mqList.page,
			},
			cli.IntFlag{
				Name:        "perPage",
				Usage:       "",
				Destination: &mqList.perPage,
			},
			cli.StringFlag{
				Name:        "filter",
				Usage:       "",
				Destination: &mqList.filter,
			},
		},
		Action: func(c *cli.Context) error {
			queues, err := mq.FilterPage(mqList.filter, mqList.page, mqList.perPage)
			if err != nil {
				return err
			}

			if common.IsPipedOut() {
				for _, q := range queues {
					fmt.Println(q.Name)
				}
			} else {
				fmt.Println(common.LINES, "Listing queues")
				for _, q := range queues {
					fmt.Println(common.BLANKS, "*", q.Name)
				}

				if tag, err := common.GetHudTag(*settings); err == nil {
					fmt.Printf("%s Go to hud-e.iron.io/mq/%s/projects/%s/queues for more info",
						common.BLANKS,
						tag,
						settings.ProjectId)
				}
				fmt.Println()
			}

			return nil
		},
	}

	return mqList
}

func (r MqList) GetCmd() cli.Command {
	return r.Command
}
