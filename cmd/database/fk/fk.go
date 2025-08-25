package fk

import (
	"fmt"
	"gorm.io/gorm"
)

func DisableFK(db *gorm.DB, table string) error {
	return db.Exec(fmt.Sprintf(`ALTER TABLE IF EXISTS %s DISABLE TRIGGER ALL;`, table)).Error
}

func EnableFK(db *gorm.DB, table string) error {
	return db.Exec(fmt.Sprintf(`ALTER TABLE IF EXISTS %s ENABLE TRIGGER ALL;`, table)).Error
}
