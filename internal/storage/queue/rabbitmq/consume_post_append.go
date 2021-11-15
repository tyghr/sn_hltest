package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tyghr/social_network/internal/model"
)

func (q *Queue) ReadPostAppendBuckets(ctx context.Context) (<-chan model.PostBucket, error) {
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
		queueAppend, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if q.logger.FailOnError(err, fmt.Sprintf("failed to declare queue: %s", queueAppend)) {
		return nil, err
	}

	deliveries, err := ch.Consume(
		queueAppend, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if q.logger.FailOnError(err, fmt.Sprintf("failed to register a consumer on: %s", queueAppend)) {
		return nil, err
	}

	posts := make(chan model.PostBucket)
	go func() {
		for d := range deliveries {
			var pb model.PostBucket
			err := json.Unmarshal(d.Body, &pb)
			if err != nil {
				q.logger.Errorw("ConsumePosts: unmarshaling", "error", err)
				break
			}

			q.logger.Debugw("ReadPostBuckets message consumed",
				"queueName", queueAppend,
				"user", pb.Post.UserName,
				"post", pb.Post.Header,
				"cut_subs", pb.Subscribers)

			posts <- pb
		}
		close(posts)
	}()

	return posts, nil
}
