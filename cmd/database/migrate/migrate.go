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
	doMigrate(db, "auctionBidders", &entities.AuctionBidder{})
	doMigrate(db, "bidderPayments", &entities.BidderPayment{})
	doMigrate(db, "bids", &entities.Bid{})

	log.Println("migration done")
	return nil
}

/**
2025/09/21 15:35:11 D:/farras/SMK_Taruna_Bhakti/Projek/Kelas_12_Tugas_Lelang_Online/Revisi_API_Lelang_Online/cmd/database/migrate/migrate.go:41 ERROR: insert or update on table "auction_bidders" violates foreign key constraint "fk_auction_bidders_user" (SQLSTATE 23503)
[3.234ms] [rows:0] ALTER TABLE "auction_bidders" ADD CONSTRAINT "fk_auction_bidders_user" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE
2025/09/21 15:35:11 failed to migrate table auctionBidders: ERROR: insert or update on table "auction_bidders" violates foreign key constraint "fk_auction_bidders_user" (SQLSTATE 23503)
2025/09/21 15:35:11 successfully migrate table auctionBidders

2025/09/21 15:35:11 D:/farras/SMK_Taruna_Bhakti/Projek/Kelas_12_Tugas_Lelang_Online/Revisi_API_Lelang_Online/cmd/database/migrate/migrate.go:41 ERROR: insert or update on table "auction_bidders" violates foreign key constraint "fk_auction_bidders_user" (SQLSTATE 23503)
[1.335ms] [rows:0] ALTER TABLE "auction_bidders" ADD CONSTRAINT "fk_auction_bidders_user" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE
*/

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
