// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rocketmq

import (
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/douyu/jupiter/pkg/client/rocketmq"
	"github.com/douyu/jupiter/pkg/util/xgo"
)

// MessageQ ...
type MessageQ struct {
	*rocketmq.ProducerConfig

	msgChan chan *primitive.Message
}

// New ...
func New() *MessageQ {
	mq := &MessageQ{
		// Producer: rocketmq.StdNewProducer("message_q"),
		msgChan: make(chan *primitive.Message, 1000),
	}

	return mq
}

// Start ...
func (mq *MessageQ) Start() error {
	xgo.Go(mq.batchSend)
	return nil
}

// Close ...
func (mq *MessageQ) Close() error {
	// mq.Producer.Close()
	return nil
}

// Push ...
func (mq *MessageQ) Push(messages ...Message) {
	for _, message := range messages {
		msg := &primitive.Message{
			Topic:         message.Topic(),
			Body:          nil,
			Flag:          0,
			TransactionId: "",
			Batch:         false,
			Queue:         nil,
		}

		mq.msgChan <- msg

		fmt.Printf("message = %+v\n", message)
		fmt.Printf("msg = %+v\n", msg)
	}
}

// batchSend ...
func (mq *MessageQ) batchSend() {
	for {
		msgs := batchMessages(mq.msgChan, 16)
		if len(msgs) <= 0 {
			continue
		}
		// if err := mq.SendOneWay(context.Background(), msgs...); err != nil {
		// 	log.Errord("send message", log.Any("msgs", msgs))
		// }
		// log.Infod("send message", log.Any("msgs", msgs))
	}
}

// batchMessages ...
func batchMessages(msgChan <-chan *primitive.Message, size int) (msgs []*primitive.Message) {
	msgs = make([]*primitive.Message, size)
	var ticker = time.NewTimer(time.Second)
	for {
		select {
		case msg := <-msgChan:
			msgs = append(msgs, msg)
			if len(msgs) >= size {
				return msgs
			}
		case <-ticker.C:
			return msgs
		}
	}
}
