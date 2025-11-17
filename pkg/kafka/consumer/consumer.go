package consumer

import (
	"context"
	"time"

	"github.com/ilovepitsa/orders/pkg/kafka/common"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
	"go.uber.org/zap"
)

type KafkaConsumer interface {
	CommitOffsets(
		ctx context.Context,
		uncommitted map[string]map[int32]kgo.EpochOffset,
		onDone func(*kgo.Client, *kmsg.OffsetCommitRequest, *kmsg.OffsetCommitResponse, error),
	)
	PollFetches(ctx context.Context) kgo.Fetches
}

type Consumer struct {
	client KafkaConsumer
	cfg    common.ConsumerConfig
	logger *zap.Logger
}

func ProvideConsumer(
	cfg common.ConsumerConfig,
	logger *zap.Logger,
	opts ...kgo.Opt,
) (*Consumer, error) {
	client, err := common.BuildKafkaClient(opts...)
	if err != nil {
		return nil, err
	}
	return NewConsumer(client, cfg, logger), nil
}

func NewConsumer(
	client KafkaConsumer,
	cfg common.ConsumerConfig,
	logger *zap.Logger,
) *Consumer {
	return &Consumer{
		cfg:    cfg,
		client: client,
		logger: logger,
	}
}

func (c *Consumer) commitCompleted(ctx context.Context, done chan common.MessageHeader) {
	commitTicker := time.NewTicker(c.cfg.CommitTimeout)
	batch := map[string]map[int32]kgo.EpochOffset{}
	for {
		select {
		case <-ctx.Done():
			return
		case <-commitTicker.C:
			if len(batch) == 0 {
				continue
			}
			c.client.CommitOffsets(ctx, batch, func(client *kgo.Client, ocr1 *kmsg.OffsetCommitRequest, ocr2 *kmsg.OffsetCommitResponse, err error) {
				c.logger.Error("failed to commit offset", zap.Any("offsets", batch))
			})
			batch = map[string]map[int32]kgo.EpochOffset{}
		case m := <-done:
			if _, ok := batch[m.Topic]; !ok {
				batch[m.Topic] = map[int32]kgo.EpochOffset{}
			}
			batch[m.Topic][m.Partition] = kgo.EpochOffset{Offset: m.Offset + 1, Epoch: -1}

		}
	}
}

func (c *Consumer) run(ctx context.Context, out chan common.Message, done chan common.MessageHeader, errOut chan kgo.FetchError) {
	go c.commitCompleted(ctx, done)
	defer func() {
		close(out)
		close(errOut)
	}()
	for {
		fetches := c.client.PollFetches(ctx)
		for _, errs := range fetches.Errors() {
			errOut <- errs
		}

		fetches.EachPartition(func(ftp kgo.FetchTopicPartition) {
			ticker := time.NewTicker(c.cfg.WriteTimeout)
			for _, records := range ftp.Records {
				message := common.Message{
					MessageHeader: common.MessageHeader{
						Topic:     records.Topic,
						Offset:    records.Offset,
						Partition: records.Partition,
					},
					Value: records.Value,
				}
				select {
				case <-ctx.Done():
					return
				case out <- message:
				case <-ticker.C:
				}
			}
		})
	}
}

func (c *Consumer) Pool(ctx context.Context) (chan common.Message, chan common.MessageHeader, chan kgo.FetchError) {
	out := make(chan common.Message, c.cfg.ChanBuff)
	done := make(chan common.MessageHeader, c.cfg.ChanBuff)
	errChan := make(chan kgo.FetchError, c.cfg.ChanBuff*5)
	go c.run(ctx, out, done, errChan)
	return out, done, errChan
}
