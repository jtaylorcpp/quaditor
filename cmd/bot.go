package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jtaylorcpp/quaditor"
	"github.com/jtaylorcpp/quaditor/auditors"
	"github.com/jtaylorcpp/quaditor/chatbots"

	log "github.com/sirupsen/logrus"
)

func init() {
	rootCmd.AddCommand(botCmd)
}

var botCmd = &cobra.Command{
	Use:   "bot",
	Short: "run irc bot",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("starting auditor of type: ", viper.GetString("audit_type"))
		var auditor quaditor.Auditor
		switch viper.GetString("audit_type") {
		case "time-series":
			log.Println("using time-series backend type: ", viper.GetString("audit_backend"))
			switch viper.GetString("audit_backend") {
			case "postgresql":
				username := viper.GetString("audit_username")
				password := viper.GetString("audit_password")
				host := viper.GetString("audit_host")
				port := viper.GetString("audit_port")

				if username == "" ||
					password == "" ||
					host == "" ||
					port == "" {
					log.Println("username, password, host, or port unset")
					return
				}

				tsaudit, err := auditors.NewTimeSeriesAuditor("postgres", username, password, host, port)
				if err != nil {
					log.Println("error with new time series auditor: ", err.Error())
					return
				} else {
					log.Printf("new time series auditor: %#v\n", tsaudit)
				}

				auditor = tsaudit
			default:
				log.Println("unknown audit backend type")
				return
			}
		default:
			log.Errorln("unknown auditor type")
			return
		}

		log.Println("starting irc bot")
		ircbot := chatbots.NewIRCBot(auditor)
		log.Printf("bot config: %#v\n", *ircbot)
		ircbot.Run()
	},
}
