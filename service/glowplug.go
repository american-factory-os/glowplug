package service

import (
	"context"
	"fmt"
	"log"

	"github.com/american-factory-os/glowplug/sparkplug"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/redis/go-redis/v9"
)

type Glowplug interface {
	Start(ctx context.Context) error
	Stop() error
}

type Opts struct {
	MQTTBrokerURL    string
	PublishBrokerURL string
	RedisURL         string
	HTTPPort         int
}

type glowplug struct {
	logger *log.Logger
	opts   Opts
	wp     Worker
	broker mqtt.Client
}

// Start will start the glowplug service
func (g *glowplug) Start(ctx context.Context) error {

	g.logger.Println("starting glowplug")

	// subscribe to all sparkplug topics
	topic := fmt.Sprintf("%s/#", sparkplug.SPB_NS)
	if token := g.broker.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		return fmt.Errorf("unable to subscribe to %s, %w", topic, token.Error())
	}

	if len(g.opts.RedisURL) == 0 && len(g.opts.PublishBrokerURL) == 0 {
		g.logger.Println("warning: no redis or publish broker is enabled, glowplug wont do anything")
	}

	if len(g.opts.PublishBrokerURL) > 0 {
		g.logger.Println("publishing human readable data to broker", g.opts.PublishBrokerURL)
	} else {
		g.logger.Println("enable publishing metrics to a broker, ex: --publish mqtt://localhost:1883")
	}

	if len(g.opts.RedisURL) > 0 {
		g.logger.Println("using redis for metric storage", g.opts.RedisURL)
	} else {
		g.logger.Println("enable publishing metric values to redis, ex: --redis redis://localhost:6379/0")
	}

	httpListenAddr := ""
	if g.opts.HTTPPort > 0 {
		g.logger.Println("http and websocket port enabled", g.opts.HTTPPort)
		httpListenAddr = fmt.Sprintf("0.0.0.0:%d", g.opts.HTTPPort)
	} else {
		g.logger.Println("enable http and websocket port to expose metrics, ex: --http 8080")
	}

	go func() {
		err := g.wp.Run(httpListenAddr)
		if err != nil {
			g.logger.Fatal(err)
		}
	}()

	g.logger.Println("glowplug started")
	<-ctx.Done()

	return nil
}

// Stop will stop the glowplug service
func (g *glowplug) Stop() error {
	g.logger.Println("stopping glowplug")
	g.wp.Stop()
	return nil
}

// msgHandler is the default message handler for the glowplug service
func (g *glowplug) msgHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {

		cap, _ := g.wp.Capacity()
		if cap == 0 {
			g.logger.Println("dropping message, worker pool full")
			return
		}

		err := g.wp.AddMessage(Message{
			topic:   msg.Topic(),
			payload: msg.Payload(),
		})

		if err != nil {
			g.logger.Println("unable to process message,", err)
		}
	}
}

// New will create a new glowplug service
func New(logger *log.Logger, opts Opts) (Glowplug, error) {
	if len(opts.RedisURL) > 0 && len(opts.RedisURL) <= 14 {
		return nil, fmt.Errorf("redis URL too short")
	}

	if len(opts.MQTTBrokerURL) <= 14 {
		return nil, fmt.Errorf("mqtt broker URL too short: %s", opts.MQTTBrokerURL)
	}

	if len(opts.PublishBrokerURL) > 0 && len(opts.PublishBrokerURL) <= 14 {
		return nil, fmt.Errorf("publish broker URL too short: %s", opts.PublishBrokerURL)
	}

	var rdb *redis.UniversalClient
	if len(opts.RedisURL) > 0 {
		logger.Println("connecting to redis", opts.RedisURL)
		redisClient, err := NewRedis(opts.RedisURL)
		if err != nil {
			return nil, err
		}
		rdb = redisClient
	}

	g := glowplug{
		logger: logger,
		opts:   opts,
	}

	logger.Println("connecting to mqtt broker", opts.MQTTBrokerURL)
	handler := g.msgHandler()
	broker, err := brokerClientFromURL(opts.MQTTBrokerURL, &handler)
	if err != nil {
		return nil, err
	}
	g.broker = broker

	var publishBroker *mqtt.Client = nil
	if len(opts.PublishBrokerURL) > 0 {
		logger.Println("connecting to mqtt publish broker", opts.PublishBrokerURL)
		pb, pErr := brokerClientFromURL(opts.PublishBrokerURL, nil)
		if pErr != nil {
			return nil, pErr
		}
		publishBroker = &pb
	}

	wss := NewWebsocketServer(logger)

	wp, err := NewWorker(logger, rdb, publishBroker, wss)
	if err != nil {
		return nil, err
	}
	g.wp = wp

	return &g, nil
}
