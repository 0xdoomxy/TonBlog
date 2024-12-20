package controller

import (
	"blog/model"
	"blog/service"
	"blog/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type airport struct {
}

var airportController = &airport{}

func GetAirport() *airport {
	return airportController
}

func (a *airport) FindFinishAirport(c *gin.Context) {
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil || page <= 0 {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	var pagesize int64
	pagesize, err = strconv.ParseInt(c.Query("pagesize"), 10, 64)
	if err != nil || pagesize <= 0 {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	var res []*model.Airport
	res, err = service.GetAirport().QueryFinishAirportWithFinishTimeByPage(c, int(page), int(pagesize))
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("查询失败"))
		return
	}
	c.JSON(200, res)
}

func (a *airport) FindRunningAirport(c *gin.Context) {
	var ok bool
	var addressAny any
	addressAny, ok = c.Get("address")
	if !ok {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	var address string
	if address, ok = addressAny.(string); !ok {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil || page <= 0 {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	var pagesize int64
	pagesize, err = strconv.ParseInt(c.Query("pagesize"), 10, 64)
	if err != nil || pagesize <= 0 {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	var res []*model.Airport
	res, err = service.GetAirport().QueryRunningAirportWithWeightByPage(c, address, int(page), int(pagesize))
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("查询失败"))
		return
	}
	c.JSON(200, res)
}
func (a *airport) DeleteAirport(c *gin.Context) {

	airportId, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil || airportId <= 0 {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	err = service.GetAirport().DeleteAirport(c, uint(airportId))
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("删除失败"))
		return
	}
	return
}
func (a *airport) CreateAirport(c *gin.Context) {
	var airport = new(model.Airport)
	err := c.Bind(airport)
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	if airport.Weight > 5 || airport.Weight <= 0 {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	if airport.ID != 0 {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	if airport.EndTime.IsZero() {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	if airport.StartTime.IsZero() {
		airport.StartTime = time.Now()
	}
	if airport.Name == "" {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	err = service.GetAirport().CreateAirport(c, airport)
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("创建失败"))
		return
	}
}

func (a *airport) UpdateAirport(c *gin.Context) {
	var ok bool
	var addressAny any
	addressAny, ok = c.Get("address")
	if !ok {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	var address string
	if address, ok = addressAny.(string); !ok {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	updateType := service.UpdateSchema(c.Query("type"))
	if updateType == "" {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	var err error
	var airportId uint64
	airportId, err = strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	params := &service.UpdateAirportTemplate{
		Schema:              updateType,
		AirportRelationship: &model.AirportRelationship{AirportId: uint(airportId), UserAddress: address},
	}
	var balance float64
	switch updateType {
	case service.UserAddressBalance:
		balance, err = strconv.ParseFloat(c.Query("balance"), 64)
		if err != nil {
			c.JSON(200, utils.NewFailedResponse("参数出错"))
			return
		}
		params.AirportRelationship.Balance = balance
	}
	err = service.GetAirport().UpdateAirport(c, params)
	if err != nil {
		if err != nil {
			c.JSON(200, utils.NewFailedResponse("修改失败"))
			return
		}
	}
	return
}
