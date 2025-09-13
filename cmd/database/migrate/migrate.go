package migrate

import (
	"fmt"
	"github.com/yayayapluto/revisi_api_lelang_online/cmd/database/fk"
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
	doMigrate(db, "item_thumbnails", &entities.ItemThumbnail{})
	doMigrate(db, "pics", &entities.PIC{})
	doMigrate(db, "auctions", &entities.Auction{})
	doMigrate(db, "users", &entities.User{})

	log.Println("migration done")
	return nil
}

func doMigrate(db *gorm.DB, tblName string, entity interface{}) {
	if err := fk.DisableFK(db, tblName); err != nil {
		log.Fatal("failed to disable fk")
	}

	if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tblName)).Error; err != nil {
		log.Printf("failed to dropping table %s: %s\n", tblName, err)
	}
	if err := db.AutoMigrate(entity); err != nil {
		log.Printf("failed to migrate table %s: %s\n", tblName, err)
	}

	if err := fk.EnableFK(db, tblName); err != nil {
		log.Fatal("failed to enable fk")
	}

	log.Printf("successfully migrate table %s\n", tblName)
}
