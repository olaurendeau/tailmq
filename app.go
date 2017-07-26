package main

import (
  "log"
  "flag"
)

type App struct {
	uri *string
	verbose *bool
	prefix *bool
	exchangeName string
	configList []AmqpConfig
}

func newApp() *App {
	return &App{}
}

func (app *App) Run() {
	app.DefineFlagsAndParse()
	tail(*app)
}

func (app *App) DefineFlagsAndParse() {
	// Define flags
	app.uri = flag.String("uri", "amqp://guest:guest@localhost:5672/", "RabbitMQ amqp uri")
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

func (app *App) info(msg string) {
  if *app.verbose {
    log.Printf(msg)
  }
}
