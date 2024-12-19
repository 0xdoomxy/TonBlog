package model

import "time"

type Airport struct {
	ID               uint      `gorm:"primary_key;AUTO_INCREMENT"`
	Name             string    `gorm:"type:varchar(255);not null"`
	StartTime        time.Time `gorm:"type:datetime"`
	EndTime          time.Time `gorm:"type:datetime"`
	FinalTime        time.Time `gorm:"type:datetime"`
	Address          string    `gorm:"type:varchar(255)"`
	Tag              string    `gorm:"type:varchar(255);not null"`
	FinancingBalance float64   `gorm:"type:float"`
	FinancingFrom    string    `gorm:"type:varchar(255)"`
	TaskType         string    `gorm:"type:varchar(255);not null"`
	AirportBalance   float64   `gorm:"type:float"`
	Teaching         string    `gorm:"type:varchar(255)"`
}

func (a *Airport) TableName() string {
	return "airport"
}
