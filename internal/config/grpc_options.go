package config

import (
	"time"
)

type GRPCOptions struct {
	Host            string
	MaxRetry        uint
	PerRetryTimeout time.Duration
}
