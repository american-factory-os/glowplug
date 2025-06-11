# Docker usage

Use docker compose to spin up a MQTT broker, Redis, and **[Glowplug](https://github.com/american-factory-os/glowplug)** for local development

```bash
docker compose up
```

Or, if you already have a MQTT broker setup, start listening:

```bash
# Listen only
docker run --rm -it -p 8000:8000 aphexddb/glowplug:latest listen --broker mqtt://localhost:1883 --http 8000

# Listen and publish to redis and MQTT
docker run --rm -it -p 8000:8000 aphexddb/glowplug:latest listen --redis redis://localhost:6379/0 --publish mqtt://localhost:1883 --broker mqtt://localhost:1883 --http 8000
```

You can also build and run the docker image locally:

```bash
docker build --tag glowplug .

# Listen only
docker run --rm -it -p 8000:8000 glowplug listen --broker mqtt://localhost:1883 --http 8000

# Listen and publish to redis and MQTT
docker run --rm -it -p 8000:8000 glowplug listen --redis redis://localhost:6379/0 --publish mqtt://localhost:1883 --broker mqtt://localhost:1883 --http 8000
```

Note: You may need to set network to "host" so docker can access local services, e.g.: `--network="host"`