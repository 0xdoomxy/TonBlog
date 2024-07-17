package cron

import (
	"blog/dao"
	"context"
	"encoding/json"
	"fmt"

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

func (acc *accessConsumerCron) Run() {
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
	accessDao := dao.GetAccess()
	acc.internal.AddJob("*/2 * * * *", cron.FuncJob(func() {
		var messages <-chan amqp.Delivery
		messages, err = channel.Consume(viper.GetString("rabbitmq.accessqueue"), "", true, false, false, false, nil)
		if err == amqp.ErrClosed {
			logrus.Errorf("the rabbitmq channel is closed")
			if connection.IsClosed() {
				connection, err = amqp.Dial(dsn)
				if err != nil {
					logrus.Fatalf("connect to rabbitmq %s failed: %s", dsn, err.Error())
				}
				logrus.Info("connect to rabbitmq success")
				channel, err = connection.Channel()
				if err != nil {
					logrus.Fatal("create the rabbitmq channel failed:", err.Error())
				}
			}
			return
		}
		if err != nil {
			logrus.Errorf("consume the rabbitmq queue failed:%s", err.Error())
			return
		}
		for {
			select {
			//TODO level=error msg="consumer unmarshal the access {{0 0}} failed: unexpected end of JSON input"
			case msg := <-messages:
				raw := msg.Body
				if len(raw) <= 0 {
					msg.Ack(true)
					return
				}
				var access = dao.Access{}
				err = json.Unmarshal(raw, &access)
				if err != nil {
					msg.Ack(false)
					logrus.Errorf("consumer unmarshal the access {%v} failed: %s", msg.Body, err.Error())
					return
				}
				err = accessDao.IncrementAccessNumToDB(context.TODO(), access)
				if err != nil {
					msg.Ack(false)
					logrus.Errorf("increment the access {%v} num to db failed: %s", access, err.Error())
					return
				}
				msg.Ack(true)
			default:
				return
			}
		}
	}))
	acc.internal.Start()
}
