**Glowplug** makes it easier to work with Sparkplug B data over MQTT for Industrial IoT (IIoT) applications.

https://github.com/american-factory-os/glowplug

## Why Glowplug?

If you're new to IIoT or industrial automation, Sparkplug data can be difficult to work with, and it's binary encoding is not human readable. Glowplug is designed to demystify Sparkplug data, by making every metric value more accessible and human readable. 

If you need to know the current value of a specific metric on your edge without talking to a MQTT broker, Glowplug makes it available for you in Redis. You can also subscribe to that value in Redis when it updates.

If you want individual metrics exposed in a UNS when they update, Glowplug can optionally publish them for you.

## Quickstart

You will need to set network to "host" so docker can access local services, e.g.:

```bash
docker run --network="host" aphexddb/glowplug:latest start --redis redis://localhost:6379/0 --publish mqtt://localhost:1883
```
