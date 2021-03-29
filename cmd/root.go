package cmd

import (
	"fmt"
	"github.com/huyhvq/eurofxref/pkg/database"
	"github.com/huyhvq/eurofxref/pkg/handler"
	"github.com/huyhvq/eurofxref/pkg/repository"
	"github.com/huyhvq/eurofxref/pkg/server"
	"github.com/huyhvq/eurofxref/pkg/service/ecb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "eurofxref",
	Short: "Euro exchange rate API",
	Long:  "Euro exchange rate API",
	Run:   serve,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.eurofxref.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.Getwd()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("eurofxref")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	viper.SetEnvPrefix("app")
	for _, cfgName := range viper.AllKeys() {
		viper.BindEnv(cfgName)
	}
}

func serve(cmd *cobra.Command, args []string) {
	db, err := database.NewDB(database.MysqlCfg{
		Username: viper.GetString("db_user"),
		Password: viper.GetString("db_pass"),
		Host:     viper.GetString("db_host"),
		Port:     viper.GetString("db_port"),
		Name:     viper.GetString("db_name"),
		Driver:   viper.GetString("db_driver"),
	})

	if err != nil {
		panic(err)
	}
	defer db.Close()
	m, err := NewMigrate(db.DB())
	if err != nil {
		log.Println("initial database migrate failed...")
	}
	if m != nil {
		if err := m.Up(); err != nil {
			log.Println("database migrate failed...", err)
		} else {
			log.Println("database migrate successful...")
		}
	}

	r := repository.NewRate(db.DB())
	e := ecb.NewService(&ecb.Config{
		Endpoint: "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml",
	})
	s := server.NewHttpServer(handler.NewHandler(r), e)
	log.Println("initial service...")
	if err := s.Initial(r); err != nil {
		panic(err)
	}
	log.Println("initial service done")
	log.Println("starting service as port 8080...")
	if err := s.Start(); err != nil {
		log.Println("starting service failed, error:", err)
	}
}
