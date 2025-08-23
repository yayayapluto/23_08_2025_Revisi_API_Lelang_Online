package migrate

import (
	"fmt"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) error {
	do_migrate(db, "object_types", &entities.ObjectType{})

	fmt.Println("migration done")
	return nil
}

func do_migrate(db *gorm.DB, tbname string, entity interface{}) {
	if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY", tbname)).Error; err != nil {
		log.Printf("failed to truncate table %s: %s\n", tbname, err)
	}
	if err := db.AutoMigrate(entity); err != nil {
		log.Printf("failed to migrate table %s: %s\n", tbname, err)
	}

	log.Printf("successfully migrate table %s\n", tbname)
}
