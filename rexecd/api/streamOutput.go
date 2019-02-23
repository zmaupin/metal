package api

import (
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/metal-go/metal/config"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type streamOutputClient struct {
	consumer sarama.Consumer
	receive  chan []byte
	ws       *websocket.Conn
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
		receive:  make(chan []byte),
		ws:       ws,
	}, nil
}

func (s *streamOutputClient) readPump()  {}
func (s *streamOutputClient) writePump() {}

func serveStreamOutputClient(w http.ResponseWriter, r *http.Request) {
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
	go client.readPump()
	go client.writePump()
}
