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
	q := &Queue{
		config: conf,
		logger: lgr,
	}

	var url string
	switch conf.QueueType {
	case config.MQRabbit:
		url = fmt.Sprintf("amqp://%s:%s@%s:%d/%s", conf.QueueUser, conf.QueuePass, conf.QueueHost, conf.QueuePort, conf.QueueVHost)
	case config.MQRabbitSecured:
		url = fmt.Sprintf("amqps://%s:%s@%s:%d/%s", conf.QueueUser, conf.QueuePass, conf.QueueHost, conf.QueuePort, conf.QueueVHost)
	default:
		q.logger.Fatalw("unknown queue type")
	}

	q.logger.Debugw("connecting to rabbitmq...")
	conn, err := connect(url)
	if err != nil {
		q.logger.Fatalw("connection to rabbitmq failed", "error", err)
	}

	q.conn = conn

	q.logger.Debugw("connection to rabbitmq established!")

	return q
}

func (q *Queue) Close() {
	q.logger.Debugf("closing mq connections...")
	q.conn.Close()
}
