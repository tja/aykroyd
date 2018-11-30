package main

import (
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Main entry point
func main() {
	// Print banner
	color.NoColor = false

	color.HiCyan("                   __    ___ __                           __    ")
	color.HiCyan(".-----.-----.-----|  |_.'  _|__|--.--._____.--._.--.-----|  |--.")
	color.HiCyan("|  -  |  -  |__ --|   _|   _|  |-   -|_____|  | |  |  -__|  -  |")
	color.HiCyan("|   __|_____|_____|____|__| |__|__.__|     |___.___|_____|_____|")
	color.HiCyan("|__|                                                            ")
	color.HiCyan("                                                                ")

	// Cobra command flags
	cmd := &cobra.Command{
		Use:     "postfix-web",
		Long:    "Web interface for PostFix mail server.",
		Args:    cobra.NoArgs,
		Version: "0.0.1",
		Run:     run,
	}

	cmd.Flags().BoolP("verbose", "v", false, "Show more progress information")
	cmd.Flags().BoolP("quiet", "q", false, "Show less progress information")

	// Viper config
	viper.BindPFlags(cmd.Flags())

	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/postfix-web")
	viper.AddConfigPath(".")

	viper.ReadInConfig()

	// Run command
	cmd.Execute()
}

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
}
