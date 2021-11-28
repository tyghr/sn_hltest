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
	bucketLen     = 100
	queueAppend   = "appendpost"
	queueRebuild  = "rebuild_feed"
	queueCounters = "counters"
)

type Queue struct {
	conn   *amqp.Connection
	config *config.QueueConfig
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
	if len(s) <= bucketLen {
		return [][]string{s}
	}
	return append([][]string{s[:bucketLen]}, cutStringSlice(s[bucketLen:])...)
}

func New(conf *config.QueueConfig, lgr logger.Logger) storage.Queue {
	q := &Queue{
		config: conf,
		logger: lgr,
	}

	var scheme string
	switch conf.Type {
	case config.MQRabbit:
		scheme = "amqp"
	case config.MQRabbitSecured:
		scheme = "amqps"
	default:
		q.logger.Fatalw("unknown queue type")
	}
	var host string
	if conf.Port == 0 {
		host = conf.Host
	} else {
		host = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	}

	url := fmt.Sprintf("%s://%s:%s@%s/%s", scheme, conf.User, conf.Pass, host, conf.VHost)

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
