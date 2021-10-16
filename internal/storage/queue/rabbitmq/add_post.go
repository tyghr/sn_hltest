package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/streadway/amqp"
	"github.com/tyghr/social_network/internal/model"
)

func (q *Queue) AddPostBuckets(ctx context.Context, post model.Post, subs []string) error {
	ch, err := q.conn.Channel()
	if err != nil {
		q.logger.Errorw("connection to rabbitmq failed", "error", err)
		return nil
	}

	if _, err := ch.QueueDeclare(
		queueAppend, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	); err != nil {
		return err
	}

	q.logger.Debugw("AddPostBuckets invoked",
		"user", post.UserName,
		"post", post.Header,
		"subs", subs)

	for _, cutSubs := range cutStringSlice(subs) {
		jsonMsg, err := json.Marshal(model.PostBacket{
			Post:        post,
			Subscribers: cutSubs,
		})
		if err != nil {
			q.logger.Errorw("queue: error while marshaling")
			return err
		}

		if err := ch.Publish(
			"",          // exchange
			queueAppend, // routing key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent, // save messages to disk
				ContentType:  "application/json",
				Body:         jsonMsg,
			},
		); err != nil {
			return err
		}
		q.logger.Debugw("AddPostBuckets message published",
			"queueName", queueAppend,
			"user", post.UserName,
			"post", post.Header,
			"cut_subs", cutSubs)
	}

	return nil
}
