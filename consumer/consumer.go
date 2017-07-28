package consumer

import (
	"log"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

type Consumer struct {
	uri string
	exchangeName string
	conn *amqp.Connection
	ch *amqp.Channel
	Deliveries <-chan amqp.Delivery
}

func New(uri string, exchangeName string) *Consumer {
	c := &Consumer{}
	c.uri = uri
	c.exchangeName = exchangeName

	return c
}

func (c *Consumer) Start() (err error) {
	c.conn, c.ch, err = c.openChannel()

	q, err := c.createExpirableQueue(c.ch)

    c.Deliveries, err = c.ch.Consume(q.Name, "", true, false, false, false, nil)

    return err
}

func (c *Consumer) Stop() {
    c.conn.Close()
    c.ch.Close()
}

func (c *Consumer) openChannel() (*amqp.Connection, *amqp.Channel, error) {
	
	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error

    log.Printf("Establishing connection... "+ c.uri)
    conn, err = amqp.Dial(c.uri)
    log.Printf("Connected")

    if err != nil {
	    log.Printf("Fail connecting")

    	return conn, ch, err
    }

    ch, err = conn.Channel()
    log.Printf("Channel opened")
    if err != nil {
	    log.Printf("Fail opening channel")
    }

    return conn, ch, err
}

func (c *Consumer) createExpirableQueue(ch *amqp.Channel) (amqp.Queue, error) {
    var args amqp.Table
    args = make(amqp.Table)
    args["x-expires"] = int32(10000)

    q, err := ch.QueueDeclare(
      "tailmq_"+uuid.NewV4().String(), // name
      false,   // durable
      true,    // delete when unused
      false,   // exclusive
      false,   // no-wait
      args,    // arguments
    )

    if err != nil {
	    log.Printf("Fail declaring queue")
    	return q, err
    }

    err = ch.QueueBind(
      q.Name, // name
      "#",   // routing key
      c.exchangeName,   // exchange name
      false,   // exclusive
      nil,     // arguments
    )

    if err != nil {
	    log.Printf("Fail binding queue")
		return q, err
    }

    err = ch.QueueBind(
      q.Name, // name
      "",   // routing key
      c.exchangeName,   // exchange name
      false,   // exclusive
      nil,     // arguments
    )

    if err != nil {
	    log.Printf("Fail binding queue")
		return q, err
    }

    log.Printf("Queue defined")

    return q, err
}