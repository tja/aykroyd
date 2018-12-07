//go:generate go run github.com/UnnoTed/fileb0x ./pkg/assets/config.yaml

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tja/postfix-web/pkg/backend"
)

// main is the main entry point of the app.
func main() {
	// Print banner
	color.NoColor = false

	color.HiCyan("                   __    ___ __                           __    ")
	color.HiCyan(".-----.-----.-----|  |_.'  _|__|--.--._____.--._.--.-----|  |--.")
	color.HiCyan("|  -  |  -  |__ --|   _|   _|  |-   -|_____|  | |  |  -__|  -  |")
	color.HiCyan("|   __|_____|_____|____|__| |__|__.__|     |___.___|_____|_____|")
	color.HiCyan("|__|                                                            ")
	color.HiCyan("                                                                ")

	// Cobra command
	cmd := &cobra.Command{
		Use:     "postfix-web",
		Long:    "Web interface for PostFix mail server.",
		Args:    cobra.NoArgs,
		Version: "0.0.1",
		Run:     run,
	}

	cmd.Flags().BoolP("verbose", "v", false, "Show more progress information")
	cmd.Flags().BoolP("quiet", "q", false, "Show less progress information")

	cmd.Flags().StringP("bind", "b", "127.0.0.1", "Interface to which the server will bind")
	cmd.Flags().IntP("port", "p", 2105, "Port on which the server will listen")
	cmd.Flags().StringP("content", "c", "./web", "Path of folder with static content")

	cmd.Flags().StringP("database", "d", "", "Database connection string")

	// Viper config
	viper.BindPFlags(cmd.Flags())

	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/postfix-web")
	viper.AddConfigPath(".")

	viper.ReadInConfig()

	// Run command
	cmd.Execute()
}

// run is called if the CLI interfaces has been satisfied.
func run(cmd *cobra.Command, args []string) {
	// Set logging level
	switch {
	case viper.GetBool("verbose"):
		logrus.SetLevel(logrus.DebugLevel)
	case viper.GetBool("quiet"):
		logrus.SetLevel(logrus.WarnLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Set up server
	server, err := backend.NewServer(
		viper.GetString("content"),
		viper.GetString("database"),
	)

	if err != nil {
		logrus.Fatal(err)
	}

	defer server.Close()

	// Start listening
	httpServer := &http.Server{
		Handler:      server.Router,
		Addr:         fmt.Sprintf("%s:%d", viper.GetString("bind"), viper.GetInt("port")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Fatal(httpServer.ListenAndServe())
}
