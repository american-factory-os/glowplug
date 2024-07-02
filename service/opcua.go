package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/american-factory-os/glowplug/json_type"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/monitor"
	"github.com/gopcua/opcua/ua"
	"github.com/redis/go-redis/v9"
)

type OpcuaClient interface {
	Start(ctx context.Context) error
	Stop() error
}

type OpcuaClientOpts struct {
	RedisURL         string
	PublishBrokerURL string
	Endpoint         string
	Policy           string
	Mode             string
	CertFile         string
	KeyFile          string
	Nodes            []string
	Interval         time.Duration
}

type OpcuaServerInfo struct {
	ApplicationURI      string
	ProductURI          string
	ApplicationName     string
	ApplicationType     string
	GatewayServerURI    string
	DiscoveryProfileURI string
	DiscoveryURLs       []string
}

type opcuaClient struct {
	logger        *log.Logger
	opts          OpcuaClientOpts
	serverInfo    OpcuaServerInfo
	productURI    string
	rdb           *redis.UniversalClient
	publishBroker *mqtt.Client
	topics        map[string]string
}

func (c *opcuaClient) Start(ctx context.Context) error {

	// var nodeParamList []string
	// if err := json.Unmarshal([]byte(opts.Nodes), &nodeParamList); err != nil {
	// 	logger.Fatalf("unable to read supplied nodes, %v", err)
	// }

	if len(c.opts.Nodes) == 0 {
		c.logger.Fatal("no nodes were specified, at least one node must be provided")
	}

	c.logger.Println("validing node ids", c.opts.Nodes)

	nodes, err := parseNodeSlice(c.opts.Nodes...)
	if err != nil {
		c.logger.Fatal(err)
	}

	c.logger.Println("attempting to monitor opcua nodes ids:", nodes)

	servers, err := findServers(context.TODO(), c.opts.Endpoint)
	if err != nil {
		c.logger.Fatal(err)
	}

	if len(servers) == 0 {
		c.logger.Fatal("no servers found")
	}

	if len(servers) > 1 {
		c.logger.Println("multiple servers found, using first")
	}

	c.logger.Println("opcua server discovered")
	c.logger.Println("server ApplicationURI:", servers[0].ApplicationURI)
	c.logger.Println("server ProductURI:", servers[0].ProductURI)
	c.logger.Println("server ApplicationName:", servers[0].ApplicationName)
	c.logger.Println("server ApplicationType:", servers[0].ApplicationType)
	// c.logger.Println("server GatewayServerURI:", servers[0].GatewayServerURI)
	// c.logger.Println("server DiscoveryProfileURI:", servers[0].DiscoveryProfileURI)
	// c.logger.Println("server DiscoveryURLs:", servers[0].DiscoveryURLs)

	// store server context

	if len(servers[0].ProductURI) == 0 {
		panic("ProductURI is empty")
	}

	c.serverInfo = servers[0]
	c.productURI = servers[0].ProductURI

	endpoints, err := opcua.GetEndpoints(ctx, c.opts.Endpoint)
	if err != nil {
		c.logger.Fatal(err)
	}

	ep := opcua.SelectEndpoint(endpoints, c.opts.Policy, ua.MessageSecurityModeFromString(c.opts.Mode))
	if ep == nil {
		c.logger.Fatal("Failed to find suitable endpoint, check your security policy and mode settings")
	}

	c.logger.Println("opcua client using security policy", ep.SecurityPolicyURI, ep.SecurityMode)

	options := []opcua.Option{
		opcua.SecurityPolicy(c.opts.Policy),
		opcua.SecurityModeString(c.opts.Mode),
		opcua.CertificateFile(c.opts.CertFile),
		opcua.PrivateKeyFile(c.opts.KeyFile),
		opcua.AuthAnonymous(),
		opcua.SecurityFromEndpoint(ep, ua.UserTokenTypeAnonymous),
	}

	client, err := opcua.NewClient(ep.EndpointURL, options...)
	if err != nil {
		c.logger.Fatal(err)
	}
	if err := client.Connect(ctx); err != nil {
		c.logger.Fatal(err)
	}

	defer client.Close(ctx)

	m, err := monitor.NewNodeMonitor(client)
	if err != nil {
		c.logger.Fatal(err)
	}

	m.SetErrorHandler(func(_ *opcua.Client, sub *monitor.Subscription, err error) {
		c.logger.Printf("error: sub=%d err=%s", sub.SubscriptionID(), err.Error())
	})

	wg := &sync.WaitGroup{}

	// start callback-based subscription
	wg.Add(1)
	for _, node := range nodes {
		key, err := keyFromUaNodeId(c.productURI, node)
		if err != nil {
			c.logger.Fatal(err)
		}
		c.logger.Println("monitoring node", node.String(), key)

		// add topic to map of keys
		if topic, err := topicFromUaNode(c.productURI, node); err != nil {
			c.logger.Fatal("unable to create topic for this node", err)
		} else {
			c.topics[key] = topic
		}

		go c.startCallbackSub(ctx, key, m, c.opts.Interval, 0, wg, node.String())
	}

	// // start channel-based subscription
	// wg.Add(1)
	// for _, node := range nodes {
	// 	b, _ := json.Marshal(node)
	// 	c.logger.Println("monitoring node", string(b), "namespace", node.Namespace(), node.Type().String(), node.IntID(), node.StringID())

	// 	go startChanSub(logger, ctx, m, *interval, 0, wg, node.String())
	// }

	<-ctx.Done()
	wg.Wait()

	return nil
}

func (c *opcuaClient) Stop() error {
	return nil
}

