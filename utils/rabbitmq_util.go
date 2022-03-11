package utils

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"seckill/common"
	"seckill/services"
)

var intiRabbitMQConn *amqp.Connection

func RabbitMQConn() *amqp.Connection {
	return intiRabbitMQConn
}

//测试发送
func TestSender() {
	conn := RabbitMQConn()
	//关闭连接
	//defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	declare, err := channel.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		panic(err)
	}

	err = channel.Publish(
		"",           // exchange
		declare.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("hello,world!!!"),
		})

	if err != nil {
		panic(err)
	}
}

//测试接收
func TestReceiver() {
	conn := RabbitMQConn()
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	declare, err := channel.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		panic(err)
	}
	msgs, err := channel.Consume(
		declare.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		panic(err)
	}
	go func() {
		for d := range msgs {
			println(d.Body)
		}
	}()
}

func NewRabbitMQConn() *amqp.Connection {
	//获取配置信息
	addr := viper.GetString("rabbitmq.addr")
	//初始化连接
	conn, err := amqp.Dial(addr)
	if err != nil {
		panic(err)
	}
	return conn
}

/**
消息发送
*/
func SendSecKillMessage(message []byte) {
	//获取RabbitMQ连接
	rabbitMQConn := NewRabbitMQConn()

	//关闭连接
	defer rabbitMQConn.Close()

	channel, err := rabbitMQConn.Channel()
	if err != nil {
		panic(err)
	}

	declare, err := channel.QueueDeclare(
		viper.GetString("rabbitmq.queue"), // name
		false,                             // durable
		false,                             // delete when unused
		false,                             // exclusive
		false,                             // no-wait
		nil,                               // arguments
	)
	if err != nil {
		panic(err)
	}

	err = channel.ExchangeDeclare(
		viper.GetString("rabbitmq.exchange"),
		"topic",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	err = channel.Publish(
		viper.GetString("rabbitmq.exchange"), // exchange
		declare.Name,                         // routing key
		false,                                // mandatory
		false,                                // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})

	if err != nil {
		panic(err)
	}
}

func ReceiveSecKillMessage() {
	//获取RabbitMQ连接
	rabbitMQConn := NewRabbitMQConn()

	//关闭连接
	//defer conn.Close()

	channel, err := rabbitMQConn.Channel()
	if err != nil {
		panic(err)
	}

	err = channel.ExchangeDeclare(
		viper.GetString("rabbitmq.exchange"),
		"topic",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	_, err = channel.QueueDeclare(
		viper.GetString("rabbitmq.queue"), // name
		false,                             // durable
		false,                             // delete when unused
		false,                             // exclusive
		false,                             // no-wait
		nil,                               // arguments
	)
	if err != nil {
		panic(err)
	}

	//队列绑定交换机
	err = channel.QueueBind(
		viper.GetString("rabbitmq.queue"),    //队列名
		viper.GetString("rabbitmq.queue"),    //key
		viper.GetString("rabbitmq.exchange"), //交换机
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	//消费者
	messages, err := channel.Consume(
		viper.GetString("rabbitmq.queue"),
		"",
		false,
		false,
		false, //设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false,
		nil,
	)

	forever := make(chan bool)

	//处理消息
	go func() {
		for v := range messages {
			var seckillMassage common.SecKillMassage
			jsoniter.Unmarshal(v.Body, &seckillMassage)
			order := services.MakeSecKillOrder(seckillMassage.User, seckillMassage.GoodId)
			if order == nil {
				v.Ack(false)
			} else {
				v.Ack(true)
			}
		}
	}()

	<-forever
}
