package dao

import (
	"blog/dao/db"
	"blog/model"
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
)

func init() {
	err := db.GetMysql().AutoMigrate(&model.Airport{})
	if err != nil {
		logrus.Panicf("auto migrate airport table err:%s", err.Error())
	}

}

var airportDao = newAirportDao()

func GetAirport() *airport {
	return airportDao
}

type airport struct {
	_  [0]func()
	sf singleflight.Group
}

func newAirportDao() *airport {
	return &airport{
		sf: singleflight.Group{},
	}
}

func (a *airport) CreateAirport(ctx context.Context, data *model.Airport) (err error) {
	storage := db.GetMysql()
	err = storage.WithContext(ctx).Model(&model.Airport{}).Create(data).Error
	if err != nil {
		logrus.Errorf("create airport %v error:%s", data, err.Error())
	}
	return
}

func (a *airport) QueryRunningAirportWithWeightByPage(ctx context.Context, address string, page int, pageSize int) (res []*model.Airport, err error) {
	storage := db.GetMysql()
	var raw any
	raw, err, _ = a.sf.Do(fmt.Sprintf("running_%s_%d_%d", address, page, pageSize), func() (interface{}, error) {
		var res []*model.Airport
		var raw *sql.Rows
		var err error
		raw, err = storage.WithContext(ctx).Raw(`SELECT
		a.id AS id,
		a.name AS name,
		a.start_time AS start_time,
		a.end_time AS end_time,
		a.final_time AS final_time,
		a.address AS address,
		a.tag AS tag,
		a.financing_balance AS financing_balance,
		a.financing_from AS financing_from,
		a.task_type AS task_type,
		a.airport_balance AS airport_balance,
		a.teaching AS teaching,
		a.weight AS weight
	FROM
		airport AS a 
	LEFT JOIN 
		airport_relationship AS ar
	ON 	
		ar.user_address = ? AND
		a.final_time is null AND
		a.id=ar.airport_id
		AND ar.airport_id IS not  NULL 
	ORDER BY a.weight
	LIMIT ?  offset ?`, address, pageSize, (page-1)*pageSize).Rows()
		if err != nil {
			return res, err
		}
		var tmp = new(model.Airport)
		for raw.Next() {
			if res == nil {
				res = make([]*model.Airport, 0)
			}
			err = storage.ScanRows(raw, tmp)
			if err != nil {
				return res, err
			}
			res = append(res, tmp)
		}
		return res, err
	})
	if err != nil {
		logrus.Errorf("query running airport (address:%s,page:%d,pagesize:%d) with weight by page error:%s ", address, page, pageSize, err.Error())
		return
	}
	return raw.([]*model.Airport), err
}

func (a *airport) QueryFinishAirportWithFinishTimeByPage(ctx context.Context, page int, pagesize int) (res []*model.Airport, err error) {
	var raw any

	storage := db.GetMysql()
	raw, err, _ = a.sf.Do(fmt.Sprintf("finish_%d_%d", page, pagesize), func() (interface{}, error) {
		var res []*model.Airport
		var err error
		err = storage.WithContext(ctx).Model(&model.Airport{}).Where("final_time is not null").Limit(pagesize).Offset((page - 1) * pagesize).Order("final_time desc").Find(&res).Error
		return res, err
	})
	if err != nil {
		logrus.Errorf("query finish airport (page:%d,pageSize:%d) with finish_time by page error:%s", page, pagesize, err.Error())
		return
	}
	return raw.([]*model.Airport), err
}
func (a *airport) UpdateAirport(ctx context.Context, data *model.Airport) (err error) {
	storage := db.GetMysql()
	return storage.WithContext(ctx).Updates(data).Error
}
