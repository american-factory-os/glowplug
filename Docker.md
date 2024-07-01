# Docker

## Build

Example:

```bash
docker buildx build --build-arg PROJECT_VERSION=0.0.0-devel -t american-factory-os/glowplug:latest .
```

Or without [BuildKit](https://docs.docker.com/build/):

```bash
docker build --build-arg PROJECT_VERSION=0.0.0-devel -t american-factory-os/glowplug:latest .
```

Push a tag to docker hub

```bash
docker tag 03b49dcb30eb american-factory-os/glowplug:latest
docker push american-factory-os/glowplug:latest
```

## Local development

You will need to set network to "host" so docker can access local services, e.g.:

```bash
docker run --rm --name glowplug --network="host" american-factory-os/glowplug:latest start --redis redis://localhost:6379/0 --publish mqtt://localhost:1883
```