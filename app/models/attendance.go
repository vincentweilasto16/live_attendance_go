package models

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	ID            string `gorm:"size:36";not null;uniqueIndex;primary_key"`
	Employee      Employee
	EmployeeID    string `gorm:"size:36;index"`
	ClockInImage  string `gorm:"size:255;not null"`
	ClockOutImage string `gorm:"size:255;not null"`
	LastClockIn   time.Time
	LastClockOut  time.Time
}

func (a *Attendance) FindByIDAndClockIn(db *gorm.DB, employeeID string) (*Attendance, error) {
	var err error
	var attendance Attendance

	currentDate := time.Now().Format("2006-01-02")

	err = db.Debug().Model(Attendance{}).Where("employee_id = ? AND DATE(last_clock_in) = ?", employeeID, currentDate).First(&attendance).Error

	if err != nil {
		return nil, err
	}

	return &attendance, nil
}

func (a *Attendance) FindByIDAndClockOut(db *gorm.DB, employeeID string) (*Attendance, error) {
	var err error
	var attendance Attendance

	currentDate := time.Now().Format("2006-01-02")

	err = db.Debug().Model(Attendance{}).Where("employee_id = ? AND DATE(last_clock_out) = ?", employeeID, currentDate).First(&attendance).Error

	if err != nil {
		return nil, err
	}

	return &attendance, nil
}
