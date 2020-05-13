package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jtaylorcpp/quaditor"
	"github.com/jtaylorcpp/quaditor/auditors"
)

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.PersistentFlags().StringVarP(&queryFileName, "file", "f", "", "json file of quads")
	queryCmd.MarkPersistentFlagRequired("file")
}

var queryFileName string

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "query quads into backend",
	Run: func(cmd *cobra.Command, args []string) {
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
					return
				} else {
					log.Printf("new time series auditor: %#v\n", tsaudit)
				}

				jsonFile, err := ioutil.ReadFile(queryFileName)
				if err != nil {
					log.Println("error reading json query: ", err.Error())
					return
				}

				var query []quaditor.Query
				err = json.Unmarshal(jsonFile, &query)
				if err != nil {
					log.Println("error reading json query: ", err.Error())
					return
				}

				log.Println("running query")
				err = tsaudit.Query(query...)
				if err != nil {
					log.Println("error running query: ", err.Error())
				}
				log.Println("query ran")
			default:
				log.Println("unknown audit backend type")
			}
		default:
			log.Println("unknown auditor type")
		}
	},
}
