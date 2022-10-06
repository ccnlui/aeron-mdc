package config

import "time"

var (
	// Media driver
	AeronDir            string        = "/Volumes/DevShm/alvantage"
	MediaDriverTimeout  time.Duration = 10 * time.Second
	ServerChannel       string        = "aeron:udp?control=localhost:40456|control-mode=dynamic"
	ClientChannelFormat string        = "aeron:udp?endpoint=localhost:%d|control=localhost:40456|control-mode=dynamic"
	TimeStream          int           = 8000
)
