package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "io/ioutil"
  "net/url"
  "strings"
  "time"
  "context"

  "github.com/olaurendeau/tailmq/consumer"
  "github.com/streadway/amqp"
)

const Help = `DESCRIPTION
  TailMQ tail AMQP exchanges and output messages in stdout
USAGE
  tailmq [options] <exchange_name>
EXAMPLES
  tailmq amp.direct - Tail exchange amp.direct on local server with default access
  tailmq -uri=amqp://user:password@tailmq.com:5672//awesome amp.topic - Tail exchange amp.topic from server tailmq.com in vhost /awesome
  tailmq -server=prod amp.fanout - Tail exchange amp.fanout from server prod configured in file ~/.tailmq
  tailmq -server=prod -vhost=/foobar amp.fanout - Tail exchange amp.fanout from server prod configured in file ~/.tailmq but use vhost /foobar
OPTIONS
`

type Config struct {
  uri *string
  server *string
  vhost *string
  verbose *bool
  prefix *bool
  help *bool
  exchangeName string
  globalConfigFilePath *string
  global GlobalConfig
}

func main() {
  var err error

  ctx := context.Background()

  config := new(Config)

  config.uri = flag.String("uri", "amqp://guest:guest@localhost:5672/", "RabbitMQ amqp uri")
  config.server = flag.String("server", "", "Use predefined server from configuration")
  config.vhost = flag.String("vhost", "", "Define vhost to tail from")
  config.prefix = flag.Bool("prefix", false, "Should output be prefixed with datetime and time")
  config.verbose = flag.Bool("verbose", false, "Do you want more informations ?")
  config.help = flag.Bool("help", false, "How does it work ?")
  config.globalConfigFilePath = flag.String("config", "", "Path of the global config file to use")
  flag.Parse()

  config.global, err = NewGlobalConfig(*config.globalConfigFilePath)
  failOnError(err, "Fail retrieving server list")

  configureLogger(config)
  displayHelp(config)
  setUri(config)
  
  if (flag.NArg() == 0) {
    log.Fatalf("Please choose an exchange to listen from")
  } else if (flag.NArg() == 1) {
    config.exchangeName = flag.Arg(0)
  } else {
    log.Fatalf("Not yet available")
  }

  c := consumer.New(*config.uri, config.exchangeName)
  go c.Start()
  defer c.Stop()

  for {
    select {
      case d := <-c.Deliveries:
        printDelivery(d, config)
      case err := <-c.Err:
        failOnError(err, "Fail consuming")
      case <-ctx.Done():
        return
    }
  }
}

func printDelivery(d amqp.Delivery, config *Config) {
  if *config.prefix {
    fmt.Printf("[%s]", time.Now().Format("2006-01-02 15:04:05"))
    if d.RoutingKey != "" {
      fmt.Printf(" %s ", d.RoutingKey)
    }
    fmt.Printf(" ")
  }
  fmt.Printf("%s\n", d.Body)
}

func configureLogger(config *Config) {
  if (!*config.verbose) {
    log.SetOutput(ioutil.Discard)
  }
}

func setUri(config *Config) {
  if (*config.server != "") {
    server, err := config.global.getServerUri(*config.server)
    failOnError(err, "Failed to find server configuration")

    *config.uri = server
  }

  if (*config.vhost != "") {
    parsedUri, err := url.Parse(*config.uri)
    failOnError(err, "Failed to parse uri")

    // If vhost start with a single slash it would be removed by Uri String() so we double it
    if (strings.Index(*config.vhost, "/") == 0) {
      parsedUri.Path = "/" + *config.vhost
    } else {
      parsedUri.Path = *config.vhost
    }

    *config.uri = parsedUri.String()
  }
}

func displayHelp(config *Config) {
  if *config.help {

    fmt.Print(Help)
    flag.PrintDefaults()

    os.Exit(0)
  }
}

func failOnError(err error, msg string) {
  if err != nil {
  log.Fatalf("%s: %s", msg, err)
  }
}
