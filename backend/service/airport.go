package service

import (
	"blog/dao"
	"blog/model"
	"context"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"time"
)

type airport struct {
	balanceLock map[uint]*redsync.Mutex
}

var airportService *airport = newAirportService()

func newAirportService() *airport {
	return &airport{
		balanceLock: make(map[uint]*redsync.Mutex),
	}
}

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

type UpdateAirportTemplate struct {
	Airport             *model.Airport             `json:"-"`
	AirportRelationship *model.AirportRelationship `json:"-"`
	Schema              UpdateSchema               `json:"-"`
}
type UpdateSchema string

const (
	UserUpdateTime     UpdateSchema = "user_update_time"
	UserAddressBalance UpdateSchema = "user_address_balance"
	UserFinishTime     UpdateSchema = "user_finish"
	UserAddIntoAddress UpdateSchema = "user_add_into_address"
)

func (a *airport) UpdateAirport(ctx context.Context, data *UpdateAirportTemplate) error {
	switch data.Schema {
	case UserUpdateTime:
		return a.updateUserUpdateTime(ctx, data)
	case UserAddressBalance:
		return a.updateUserAddressBalance(ctx, data)
	case UserFinishTime:
		return a.updateUserFinishTime(ctx, data)
	case UserAddIntoAddress:
		return a.createUserAddIntoAddress(ctx, data)
	default:
		return fmt.Errorf("unknow airport schema: %s", data.Schema)
	}
	return nil
}

func (a *airport) createUserAddIntoAddress(ctx context.Context, data *UpdateAirportTemplate) error {
	if data == nil || data.AirportRelationship == nil || data.AirportRelationship.AirportId <= 0 || len(data.AirportRelationship.UserAddress) <= 0 {
		return fmt.Errorf("参数出错")
	}
	airportRelationshipDao := dao.GetAirportRelationship()
	return airportRelationshipDao.CreateAirportRelationship(ctx, &model.AirportRelationship{AirportId: data.AirportRelationship.AirportId, UserAddress: data.AirportRelationship.UserAddress, CreateTime: time.Now()})
}
func (a *airport) updateUserFinishTime(ctx context.Context, data *UpdateAirportTemplate) (err error) {
	if data == nil || data.AirportRelationship == nil || data.AirportRelationship.AirportId <= 0 || len(data.AirportRelationship.UserAddress) <= 0 {
		return fmt.Errorf("参数出错")
	}
	airportRelationshipDao := dao.GetAirportRelationship()
	now := time.Now()
	return airportRelationshipDao.UpdateAirportRelationship(ctx, &model.AirportRelationship{AirportId: data.AirportRelationship.AirportId, UserAddress: data.AirportRelationship.UserAddress, UpdateTime: &now})
}
func (a *airport) updateUserAddressBalance(ctx context.Context, data *UpdateAirportTemplate) (err error) {
	if data == nil || data.AirportRelationship == nil || data.AirportRelationship.Balance < 0 || data.AirportRelationship.AirportId <= 0 || len(data.AirportRelationship.UserAddress) <= 0 {
		return fmt.Errorf("参数出错")
	}
	lock, ok := a.balanceLock[data.AirportRelationship.AirportId]
	if !ok {
		err = fmt.Errorf("系统出错啦")
		return
	}
	err = lock.Lock()
	if err != nil {
		return err
	}
	defer lock.Unlock()
	airportRelationshipDao := dao.GetAirportRelationship()
	airportDao := dao.GetAirport()
	var oldAirportRelationship *model.AirportRelationship
	oldAirportRelationship, err = airportRelationshipDao.FindAirportRelationshipByAddressAndId(ctx, data.AirportRelationship.UserAddress, data.AirportRelationship.AirportId)
	if err != nil {
		return err
	}
	incr := oldAirportRelationship.Balance - data.AirportRelationship.Balance
	err = airportDao.UpdateAirportBalance(ctx, oldAirportRelationship.AirportId, incr)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			airportDao.UpdateAirportBalance(ctx, oldAirportRelationship.AirportId, -incr)
		}
	}()
	err = airportRelationshipDao.UpdateAirportRelationship(ctx, &model.AirportRelationship{
		AirportId:   data.AirportRelationship.AirportId,
		UserAddress: data.AirportRelationship.UserAddress,
		Balance:     data.AirportRelationship.Balance,
	})
	return
}

func (a *airport) updateUserUpdateTime(ctx context.Context, data *UpdateAirportTemplate) (err error) {
	airportRelationshipDao := dao.GetAirportRelationship()
	updateTime := time.Now()
	return airportRelationshipDao.UpdateAirportRelationship(ctx, &model.AirportRelationship{
		AirportId:   data.AirportRelationship.AirportId,
		UserAddress: data.AirportRelationship.UserAddress,
		UpdateTime:  &updateTime,
	})
}

func (a *airport) DeleteAirport(ctx context.Context, airportId uint) (err error) {
	return dao.GetAirport().DeleteAirport(ctx, airportId)
}
