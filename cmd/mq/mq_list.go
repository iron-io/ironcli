package mq

import (
	"fmt"

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

func NewMqList(settings *common.Settings) *MqList {
	mqList := &MqList{}

	mqList.Command = cli.Command{
		Name:  "list",
		Usage: "list of queues",
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
		Before: func(c *cli.Context) error {
			if err := common.SetSettings(settings); err != nil {
				return err
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			err := mqList.Action(settings)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return mqList
}

func (m MqList) GetCmd() cli.Command {
	return m.Command
}

func (m *MqList) Action(settings *common.Settings) error {
	queues, err := mq.FilterPage(m.filter, m.page, m.perPage)
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

		if tag, err := common.GetHudTag(settings.Worker); err == nil {
			fmt.Printf("%s Go to hud-e.iron.io/mq/%s/projects/%s/queues for more info",
				common.BLANKS,
				tag,
				settings.Worker.ProjectId)
		}
		fmt.Println()
	}

	return nil
}
