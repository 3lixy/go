package common

import (
	"strings"

	"github.com/Shopify/sarama"
)

const (
	//SyncProducerType 同步
	SyncProducerType = 1
	//AsyncProductType 异步
	AsyncProductType = 2
)

//Producer kafa生产者
//brokerList ip:端口  "127.0.0.1:3306"
//message  传递消息体
//producerType 生产发送  1同步  2异步
//返回nil成功
func Producer(brokerList string, Topic string, message string, producerType int) error {

	var err error
	switch producerType {
	case SyncProducerType:
		err = SyncProducer(brokerList, Topic, message)
	case AsyncProductType:
		err = ASyncProducer(brokerList, Topic, message)
	}
	return err
}

//SyncProducer 同步生产
func SyncProducer(brokerList string, Topic string, message string) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(strings.Split(brokerList, ","), config)
	if err != nil {
		return err
	}

	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: Topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

//ASyncProducer 异步生产
func ASyncProducer(brokerList string, Topic string, message string) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(strings.Split(brokerList, ","), config)

	defer producer.Close()

	if err != nil {
		return err
	}

	go func(producer sarama.AsyncProducer) {
		errors := producer.Errors()
		success := producer.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {

				}
			case <-success:
			}
		}
	}(producer)

	msg := &sarama.ProducerMessage{
		Topic: Topic,
		Value: sarama.ByteEncoder(message),
	}
	producer.Input() <- msg
	return nil
}
