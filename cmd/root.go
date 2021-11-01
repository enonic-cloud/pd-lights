package cmd

import (
	"context"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "pd-lights",
	Short: "A small app that controls the office traffic lights",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Starting ...")

		if viper.GetString("token") == "" {
			log.Fatalf("pagerduty token not set")
		}

		if viper.GetString("ip") == "" {
			log.Fatalf("npc ip not set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		log.Infof("Rolling throuh all the lights")
		if err := SetLights(ctx, On, On, On); err != nil {
			log.Fatalf("Failed to turn off all lights")
		}
		if err := SetLights(ctx, On, Off, Off); err != nil {
			log.Fatalf("Failed to turn on red light")
		}
		if err := SetLights(ctx, Off, On, Off); err != nil {
			log.Fatalf("Failed to turn on yellow light")
		}
		if err := SetLights(ctx, Off, Off, On); err != nil {
			log.Fatalf("Failed to turn on green light")
		}
		log.Printf("Lights working, creating pagerduty client")
		cancel()

		client := pagerduty.NewClient(viper.GetString("token"))
		for {
			if err := checkIncidents(client); err != nil {
				log.Errorf("Incident check loop failed: %s", err)
			} else {
				log.Infof("Incident check loop successful")
			}
			time.Sleep(viper.GetDuration("loop"))
		}
	},
}

func checkIncidents(client *pagerduty.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
	defer cancel()

	res, err := client.ListIncidentsWithContext(ctx, pagerduty.ListIncidentsOptions{
		Since: time.Now().Add(-time.Hour * 12).Format("2006-01-02"),
	})
	if err != nil {
		log.Errorf("Failed calling PD api: %s", err)
		return SetLights(ctx, On, Off, On)
	}
	for _, i := range res.Incidents {
		switch i.Status {
		case "triggered":
			return SetLights(ctx, On, Off, Off)
		case "acknowledged":
			return SetLights(ctx, Off, On, Off)
		case "resolved":
			// Do nothing
		default:
			log.Warnf("Unknown status: %s", i.Status)
			return SetLights(ctx, On, On, Off)
		}
	}
	return SetLights(ctx, Off, Off, On)
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().String("token", "", "pagerduty token")
	_ = viper.BindPFlag("token", rootCmd.Flags().Lookup("token"))

	rootCmd.Flags().String("ip", "", "npc ip")
	_ = viper.BindPFlag("ip", rootCmd.Flags().Lookup("ip"))

	rootCmd.Flags().Duration("loop", time.Second*30, "loop interval")
	_ = viper.BindPFlag("loop", rootCmd.Flags().Lookup("loop"))

	rootCmd.Flags().Duration("timeout", time.Second*15, "request timeout")
	_ = viper.BindPFlag("timeout", rootCmd.Flags().Lookup("timeout"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".pd-lights")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}
}
