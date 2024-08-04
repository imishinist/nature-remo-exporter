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
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tenntenn/natureremo"
)

var (
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
			client := natureremo.NewClient(accessToken)

			devices, err := client.DeviceService.GetAll(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get all devices from Nature Remo API: %v", err)
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(devices); err != nil {
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nature-remo-exporter.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVar(&accessToken, "token", "", "Nature Remo access token")
}
