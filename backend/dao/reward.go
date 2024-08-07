package dao

import (
	"blog/dao/db"
	"encoding/json"
)

func GetReward() *reward {
	return rewardDao
}

type reward struct {
	_ [0]func()
}

var rewardDao = &reward{}

func (r *reward) CreateReward(reward *Reward) (err error) {
	err = db.GetMysql().Model(&Reward{}).Create(reward).Error
	if err != nil {
		return
	}
	return
}

func (r *reward) FindRewardById(id uint) (reward Reward, err error) {
	err = db.GetMysql().Model(&Reward{}).Where("id = ?", id).First(&reward).Error
	return
}

func (r *reward) DeleteReward(id uint) (err error) {
	err = db.GetMysql().Model(&Reward{}).Where("id = ?", id).Delete(&Reward{}).Error
	return
}

func (r *reward) UpdateReward(reward *Reward) (err error) {
	err = db.GetMysql().Model(&Reward{}).Where("id = ?", reward.ID).Updates(reward).Error
	return
}

func (r *reward) FindRewardByArticleId(articleId uint) (rewards []Reward, err error) {
	err = db.GetMysql().Model(&Reward{}).Where("article_id = ?", articleId).Find(&rewards).Error
	return
}

func (r *reward) FindRewardByUserId(userId uint) (rewards []Reward, err error) {
	err = db.GetMysql().Model(&Reward{}).Where("user_address = ?", userId).Find(&rewards).Error
	return
}

func init() {
	db.GetMysql().AutoMigrate(&Reward{})
}

// should replace the origin cacheKey which should assign the value by user. then we pass the tag table name to assign the cache prefix
var _r = &Reward{}

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
