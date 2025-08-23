package seeder

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"gorm.io/gorm"
	"log"
)

func Seed(db *gorm.DB) error {

	// data master
	SeedObjectType(db, 100)

	fmt.Println("Seeding done")
	return nil
}

func SeedObjectType(db *gorm.DB, total int) {
	for i := 0; i < total; i++ {
		if err := db.Exec("TRUNCATE table object_types RESTART IDENTITY").Error; err != nil {
			panic(err)
		}
		for i := 0; i < total; i++ {
			data := &entities.ObjectType{
				Name: gofakeit.LoremIpsumWord(),
			}
			if err := db.Create(data).Error; err != nil {
				log.Printf("skipped entry %s: %s", data.Name, err)
				continue
			}
		}
		log.Println("seeding object type done")
	}
}
