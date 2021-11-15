package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/streadway/amqp"
	"github.com/tyghr/social_network/internal/model"
)

func (q *Queue) incCounters(ctx context.Context, command string, subs []string) error {
	ch, err := q.conn.Channel()
	if err != nil {
		q.logger.Errorw("connection to rabbitmq failed", "error", err)
		return nil
	}

	queueName := queueCounters

	if _, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	); err != nil {
		return err
	}

	q.logger.Debugw("incCounters invoked",
		"command", command,
		"subs", subs)

	for _, cutSubs := range cutStringSlice(subs) {
		jsonMsg, err := json.Marshal(model.CounterCmd{
			Command:     command,
			Subscribers: cutSubs,
		})
		if err != nil {
			q.logger.Errorw("queue: error while marshaling")
			return err
		}

		if err := ch.Publish(
			"",        // exchange
			queueName, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent, // save messages to disk
				ContentType:  "application/json",
				Body:         jsonMsg,
			},
		); err != nil {
			return err
		}
		q.logger.Debugw("incCounters message published",
			"command", command,
			"queueName", queueName,
			"cut_subs", cutSubs,
		)
	}

	return nil
}

func (q *Queue) IncTotalCounters(ctx context.Context, subs []string) error {
	return q.incCounters(ctx, model.CounterCmdTotalInc, subs)
}

func (q *Queue) IncCursorCounters(ctx context.Context, subs []string) error {
	return q.incCounters(ctx, model.CounterCmdCursorInc, subs)
}

func (q *Queue) UpdateCursorCounter(ctx context.Context, user string) error {
	ch, err := q.conn.Channel()
	if err != nil {
		q.logger.Errorw("connection to rabbitmq failed", "error", err)
		return nil
	}

	queueName := queueCounters

	if _, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	); err != nil {
		return err
	}

	q.logger.Debugw("UpdateCursorCounter invoked",
		"user", user)

	jsonMsg, err := json.Marshal(model.CounterCmd{
		Command:     model.CounterCmdCursorUpdate,
		Subscribers: []string{user},
	})
	if err != nil {
		q.logger.Errorw("queue: error while marshaling")
		return err
	}

	if err := ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // save messages to disk
			ContentType:  "application/json",
			Body:         jsonMsg,
		},
	); err != nil {
		return err
	}
	q.logger.Debugw("UpdateCursorCounter message published",
		"user", user,
	)

	return nil
}
