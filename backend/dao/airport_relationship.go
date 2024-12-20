package dao

import (
	"blog/dao/db"
	"blog/model"
	"context"

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

func (ar *airportRelationship) UpdateAirportRelationship(ctx context.Context)
