package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fuhrmannb/duckmail"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgPath string

var rootCmd = &cobra.Command{
	Use:   "duckmail",
	Short: "Duckmail is an app that send notification when a mail has been received",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read config
		var cfg duckmail.RootCfg
		if err := viper.Unmarshal(&cfg); err != nil {
			cfgError(err)
		}

		// Start Duckmail controller
		if err := duckmail.StartController(&cfg); err != nil {
			return err
		}
		return nil
	},
}

func main() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "",
		`config file (by default, duckmail.yaml file is located in /etc/duckmail, $HOME/.duckmail or at current path)`)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgPath != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgPath)
	} else {
		// Find home directory
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.SetConfigName("duckmail")
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(home, ".duckmail"))
		viper.AddConfigPath("/etc/duckmail")
	}

	if err := viper.ReadInConfig(); err != nil {
		cfgError(err)
	}
}

func cfgError(err error) {
	fmt.Printf("Can't read config:%v\n, ", err)
	os.Exit(1)
}
