package migration

import (
	"gorm.io/gorm"
	"tele-sticker-finder/model"
)

var models = []interface{}{
	&model.Sticker{},
}

func AutoMigration(db *gorm.DB) {
	db.Set("gorm:table_options", "").AutoMigrate(models...)
}
