package rabbitmq

import (
	"context"

	"github.com/streadway/amqp"
)

func (q *Queue) PostRebuildSubsFeedRequest(ctx context.Context, subs []string) error {
	ch, err := q.conn.Channel()
	if err != nil {
		q.logger.Errorw("connection to rabbitmq failed", "error", err)
		return nil
	}

	if _, err := ch.QueueDeclare(
		queueRebuild, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	); err != nil {
		return err
	}

	q.logger.Debugw("PostRebuildSubsFeedRequest invoked",
		"subs", subs)

	for _, sub := range subs {
		if err := ch.Publish(
			"",           // exchange
			queueRebuild, // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent, // save messages to disk
				ContentType:  "application/json",
				Body:         []byte(sub),
			},
		); err != nil {
			return err
		}
		q.logger.Debugw("PostRebuildSubsFeedRequest message published",
			"queueName", queueRebuild,
			"sub", sub)
	}

	return nil
}
