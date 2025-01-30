package network

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/RustyDaemon/go-dsn-now/model/response"
)

var (
	DSNConfigUrl = "https://eyes.nasa.gov/apps/dsn-now/config.xml"
	DSNDataUrl   = "https://eyes.nasa.gov/dsn/data/dsn.xml?r=%v"
)

func LoadDSNConfig(result chan response.DSNConfig, ce chan error) {
	resp, err := http.Get(DSNConfigUrl)
	if err != nil {
		ce <- err
	}
	defer resp.Body.Close()

	dsnConfig := response.DSNConfig{}
	err = xml.NewDecoder(resp.Body).Decode(&dsnConfig)
	if err != nil {
		ce <- err
	}

	result <- dsnConfig
}

func LoadDSNData(result chan response.DSN, ce chan error) {
	r := time.Now().UnixMilli() / 5000
	url := fmt.Sprintf(DSNDataUrl, r)
	resp, err := http.Get(url)
	if err != nil {
		ce <- err
	}
	defer resp.Body.Close()

	var dsn response.DSN
	err = xml.NewDecoder(resp.Body).Decode(&dsn)
	if err != nil {
		ce <- err
	}

	result <- dsn
}
