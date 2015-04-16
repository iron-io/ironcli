package common

import (
	"strconv"

	"github.com/iron-io/ironcli/Godeps/_workspace/src/github.com/iron-io/iron_go3/config"
)

func NewIronConfig(g *GlobalFlags) *config.Settings {
	return &config.Settings{
		Token:      g.Token,
		ProjectId:  g.ProjID,
		Host:       g.Host,
		Scheme:     "https",
		Port:       443,
		ApiVersion: strconv.Itoa(g.Version),
		UserAgent:  "ironcli",
	}
}
