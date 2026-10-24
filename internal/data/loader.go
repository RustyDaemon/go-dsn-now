package data

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/RustyDaemon/go-dsn-now/internal/config"
	"github.com/RustyDaemon/go-dsn-now/internal/model/response"
)

func NewHTTPClient(cfg *config.Config) *http.Client {
	return &http.Client{
		Timeout: cfg.HTTPTimeout,
	}
}

func LoadDSNConfig(ctx context.Context, client *http.Client, cfg *config.Config, result chan response.DSNConfig, ce chan error) {
	req, err := http.NewRequestWithContext(ctx, "GET", cfg.DSNConfigURL, nil)
	if err != nil {
		ce <- err
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		ce <- err
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		ce <- fmt.Errorf("HTTP %s", resp.Status)
		return
	}

	dsnConfig := response.DSNConfig{}
	err = xml.NewDecoder(resp.Body).Decode(&dsnConfig)
	if err != nil {
		ce <- fmt.Errorf("invalid configuration XML: %w", err)
		return
	}
	dishes := 0
	for _, site := range dsnConfig.Sites {
		dishes += len(site.Dishes)
	}
	if dishes == 0 {
		ce <- fmt.Errorf("configuration contains no dishes")
		return
	}

	result <- dsnConfig
}

func LoadDSNData(ctx context.Context, client *http.Client, cfg *config.Config, result chan response.DSN, ce chan error) {
	r := time.Now().UnixMilli() / 5000
	url := fmt.Sprintf(cfg.DSNDataURL, r)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		ce <- err
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		ce <- err
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		ce <- fmt.Errorf("HTTP %s", resp.Status)
		return
	}

	var dsn response.DSN
	err = xml.NewDecoder(resp.Body).Decode(&dsn)
	if err != nil {
		ce <- fmt.Errorf("invalid live data XML: %w", err)
		return
	}
	if len(dsn.D) == 0 {
		ce <- fmt.Errorf("live data contains no dishes")
		return
	}

	result <- dsn
}
