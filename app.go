package main

import (
  "log"
  "flag"
  "fmt"
  "net/url"
  "strings"
)

type App struct {
	uri *string
	server *string
	vhost *string
	verbose *bool
	prefix *bool
	exchangeName string
	configList []AmqpConfig
}

func newApp() *App {
	return &App{}
}

func (app *App) run() {
	// tailmq -uri=amqp://... -prefix <exchange> [routing_pattern]
	// tailmq -server=prod -vhost=/event <exchange> [routing_pattern]
	// tailmq configure <preprod> amqp://...

	app.defineFlagsAndParse()
	app.configList = getAmqpConfigList()
	app.setUri()
	tail(*app)
}

func (app *App) defineFlagsAndParse() {
	// Define flags
	app.uri = flag.String("uri", "amqp://guest:guest@localhost:5672/", "RabbitMQ amqp uri")
	app.server = flag.String("server", "", "Use predefined server from configuration")
	app.vhost = flag.String("vhost", "", "Define vhost to tail from")
	app.prefix = flag.Bool("prefix", false, "Should output be prefixed with datetime and time")
	app.verbose = flag.Bool("verbose", false, "Do you want more informations ?")
	flag.Parse()

	if (flag.NArg() == 0) {
		log.Fatalf("Please choose an exchange to listen from")
	} else if (flag.NArg() == 1) {
		app.exchangeName = flag.Arg(0)
	} else {
		log.Fatalf("Not yet available")
	}	
}

func (app *App) setUri() {
	if (*app.server != "" && len(app.configList) > 0) {
		config, err := getConfig(*app.server, app.configList)
		app.failOnError(err, "Failed to find config")

		*app.uri = config.Uri
	}

	if (*app.vhost != "") {
		parsedUri, err := url.Parse(*app.uri)
		app.failOnError(err, "Failed to parse uri")

		// If vhost start with a single slash it would be removed by Uri String() so we double it
		if (strings.Index(*app.vhost, "/") == 0) {
			parsedUri.Path = "/" + *app.vhost
		} else {
			parsedUri.Path = *app.vhost
		}

		*app.uri = parsedUri.String()
	}
}

func (app *App) info(msg string) {
  if *app.verbose {
    log.Printf(msg)
  }
}

func (app *App) failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
    panic(fmt.Sprintf("%s: %s", msg, err))
  }
}

