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

	SeedObjectType(db, 5)
	SeedOrganizer(db, 20)
	SeedFile(db, 50)
	SeedItem(db, 50)
	SeedItemDetail(db, 50)
	SeedItemDocument(db, 50)
	SeedItemGrade(db, 50)
	SeedItemThumbnail(db, 50)
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
		db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", t))
	}
}

func SeedObjectType(db *gorm.DB, total int) {
	var data []entities.ObjectType
	for i := 0; i < total; i++ {
		name := fmt.Sprintf("Type-%d-%s", i, gofakeit.Word())
		data = append(data, entities.ObjectType{Name: name})
	}
	db.CreateInBatches(&data, 10)
	log.Println("seeding object type done")
}

func SeedOrganizer(db *gorm.DB, total int) {
	var data []entities.Organizer
	for i := 0; i < total; i++ {
		addr := gofakeit.Address()
		data = append(data, entities.Organizer{
			Name:          fmt.Sprintf("Org-%d-%s", i, gofakeit.Company()),
			Address:       addr.Address,
			BankName:      gofakeit.BankName(),
			AccountNumber: gofakeit.Numerify("############"),
			AccountName:   gofakeit.Name(),
		})
	}
	db.CreateInBatches(&data, 10)
	log.Println("seeding organizer done")
}

func SeedFile(db *gorm.DB, total int) {
	var data []entities.File
	for i := 0; i < total; i++ {
		path := fmt.Sprintf("https://placehold.co/%dx%d", gofakeit.Number(300, 800), gofakeit.Number(300, 800))
		data = append(data, entities.File{Path: path})
	}
	db.CreateInBatches(&data, 10)
	log.Println("seeding file done")
}

func SeedItem(db *gorm.DB, total int) {
	var items []entities.Item
	for i := 0; i < total; i++ {
		desc := gofakeit.Sentence(5)
		items = append(items, entities.Item{
			ObjectTypeID: uint(gofakeit.Number(1, 5)),
			Name:         fmt.Sprintf("Item-%d-%s", i, gofakeit.Product().Suffix),
			Price:        float64(gofakeit.Number(1e6, 1e8)),
			DepositPrice: float64(gofakeit.Number(1e5, 1e7)),
			Description:  &desc,
			FileID:       uint(gofakeit.Number(1, 50)),
		})
	}
	db.CreateInBatches(&items, 10)
	log.Println("seeding item done")
}

func SeedItemDetail(db *gorm.DB, total int) {
	var details []entities.ItemDetail
	for i := 0; i < total; i++ {
		itemID := uint(i%50 + 1)
		plate := gofakeit.Regex(`[A-Z]{1,2}\d{1,4}[A-Z]{2}`)
		series := gofakeit.CarModel()
		cc := gofakeit.Float64Range(1000, 5000)
		itype := gofakeit.CarType()
		transmission := gofakeit.RandomString([]string{"Manual", "Automatic"})
		model := gofakeit.CarModel()
		frameNum := gofakeit.Regex(`[A-HJ-NPR-Z0-9]{17}`)
		machineNum := gofakeit.Regex(`[A-HJ-NPR-Z0-9]{17}`)
		km := gofakeit.Number(1000, 200000)
		fuel := gofakeit.RandomString([]string{"Petrol", "Diesel", "Electric", "Hybrid"})
		driveType := gofakeit.RandomString([]string{"FWD", "RWD", "AWD"})
		stnkDate := gofakeit.Date()

		details = append(details, entities.ItemDetail{
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
		})
	}
	db.CreateInBatches(&details, 10)
	log.Println("seeding item_detail done")
}

func SeedItemDocument(db *gorm.DB, total int) {
	var docs []entities.ItemDocument
	for i := 0; i < total; i++ {
		itemID := uint(i%50 + 1)
		Bpkb := gofakeit.Bool()
		Stnk := gofakeit.Bool()
		Facture := gofakeit.Bool()
		Receipt := gofakeit.Bool()
		OwnershipRelease := gofakeit.Bool()
		Warranty := gofakeit.Bool()
		Box := gofakeit.Bool()

		docs = append(docs, entities.ItemDocument{
			ItemID:           itemID,
			Bpkb:             &Bpkb,
			Stnk:             &Stnk,
			Facture:          &Facture,
			Receipt:          &Receipt,
			OwnershipRelease: &OwnershipRelease,
			Warranty:         &Warranty,
			Box:              &Box,
		})
	}
	db.CreateInBatches(&docs, 10)
	log.Println("seeding item_document done")
}

func SeedItemGrade(db *gorm.DB, total int) {
	var grades []entities.ItemGrade
	for i := 0; i < total; i++ {
		itemID := uint(i%50 + 1)
		grades = append(grades, entities.ItemGrade{
			ItemID:   itemID,
			Interior: gofakeit.RandomString([]string{"a", "b", "c", "d", "e", "f"}),
			Exterior: gofakeit.RandomString([]string{"a", "b", "c", "d", "e", "f"}),
			Frame:    gofakeit.RandomString([]string{"a", "b", "c", "d", "e", "f"}),
			Machine:  gofakeit.RandomString([]string{"a", "b", "c", "d", "e", "f"}),
		})
	}
	db.CreateInBatches(&grades, 10)
	log.Println("seeding item_grade done")
}

func SeedItemThumbnail(db *gorm.DB, total int) {
	var thumbs []entities.ItemThumbnail
	for i := 0; i < total; i++ {
		thumbs = append(thumbs, entities.ItemThumbnail{
			Name:   fmt.Sprintf("Thumb-%d-%s", i, gofakeit.Word()),
			ItemID: uint(gofakeit.Number(1, 50)),
			FileID: uint(gofakeit.Number(1, 50)),
		})
	}
	db.CreateInBatches(&thumbs, 10)
	log.Println("seeding item_thumbnail done")
}

func SeedPIC(db *gorm.DB, total int) {
	var pics []entities.PIC
	for i := 0; i < total; i++ {
		pics = append(pics, entities.PIC{
			Name:        fmt.Sprintf("PIC-%d-%s", i, gofakeit.Name()),
			PhoneNumber: gofakeit.Phone(),
		})
	}
	db.CreateInBatches(&pics, 10)
	log.Println("seeding pic done")
}

func SeedAuction(db *gorm.DB, total int) {
	var auctions []entities.Auction
	for i := 0; i < total; i++ {
		start := gofakeit.Date()
		end := gofakeit.DateRange(start, start.AddDate(0, 0, 14))
		auctions = append(auctions, entities.Auction{
			ItemID:      uint(gofakeit.Number(1, 50)),
			OrganizerID: uint(gofakeit.Number(1, 20)),
			PicID:       uint(gofakeit.Number(1, 10)),
			StartDate:   start,
			EndDate:     end,
		})
	}
	db.CreateInBatches(&auctions, 10)
	log.Println("seeding auction done")
}
