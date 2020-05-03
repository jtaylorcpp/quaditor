package main

import (
	//"log"

	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
	//"github.com/jtaylorcpp/quaditor/auditors"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the server",
	/*Run: func(cmd *cobra.Command, args []string) {
		log.Println("starting auditor of type: ", viper.GetString("audit_type"))
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
				} else {
					log.Printf("new time series auditor: %#v\n", tsaudit)
				}

			default:
				log.Println("unknown audit backend type")
			}
		default:
			log.Println("unknown auditor type")
		}
	},*/
}
