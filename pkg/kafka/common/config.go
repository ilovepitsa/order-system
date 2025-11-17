package common

import "time"

type ConsumerConfig struct {
	ChanBuff      int           `json:"chan_buff" yaml:"chan_buff"`
	WriteTimeout  time.Duration `json:"write_timeout" yaml:"write_timeout"`
	CommitTimeout time.Duration `json:"commit_timeout" yaml:"commit_timeout"`
}

type ProducerConfig struct {
	WorkerAmount      int           `json:"worker_amount" yaml:"worker_amount"`
	QueueWriteTimeout time.Duration `json:"queue_write_timeout" yaml:"queue_write_timeout"`
	SendBuff          int           `json:"send_buff" yaml:"send_buff"`
	SendPeriod        time.Duration `json:"send_period" yaml:"send_period"`
}
