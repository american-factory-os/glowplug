package cmd

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/american-factory-os/glowplug/service"
	"github.com/american-factory-os/glowplug/version"
	"github.com/spf13/cobra"
)

// listenCmd represents the start command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for Sparkplug messages",
	Long: `This command will start glowplug. 

It will read Sparkplug B data over MQTT and optionally publish human
readable data to MQTT.`,

	Run: func(cmd *cobra.Command, args []string) {

		logger := service.NewLogger()
		if len(version.Version) > 0 || len(version.Revision) > 0 {
			logger.Printf("glowplug version %s %s", version.Version, version.Revision)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

		httpPort := 0
		httpPortStr := cmd.Flag("http").Value.String()
		if len(httpPortStr) > 0 {
			var err error
			httpPort, err = strconv.Atoi(httpPortStr)
			if err != nil {
				logger.Fatalf("invalid http port: %v", err)
			}
		}

		svc, err := service.New(logger, service.Opts{
			MQTTBrokerURL:    cmd.Flag("broker").Value.String(),
			RedisURL:         cmd.Flag("redis").Value.String(),
			PublishBrokerURL: cmd.Flag("publish").Value.String(),
			HTTPPort:         httpPort,
		})

		if err != nil {
			logger.Fatal(err)
		}

		go func(c context.Context) {
			if err := svc.Start(c); err != nil {
				logger.Fatal(err)
			}
		}(ctx)

	FOR:
		for {
			select {
			case <-ctx.Done():
				svc.Stop()

				// wait for glowplug to stop
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
	rootCmd.AddCommand(listenCmd)
	listenCmd.PersistentFlags().StringP("broker", "b", "mqtt://localhost:1883", "MQTT broker URL to listen for Sparkplug messages")
	listenCmd.PersistentFlags().StringP("publish", "p", "", "Publish human readable Sparkplug metrics values to this MQTT broker, e.g. mqtt://localhost:1883")
	listenCmd.PersistentFlags().StringP("redis", "r", "", "Redis URL to store Sparkplug data, e.g. redis://localhost:6379/0")
	listenCmd.PersistentFlags().IntP("http", "w", 0, "HTTP port that exposes Sparkplug data over websockets")
}
