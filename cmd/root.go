package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/quaditor")
	viper.AddConfigPath("$HOME/.quaditor")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("quad")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			// noop
			log.Println("no config file found")
		} else {
			// Config file was found but another error was produced
			panic(err.Error())
		}
	}

	rootCmd.PersistentFlags().StringVarP(&auditType, "audit-type", "", "time-series", "type of auditor to use")
	viper.BindPFlag("audit_type", rootCmd.PersistentFlags().Lookup("audit-type"))
	rootCmd.PersistentFlags().StringVarP(&auditBackend, "audit-backend", "", "postgresql", "type of auditor to use")
	viper.BindPFlag("audit_backend", rootCmd.PersistentFlags().Lookup("audit-backend"))
	rootCmd.PersistentFlags().StringVarP(&auditUsername, "audit-username", "", "euler", "type of auditor to use")
	viper.BindPFlag("audit_username", rootCmd.PersistentFlags().Lookup("audit-username"))
	rootCmd.PersistentFlags().StringVarP(&auditPassword, "audit-password", "", "euler", "type of auditor to use")
	viper.BindPFlag("audit_password", rootCmd.PersistentFlags().Lookup("audit-password"))
	rootCmd.PersistentFlags().StringVarP(&auditHost, "audit-host", "", "localhost", "type of auditor to use")
	viper.BindPFlag("audit_host", rootCmd.PersistentFlags().Lookup("audit-host"))
	rootCmd.PersistentFlags().StringVarP(&auditPort, "audit-port", "", "5432", "type of auditor to use")
	viper.BindPFlag("audit_port", rootCmd.PersistentFlags().Lookup("audit-port"))
}

var (
	auditType     string
	auditBackend  string
	auditUsername string
	auditPassword string
	auditHost     string
	auditPort     string
)

var rootCmd = &cobra.Command{
	Use:   "quaditor",
	Short: "the quaditor",
}

func main() {
	log.Println(rootCmd.Execute())
}
