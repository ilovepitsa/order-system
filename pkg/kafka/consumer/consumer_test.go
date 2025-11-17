package consumer

import (
	"github.com/ilovepitsa/orders/pkg/kafka/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
	"testing"
	"time"
)

func MakeFetches(topic string, partition int32, offset int64, value []byte) []kgo.Fetch {
	return []kgo.Fetch{
		{
			Topics: []kgo.FetchTopic{
				{
					Topic: topic,
					Partitions: []kgo.FetchPartition{
						{
							Partition: partition,
							Records: []*kgo.Record{
								{
									Value:  value,
									Offset: offset,
									Topic:  topic,
								},
							},
						},
					},
				},
			},
		},
	}
}

func TestConsumer_Pool(t *testing.T) {
	type fields struct {
		cfg    common.ConsumerConfig
		logger *zap.Logger
	}
	type mockOutput struct {
		fetches kgo.Fetches
	}
	type output struct {
		out []common.Message
		err []error
	}
	tests := []struct {
		name        string
		fields      fields
		mockOut     mockOutput
		out         output
		testTimeout time.Duration
	}{
		// TODO: Add test cases.
		{
			name: "return_1_message_no_error ",
			fields: fields{
				cfg: common.ConsumerConfig{
					ChanBuff:      4,
					WriteTimeout:  time.Second,
					CommitTimeout: 3600 * time.Second,
				},
				logger: zap.NewNop(),
			},
			mockOut: mockOutput{
				fetches: MakeFetches("test", 0, 1, []byte("test")),
			},
			out: output{
				out: []common.Message{
					{
						MessageHeader: common.MessageHeader{
							Topic:     "test",
							Offset:    1,
							Partition: 0,
						},
						Value: []byte("test"),
					},
				},
				err: []error{},
			},
			testTimeout: time.Second * 60,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewMockKafkaConsumer(t)
			client.
				On("PollFetches", mock.Anything).
				Return(tt.mockOut.fetches)
			c := &Consumer{
				client: client,
				cfg:    tt.fields.cfg,
				logger: tt.fields.logger,
			}
			out, done, errChan := c.Pool(t.Context())
			close(done)
			select {
			case err := <-errChan:
				if len(tt.out.err) == 0 {
					assert.NoError(t, err.Err)
					return
				}
				assert.Equal(t, tt.out.err[0], err.Err)
			case out := <-out:
				if len(tt.out.err) != 0 {
					t.Fatal("expected error")
					return
				}
				assert.Equal(t, tt.out.out[0], out)
			case <-time.After(tt.testTimeout):
				t.Fatal("timeout")
			}

		})
	}
}
