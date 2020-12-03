package hub

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	"sync"
	"time"
)

var RecheckExchange = "recheck_exchange_for_rebalancer"

type RecheckMessage struct {
	QueueName string
	Exchange  string
	Kind      string
	Host      string
}

type Rebalancer struct {
	hub         *Hub
	syncTimeMap sync.Map
}

func NewRebalancer(hub *Hub) *Rebalancer {
	rb := &Rebalancer{hub: hub}

	rb.ListenQueue()

	return rb
}

func (r *Rebalancer) ReBalance(queueName, kind, exchange string) {
	var (
		producer ProducerInterface
		err      error
		hostname string
	)

	// 无需持久化
	if producer, err = r.hub.pm.GetProducer("", amqp.ExchangeFanout, RecheckExchange); err != nil {
		log.Error(err)
		return
	}

	if hostname, err = os.Hostname(); err != nil {
		log.Error(err)
		return
	}

	if marshal, err := json.Marshal(&RecheckMessage{QueueName: queueName, Host: hostname, Kind: kind, Exchange: exchange}); err != nil {
		log.Error(err)
	} else {
		if err = producer.Publish(Message{Data: string(marshal)}); err != nil {
			log.Error(err)
		}
	}
}

func (r *Rebalancer) ListenQueue() {
	log.Info("LISTEN Rebalance!")
	var (
		consumer ConsumerInterface
		err      error
		hostname string
	)

	if hostname, err = os.Hostname(); err != nil {
		log.Error(err)
		return
	}

	if consumer, err = r.hub.cm.GetConsumer(hostname, amqp.ExchangeFanout, RecheckExchange, WithQueueAutoDelete(true)); err != nil {
		log.Error(err)
		return
	}

	// todo 30秒内不再触发重平衡
	go func() {
		for {
			select {
			case <-r.hub.Done():
				log.Error("rb hub done。")
				return
			case delivery, ok := <-consumer.Delivery():
				var (
					msg        Message
					recheckMsg RecheckMessage
				)
				delivery.Ack(true)
				if !ok {
					log.Info("not okkk")
					break
				}
				if err = json.Unmarshal(delivery.Body, &msg); err != nil {
					log.Error(err)
					break
				}
				if err = json.Unmarshal([]byte(msg.GetData()), &recheckMsg); err != nil {
					log.Error(err)
					break
				}
				key := getKey(recheckMsg.QueueName, recheckMsg.Kind, recheckMsg.Exchange)

				log.Infof("收到重平衡消息 queueName: %s 队列长度：%d host: %s", recheckMsg.QueueName, len(r.hub.cm.Delivery(key)), hostname)
				if len(r.hub.cm.Delivery(key)) > 0 {
					if load, ok := r.syncTimeMap.Load(key); ok {
						if load.(time.Time).After(time.Now().Add(3 * time.Second)) {
							log.Warn("3 秒内已经触发过重平衡了")
							return
						}
					}
					log.Warnf("触发重平衡 %s host: %s", recheckMsg.QueueName, hostname)
					for {
						select {
						case d := <-r.hub.cm.Delivery(key):
							d.Nack(false, true)
						default:
							r.syncTimeMap.Store(key, time.Now())
							break
						}
					}
				}
			}
		}
	}()
}