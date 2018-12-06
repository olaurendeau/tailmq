package consumer

import (
  "log"
  "github.com/satori/go.uuid"
  "github.com/streadway/amqp"
)

type Consumer struct {
  uri string
  exchangeName string
  routingKey string
  conn *amqp.Connection
  ch *amqp.Channel
  Deliveries <-chan amqp.Delivery
  Err chan error
}

func New(uri string, exchangeName string, routingKey string) *Consumer {
  c := &Consumer{}
  c.uri = uri
  c.exchangeName = exchangeName
  c.routingKey = routingKey
  c.Err = make(chan error)

  return c
}

func (c *Consumer) Start() () {
  var err error

  c.conn, c.ch = c.openChannel()

  q := c.createExpirableQueue(c.ch)

  c.Deliveries, err = c.ch.Consume(q.Name, "", true, false, false, false, nil)
  c.Err <- err
}

func (c *Consumer) Stop() {
  c.conn.Close()
  c.ch.Close()
}

func (c *Consumer) openChannel() (*amqp.Connection, *amqp.Channel) {
  log.Printf("Establishing connection... "+ c.uri)
  conn, err := amqp.Dial(c.uri)
  c.Err <- err
  log.Printf("Connected")

  ch, err := conn.Channel()
  c.Err <- err
  log.Printf("Channel opened")

  return conn, ch
}

func (c *Consumer) createExpirableQueue(ch *amqp.Channel) (amqp.Queue) {
  var args amqp.Table
  args = make(amqp.Table)
  args["x-expires"] = int32(10000)

  q, err := ch.QueueDeclare("tailmq_"+uuid.NewV4().String(), false, true, false, false, args)
  c.Err <- err
  log.Printf("Queue defined")
  
  if c.routingKey != "" {
    err = ch.QueueBind(q.Name, c.routingKey, c.exchangeName, false, nil)
    c.Err <- err
    log.Printf("Queue " + q.Name + " binded to exchange "+ c.exchangeName + " on routing key " + c.routingKey)
  } else {
    err = ch.QueueBind(q.Name, "#", c.exchangeName, false, nil)
    c.Err <- err
    err = ch.QueueBind(q.Name, "", c.exchangeName, false, nil)
    c.Err <- err
    log.Printf("Queue " + q.Name + " binded to exchange "+ c.exchangeName)
  }

  return q
}
