package counters

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tyghr/social_network/internal/model"
)

func (q *Queue) ReadCountersBuckets(ctx context.Context) (<-chan model.CounterCmd, error) {
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

	queueName := queueCounters

	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if q.logger.FailOnError(err, fmt.Sprintf("failed to declare queue: %s", queueName)) {
		return nil, err
	}

	deliveries, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if q.logger.FailOnError(err, fmt.Sprintf("failed to register a consumer on: %s", queueName)) {
		return nil, err
	}

	cmds := make(chan model.CounterCmd)
	go func() {
		for d := range deliveries {
			var pb model.CounterCmd
			err := json.Unmarshal(d.Body, &pb)
			if err != nil {
				q.logger.Errorw("ConsumePosts: unmarshaling", "error", err)
				break
			}

			q.logger.Debugw("ReadCountersBuckets message consumed",
				"queueName", queueName,
				"command", pb.Command,
				"cut_subs", pb.Subscribers)

			cmds <- pb
		}
		close(cmds)
	}()

	return cmds, nil
}
