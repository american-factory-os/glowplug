/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/american-factory-os/glowplug/service"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start glowplug",
	Long: `This command will start the glowplug service. 

It will read Sparkplug B data over MQTT and optionally publish human
readable data to MQTT.`,

	Run: func(cmd *cobra.Command, args []string) {

		logger := service.NewLogger()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

		svc, err := service.New(logger, service.Opts{
			MQTTBrokerURL:    cmd.Flag("mqtt").Value.String(),
			RedisURL:         cmd.Flag("redis").Value.String(),
			PublishBrokerURL: cmd.Flag("publish").Value.String(),
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
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringP("mqtt", "m", "mqtt://localhost:1883", "MQTT broker URL to listen for Sparkplug messages")
	startCmd.PersistentFlags().StringP("publish", "p", "", "Publish human readable Sparkplug metrics values to this MQTT broker, e.g. mqtt://localhost:1883")
	startCmd.PersistentFlags().StringP("redis", "r", "", "Redis URL to store Sparkplug data, e.g. redis://localhost:6379/0")
}
