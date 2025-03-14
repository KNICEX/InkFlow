package dao

import "gorm.io/gorm"

func InitTable(db *gorm.DB) error {
	if err := db.AutoMigrate(&DraftInk{}, &LiveInk{}); err != nil {
		return err
	}
	return nil
}
