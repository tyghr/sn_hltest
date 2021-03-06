package counters

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"github.com/tyghr/logger"
	config "github.com/tyghr/social_network/internal/config/counters"
	"github.com/tyghr/social_network/internal/storage"
)

var (
	queueCounters = "counters"
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

func New(conf *config.Config, lgr logger.Logger) storage.CountersQueue {
	q := &Queue{
		config: conf,
		logger: lgr,
	}

	var scheme string
	switch conf.QueueType {
	case config.MQRabbit:
		scheme = "amqp"
	case config.MQRabbitSecured:
		scheme = "amqps"
	default:
		q.logger.Fatalw("unknown queue type")
	}
	var host string
	if conf.QueuePort == 0 {
		host = conf.QueueHost
	} else {
		host = fmt.Sprintf("%s:%d", conf.QueueHost, conf.QueuePort)
	}

	url := fmt.Sprintf("%s://%s:%s@%s/%s", scheme, conf.QueueUser, conf.QueuePass, host, conf.QueueVHost)

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
