# tailmq
Tail messages from a RabbitMQ exchange into your CLI console

# Installation

## Linux

```bash
curl -O https://github.com/olaurendeau/tailmq/releases/download/v1.0.0/tailmq-linux-amd64
mv tailmq-linux-amd64 /usr/local/bin/tailmq
rm tailmq-linux-amd64
```

## MacOS

```bash
curl -O https://github.com/olaurendeau/tailmq/releases/download/v1.0.0/tailmq-darwin-amd64
sudo mv tailmq-darwin-amd64 /usr/local/bin/tailmq
rm tailmq-darwin-amd64
```

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
DESCRIPTION
  TailMQ tail AMQP exchanges and output messages in stdout
USAGE
  tailmq [options] <exchange_name>
EXAMPLES
  tailmq amp.direct - Tail exchange amp.direct on local server with default access
  tailmq -uri=amqp://user:password@tailmq.com:5672//awesome amp.topic - Tail exchange amp.topic from server tailmq.com in vhost /awesome
  tailmq -server=prod amp.fanout - Tail exchange amp.fanout from server prod configured in file ~/.tailmq
  tailmq -server=prod -vhost=/foobar amp.fanout - Tail exchange amp.fanout from server prod configured in file ~/.tailmq but use vhost /foobar
OPTIONS
  -config string
    	Path of the global config file to use
  -help
    	How does it work ?
  -prefix
    	Should output be prefixed with datetime and time
  -server string
    	Use predefined server from configuration
  -uri string
    	RabbitMQ amqp uri (default "amqp://guest:guest@localhost:5672/")
  -verbose
    	Do you want more informations ?
  -vhost string
    	Define vhost to tail from
```

# Config file

## Format

`~/.tailmq`
```yaml
servers:
    server_name: amqp_uri    
```

## Sample

`~/.tailmq`
```yaml
servers:
    local: amqp://localhost:5672/
    staging: amqp://staging.tailmq.io:5672/
    staging_the_vhost: amqp://staging.tailmq.io:5672/the_vhost
    prod: amqp://tailmq.io:5672/
```
