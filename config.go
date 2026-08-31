package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
)

type clusterMode string

const (
	singleProcess clusterMode = "single_process"
	singleNode    clusterMode = "single_node"
)

type RaftConfig struct {
	HeartBeatInterval uint `toml:"heartbeat_interval"`
	ElectionMin       uint `toml:"election_min"`
	ElectionMax       uint `toml:"election_max"`
}

type RawConfig struct {
	Addrs         []string    `toml:"addrs"`
	Mode          clusterMode `toml:"mode"`
	LogLevel      int         `toml:"log_level"`
	LogFormat     int         `toml:"log_format"`
	HTTPPprofAddr string      `toml:"http_pprof_addr"`
	RaftConfig    RaftConfig  `toml:"raft_config"`
	Persist       bool        `toml:"persist"`
}

const configFilePath = "config.toml"

func parseConfig() (*RawConfig, error) {
	var config RawConfig
	_, err := toml.DecodeFile(configFilePath, &config)
	if err != nil {
		return nil, fmt.Errorf("could not read config file. %w", err)
	}

	return &config, nil
}

func (rc *RawConfig) ToConfig() (*Configuration, error) {
	var configuration Configuration
	switch rc.LogLevel {
	case -4:
		configuration.LogLevel = slog.LevelDebug
	case 0:
		configuration.LogLevel = slog.LevelInfo
	case 4:
		configuration.LogLevel = slog.LevelWarn
	case 8:
		configuration.LogLevel = slog.LevelError
	default:
		return nil, fmt.Errorf("unsupported log level: %d", rc.LogLevel)

	}

	switch rc.LogFormat {
	case 0, 1:
		configuration.LogFormat = rc.LogFormat
	default:
		return nil, fmt.Errorf("unsupported log format: %d", rc.LogLevel)
	}

	if len(rc.Addrs) == 0 {
		return nil, fmt.Errorf("addr cannot be empty")
	}

	HEARTBEAT_INTERVAL = rc.RaftConfig.HeartBeatInterval
	ELECTION_INTERVAL_MIN = rc.RaftConfig.ElectionMin
	ELECTION_INTERVAL_MAX = rc.RaftConfig.ElectionMax

	configuration.Addrs = rc.Addrs
	// TODO: create a way for writing to file to instead
	configuration.Out = os.Stdout
	return &configuration, nil
}
