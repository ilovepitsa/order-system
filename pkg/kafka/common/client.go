package common

import (
	"github.com/twmb/franz-go/pkg/kgo"
)

func BuildKafkaClient(opts ...kgo.Opt) (*kgo.Client, error) {
	opts = append(opts, kgo.DisableAutoCommit())
	return kgo.NewClient(opts...)
}
