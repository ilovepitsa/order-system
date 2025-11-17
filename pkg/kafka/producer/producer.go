package producer

import (
	"context"
	"github.com/ilovepitsa/orders/pkg/kafka/common"
	utils "github.com/ilovepitsa/orders/pkg/utils/sync"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

type Producer struct {
	client  *kgo.Client
	cfg     common.ProducerConfig
	logger  *zap.Logger
	toSend  chan common.Message
	backlog *utils.SyncSlice[common.Message]
}

func NewProducer(
	cfg common.ProducerConfig,
	logger *zap.Logger,
	opts ...kgo.Opt,
) (*Producer, error) {
	client, err := common.BuildKafkaClient()
	if err != nil {
		return nil, err
	}
	return &Producer{
		toSend:  make(chan common.Message, 0),
		backlog: utils.NewSyncSlice[common.Message](0, cfg.SendBuff),
		client:  client,
		cfg:     cfg,
		logger:  logger,
	}, nil
}

func (p *Producer) send(ctx context.Context, message *common.Message) {
	record := &kgo.Record{
		Topic: message.Topic,
		Key:   []byte(message.Key),
		Value: message.Value,
	}
	p.client.Produce(ctx, record, func(record *kgo.Record, err error) {
		p.logger.Debug("Producer send message",
			zap.String("topic", record.Topic),
			zap.String("key", string(record.Key)),
			zap.String("value", string(record.Value)),
			zap.Int32("partition", record.Partition),
			zap.Error(err))
		if err != nil {
			message.OnFailure(message, err)
		}
	})
}

func (p *Producer) Produce(ctx context.Context, topic, key string, value []byte, onFailure common.FailureCallback) {
	message := &common.Message{
		MessageHeader: common.MessageHeader{
			Topic: topic,
			Key:   key,
		},
		Value: value,
		MessageAdditional: common.MessageAdditional{
			OnFailure: onFailure,
		},
	}
	go p.send(ctx, message)
}
