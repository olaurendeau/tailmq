# tailmq
Tail messages from a RabbitMQ exchange into your CLI console

# Installation

# Usage examples

Dump messages from an exchange to your console
```bash
$ tailmq -u amqp://guest:guest@localhost:5672/ amq.direct
{"id":"592fd0066fa63","method":"createDocument","params":{"email":"john.doe@tutu.com"}}
{"id":"592fd00670cbf","method":"createDocument","params":{"email":"john.doe@tutu.com"}}
{"id":"592fd00674483","method":"createDocument","params":{"email":"john.doe@tutu.com"}}
{"id":"592fd006758ad","method":"createDocument","params":{"email":"john.doe@tutu.com"}}
```

Retrieve a specific value from messages
```bash
$ tailmq -u amqp://guest:guest@localhost:5672/ amq.direct | jq '.params.email'
"john.doe@tutu.com"
"john.doe@tutu.com"
```

Look for specific messages
```bash
$ tailmq -u amqp://guest:guest@localhost:5672/ amq.direct | grep 74
{"id":"592fe92c06747","method":"createDocument","params":{"email":"john.doe@tutu.com"}}
{"id":"592fe92c174e1","method":"createDocument","params":{"email":"john.doe@tutu.com"}}
```

# Usage

```bash
$ tailmq --help
NAME:
   tailmq - Tail a RabbitMQ exchange

USAGE:
   tailmq [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --uri value, -u value  RabbitMQ amqp uri (default: "amqp://guest:guest@localhost:5672/vhost")
   --prefix               Should output be prefixed with date and time
   --verbose              Do you want more informations ?
   --help, -h             show help
   --version, -v          print the version```