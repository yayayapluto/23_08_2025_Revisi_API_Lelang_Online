package seeder

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/yayayapluto/revisi_api_lelang_online/cmd/database/fk"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"gorm.io/gorm"
	"log"
)

func Seed(db *gorm.DB) error {
	tables := []string{
		"object_types", "organizers", "files",
		"items", "item_details", "item_documents",
		"item_grades", "item_thumbnails", "pics", "auctions",
	}

	if err := toggleFK(db, tables, false); err != nil {
		return err
	}

	truncateTables(db, tables)

	SeedObjectType(db, 10)
	SeedOrganizer(db, 20)
	SeedFile(db, 30)
	SeedItem(db, 30)
	SeedItemDetail(db, 30)
	SeedItemDocument(db, 30)
	SeedItemGrade(db, 30)
	SeedItemThumbnail(db, 30)
	SeedPIC(db, 10)
	SeedAuction(db, 100)

	if err := toggleFK(db, tables, true); err != nil {
		return err
	}

	log.Println("Seeding done")
	return nil
}

func toggleFK(db *gorm.DB, tables []string, enable bool) error {
	for _, t := range tables {
		var err error
		if enable {
			err = fk.EnableFK(db, t)
		} else {
			err = fk.DisableFK(db, t)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func truncateTables(db *gorm.DB, tables []string) {
	for _, t := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", t)).Error; err != nil {
			panic(err)
		}
	}
}

func SeedObjectType(db *gorm.DB, total int) {
	exists := map[string]bool{}
	for i := 0; i < total; {
		name := gofakeit.HipsterWord()
		if exists[name] {
			continue
		}
		data := &entities.ObjectType{Name: name}
		if err := db.Create(data).Error; err == nil {
			exists[name] = true
			i++
		}
	}
	log.Println("seeding object type done")
}

func SeedOrganizer(db *gorm.DB, total int) {
	exists := map[string]bool{}
	for i := 0; i < total; {
		name := gofakeit.HipsterWord()
		if exists[name] {
			continue
		}
		data := &entities.Organizer{
			Name:          name,
			Address:       gofakeit.Address().Address,
			BankName:      gofakeit.BankName(),
			AccountNumber: gofakeit.Numerify("############"),
			AccountName:   gofakeit.BuzzWord(),
		}
		if err := db.Create(data).Error; err == nil {
			exists[name] = true
			i++
		}
	}
	log.Println("seeding organizer done")
}

func SeedFile(db *gorm.DB, total int) {
	exists := map[string]bool{}
	for i := 0; i < total; {
		path := gofakeit.URL()
		if exists[path] {
			continue
		}
		data := &entities.File{Path: path}
		if err := db.Create(data).Error; err == nil {
			exists[path] = true
			i++
		}
	}
	log.Println("seeding file done")
}

func SeedItem(db *gorm.DB, total int) {
	exists := map[string]bool{}
	for i := 0; i < total; {
		name := gofakeit.HipsterWord()
		if exists[name] {
			continue
		}
		desc := gofakeit.SentenceSimple()
		data := &entities.Item{
			ObjectTypeID: uint(gofakeit.Number(1, 10)),
			Name:         name,
			Price:        float64(gofakeit.Number(1000000, 100000000)),
			DepositPrice: float64(gofakeit.Number(100000, 10000000)),
			Description:  &desc,
			FileID:       uint(gofakeit.Number(1, 30)),
		}
		if err := db.Create(data).Error; err == nil {
			exists[name] = true
			i++
		}
	}
	log.Println("seeding item done")
}

func SeedItemDetail(db *gorm.DB, total int) {
	exists := map[string]bool{}
	for i := 0; i < total; {
		itemID := uint(gofakeit.Number(1, 30))
		if exists[fmt.Sprint(itemID)] {
			continue
		}
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
			ItemID:        itemID,
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
		if err := db.Create(data).Error; err == nil {
			exists[fmt.Sprint(itemID)] = true
			i++
		}
	}
	log.Println("seeding item_detail done")
}

func SeedItemDocument(db *gorm.DB, total int) {
	exists := map[string]bool{}
	for i := 0; i < total; {
		itemID := uint(gofakeit.Number(1, 30))
		if exists[fmt.Sprint(itemID)] {
			continue
		}

		Bpkb := gofakeit.Bool()
		Stnk := gofakeit.Bool()
		Facture := gofakeit.Bool()
		Receipt := gofakeit.Bool()
		OwnershipRelease := gofakeit.Bool()
		Warranty := gofakeit.Bool()
		Box := gofakeit.Bool()

		data := &entities.ItemDocument{
			ItemID:           itemID,
			Bpkb:             &Bpkb,
			Stnk:             &Stnk,
			Facture:          &Facture,
			Receipt:          &Receipt,
			OwnershipRelease: &OwnershipRelease,
			Warranty:         &Warranty,
			Box:              &Box,
		}
		if err := db.Create(data).Error; err == nil {
			exists[fmt.Sprint(itemID)] = true
			i++
		}
	}
	log.Println("seeding item_document done")
}

func SeedItemGrade(db *gorm.DB, total int) {
	exists := map[string]bool{}
	for i := 0; i < total; {
		itemID := uint(gofakeit.Number(1, 30))
		if exists[fmt.Sprint(itemID)] {
			continue
		}
		data := &entities.ItemGrade{
			ItemID:   itemID,
			Interior: gofakeit.RandomString([]string{"a", "b", "c", "d", "e", "f"}),
			Exterior: gofakeit.RandomString([]string{"a", "b", "c", "d", "e", "f"}),
			Frame:    gofakeit.RandomString([]string{"a", "b", "c", "d", "e", "f"}),
			Machine:  gofakeit.RandomString([]string{"a", "b", "c", "d", "e", "f"}),
		}
		if err := db.Create(data).Error; err == nil {
			exists[fmt.Sprint(itemID)] = true
			i++
		}
	}
	log.Println("seeding item_grade done")
}

func SeedItemThumbnail(db *gorm.DB, total int) {
	exists := map[string]bool{}
	for i := 0; i < total; {
		name := gofakeit.BuzzWord()
		if exists[name] {
			continue
		}
		data := &entities.ItemThumbnail{
			Name:   name,
			ItemID: uint(gofakeit.Number(1, 30)),
			FileID: uint(gofakeit.Number(1, 30)),
		}
		if err := db.Create(data).Error; err == nil {
			exists[name] = true
			i++
		}
	}
	log.Println("seeding item_thumbnail done")
}

func SeedPIC(db *gorm.DB, total int) {
	exists := map[string]bool{}
	for i := 0; i < total; {
		name := gofakeit.BuzzWord()
		if exists[name] {
			continue
		}
		data := &entities.PIC{
			Name:        name,
			PhoneNumber: gofakeit.Phone(),
		}
		if err := db.Create(data).Error; err == nil {
			exists[name] = true
			i++
		}
	}
	log.Println("seeding pic done")
}

func SeedAuction(db *gorm.DB, total int) {
	exists := map[string]bool{}
	for i := 0; i < total; {
		id := fmt.Sprintf("%d-%d-%d", gofakeit.Number(1, 30), gofakeit.Number(1, 30), gofakeit.Number(1, 30))
		if exists[id] {
			continue
		}
		data := &entities.Auction{
			ItemID:      uint(gofakeit.Number(1, 30)),
			OrganizerID: uint(gofakeit.Number(1, 30)),
			PicID:       uint(gofakeit.Number(1, 30)),
			StartDate:   gofakeit.Date(),
			EndDate:     gofakeit.Date(),
		}
		if err := db.Create(data).Error; err == nil {
			exists[id] = true
			i++
		}
	}
	log.Println("seeding auction done")
}
