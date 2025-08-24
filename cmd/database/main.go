package main

import (
	"flag"
	"github.com/yayayapluto/revisi_api_lelang_online/cmd/config"
	"github.com/yayayapluto/revisi_api_lelang_online/cmd/database/migrate"
	"github.com/yayayapluto/revisi_api_lelang_online/cmd/database/seeder"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"gorm.io/gorm"
	"log"
)

func DatabaseSetUp() (*gorm.DB, error) {
	env, err := utils.LoadEnv()
	if err != nil {
		return nil, err
	}

	db, err := config.ConnectDB(env.DBHOST, env.DBUSER, env.DBPASSWORD, env.DBNAME, env.DBPORT)
	if db == nil || err != nil {
		return nil, err
	}

	migrateFlag := flag.Bool("migrate", false, "migrating the database")
	seedFlag := flag.Bool("seed", false, "seeding the database")

	flag.Parse()

	if *migrateFlag {
		if err := migrate.Migrate(db); err != nil {
			return nil, err
		}
	}
	if *seedFlag {
		if err := seeder.Seed(db); err != nil {
			return nil, err
		}
	}
	return db, nil
}

func main() {
	_, err := DatabaseSetUp()
	if err != nil {
		log.Fatalf("Error setting up database : %v", err)
	}
}
