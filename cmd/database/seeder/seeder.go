package seeder

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"gorm.io/gorm"
	"log"
)

func Seed(db *gorm.DB) error {
	// Disable foreign key check untuk tiap tabel
	if err := disableFK(db, "items"); err != nil {
		return err
	}
	if err := disableFK(db, "object_types"); err != nil {
		return err
	}
	if err := disableFK(db, "organizers"); err != nil {
		return err
	}
	if err := disableFK(db, "files"); err != nil {
		return err
	}

	// Data master
	SeedObjectType(db, 25)
	SeedOrganizer(db, 50)

	// Auction-related
	SeedFile(db, 100)
	SeedItem(db, 100)

	// Enable foreign key check lagi
	if err := enableFK(db, "items"); err != nil {
		return err
	}
	if err := enableFK(db, "object_types"); err != nil {
		return err
	}
	if err := enableFK(db, "organizers"); err != nil {
		return err
	}
	if err := enableFK(db, "files"); err != nil {
		return err
	}

	fmt.Println("Seeding done")
	return nil
}

func disableFK(db *gorm.DB, table string) error {
	return db.Exec(fmt.Sprintf(`ALTER TABLE %s DISABLE TRIGGER ALL;`, table)).Error
}

func enableFK(db *gorm.DB, table string) error {
	return db.Exec(fmt.Sprintf(`ALTER TABLE %s ENABLE TRIGGER ALL;`, table)).Error
}

func SeedObjectType(db *gorm.DB, total int) {
	if err := db.Exec("TRUNCATE TABLE object_types RESTART IDENTITY CASCADE").Error; err != nil {
		panic(err)
	}
	for i := 0; i < total; i++ {
		data := &entities.ObjectType{
			Name: gofakeit.HipsterWord(),
		}
		if err := db.Create(data).Error; err != nil {
			log.Printf("skipped entry %s: %s", data.Name, err)
			continue
		}
	}
	log.Println("seeding object type done")
}

func SeedOrganizer(db *gorm.DB, total int) {
	if err := db.Exec("TRUNCATE TABLE organizers RESTART IDENTITY CASCADE").Error; err != nil {
		panic(err)
	}
	for i := 0; i < total; i++ {
		data := &entities.Organizer{
			Name:          gofakeit.HipsterWord(),
			Address:       gofakeit.Address().Address,
			BankName:      gofakeit.BankName(),
			AccountNumber: gofakeit.Numerify("############"),
			AccountName:   gofakeit.BuzzWord(),
		}
		if err := db.Create(data).Error; err != nil {
			log.Printf("skipped entry %s: %s", data.Name, err)
			continue
		}
	}
	log.Println("seeding organizer done")
}

func SeedFile(db *gorm.DB, total int) {
	if err := db.Exec("TRUNCATE TABLE files RESTART IDENTITY CASCADE").Error; err != nil {
		panic(err)
	}
	for i := 0; i < total; i++ {
		data := &entities.File{
			Path: gofakeit.URL(),
		}
		if err := db.Create(data).Error; err != nil {
			log.Printf("skipped entry %s: %s", data.Path, err)
			continue
		}
	}
	log.Println("seeding file done")
}

func SeedItem(db *gorm.DB, total int) {
	if err := db.Exec("TRUNCATE TABLE items RESTART IDENTITY CASCADE").Error; err != nil {
		panic(err)
	}
	for i := 0; i < total; i++ {
		desc := gofakeit.SentenceSimple()
		data := &entities.Item{
			ObjectTypeID: uint(gofakeit.Number(1, 25)),
			Name:         gofakeit.HipsterWord(),
			Price:        float64(gofakeit.Number(1000000, 100000000)),
			DepositPrice: float64(gofakeit.Number(100000, 10000000)),
			Description:  &desc,
			FileID:       uint(gofakeit.Number(1, 10)),
		}
		if err := db.Create(data).Error; err != nil {
			log.Printf("skipped entry %s: %s", data.Name, err)
			continue
		}
	}
	log.Println("seeding item done")
}
