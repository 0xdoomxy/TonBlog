package model

type Token struct {
	AirportId   uint
	ChainId     uint    `gorm:"chain_id"`
	TokenName   string  `gorm:"token_name;type:varchar(255)"`
	Address     string  `gorm:"address;type:varchar(255)"`
	PriceOracle float64 `gorm:"price_oracle;type:float"`
}

func (t *Token) TableName() string {
	return "token"
}
