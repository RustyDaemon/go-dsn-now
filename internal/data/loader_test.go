package data

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/RustyDaemon/go-dsn-now/internal/config"
	"github.com/RustyDaemon/go-dsn-now/internal/model/response"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestLoadersValidateResponses(t *testing.T) {
	for _, tc := range []struct {
		name                          string
		status                        int
		configXML, dataXML, wantError string
	}{
		{"valid", 200, `<dsn><sites><site name="gdscc"><dish name="dss14"/></site></sites></dsn>`, `<dsn><dish name="dss14"/></dsn>`, ""},
		{"HTTP failure", 503, `<html/>`, `<html/>`, "HTTP 503"},
		{"malformed XML", 200, `<dsn>`, `<dsn>`, "XML"},
		{"unexpected document", 200, `<html/>`, `<html/>`, "no dishes"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, configRequest := range []bool{true, false} {
				body := tc.dataXML
				if configRequest {
					body = tc.configXML
				}
				client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					if req.Context() == nil {
						t.Fatal("request missing context")
					}
					return &http.Response{StatusCode: tc.status, Status: fmt.Sprintf("%d %s", tc.status, http.StatusText(tc.status)), Body: io.NopCloser(strings.NewReader(body))}, nil
				})}
				cfg := config.Load()
				errors := make(chan error, 1)
				configResult, dataResult := make(chan response.DSNConfig, 1), make(chan response.DSN, 1)
				if configRequest {
					LoadDSNConfig(context.Background(), client, cfg, configResult, errors)
				} else {
					LoadDSNData(context.Background(), client, cfg, dataResult, errors)
				}
				select {
				case err := <-errors:
					if tc.wantError == "" || !strings.Contains(err.Error(), tc.wantError) {
						t.Fatalf("error = %v; want %q", err, tc.wantError)
					}
				case <-configResult:
					if tc.wantError != "" {
						t.Fatalf("configuration accepted; want %q", tc.wantError)
					}
				case <-dataResult:
					if tc.wantError != "" {
						t.Fatalf("data accepted; want %q", tc.wantError)
					}
				default:
					t.Fatal("loader returned no result")
				}
			}
		})
	}
}

func TestLoadersHonorCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) { return nil, req.Context().Err() })}
	for _, configRequest := range []bool{true, false} {
		errs := make(chan error, 1)
		if configRequest {
			LoadDSNConfig(ctx, client, config.Load(), make(chan response.DSNConfig, 1), errs)
		} else {
			LoadDSNData(ctx, client, config.Load(), make(chan response.DSN, 1), errs)
		}
		if err := <-errs; !errors.Is(err, context.Canceled) {
			t.Fatalf("error = %v; want cancellation", err)
		}
	}
}
