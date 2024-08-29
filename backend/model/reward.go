package model

import "encoding/json"

/*打赏表*/
type Reward struct {
	ID          uint   `gorm:"primaryKey"`
	CreateAt    uint64 `gorm:"autoCreateTime:milli"`
	ArticleID   uint   `gorm:"not null;index:searchforarticle"`
	UserAddress string `gorm:"varchar(64);not null;index:searchforuser"`
	Amount      uint   `gorm:"not null"`
}

func (reward *Reward) TableName() string {
	return "reward"
}

func (reward *Reward) MarshalBinary() ([]byte, error) {
	return json.Marshal(reward)
}

func (reward *Reward) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, reward)
}
