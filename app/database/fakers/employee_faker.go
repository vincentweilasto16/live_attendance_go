package fakers

import (
	"time"

	"live_attendance/main/app/models"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func EmployeeFaker(db *gorm.DB) *models.Employee {
	return &models.Employee{
		ID:        uuid.New().String(),
		Name:      faker.Name(),
		Email:     faker.Email(),
		Password:  faker.Password(),
		Image:     faker.URL(),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}
}
