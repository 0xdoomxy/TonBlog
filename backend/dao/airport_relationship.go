package dao

import (
	"blog/dao/db"
	"blog/model"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
)

func init() {
	err := db.GetMysql().AutoMigrate(&model.AirportRelationship{})
	if err != nil {
		logrus.Panicf("auto migrate airport_relationship table err:%s", err.Error())
	}

}

var airportRelationshipDao = newAirportRelationshipDao()

func GetAirportRelationship() *airportRelationship {
	return airportRelationshipDao
}

type airportRelationship struct {
	_  [0]func()
	sf singleflight.Group
}

func newAirportRelationshipDao() *airportRelationship {
	return &airportRelationship{
		sf: singleflight.Group{},
	}
}

func (ar *airportRelationship) UpdateAirportRelationship(ctx context.Context, data *model.AirportRelationship) (err error) {
	cache := db.GetRedis()
	err = cache.Del(ctx, fmt.Sprintf("airport_relationship_%d_%s", data.AirportId, data.UserAddress)).Err()
	if err != nil {
		logrus.Panicf("delete airport_relationship cache (%v) when update airport_relationship err:%s", data, err.Error())
		return
	}
	storage := db.GetMysql()
	err = storage.WithContext(ctx).UpdateColumns(data).Error
	if err != nil {
		logrus.Errorf("update airport_relationship (%v) err:%s", data, err.Error())
	}
	return
}

func (ar *airportRelationship) CreateAirportRelationship(ctx context.Context, data *model.AirportRelationship) (err error) {
	storage := db.GetMysql()
	err = storage.WithContext(ctx).Create(data).Error
	if err != nil {
		logrus.Errorf("create airport_relationship(%v) err:%s", data, err.Error())
		return
	}
	cache := db.GetRedis()
	ignoreErr := cache.SetEx(ctx, fmt.Sprintf("airport_relationship_%d_%s", data.AirportId, data.UserAddress), data, 12*time.Hour).Err()
	if ignoreErr != nil {
		logrus.Errorf("set airport_relationship(%v) err:%s", data, ignoreErr.Error())
	}
	return
}

func (ar *airportRelationship) DeleteAirportRelationship(ctx context.Context, data *model.AirportRelationship) (err error) {
	cache := db.GetRedis()
	err = cache.Del(ctx, fmt.Sprintf("airport_relationship_%d_%s", data.AirportId, data.UserAddress)).Err()
	if err != nil {
		logrus.Errorf("delete airport_relationship(%v) in redis err:%s", data, err.Error())
		return
	}
	storage := db.GetMysql()
	res := storage.WithContext(ctx).Model(&model.AirportRelationship{}).Where("airport_id = ? and  user_address = ? and delete_time is null", data.AirportId, data.UserAddress).Update("delete_time", time.Now())
	if res.RowsAffected <= 0 || res.Error == nil {
		err = res.Error
		if err == nil {
			err = fmt.Errorf("删除失败")
		}
		return err
	}
	return nil
}

func (ar *airportRelationship) FindAirportRelationshipByAddressAndId(ctx context.Context, address string, airportId uint) (res *model.AirportRelationship, err error) {
	cache := db.GetRedis()
	err = cache.Get(ctx, fmt.Sprintf("airport_relationship_%d_%s", airportId, address)).Scan(res)
	if !errors.Is(err, redis.Nil) {
		if err != nil {
			logrus.Errorf("find airport_relationship (airportId:%d,address:%s) by redis err:%s", airportId, address, err.Error())
		}
		return
	}
	var raw any
	raw, err, _ = ar.sf.Do(fmt.Sprintf("airport_relationship_%d_%s", airportId, address), func() (interface{}, error) {
		storage := db.GetMysql()
		var closureRes *model.AirportRelationship
		var closureErr error
		closureErr = storage.WithContext(ctx).Model(&model.AirportRelationship{}).Where("airport_id = ? and user_address = ? and delete_time is null", airportId, address).First(closureRes).Error
		return closureRes, closureErr
	})
	if err != nil {
		logrus.Errorf("find airport_relationship  (airportId:%d,address:%s) by mysql err:%s", airportId, address, err.Error())
		return
	}
	ignoreErr := cache.SetEx(ctx, fmt.Sprintf("airport_relationship_%d_%s", airportId, address), raw, 12*time.Hour).Err()
	if ignoreErr != nil {
		logrus.Errorf("set airport_relationship in redis (airportId:%d,address:%s)  error:%s", airportId, address, ignoreErr.Error())
	}
	return raw.(*model.AirportRelationship), nil
}
