package ratelimiter

import "time"

type Limiter interface {
	Permit(ip string) (bool, time.Duration)
}

type Config struct {
	RequestPerTimeFrame int
	TimeFrame           time.Duration
	Enabled             bool
}
