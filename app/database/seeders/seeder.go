package seeders

import (
	"live_attendance/main/app/database/fakers"

	"gorm.io/gorm"
)

type Seeder struct {
	Seeder interface{}
}

func RegisterSeeder(db *gorm.DB) []Seeder {
	return []Seeder{
		{Seeder: fakers.EmployeeFaker(db)},
	}
}

func DBSeed(db *gorm.DB) error {
	for _, seeder := range RegisterSeeder(db) {
		err := db.Debug().Create(seeder.Seeder).Error
		if err != nil {
			return err
		}
	}

	return nil
}
