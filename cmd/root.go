/*
Copyright © 2024 Taisuke Miyazaki <imishinist@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/tenntenn/natureremo"
)

type Metrics struct {
	Temperature  *prometheus.GaugeVec
	Humidity     *prometheus.GaugeVec
	Illumination *prometheus.GaugeVec
	Movement     *prometheus.GaugeVec
}

func NewMetrics() *Metrics {
	namespace := "nature_remo"
	deviceLabels := []string{
		"id",
		"name",
		"firmware_version",
		"mac_address",
		"serial_number",
	}

	temperature := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "temperature",
		Help:      "current temperature",
	}, deviceLabels)
	humidity := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "humidity",
		Help:      "current humidity",
	}, deviceLabels)
	illumination := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "illumination",
		Help:      "current illumination",
	}, deviceLabels)
	movement := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "movement",
		Help:      "current movement",
	}, deviceLabels)
	return &Metrics{
		Temperature:  temperature,
		Humidity:     humidity,
		Illumination: illumination,
		Movement:     movement,
	}
}

func (m *Metrics) Set(devices []*natureremo.Device) error {
	for _, device := range devices {
		labels := prometheus.Labels{
			"id":               device.ID,
			"name":             device.Name,
			"firmware_version": device.FirmwareVersion,
			"mac_address":      device.MacAddress,
			"serial_number":    device.SerialNumber,
		}
		m.Temperature.With(labels).Set(device.NewestEvents[natureremo.SensorTypeTemperature].Value)
		m.Humidity.With(labels).Set(device.NewestEvents[natureremo.SensorTypeHumidity].Value)
		m.Illumination.With(labels).Set(device.NewestEvents[natureremo.SensorTypeIllumination].Value)
		m.Movement.With(labels).Set(device.NewestEvents[natureremo.SensorTypeMovement].Value)
	}
	return nil
}

var (
	port     int
	interval time.Duration

	accessToken string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "nature-remo-exporter",
		Short: "A Prometheus exporter for Nature Remo",
		Long: `Nature Remo Exporter is a Prometheus exporter for Nature Remo smart devices.

This tool collects metrics from Nature Remo Cloud API and exposes them in a format 
that Prometheus can scrape. It is designed to help monitor and analyze 
the performance and data from Nature Remo devices`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

			client := natureremo.NewClient(accessToken)
			metrics := NewMetrics()

			update := func(ctx context.Context) error {
				devices, err := client.DeviceService.GetAll(ctx)
				if err != nil {
					return fmt.Errorf("failed to get all devices from Nature Remo API: %v", err)
				}
				if err := metrics.Set(devices); err != nil {
					return fmt.Errorf("failed to set metrics: %v", err)
				}
				return nil
			}

			go func() {
				if err := update(cmd.Context()); err != nil {
					logger.Error(err.Error())
				}

				ticker := time.NewTicker(interval)
				defer ticker.Stop()
				for {
					select {
					case <-cmd.Context().Done():
						logger.Info("shutting down")
						return
					case <-ticker.C:
						if err := update(cmd.Context()); err != nil {
							logger.Error(err.Error())
							return
						}
						logger.Debug("metrics updated")
					}
				}
			}()

			reg := prometheus.NewRegistry()
			reg.MustRegister(collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
			reg.MustRegister(metrics.Temperature, metrics.Humidity, metrics.Illumination, metrics.Movement)
			http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

			logger.Info(fmt.Sprintf("Listening on port %d", port))
			if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
				return err
			}
			return nil
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVar(&port, "port", 9199, "Port to listen on")
	rootCmd.PersistentFlags().DurationVar(&interval, "interval", time.Second*30, "Interval between metrics refresh")
	rootCmd.PersistentFlags().StringVar(&accessToken, "token", "", "Nature Remo access token")
}
