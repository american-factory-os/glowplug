/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/american-factory-os/glowplug/service"
	"github.com/american-factory-os/glowplug/version"
	"github.com/gopcua/opcua"
	"github.com/spf13/cobra"
)

// opcuaCmd represents the opcua command
var opcuaCmd = &cobra.Command{
	Use:   "opcua",
	Short: "Start monitoring opcua with glowplug",
	Long: `This command will start the glowplug opcua client. 

It will read ths specified OPC UA nodes and optionally publish human
readable data to Redis and MQTT.`,
	Run: func(cmd *cobra.Command, args []string) {

		logger := service.NewLogger()
		if len(version.Version) > 0 || len(version.Revision) > 0 {
			logger.Printf("glowplug version %s %s", version.Version, version.Revision)
		}

		interval := cmd.Flag("interval").Value.String()
		intervalValue, err := time.ParseDuration(interval)
		if err != nil {
			panic(err)
		}

		nodes := cmd.Flag("nodes").Value.String()
		nodeList := []string{}
		if err := json.Unmarshal([]byte(nodes), &nodeList); err != nil {
			logger.Fatalf("nodes flag: %v, expecting something like '[\"ns=3;i=1005\"]'", err)
		}

		client, err := service.NewOpcuaClient(logger, service.OpcuaClientOpts{
			RedisURL:         cmd.Flag("redis").Value.String(),
			PublishBrokerURL: cmd.Flag("mqtt").Value.String(),
			Endpoint:         cmd.Flag("endpoint").Value.String(),
			Policy:           cmd.Flag("policy").Value.String(),
			Mode:             cmd.Flag("mode").Value.String(),
			CertFile:         cmd.Flag("cert").Value.String(),
			KeyFile:          cmd.Flag("key").Value.String(),
			Nodes:            nodeList,
			Interval:         intervalValue,
		})

		if err != nil {
			logger.Fatalf("unable to create OPC UA client, %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

		go func(c context.Context) {
			if err := client.Start(c); err != nil {
				logger.Fatal(err)
			}
		}(ctx)

	FOR:
		for {
			select {
			case <-ctx.Done():
				client.Stop()

				// wait for client to stop
				time.Sleep(500 * time.Millisecond)

				break FOR
			case caught := <-sig:
				logger.Println(caught, "signal caught")
				cancel()
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(opcuaCmd)
	opcuaCmd.Flags().StringP("redis", "r", "redis://localhost:6379/0", "Redis URL to store OPC UA node data")
	opcuaCmd.Flags().StringP("nodes", "n", "", "JSON array of node id's to subscribe to, e.g. '[\"ns=3;i=1005\"]'")
	opcuaCmd.Flags().StringP("endpoint", "e", "", "OPC UA server, e.g. opc.tcp://localhost:53530/OPCUA/SimulationServer")
	opcuaCmd.Flags().StringP("policy", "p", "Basic256Sha256", "Security policy: None, Basic128Rsa15, Basic256, Basic256Sha256")
	opcuaCmd.Flags().DurationP("interval", "i", opcua.DefaultSubscriptionInterval, "opcua subscription interval in milliseconds.")
	// TODO: remove default cert/key values
	opcuaCmd.Flags().StringP("cert", "c", "cert/public.der", "Path to cert.pem. Required for security mode/policy != None")
	opcuaCmd.Flags().StringP("key", "k", "cert/default_pk.pem", "Path to private key.pem. Required for security mode/policy != None")
	opcuaCmd.Flags().StringP("mode", "m", "auto", "Security mode: auto, None, Sign, SignAndEncrypt.")
	opcuaCmd.Flags().StringP("mqtt", "q", "", "Publish human readable OPC UA node data to this MQTT broker, e.g. mqtt://localhost:1883")
}
