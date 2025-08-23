package migrate

import (
	"fmt"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {

	fmt.Println("Migration done")
	return nil
}
