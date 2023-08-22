package models

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	ID        string `gorm:"size:36";not null;uniqueIndex;primary_key"`
	Name      string `gorm:"size:255;not null"`
	Email     string `gorm:"size:255;not null"`
	Password  string `gorm:"size:255;not null"`
	Image     string `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (e *Employee) FindByEmail(db *gorm.DB, email string) (*Employee, error) {
	var err error
	var employee Employee

	err = db.Debug().Model(Employee{}).Where("email = ?", email).First(&employee).Error

	if err != nil {
		return nil, err
	}

	return &employee, nil
}

func (e *Employee) FindByID(db *gorm.DB, id string) (*Employee, error) {
	var err error
	var employee Employee

	err = db.Debug().Model(Employee{}).Where("id = ?", id).First(&employee).Error

	if err != nil {
		return nil, err
	}

	return &employee, nil
}
