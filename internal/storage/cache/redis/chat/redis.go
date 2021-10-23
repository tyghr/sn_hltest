package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/tyghr/social_network/internal/storage"
)

var (
	txMaxRetries = 1000
)

// <-chan *redis.Message, []redis.XMessage, func() error
func (chat *Chat) GetMsgs(ctx context.Context, chatID string) (<-chan storage.ChannelMessage, []storage.ChatMessage, func() error) {
	chat.rc.Ping(ctx)

	// subscribe (online events)
	pubsub := chat.rc.Subscribe(ctx,
		getChannelName(chatID),
	)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}

	// get user list
	userList, err := chat.rc.Get(ctx,
		fmt.Sprintf("userlist:%s", chatID),
	).Result()
	if err == redis.Nil {
		log.Printf("userlist for chat %s does not exist", chatID)
	} else if err != nil {
		panic(err)
	}

	// check each users counters
	streams := []string{}
	for _, u := range strings.Split(userList, ":") {
		if u == "" {
			continue
		}
		cnt, err := chat.rc.Get(ctx,
			fmt.Sprintf("msgcount:%s", u),
		).Int()
		if err == redis.Nil {
			log.Printf("msgcount for user %s does not exist", u)
		} else if err != nil {
			panic(err)
		}

		for i := 0; i <= cnt; i = i + 1000 {
			streams = append(streams, getStreamName(chatID, u, i))
		}
	}
	oldMsgs := []storage.ChatMessage{}
	// range all matched users chat streams for old messages
	for _, s := range streams {
		log.Println("reading stream:", s)

		res, err := chat.rc.XRead(ctx, &redis.XReadArgs{
			Streams: []string{s, "0"},
			Block:   -1,
		}).Result()
		if err != nil && err != redis.Nil {
			panic(err)
		}

		for _, stream := range res {
			for _, sm := range stream.Messages {
				oldMsgs = append(oldMsgs, &ChatMessage{sm})
			}
		}
	}
	sort.Sort(storage.ByMsgID(oldMsgs))

	retCh := make(chan storage.ChannelMessage)

	go func() {
		for chMsg := range pubsub.Channel() {
			retCh <- &ChannelMessage{
				cm: chMsg,
			}
		}
		close(retCh)
	}()

	// When pubsub is closed channel is closed too.
	return retCh, oldMsgs, pubsub.Close
}

func (chat *Chat) updateUserList(ctx context.Context, chatID, userName string) error {
	// transaction func
	txf := func(tx *redis.Tx) error {
		// check/add user to chat user list
		userList, err := tx.Get(ctx,
			fmt.Sprintf("userlist:%s", chatID),
		).Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if !strings.Contains(userList+":", ":"+userName+":") {
			// Operation is committed only if the watched keys remain unchanged.
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				return pipe.Append(ctx,
					fmt.Sprintf("userlist:%s", chatID),
					":"+userName,
				).Err()
			})
			return err
		}
		return nil
	}

	for i := 0; i < txMaxRetries; i++ {
		err := chat.rc.Watch(
			ctx,
			txf,
			fmt.Sprintf("userlist:%s", chatID),
		)
		if err == nil {
			return nil
		}
		if err == redis.TxFailedErr {
			continue
		}
		return err
	}

	return errors.New("tx max retries reached")
}

func (chat *Chat) incrementMsgCount(ctx context.Context, userName string) (int, error) {
	var msgCount int

	// transaction func
	txf := func(tx *redis.Tx) error {
		// get user msg_count
		cnt, err := tx.Get(ctx,
			fmt.Sprintf("msgcount:%s", userName),
		).Int()
		if err != nil && err != redis.Nil {
			return fmt.Errorf("get user msg_count %w", err)
		}
		cnt++
		msgCount = cnt

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			// increment user msg_count
			err := pipe.Set(ctx,
				fmt.Sprintf("msgcount:%s", userName),
				cnt,
				0,
			).Err()
			if err != nil {
				return fmt.Errorf("increment user msg_count %w", err)
			}
			return nil
		})
		return err
	}

	for i := 0; i < txMaxRetries; i++ {
		err := chat.rc.Watch(
			ctx,
			txf,
			fmt.Sprintf("msgcount:%s", userName),
		)
		if err == redis.TxFailedErr {
			continue
		} else if err != nil {
			return 0, fmt.Errorf("Watch: %w", err)
		}
		return msgCount, nil
	}
	return 0, errors.New("tx max retries reached")
}

func (chat *Chat) WriteMessage(ctx context.Context, chatID, userName, message string) error {
	jMsg, _ := json.Marshal(struct {
		User string `json:"user"`
		Text string `json:"text"`
	}{
		User: userName,
		Text: message,
	})

	err := chat.updateUserList(ctx, chatID, userName)
	if err != nil {
		return fmt.Errorf("updateUserList: %w", err)
	}

	cnt, err := chat.incrementMsgCount(ctx, userName)
	if err != nil {
		return fmt.Errorf("incrementMsgCount: %w", err)
	}

	// add to stream
	err = chat.rc.XAdd(ctx, &redis.XAddArgs{
		Stream: getStreamName(chatID, userName, cnt),
		ID:     "*",
		Values: []string{
			"user", userName,
			"text", message,
		},
	}).Err()
	if err != nil {
		return fmt.Errorf("add to stream %w", err)
	}

	// publish
	err = chat.rc.Publish(
		context.TODO(),
		getChannelName(chatID),
		jMsg,
	).Err()
	if err != nil {
		return fmt.Errorf("publish %w", err)
	}

	return nil
}
