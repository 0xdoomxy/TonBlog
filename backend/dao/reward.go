package dao

import (
	"blog/dao/db"
	"blog/model"
)

func GetReward() *reward {
	return rewardDao
}

type reward struct {
	_ [0]func()
}

var rewardDao = &reward{}

func (r *reward) CreateReward(reward *model.Reward) (err error) {
	err = db.GetMysql().Model(&model.Reward{}).Create(reward).Error
	if err != nil {
		return
	}
	return
}

func (r *reward) FindRewardById(id uint) (reward model.Reward, err error) {
	err = db.GetMysql().Model(&model.Reward{}).Where("id = ?", id).First(&reward).Error
	return
}

func (r *reward) DeleteReward(id uint) (err error) {
	err = db.GetMysql().Model(&model.Reward{}).Where("id = ?", id).Delete(&model.Reward{}).Error
	return
}

func (r *reward) UpdateReward(reward *model.Reward) (err error) {
	err = db.GetMysql().Model(&model.Reward{}).Where("id = ?", reward.ID).Updates(reward).Error
	return
}

func (r *reward) FindRewardByArticleId(articleId uint) (rewards []model.Reward, err error) {
	err = db.GetMysql().Model(&model.Reward{}).Where("article_id = ?", articleId).Find(&rewards).Error
	return
}

func (r *reward) FindRewardByUserId(userId uint) (rewards []model.Reward, err error) {
	err = db.GetMysql().Model(&model.Reward{}).Where("user_address = ?", userId).Find(&rewards).Error
	return
}

func init() {
	db.GetMysql().AutoMigrate(&model.Reward{})
}

// should replace the origin cacheKey which should assign the value by user. then we pass the tag table name to assign the cache prefix
var _r = &model.Reward{}
