package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/metal-go/metal/config"
	"github.com/metal-go/metal/util/queue"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type streamOutputClient struct {
	consumer sarama.Consumer
	done     chan struct{}
	receive  chan []byte
	send     *queue.Bytes
	ws       *websocket.Conn
}

type streamOutputRequest struct {
	topic string
}

func newStreamoutputClient(ws *websocket.Conn) (*streamOutputClient, error) {
	kafkaConfig := sarama.NewConfig()

	kafkaVersion, err := sarama.ParseKafkaVersion(config.RexecdGlobal.KafkaVersion)
	if err != nil {
		return nil, err
	}
	kafkaConfig.Version = kafkaVersion

	consumer, err := sarama.NewConsumer(config.RexecdGlobal.KafkaAddress, kafkaConfig)
	if err != nil {
		return nil, err
	}

	return &streamOutputClient{
		consumer: consumer,
		done:     make(chan struct{}),
		receive:  make(chan []byte),
		send:     queue.NewBytes(),
		ws:       ws,
	}, nil
}

// readPump gets the message request from the websocket client and sends it on
// the receive channel. It only accepts one message of type BinaryMessage or
// a CloseMessage
func (s *streamOutputClient) readPump(ctx context.Context) {
	for {
		select {
		case <-s.done:
			return
		case <-ctx.Done():
			s.done <- struct{}{}
			return
		default:
			messageType, b, err := s.ws.ReadMessage()
			if err != nil {
				log.Error(err)
				s.done <- struct{}{}
				return
			}
			switch messageType {
			case websocket.TextMessage, websocket.PingMessage, websocket.PongMessage:
				log.Errorf("unsupported message type %v", messageType)
				s.done <- struct{}{}
				return
			case websocket.CloseMessage:
				s.done <- struct{}{}
				return
			case websocket.BinaryMessage:
				s.receive <- b
				close(s.receive)
				return
			}
		}
	}
}

// writePump consumes all messages on the requested Kafka topic and sends
// them to the client. If there is an error at any point, we log it and return.
// It creates a go routine that consumes all messages and places each []byte
// on the queue.Bytes. This allows this func to consume one message per
// iteration
func (s *streamOutputClient) writePump(ctx context.Context) {
	var request = new(streamOutputRequest)
	var partitionConsumer sarama.PartitionConsumer
	defer func() {
		if partitionConsumer != nil {
			partitionConsumer.Close()
		}
	}()

	for {
		select {
		case msg := <-s.receive:
			err := json.Unmarshal(msg, request)
			if err != nil {
				log.Error(err)
				s.done <- struct{}{}
				return
			}
			partitionConsumer, err = s.consumer.ConsumePartition(request.topic, 0, 0)
			if err != nil {
				log.Error(err)
				s.done <- struct{}{}
				return
			}
			go s.consumePartition(ctx, partitionConsumer)
		case <-s.done:
			return
		case <-ctx.Done():
			s.done <- struct{}{}
			return
		default:
			if request.topic == "" {
				continue
			}
			b := s.send.Dequeue()
			if b == nil {
				s.done <- struct{}{}
				return
			}

			if err := s.ws.WriteMessage(websocket.BinaryMessage, b); err != nil {
				if err != nil {
					log.Error(err)
					s.done <- struct{}{}
					return
				}
			}
			if string(b) == "EOF" {
				s.done <- struct{}{}
				return
			}
		}
	}
}

// consumeParition will consume all messages from the given
// sarama.PartitionConsumer.
func (s *streamOutputClient) consumePartition(ctx context.Context, p sarama.PartitionConsumer) {
	for {
		select {
		case msg := <-p.Messages():
			s.send.Enqueue(msg.Value)
		case <-ctx.Done():
			return
		case <-s.done:
			return
		}
	}
}

func serveStreamOutputClient(timeout time.Duration, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(500)
		log.Error(err)
		return
	}
	client, err := newStreamoutputClient(ws)
	if err != nil {
		w.WriteHeader(500)
		log.Error(err)
		return
	}
	defer client.ws.Close()
	defer close(client.done)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	go client.readPump(ctx)
	go client.writePump(ctx)

	select {
	case <-client.done:
	case <-ctx.Done():
	}
}
