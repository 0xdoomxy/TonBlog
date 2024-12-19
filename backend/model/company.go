package model

type Company struct {
	Name string `gorm:"type:varchar(255);primary_key"`
}

func (c *Company) TableName() string {
	return "company"
}
