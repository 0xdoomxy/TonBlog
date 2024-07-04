package db

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

func init() {
	var err error
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/", viper.GetString("rabbitmq.username"), viper.GetString("rabbitmq.password"), viper.GetString("rabbitmq.host"), viper.GetString("rabbitmq.port"))
	rabbitmq, err = amqp.Dial(dsn)
	if err != nil {
		logrus.Fatalf("connect to rabbitmq %s failed: %s", dsn, err.Error())
	}
	logrus.Info("connect to rabbitmq success")

}

var rabbitmq *amqp.Connection

func GetRabbitmq() *amqp.Connection {
	return rabbitmq
}

func GetRabbitmqChannel() (*amqp.Channel, error) {
	return rabbitmq.Channel()

}
