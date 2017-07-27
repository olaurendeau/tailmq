package main

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

func openChannel(app App) (*amqp.Connection, *amqp.Channel) {
	
    app.info("Establishing connection... "+ *app.uri)
    conn, err := amqp.Dial(*app.uri)
    app.failOnError(err, "Failed to connect to RabbitMQ")
    app.info("Connected")

    ch, err := conn.Channel()
    app.failOnError(err, "Failed to open a channel")
    app.info("Channel opened")

    return conn, ch
}

func createExpirableQueue(app App, ch *amqp.Channel) amqp.Queue {
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

    app.failOnError(err, "Failed to declare a queue")

    err = ch.QueueBind(
      q.Name, // name
      "#",   // routing key
      app.exchangeName,   // exchange name
      false,   // exclusive
      nil,     // arguments
    )
    app.failOnError(err, "Failed to bind topic / fanout queue")

    err = ch.QueueBind(
      q.Name, // name
      "",   // routing key
      app.exchangeName,   // exchange name
      false,   // exclusive
      nil,     // arguments
    )
    app.failOnError(err, "Failed to bind direct queue")

    app.info("Queue defined")

    return q
}

func tail(app App) {

	conn, ch := openChannel(app)
    defer conn.Close()
    defer ch.Close()

    q := createExpirableQueue(app, ch)

    msgs, err := ch.Consume(
      q.Name, // queue
      "",     // consumer
      true,   // auto-ack
      false,  // exclusive
      false,  // no-local
      false,  // no-wait
      nil,    // args
    )
    app.failOnError(err, "Failed to register a consumer")

    forever := make(chan bool)

    go func() {
      for d := range msgs {
        if *app.prefix {
          fmt.Printf("[%s]", time.Now().Format("2006-01-02 15:04:05"))
          if d.RoutingKey != "" {
            fmt.Printf(" %s ", d.RoutingKey)
          }
          fmt.Printf(" ")
        }
        fmt.Printf("%s\n", d.Body)
      }
    }()

    app.info(" [*] Waiting for messages. To exit press CTRL+C")
    <-forever
}