func (c *opcuaClient) publishToRedis(key string, value interface{}) error {
	// pipeline redis commands
	if c.rdb != nil {
		rdb := *c.rdb
		if cmds, err := rdb.Pipelined(context.TODO(), func(pipeliner redis.Pipeliner) error {

			// store the metric value in a redis set
			pipeliner.Set(context.TODO(), key, value, 0)

			// publish metric value to redis channel
			pipeliner.Publish(context.TODO(), key, value)

			return nil
		}); err != nil {
			return err
		} else {
			for _, cmd := range cmds {
				if cmd.Err() != nil && cmd.Err() != redis.Nil {
					return fmt.Errorf("redis cmd error %w", cmd.Err())
				}
			}
		}
	}

	return nil
}

func NewOpcuaClient(logger *log.Logger, opts OpcuaClientOpts) (OpcuaClient, error) {

	var rdb *redis.UniversalClient
	if len(opts.RedisURL) > 0 {
		logger.Println("connecting to redis", opts.RedisURL)
		redisClient, err := NewRedis(opts.RedisURL)
		if err != nil {
			return nil, err
		}
		rdb = redisClient
	}

	var publishBroker *mqtt.Client = nil
	if len(opts.PublishBrokerURL) > 0 {
		logger.Println("connecting to mqtt publish broker", opts.PublishBrokerURL)
		pb, pErr := brokerClientFromURL(opts.PublishBrokerURL, nil)
		if pErr != nil {
			return nil, pErr
		}
		publishBroker = &pb
	}

	return &opcuaClient{
		logger:        logger,
		opts:          opts,
		rdb:           rdb,
		publishBroker: publishBroker,
		topics:        make(map[string]string),
	}, nil
}

func findServers(ctx context.Context, endpoint string) ([]OpcuaServerInfo, error) {
	servers, err := opcua.FindServers(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	resp := []OpcuaServerInfo{}

	for _, server := range servers {
		resp = append(resp, OpcuaServerInfo{
			ApplicationURI:      server.ApplicationURI,
			ProductURI:          server.ProductURI,
			ApplicationName:     server.ApplicationName.Text,
			ApplicationType:     server.ApplicationType.String(),
			GatewayServerURI:    server.GatewayServerURI,
			DiscoveryProfileURI: server.DiscoveryProfileURI,
			DiscoveryURLs:       server.DiscoveryURLs,
		})
	}

	return resp, nil
}

func (c *opcuaClient) startCallbackSub(ctx context.Context, key string, m *monitor.NodeMonitor, interval, lag time.Duration, wg *sync.WaitGroup, nodes ...string) {
	sub, err := m.Subscribe(
		ctx,
		&opcua.SubscriptionParameters{
			Interval: interval,
		},
		func(s *monitor.Subscription, msg *monitor.DataChangeMessage) {
			if msg.Error != nil {
				c.logger.Printf("opcua callback error sub=%d error=%s", s.SubscriptionID(), msg.Error)
			} else {
				go func() {
					if err := c.publishToRedis(key, msg.Value.Value()); err != nil {
						c.logger.Println("redis publish err", err)
					} else {
						// c.logger.Printf("opcua callback sub=%d node=%s value=%v", s.SubscriptionID(), msg.NodeID, msg.Value.Value())
					}

					if c.publishBroker != nil {
						if jsonType, err := json_type.NodeValueToJsonType(msg.Value); err != nil {
							c.logger.Println(msg.NodeID.String(), err)
						} else {
							topic, ok := c.topics[key]
							if !ok {
								c.logger.Println("topic not found for key", key)
							} else {
								if token := (*c.publishBroker).Publish(topic, 0, false, jsonType.Bytes()); token.Wait() && token.Error() != nil {
									c.logger.Println("mqtt publish err", token.Error())
								}
							}
						}
					}
				}()
			}
			time.Sleep(lag)
		},
		nodes...)

	if err != nil {
		c.logger.Fatal(err)
	}

	defer cleanup(c.logger, ctx, sub, wg)

	<-ctx.Done()
}

func startChanSub(logger *log.Logger, ctx context.Context, m *monitor.NodeMonitor, interval, lag time.Duration, wg *sync.WaitGroup, nodes ...string) {
	ch := make(chan *monitor.DataChangeMessage, 16)
	sub, err := m.ChanSubscribe(ctx, &opcua.SubscriptionParameters{Interval: interval}, ch, nodes...)

	if err != nil {
		logger.Fatal(err)
	}

	defer cleanup(logger, ctx, sub, wg)

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			if msg.Error != nil {
				logger.Printf("[channel ] sub=%d error=%s", sub.SubscriptionID(), msg.Error)
			} else {
				logger.Printf("[channel ] sub=%d ts=%s node=%s value=%v", sub.SubscriptionID(), msg.SourceTimestamp.UTC().Format(time.RFC3339), msg.NodeID, msg.Value.Value())
			}
			time.Sleep(lag)
		}
	}
}

func cleanup(logger *log.Logger, ctx context.Context, sub *monitor.Subscription, wg *sync.WaitGroup) {
	logger.Printf("stats: sub=%d delivered=%d dropped=%d", sub.SubscriptionID(), sub.Delivered(), sub.Dropped())
	sub.Unsubscribe(ctx)
	wg.Done()
}

func parseNodeSlice(nodes ...string) ([]*ua.NodeID, error) {
	var err error

	nodeIDs := make([]*ua.NodeID, len(nodes))

	for i, node := range nodes {
		if nodeIDs[i], err = ua.ParseNodeID(node); err != nil {
			return nil, err
		}
	}

	return nodeIDs, nil
}
