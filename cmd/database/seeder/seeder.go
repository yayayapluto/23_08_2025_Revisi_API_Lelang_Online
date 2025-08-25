package seeder

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"gorm.io/gorm"
	"log"
)

func Seed(db *gorm.DB) error {
	if err := disableFK(db, "object_types"); err != nil {
		return err
	}
	if err := disableFK(db, "organizers"); err != nil {
		return err
	}
	if err := disableFK(db, "files"); err != nil {
		return err
	}
	if err := disableFK(db, "items"); err != nil {
		return err
	}
	if err := disableFK(db, "item_details"); err != nil {
		return err
	}
	if err := disableFK(db, "item_documents"); err != nil {
		return err
	}
	if err := disableFK(db, "item_grades"); err != nil {
		return err
	}
	if err := disableFK(db, "item_thumbnails"); err != nil {
		return err
	}

	// Data master
	SeedObjectType(db, 10)
	SeedOrganizer(db, 20)

	// Auction-related
	SeedFile(db, 30)
	SeedItem(db, 30)
	SeedItemDetail(db, 30)
	SeedItemDocument(db, 30)
	SeedItemGrade(db, 30)
	SeedItemThumbnail(db, 30)

	if err := enableFK(db, "object_types"); err != nil {
		return err
	}
	if err := enableFK(db, "organizers"); err != nil {
		return err
	}
	if err := enableFK(db, "files"); err != nil {
		return err
	}
	if err := enableFK(db, "items"); err != nil {
		return err
	}
	if err := enableFK(db, "item_details"); err != nil {
		return err
	}
	if err := enableFK(db, "item_documents"); err != nil {
		return err
	}
	if err := enableFK(db, "item_grades"); err != nil {
		return err
	}
	if err := enableFK(db, "item_thumbnails"); err != nil {
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
			ObjectTypeID: uint(gofakeit.Number(1, 10)),
			Name:         gofakeit.HipsterWord(),
			Price:        float64(gofakeit.Number(1000000, 100000000)),
			DepositPrice: float64(gofakeit.Number(100000, 10000000)),
			Description:  &desc,
			FileID:       uint(gofakeit.Number(1, 30)),
		}
		if err := db.Create(data).Error; err != nil {
			log.Printf("skipped entry %s: %s", data.Name, err)
			continue
		}
	}
	log.Println("seeding item done")
}

func SeedItemDetail(db *gorm.DB, total int) {
	if err := db.Exec("TRUNCATE TABLE item_details RESTART IDENTITY CASCADE").Error; err != nil {
		panic(err)
	}

	for i := 0; i < total; i++ {
		plate := gofakeit.LetterN(2) + gofakeit.Numerify("####") + gofakeit.LetterN(2)
		series := gofakeit.Word()
		cc := gofakeit.Float64Range(1000, 5000)
		itype := gofakeit.CarType()
		transmission := gofakeit.RandomString([]string{"Manual", "Automatic"})
		model := gofakeit.Word()
		frameNum := gofakeit.UUID()
		machineNum := gofakeit.UUID()
		km := gofakeit.Number(1000, 200000)
		fuel := gofakeit.RandomString([]string{"Petrol", "Diesel", "Electric", "Hybrid"})
		driveType := gofakeit.RandomString([]string{"FWD", "RWD", "AWD"})
		stnkDate := gofakeit.Date()

		data := &entities.ItemDetail{
			ItemID:        uint(gofakeit.Number(1, 30)),
			PlateNumber:   &plate,
			Brand:         gofakeit.CarMaker(),
			Series:        &series,
			CC:            &cc,
			Type:          &itype,
			Transmission:  &transmission,
			Model:         &model,
			Year:          gofakeit.Year(),
			FrameNumber:   &frameNum,
			MachineNumber: &machineNum,
			Kilometer:     &km,
			Fuel:          &fuel,
			Color:         gofakeit.Color(),
			DriveType:     &driveType,
			StnkDate:      &stnkDate,
		}

		if err := db.Create(data).Error; err != nil {
			log.Printf("skipped entry %d: %s", data.ItemID, err)
			continue
		}
	}
	log.Println("seeding item_detail done")
}

func SeedItemDocument(db *gorm.DB, total int) {
	if err := db.Exec("TRUNCATE TABLE item_documents RESTART IDENTITY CASCADE").Error; err != nil {
		panic(err)
	}
	for i := 0; i < total; i++ {

		Bpkb := gofakeit.Bool()
		Stnk := gofakeit.Bool()
		Facture := gofakeit.Bool()
		Receipt := gofakeit.Bool()
		OwnershipRelease := gofakeit.Bool()
		Warranty := gofakeit.Bool()
		Box := gofakeit.Bool()

		data := &entities.ItemDocument{
			ItemID:           uint(gofakeit.Number(1, 30)),
			Bpkb:             &Bpkb,
			Stnk:             &Stnk,
			Facture:          &Facture,
			Receipt:          &Receipt,
			OwnershipRelease: &OwnershipRelease,
			Warranty:         &Warranty,
			Box:              &Box,
		}
		if err := db.Create(data).Error; err != nil {
			log.Printf("skipped entry %d: %s", data.ItemID, err)
			continue
		}
	}
	log.Println("seeding item done")
}

func SeedItemGrade(db *gorm.DB, total int) {
	if err := db.Exec("TRUNCATE TABLE item_grades RESTART IDENTITY CASCADE").Error; err != nil {
		panic(err)
	}
	for i := 0; i < total; i++ {
		data := &entities.ItemGrade{
			ItemID:   uint(gofakeit.Number(1, 30)),
			Interior: gofakeit.RandomString([]string{"a", "b", "c", "d", "e", "f"}),
			Exterior: gofakeit.RandomString([]string{"a", "b", "c", "d", "e", "f"}),
			Frame:    gofakeit.RandomString([]string{"a", "b", "c", "d", "e", "f"}),
			Machine:  gofakeit.RandomString([]string{"a", "b", "c", "d", "e", "f"}),
		}
		if err := db.Create(data).Error; err != nil {
			log.Printf("skipped entry %d: %s", data.ItemID, err)
			continue
		}
	}
	log.Println("seeding item done")
}

func SeedItemThumbnail(db *gorm.DB, total int) {
	if err := db.Exec("TRUNCATE TABLE item_thumbnails RESTART IDENTITY CASCADE").Error; err != nil {
		panic(err)
	}
	for i := 0; i < total; i++ {
		data := &entities.ItemThumbnail{
			Name:   gofakeit.BuzzWord(),
			ItemID: uint(gofakeit.Number(1, 30)),
			FileID: uint(gofakeit.Number(1, 30)),
		}
		if err := db.Create(data).Error; err != nil {
			log.Printf("skipped entry %d: %s", data.ItemID, err)
			continue
		}
	}
	log.Println("seeding item_thumbnail done")
}
