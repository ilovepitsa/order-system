package consumer

import (
	"context"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
	"go.uber.org/zap"
)

type ConsumerConfig struct {
	ChanBuff      int           `json:"chan_buff" yaml:"chan_buff"`
	WriteTimeout  time.Duration `json:"write_timeout" yaml:"write_timeout"`
	CommitTimeout time.Duration `json:"commit_timeout" yaml:"commit_timeout"`
}

type MessageHeader struct {
	Topic     string
	Offset    int64
	Partition int32
}

type Message struct {
	MessageHeader
	Value []byte
}

type Consumer struct {
	client *kgo.Client
	cfg    ConsumerConfig
	logger *zap.Logger
}

func NewConsumer(
	cfg ConsumerConfig,
	logger *zap.Logger,
	opts ...kgo.Opt,
) (*Consumer, error) {
	client, err := kgo.NewClient()
	if err != nil {
		return nil, err
	}
	return &Consumer{
		client: client,
		logger: logger,
	}, nil
}

func (c *Consumer) commitCompleted(ctx context.Context, done chan MessageHeader) {
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

func (c *Consumer) run(ctx context.Context, out chan Message, done chan MessageHeader, errOut chan kgo.FetchError) {
	go c.commitCompleted(ctx, done)
	for {
		fetches := c.client.PollFetches(ctx)
		for _, errs := range fetches.Errors() {
			errOut <- errs
		}

		fetches.EachPartition(func(ftp kgo.FetchTopicPartition) {
			ticker := time.NewTicker(c.cfg.WriteTimeout)
			for _, records := range ftp.Records {

				message := Message{
					MessageHeader: MessageHeader{
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

func (c *Consumer) Pool(ctx context.Context) (chan Message, chan MessageHeader, chan kgo.FetchError) {
	out := make(chan Message, c.cfg.ChanBuff)
	done := make(chan MessageHeader, c.cfg.ChanBuff)
	errChan := make(chan kgo.FetchError, c.cfg.ChanBuff*5)
	go c.run(ctx, out, done, errChan)
	return out, done, errChan
}

func ProvideKafkaConsumer() {
	client, err := kgo.NewClient()
	if err != nil {
		panic(err)
	}

	fetches := client.PollFetches(context.TODO())
	iter := fetches.RecordIter()
	if !iter.Done() {

	}
}
