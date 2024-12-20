package model

import "time"

type Airport struct {
	ID               uint      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name             string    `gorm:"type:varchar(255);not null" json:"name"`
	StartTime        time.Time `gorm:"type:datetime" json:"start_time"`
	EndTime          time.Time `gorm:"type:datetime" json:"end_time"`
	FinalTime        time.Time `gorm:"type:datetime" json:"final_time"`
	Address          string    `gorm:"type:varchar(255)" json:"address"`
	Tag              string    `gorm:"type:varchar(255)" json:"tag"`
	FinancingBalance float64   `gorm:"type:float" json:"financing_balance"`
	FinancingFrom    string    `gorm:"type:varchar(255)" json:"financing_from"`
	TaskType         string    `gorm:"type:varchar(255)" json:"task_type"`
	AirportBalance   float64   `gorm:"type:float" json:"airport_balance"`
	Teaching         string    `gorm:"type:varchar(255)" json:"teaching"`
	Weight           uint      `gorm:"type:bigint" json:"weight"`
}

func (a *Airport) TableName() string {
	return "airport"
}
