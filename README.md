# tailmq
Tail messages from a RabbitMQ exchange into your CLI console

# Installation

# Usage examples

Dump messages from an exchange to your console
```bash
$ tailmq amq.topic
{"id":"592fec584a25f","request_id":"592fd015bba42","channel":null,"message":"[10:28:40] 592fd015bba42 [worker] Invoice sent by email\n"}
{"id":"592fec584b798","request_id":"592fd015bba42","channel":null,"message":"[10:28:40] 592fd015bba42 [worker] Message processed\n"}
{"id":"592fec584cd07","request_id":"592fd015cad15","channel":null,"message":"[10:28:40] 592fd015cad15 [worker] Generating invoice\n"}
{"id":"592fec5901213","request_id":"592fd015c9bc8","channel":null,"message":"[10:28:41] 592fd015c9bc8 [worker] Invoice sent by email\n"}
{"id":"592fec59023a7","request_id":"592fd015c9bc8","channel":null,"message":"[10:28:41] 592fd015c9bc8 [worker] Message processed\n"}```

Retrieve a specific value from messages
```bash
$ tailmq amq.topic | jq '.message'
"[10:29:13] 592fd0166f482 [worker] Invoice sent by email\n"
"[10:29:13] 592fd0166f482 [worker] Message processed\n"
"[10:29:13] 592fd016780ce [worker] Generating invoice\n"
"[10:29:14] 592fd01670bc3 [worker] Invoice sent by email\n"
```

Look for specific messages
```bash
$ tailmq amq.topic | grep sent
{"id":"592fec63717ac","request_id":"592fd015f36f8","channel":null,"message":"[10:28:51] 592fd015f36f8 [worker] Invoice sent by email\n"}
{"id":"592fec6414d9d","request_id":"592fcf0447069","channel":null,"message":"[10:28:52] 592fcf0447069 [worker] Invoice sent by email\n"}
{"id":"592fec6448ae9","request_id":"592fd01600770","channel":null,"message":"[10:28:52] 592fd01600770 [worker] Invoice sent by email\n"}
```

# Usage

```bash
$ tailmq --help
NAME:
   tailmq - Tail a RabbitMQ exchange

USAGE:
   tailmq [global options] command [command options] [exchangeName]

VERSION:
   0.1.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --uri value, -u value  RabbitMQ amqp uri (default: "amqp://guest:guest@localhost:5672/")
   --prefix               Should output be prefixed with date and time
   --verbose              Do you want more informations ?
   --help, -h             show help
   --version              print only the version
```