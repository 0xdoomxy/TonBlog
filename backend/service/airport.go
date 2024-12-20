package service

import (
	"blog/dao"
	"blog/model"
	"context"
)

type airport struct {
}

var airportService *airport

func GetAirport() *airport {
	return airportService
}

func (a *airport) QueryFinishAirportWithFinishTimeByPage(ctx context.Context, page int, pageSize int) ([]*model.Airport, error) {
	return dao.GetAirport().QueryFinishAirportWithFinishTimeByPage(ctx, page, pageSize)
}
func (a *airport) QueryRunningAirportWithWeightByPage(ctx context.Context, address string, page int, pageSize int) ([]*model.Airport, error) {
	return dao.GetAirport().QueryRunningAirportWithWeightByPage(ctx, address, page, pageSize)
}
func (a *airport) CreateAirport(ctx context.Context, data *model.Airport) (err error) {
	return dao.GetAirport().CreateAirport(ctx, data)
}
