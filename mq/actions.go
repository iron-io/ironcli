package mq

import (
	"fmt"

	"github.com/iron-io/ironcli/Godeps/_workspace/src/github.com/iron-io/iron_go3/mq"
	"github.com/iron-io/ironcli/common"
)

func enqueue(g *common.GlobalFlags) {
	qName := g.StringOrFail("queue")
	body := g.StringOrFail("body")
	delay := int64(g.IntOrFail("delay", -1))
	q := mq.ConfigNew(qName, common.NewIronConfig(g))
	msg := mq.Message{Body: body, Delay: delay}
	id, err := q.PushMessage(msg)
	common.FailErr(err)
	fmt.Println("success. message id = ", id)
}

func delete(g *common.GlobalFlags) {
	qName := g.StringOrFail("queue")
	id := g.StringOrFail("id")
	resID := g.StringOrFail("reservation_id")
	q := mq.ConfigNew(qName, common.NewIronConfig(g))
	err := q.DeleteMessage(id, resID)
	common.FailErr(err)
	fmt.Println("success")
}
