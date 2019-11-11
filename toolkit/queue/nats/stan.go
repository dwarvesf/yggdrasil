package nats

import (
	"time"

	"github.com/nats-io/stan.go"

	"github.com/dwarvesf/yggdrasil/toolkit/queue"
)

// Queue implementation
type nats struct {
	clusterID string
	clientID  string
	natsURL   string
}

// New return a new nats information
func New(clusterID, clientID, natsURL string) queue.Queue {
	return &nats{
		clusterID: clusterID,
		clientID:  clientID,
		natsURL:   natsURL,
	}
}

// Message writer implementation
type natsWriter struct {
	Writer stan.Conn
}

// NewWriter return a new writer for NATS-streaming
func (n *nats) NewWriter() (queue.Writer, error) {
	sc, err := stan.Connect(
		n.clusterID,
		n.clientID+"-writer",
		stan.NatsURL(n.natsURL),
	)
	if err != nil {
		return nil, err
	}
	return &natsWriter{
		Writer: sc,
	}, nil
}

func (w *natsWriter) Write(channel string, msg []byte) error {
	return w.Writer.Publish(channel, msg)
}

func (w *natsWriter) Close() error {
	return w.Writer.Close()
}

// Message reader implementation
type natsReader struct {
	channel string
	Reader  stan.Conn
}

// NewReader return a new reader for NATS-streaming
func (n *nats) NewReader(channel string) (queue.Reader, error) {
	sc, err := stan.Connect(
		n.clusterID,
		n.clientID+"-reader",
		stan.NatsURL(n.natsURL),
	)
	if err != nil {
		return nil, err
	}
	return &natsReader{
		channel: channel,
		Reader:  sc,
	}, nil
}

func (r *natsReader) Read(handleMsg func(m *stan.Msg)) error {
	_, err := r.Reader.QueueSubscribe(
		r.channel,                    // channel name
		r.channel,                    // queue name
		handleMsg,                    // business logic of handle msg
		stan.DurableName(r.channel),  // durable name
		stan.DeliverAllAvailable(),   // deliver all message available
		stan.SetManualAckMode(),      // manually send ack
		stan.AckWait(20*time.Second), // server ack wait time
		stan.MaxInflight(25),         // maximum msg sending in flight
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *natsReader) Close() error {
	return r.Reader.Close()
}
