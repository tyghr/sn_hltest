package rabbitmq

import (
	"context"
	"fmt"
)

func (q *Queue) ReadFeedRebuild(ctx context.Context) (<-chan string, error) {
	ch, err := q.conn.Channel()
	if err != nil {
		q.logger.Errorw("connection to rabbitmq failed", "error", err)
		return nil, err
	}

	//err := c.Chan.Qos(
	//	1,     // prefetch count
	//	0,     // prefetch size
	//	false, // global
	//)

	_, err = ch.QueueDeclare(
		queueRebuild, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if q.logger.FailOnError(err, fmt.Sprintf("failed to declare queue: %s", queueRebuild)) {
		return nil, err
	}

	deliveries, err := ch.Consume(
		queueRebuild, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if q.logger.FailOnError(err, fmt.Sprintf("failed to register a consumer on: %s", queueRebuild)) {
		return nil, err
	}

	subs := make(chan string)
	go func() {
		for d := range deliveries {
			sub := string(d.Body)

			q.logger.Debugw("ReadFeedRebuild message consumed",
				"queueName", queueRebuild,
				"user", sub,
			)

			subs <- sub
		}
		close(subs)
	}()

	return subs, nil
}
