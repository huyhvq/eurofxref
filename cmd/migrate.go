package cmd

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/huyhvq/eurofxref/pkg/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate",
	Long:  `migrate.`,
	Run:   migrateExecute,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func migrateExecute(cmd *cobra.Command, args []string) {
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
		panic(err)
	}

	if len(args) > 0 && args[0] == "down" {
		if err := m.Down(); err != nil {
			log.Println(err)
		}
		return
	}

	if err := m.Up(); err != nil {
		log.Println(err)
	}
}

func NewMigrate(db *sql.DB) (*migrate.Migrate, error) {
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}
