package cron

import (
	"blog/dao"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type accessConsumerCron struct {
	internal *cron.Cron
}

func NewAccessConsumerCron() *accessConsumerCron {
	return &accessConsumerCron{
		internal: cron.New(),
	}

}

func (acc *accessConsumerCron) Run(Done func()) {
	var err error
	var connection *amqp.Connection
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/", viper.GetString("rabbitmq.username"), viper.GetString("rabbitmq.password"), viper.GetString("rabbitmq.host"), viper.GetString("rabbitmq.port"))
	connection, err = amqp.Dial(dsn)
	if err != nil {
		logrus.Fatalf("connect to rabbitmq %s failed: %s", dsn, err.Error())
	}
	logrus.Info("connect to rabbitmq success")
	var channel *amqp.Channel
	channel, err = connection.Channel()
	if err != nil {
		logrus.Fatal("create the rabbitmq channel failed:", err.Error())
	}
	// once exactly
	onceDone := sync.OnceFunc(Done)
	accessDao := dao.GetAccess()
	acc.internal.AddJob("*/2 * * * *", cron.FuncJob(func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("access consumer panic: %v", r)
				onceDone()
			}
		}()
		var messages <-chan amqp.Delivery
		messages, err = channel.Consume(viper.GetString("rabbitmq.accessqueue"), "access_cron_consumer", true, false, false, false, nil)
		if err != nil {
			logrus.Panic("consume the rabbitmq queue failed:", err.Error())
		}
		for {
			select {
			//TODO level=error msg="consumer unmarshal the access {{0 0}} failed: unexpected end of JSON input"
			case msg := <-messages:
				raw := msg.Body
				var access = dao.Access{}
				err = json.Unmarshal(raw, &access)
				if err != nil {
					logrus.Errorf("consumer unmarshal the access {%v} failed: %s", access, err.Error())
					continue
				}
				err = accessDao.IncrementAccessNumToDB(access)
				if err != nil {
					logrus.Errorf("increment the access {%v} num to db failed: %s", access, err.Error())
					continue
				}
			default:
			}
		}
	}))
	acc.internal.Start()
}
