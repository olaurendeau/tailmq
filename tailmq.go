package main

import (
  "fmt"
  "os"
  "log"
  "time"

  "github.com/satori/go.uuid"
  "github.com/urfave/cli"
  "github.com/streadway/amqp"
)

var verbose bool

func main() {

  var uri string
  var exchangeName string
  var prefix bool
  
  cli.VersionFlag = cli.BoolFlag{
    Name: "version",
    Usage: "print only the version",
  }

  app := cli.NewApp()
  app.Name = "tailmq"
  app.Usage = "Tail a RabbitMQ exchange"
  app.Version = "0.1.0"
  app.Compiled = time.Now()
  app.ArgsUsage = "[exchangeName]"

  app.Flags = []cli.Flag {

    cli.StringFlag{
      Name: "uri, u",
      Value: "amqp://guest:guest@localhost:5672/",
      Usage: "RabbitMQ amqp uri",
      Destination: &uri,
    },

    cli.BoolFlag{
      Name: "prefix",
      Usage: "Should output be prefixed with date and time",
      Destination: &prefix,
    },

    cli.BoolFlag{
      Name: "verbose",
      Usage: "Do you want more informations ?",
      Destination: &verbose,
    },
  }

  app.Action = func(c *cli.Context) error {

    if c.NArg() <= 0 {
      log.Fatalf("Please choose an exchange to listen from")
    }
    exchangeName = c.Args().Get(0)

    info("Establishing connection... "+uri)
    conn, err := amqp.Dial(uri)
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()
    info("Connected")

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()
    info("Channel opened")

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

    failOnError(err, "Failed to declare a queue")
    err = ch.QueueBind(
      q.Name, // name
      "#",   // routing key
      exchangeName,   // exchange name
      false,   // exclusive
      nil,     // arguments
    )
    failOnError(err, "Failed to bind topic / fanout queue")

    err = ch.QueueBind(
      q.Name, // name
      "",   // routing key
      exchangeName,   // exchange name
      false,   // exclusive
      nil,     // arguments
    )
    failOnError(err, "Failed to bind direct queue")

    info("Queue defined")

    msgs, err := ch.Consume(
      q.Name, // queue
      "",     // consumer
      true,   // auto-ack
      false,  // exclusive
      false,  // no-local
      false,  // no-wait
      nil,    // args
    )
    failOnError(err, "Failed to register a consumer")

    forever := make(chan bool)

    go func() {
      for d := range msgs {
        if prefix {
          fmt.Printf("[%s]", time.Now().Format("2006-01-02 15:04:05"))
          if d.RoutingKey != "" {
            fmt.Printf(" %s ", d.RoutingKey)
          }
          fmt.Printf(" ")
        }
        fmt.Printf("%s\n", d.Body)
      }
    }()

    info(" [*] Waiting for messages. To exit press CTRL+C")
    <-forever
    
    return nil
  }

  app.Run(os.Args)
}

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
    panic(fmt.Sprintf("%s: %s", msg, err))
  }
}

func info(msg string) {
  if verbose {
    log.Printf(msg)
  }
}
