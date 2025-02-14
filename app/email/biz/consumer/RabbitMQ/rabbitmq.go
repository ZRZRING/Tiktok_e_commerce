package RabbitMQ

import (
	"fmt"
	"github.com/Blue-Berrys/Tiktok_e_commerce/app/email/conf"
	"github.com/Blue-Berrys/Tiktok_e_commerce/app/email/infra/rpc"
	"github.com/cloudwego/kitex/pkg/klog"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

var (
	//定时取消订单
	ORDER_DELAY_QUEUE_NAME = "order_dlq"
)

func InitRabbitMQConsumer() {

}

func StartDelayQueueConsumer() {
	//1.连接到rabbitmq server
	conn, err := connectRabbitMQ()
	if err != nil {
		klog.Fatal(err)
	}
	defer func() {
		_ = conn.Close()
	}()
	klog.Info("started delay queue consumer")

	//2.打开channel
	ch, err := conn.Channel()
	if err != nil {
		klog.Errorf("打开channel失败:%v", err)
	}
	defer func() {
		_ = ch.Close()
	}()

	//3.声明queue
	_, err = ch.QueueDeclare(
		ORDER_DELAY_QUEUE_NAME,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		klog.Fatalf("延迟队列声明失败: %v", err)
	}

	//4.消息延迟队列中的消息.事实上就是将订单状态修改为canceled状态
	go func() {
		dlqMsg, err := ch.Consume(
			ORDER_DELAY_QUEUE_NAME,
			"",
			false,
			false,
			false,
			false,
			nil)
		if err != nil {
			klog.Fatalf("注册延迟队列消费者失败:%v", err)
		}

		for msg := range dlqMsg {
			orderId := msg.Body
			rpc.OrderClient.
		}
	}()
}

func connectRabbitMQ() (*amqp.Connection, error) {
	rabbitMQAddr := fmt.Sprintf(conf.GetConf().RabbitMQ.Addr,
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_SWITCH"),
	)
	return amqp.Dial(rabbitMQAddr)
}
