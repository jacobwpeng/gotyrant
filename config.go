package gotyrant

import "time"

type Config struct {
	Addr    string
	Timeout time.Duration
}
