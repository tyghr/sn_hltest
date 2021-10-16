package rabbitmq

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"github.com/tyghr/logger"
	"github.com/tyghr/social_network/internal/config"
	"github.com/tyghr/social_network/internal/storage"
)

var (
	backetLen    = 100
	queueAppend  = "appendpost"
	queueRebuild = "rebuild_feed"
)

type Queue struct {
	conn   *amqp.Connection
	config *config.Config
	logger logger.Logger
}

func connect(url string) (conn *amqp.Connection, err error) {
	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			return
		}
		time.Sleep(2 * time.Second)
	}
	return
}

func cutStringSlice(s []string) [][]string {
	if len(s) <= backetLen {
		return [][]string{s}
	}
	return append([][]string{s[:backetLen]}, cutStringSlice(s[backetLen:])...)
}

func New(conf *config.Config, lgr logger.Logger) storage.Queue {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", conf.QueueUser, conf.QueuePass, conf.QueueHost, conf.QueuePort)
	q := &Queue{
		config: conf,
		logger: lgr,
	}

	q.logger.Debugw("connecting to rabbitmq...")
	conn, err := connect(url)
	if err != nil {
		q.logger.Errorw("connection to rabbitmq failed", "error", err)
		return nil
	}

	q.conn = conn

	q.logger.Debugw("connection to rabbitmq established!")

	return q
}

func (q *Queue) Close() {
	q.logger.Debugf("closing mq connections...")
	q.conn.Close()
}
