# Glowplug
**Glowplug** makes it easier to work with Sparkplug B data over MQTT for Industrial IoT (IIoT) applications.

## Why Glowplug?
If you're new to IIoT or industrial automation, Sparkplug data can be difficult to work with, and it's binary encoding is not human readable. Glowplug is designed to demystify Sparkplug data, by making every metric value more accessible and human readable. 

If you need to know the current value of a specific metric on your edge without talking to a MQTT broker, Glowplug makes it available for you in a browser, redis, or JSON in your MQTT broker. 

## Quickstart

Connect to a UNS and view SparkplugB metrics on a webpage (handy for exploring):
```bash
docker run -p 8000:8000 aphexddb/glowplug:latest listen -b mqtt://localhost:1883 -w 8000
```

Next, add a redis url to store current values...
```bash
docker run -p 8000:8000 aphexddb/glowplug:latest listen -b mqtt://localhost:1883 -w 8000 --redis redis://localhost:6379/0
```

...And then you can publish simplified metrics to your broker:
```bash
docker run -p 8000:8000 aphexddb/glowplug:latest listen -b mqtt://localhost:1883 -w 8000 --redis redis://localhost:6379/0 --publish mqtt://localhost:1883
```

Choose whatever flags suit your neeeds! Note: You may need to set network to "host" so docker can access local services, e.g.: `--network="host"`

Run 

## Features

* MQTT Integration: Connect to your MQTT broker and listen for Sparkplug B messages.
* Sparkplug B Support: Consumes all Sparkplug v3.0 messages and topics
* HTTP and websocket support: View your UNS data live in a browser
* Redis Support: Store the last known values of Sparkplug metrics in Redis, and optionally subscribe to changes.
* UNS compatible: Publishes to unique and consistent Redis keys and MQTT topic names.
* Human-readable: Optionally publish Sparkplug metrics to human-readable MQTT topics.

## HTTP and Websockets
* The flag `--http` or `-w` enables an HTTP server on the specified port. HTTP is disabled unless the flag is set.
* Sparkplug metrics are published on the path`/ws`, e.g. `ws://localhost:8000/ws`
Sparkplug metrics are published over HTTP and are viewable on the specified port, and are available via websocket at http://localhost:8000.

## MQTT
* The flag `--broker` or `-b` contains the MQTT broker glowplug will listen for Sparkplug messages.
  * The value defaults to `mqtt://localhost:1883` (commonly used for [mosquitto](https://github.com/eclipse/mosquitto)).
* The flag `--publish` or `-p` will publish each metrics to a unique topic in a UNS (more below on this).   
  * **Note:** this flag will generate a new topic in your broker for each Sparkplug metric published. If you have 100k's of tags there may be a compute impact.

View your MQTT broker directly with [MQTT Explorer](https://mqtt-explorer.com/).

## Redis
* The flag `--redis` or `-r` will specify the redis server for glowplug to store all Sparkplug metrics from birth and data messages in a [SET](https://redis.io/docs/latest/commands/set/), and publish them to a [channel](https://redis.io/docs/latest/commands/pubsub-channels/) of the same key as the set.

You can explore glowplug data in Redis with [Redis Insight](https://redis.io/insight/). All data is prefixed with the value `glowplug`, so you can search for `glowplug:*` to see all the keys. Glowplug also stores all the metric data types it has seen in the hash `glowplug:metric_types`. 

If you are using [Node Red](https://nodered.org/), the [node-red-contrib-redis](https://flows.nodered.org/node/node-red-contrib-redis) module makes it easy to consume a redis channel containing Sparkplug metric data using [SUBSCRIBE](https://redis.io/docs/latest/commands/subscribe/).

<img src="example/redis-in-node-red.png" />

PSA: You can also [subscribe](https://redis.io/docs/latest/develop/use/keyspace-notifications/) to keys in Redis when they update.

## UNS Namespace

Glowplug ensures a consistent and unique namespace for all Redis keys, and MQTT topics. For example, given a Sparkplug payload (using the [parris method](https://www.hivemq.com/blog/implementing-unified-namespace-uns-mqtt-sparkplug/)) that contains a device metric of data type "Float" named "Current/Celsius":

* `spBv1.0/Plant1:Area3:Line4:Cell2/DDATA/Heater/TempSensor`

will be mapped as follows:

- **Redis SET key**:
    ```txt
    glowplug:plant1:area3:line4:cell2:heater:tempsensor:current:celsius
    ```
- **Redis channel key**:
    ```txt
    glowplug:plant1:area3:line4:cell2:heater:tempsensor:current:celsius
    ```    
- **MQTT Topic**:
    ```txt
    glowplug/Plant1/Area3/Line4/Cell2/Heater/TempSensor/Current/Celsius
    ```

That same metric is published via websocket as follows:
```json
{
    "topic": {
        "command": "DDATA",
        "group_id": "Plant1:Area3:Line4:Cell2",
        "edge_node_id": "Heater",
        "device_id": "TempSensor",
        "has_device": true
    },
    "alias": 7,
    "name": "Current/Celsius",
    "value": 42.0
}
```

## Sparkplug Go Module

Glowplug also includes a Go module for interacting with Sparkplug B protocol buffers. Install it with:

```bash
go get github.com/american-factory-os/glowplug/sparkplug
```

See `example/sparkplug_and_mqtt` for example usage.

