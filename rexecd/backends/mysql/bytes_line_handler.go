package mysql

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io/ioutil"

	"github.com/Shopify/sarama"
)

// BytesLineHandlerType type
type BytesLineHandlerType string

const (
	// MySQLStdout type
	MySQLStdout BytesLineHandlerType = "stdout"
	// MySQLStderr type
	MySQLStderr BytesLineHandlerType = "stderr"
)

// BytesLineHandler handles
type BytesLineHandler struct {
	command            *Command
	kafkaClusterAdmin  sarama.ClusterAdmin
	kafkaAsyncProducer sarama.AsyncProducer
	handlerType        BytesLineHandlerType
	lineNo             uint64
	topic              string
	init               bool
}

// NewBytesLineHandler returns a new BytesLineHandler
func NewBytesLineHandler(command *Command, handlerType BytesLineHandlerType, kafkaClusterAdmin sarama.ClusterAdmin, kafkaAsyncProducer sarama.AsyncProducer) *BytesLineHandler {
	return &BytesLineHandler{
		command:            command,
		handlerType:        handlerType,
		kafkaAsyncProducer: kafkaAsyncProducer,
		kafkaClusterAdmin:  kafkaClusterAdmin,
	}
}

// Handle satisfies rexecd.BytesLineHandler
func (b *BytesLineHandler) Handle(ctx context.Context, data []byte) error {
	if !b.init {
		if err := b.initFunc(); err != nil {
			return err
		}
	}
	b.lineNo++

	key, err := b.uint64ToKey(b.lineNo)
	if err != nil {
		return err
	}

	b.kafkaAsyncProducer.Input() <- &sarama.ProducerMessage{
		Key:   sarama.ByteEncoder(key),
		Topic: b.topic,
		Value: sarama.ByteEncoder(data),
	}
	return nil
}

// Finish wraps up the handling of bytes
func (b *BytesLineHandler) Finish(ctx context.Context) error {
	key, err := b.uint64ToKey(0)
	if err != nil {
		return err
	}
	b.kafkaAsyncProducer.Input() <- &sarama.ProducerMessage{
		Key:   sarama.ByteEncoder(key),
		Topic: b.topic,
		Value: sarama.ByteEncoder([]byte("EOF")),
	}
	return nil
}

func (b *BytesLineHandler) uint64ToKey(i uint64) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := binary.Write(buf, binary.LittleEndian, b.lineNo); err != nil {
		return nil, err
	}
	key, err := ioutil.ReadAll(buf)
	return key, err
}

func (b *BytesLineHandler) initFunc() error {
	topic := fmt.Sprintf("%d-%v", b.command.ID, b.handlerType)
	b.topic = topic
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1, // TODO: Make this num of brokers
	}
	if err := b.kafkaClusterAdmin.CreateTopic(topic, topicDetail, false); err != nil {
		return err
	}

	b.init = true
	return nil
}
