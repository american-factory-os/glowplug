# Glowplug
**Glowplug** makes it easier to work with Sparkplug B data over MQTT for Industrial IoT (IIoT) applications.

## Why Glowplug?
If you're new to IIoT or industrial automation, Sparkplug data can seem overwhelming and it's binary encoding isn't human readable. Every Sparkplug payload has a list of metrics in various data types,  and Glowplug is designed to demystify Sparkplug data, by making metric values more accessible and human readable by extracting them from a payload into a Unified Namespace (UNS) context.

If you need to know the current value of a specific metric without talking to a MQTT broker, Glowplug makes it available for you in Redis. You can also subscribe to this value in Redis.

If you want to republish individual metrics to a UNS, Glowplug can optionally publish them for you.

## Quickstart

If you have Go installed, install glowplug from the command line:

```bash
go install github.com/american-factory-os/glowplug@latest

# start glowplug listening and publishing to your local broker
glowplug start --mqtt mqtt://localhost:1883 --publish mqtt://localhost:1883
```

Or, try providing a redis url:
```bash
glowplug start --mqtt mqtt://localhost:1883 --redis redis://localhost:6379 
```

## Features

* MQTT Integration: Connect to your MQTT broker and listen for Sparkplug B messages.
* Sparkplug B Support: Consumes all Sparkplug v3.0 messages and topics
* Redis Support: Store the last known values of Sparkplug metrics in Redis, and optionally subscribe to changes.
* UNS compatible: Publishes to unique and consistent Redis keys and MQTT topic names.
* Human-readable: Optionally publish Sparkplug metrics to human-readable MQTT topics.

## MQTT
The flag `--mqtt` contains the MQTT broker URL (e.g. [Mosquitto](https://github.com/eclipse/mosquitto)) to connect to and listen for Sparkplug B messages. When providing a MQTT url, it should start with "ws", "wss", or "mqtt".

Enable the flag `--publish` with another broker URL (could be the same one!) to publish all metrics to unique MQTT topics. 

*Note: the `--publish` flag will generate new topics in your broker for each Sparkplug metric published!*

### Example Usage
```bash
glowplug start --publish mqtt://localhost:1883
```
View your MQTT broker directly with [MQTT Explorer](https://mqtt-explorer.com/).

## Redis
If the flag `--redis` contains a valid redis URL, glowplug will store all Sparkplug metrics from birth and data messages in a set, and publish them to a channel of the same key as the set.

### Example Usage

```bash
glowplug start --redis redis://localhost:6379
```

You can explore glowplug data in Redis with [Redis Insight](https://redis.io/insight/). All data is prefixed with the value `glowplug`, so you can search for `glowplug:*`. Glowplug also stores all the metric data types it has seen in the hash `glowplug:metric_types`. 

If you are using [Node Red](https://nodered.org/), the [node-red-contrib-redis](https://flows.nodered.org/node/node-red-contrib-redis) module makes it easy to consume a redis channel containing Sparkplug metric data using the SUBSCRIBE method.

<img src="example/redis-in-node-red.png" />

## UNS Namespace

Glowplug ensures a consistent and unique namespace for all Redis keys and MQTT topics. For example, given a Sparkplug payload (using the [parris method](https://www.hivemq.com/blog/implementing-unified-namespace-uns-mqtt-sparkplug/)) that contains a device metric of data type "Float" named "Current/Celsius":

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

## Sparkplug Go Module

Glowplug also includes a Go module for interacting with Sparkplug B protocol buffers. Install it with:

```bash
go get github.com/american-factory-os/glowplug/sparkplug
```

See `example/sparkplug_and_mqtt` for example usage.

