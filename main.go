//go:generate go run github.com/UnnoTed/fileb0x ./pkg/assets/config.yaml

package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tja/aykroyd/pkg/backend"
)

// main is the main entry point of the app.
func main() {
	// Print banner
	color.NoColor = false

	color.HiCyan("             __                        __ ")
	color.HiCyan(".---.-.--.--|  |--.-.--.-----.--.--.--|  |")
	color.HiCyan("|  -  |  |  |    <|  .-|  -  |  |  |  -  |")
	color.HiCyan("|___._|___  |__|__|__| |_____|___  |_____|")
	color.HiCyan("         |__|                   |__|      ")
	color.HiCyan("                                          ")

	// Cobra command
	cmd := &cobra.Command{
		Use:     "aykroyd",
		Long:    "Email forwards via PostFix.",
		Args:    cobra.NoArgs,
		Version: "2.1.0",
		Run:     aykroyd,
	}

	cmd.Flags().BoolP("verbose", "v", false, "Write more")

	cmd.Flags().StringP("listen", "l", "0.0.0.0:80", "IP and port on which the server will listen")
	cmd.Flags().StringP("assets", "a", "", "Path to static web assets")

	cmd.Flags().StringP("db-host", "H", "localhost", "MySQL host")
	cmd.Flags().StringP("db-database", "d", "postfix", "MySQL database")
	cmd.Flags().StringP("db-username", "u", "postfix", "MySQL username")
	cmd.Flags().StringP("db-password", "p", "", "MySQL password")

	// Viper config
	viper.SetEnvPrefix("AYKROYD")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	viper.BindPFlags(cmd.Flags())

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/aykroyd")
	viper.AddConfigPath("$HOME/.config/aykroyd")
	viper.AddConfigPath(".")

	viper.ReadInConfig()

	// Run command
	cmd.Execute()
}

// aykroyd is called if the CLI interfaces has been satisfied.
func aykroyd(cmd *cobra.Command, args []string) {
	// Set logging level
	if viper.GetBool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Set up server
	server, err := backend.NewServer(
		viper.GetString("assets"),
		viper.GetString("db-host"),
		viper.GetString("db-database"),
		viper.GetString("db-username"),
		viper.GetString("db-password"),
	)

	if err != nil {
		logrus.Fatal(err)
	}

	defer server.Close()

	// Start listening
	httpServer := &http.Server{
		Handler:      server.Router,
		Addr:         viper.GetString("listen"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Infof("Listening on %s", httpServer.Addr)

	logrus.Fatal(httpServer.ListenAndServe())
}
