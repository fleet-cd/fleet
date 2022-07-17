/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"github.com/spf13/cobra"
	"github.com/tgs266/fleet/config"
	"github.com/tgs266/fleet/fleet/auth"
	"github.com/tgs266/fleet/fleet/persistence"
)

// createAdminCmd represents the createAdmin command
var createAdminCmd = &cobra.Command{
	Use:   "createAdmin",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		email, err := cmd.Flags().GetString("email")
		if err != nil {
			panic(err)
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			panic(err)
		}
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			panic(err)
		}
		config := config.Read(configPath)
		persistence.Connect(config)

		auth.CreateAdminUser(email, password)
	},
}

func init() {
	rootCmd.AddCommand(createAdminCmd)
	createAdminCmd.Flags().StringP("config", "c", "config.yaml", "Fleet config file path")
	createAdminCmd.Flags().String("email", "", "")
	createAdminCmd.Flags().String("password", "", "")
}
