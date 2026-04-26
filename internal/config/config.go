package config

import "time"

type Config struct {
	DSNConfigURL    string
	DSNDataURL      string
	RefreshInterval time.Duration
	HTTPTimeout     time.Duration
	AppVersion      string
	AppGithubURL    string
	Theme           string
}

func Load() *Config {
	return &Config{
		DSNConfigURL:    "https://eyes.nasa.gov/apps/dsn-now/config.xml",
		DSNDataURL:      "https://eyes.nasa.gov/dsn/data/dsn.xml?r=%v",
		RefreshInterval: 10 * time.Second,
		HTTPTimeout:     10 * time.Second,
		AppVersion:      "0.1.3",
		AppGithubURL:    "https://github.com/RustyDaemon/go-dsn-now",
		Theme:           "dark",
	}
}
