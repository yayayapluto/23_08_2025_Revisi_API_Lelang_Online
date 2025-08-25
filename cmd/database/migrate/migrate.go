package migrate

import (
	"fmt"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) error {
	doMigrate(db, "object_types", &entities.ObjectType{})
	doMigrate(db, "organizers", &entities.Organizer{})
	doMigrate(db, "files", &entities.File{})
	doMigrate(db, "items", &entities.Item{})
	doMigrate(db, "item_details", &entities.ItemDetail{})
	doMigrate(db, "item_documents", &entities.ItemDocument{})
	doMigrate(db, "item_grades", &entities.ItemGrade{})

	fmt.Println("migration done")
	return nil
}

func doMigrate(db *gorm.DB, tbname string, entity interface{}) {
	if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY", tbname)).Error; err != nil {
		log.Printf("failed to truncate table %s: %s\n", tbname, err)
	}
	if err := db.AutoMigrate(entity); err != nil {
		log.Printf("failed to migrate table %s: %s\n", tbname, err)
	}

	log.Printf("successfully migrate table %s\n", tbname)
}
