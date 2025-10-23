package migrate

import (
	"fmt"
	"github.com/yayayapluto/revisi_api_lelang_online/cmd/database/fk"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) error {

	doMigrate(db, "object_types", &entities.ObjectType{}) // Y
	doMigrate(db, "organizers", &entities.Organizer{})    // Y
	doMigrate(db, "files", &entities.File{})              // Y
	doMigrate(db, "users", &entities.User{})              // Y

	doMigrate(db, "items", &entities.Item{})                    // Y
	doMigrate(db, "item_details", &entities.ItemDetail{})       // Y
	doMigrate(db, "item_documents", &entities.ItemDocument{})   // Y
	doMigrate(db, "item_grades", &entities.ItemGrade{})         // Y
	doMigrate(db, "item_thumbnails", &entities.ItemThumbnail{}) // Y
	doMigrate(db, "pics", &entities.PIC{})                      // Y
	doMigrate(db, "auctions", &entities.Auction{})              // Y
	doMigrate(db, "auctionBidders", &entities.AuctionBidder{})  // Y
	doMigrate(db, "bids", &entities.Bid{})
	doMigrate(db, "bidderPayments", &entities.BidderPayment{}) // Y

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
