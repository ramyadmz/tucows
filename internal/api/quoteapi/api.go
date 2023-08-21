// Package quoteapi provides functions for interacting with quote APIs and handling quote configurations.
package quoteapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/avast/retry-go"
	"github.com/ramyad/tucows/internal/shared"
)

// DefaultBaseUrl is the default base URL for the quote API.
var DefaultBaseUrl = "http://api.forismatic.com/api/1.0/"

const (
	// DefaultMethod and DefaultFormat represent default values for the quote API.
	DefualtMethod   = "getQuote"
	DefualtFormat   = "json"
	DefaultLanguage = "en"
	MaxKeyValue     = 999999
	RetryAttempts   = 4
	RetryDelay      = time.Second
)

// QuoteApiBuilder provides methods for building a quoteAPI instance.
type QuoteApiBuilder struct {
	api *quoteAPI
}

// quoteAPI represents a quote API with various properties.
type quoteAPI struct {
	baseURL  string
	method   string
	format   string
	language string
}

// QuoteProvider is an interface that defines the contract for fetching random quote.
type QuoteProvider interface {
	GetRandomQuote(qtCnfgBldr *QuoteConfigBuilder) (string, error)
}

// QuoteConfigBuilder provides methods for building a quoteConfig instance.
type QuoteConfigBuilder struct {
	config quoteConfig
}

// quoteConfig represents configuration options for fetching quotes.
type quoteConfig struct {
	Key int
}

// NewQuoteApiBuilder creates a new QuoteApiBuilder instance with default properties.
func NewQuoteApiBuilder() *QuoteApiBuilder {
	return &QuoteApiBuilder{
		api: &quoteAPI{
			baseURL:  DefaultBaseUrl,
			method:   DefualtMethod,
			format:   DefualtFormat,
			language: DefaultLanguage,
		},
	}
}

// WithBaseURL sets the base URL for the quote API and returns the builder instance.
func (tab *QuoteApiBuilder) WithBaseURL(baseURL string) *QuoteApiBuilder {
	tab.api.baseURL = baseURL
	return tab
}

// WithMethod sets the method for the quote API and returns the builder instance.
func (tab *QuoteApiBuilder) WithMethod(method string) *QuoteApiBuilder {
	tab.api.method = method
	return tab
}

// WithFormat sets the format for the quote API and returns the builder instance.
func (tab *QuoteApiBuilder) WithFormat(format string) *QuoteApiBuilder {
	tab.api.format = format
	return tab
}

// WithLanguage sets the language for the quote API and returns the builder instance.
func (tab *QuoteApiBuilder) WithLanguage(language string) *QuoteApiBuilder {
	tab.api.language = language
	return tab
}

// Build constructs and returns a QuoteProvider interface.
func (tab *QuoteApiBuilder) Build() QuoteProvider {
	return tab.api
}

// NewQuoteConfigBuilder creates a new QuoteConfigBuilder instance.
func NewQuoteConfigBuilder() *QuoteConfigBuilder {
	return &QuoteConfigBuilder{
		config: quoteConfig{},
	}
}

// WithKey sets the key in the configuration and returns the builder instance.
func (tcb *QuoteConfigBuilder) WithKey(key int) *QuoteConfigBuilder {
	tcb.config.Key = min(key, MaxKeyValue)
	return tcb
}

// Build constructs and returns a quoteConfig instance.
func (tcb *QuoteConfigBuilder) Build() quoteConfig {
	return tcb.config
}

// Data represents the structure of the JSON response data containing a quote.
type Data struct {
	QuoteText string `json:"quoteText"`
}

// GetRandomQuote fetches a random quote using the provided configuration from the quote API.
func (api *quoteAPI) GetRandomQuote(qtCnfgBldr *QuoteConfigBuilder) (string, error) {
	data := &Data{}
	path := api.buildPath(qtCnfgBldr.Build())

	var resp *http.Response
	err := retry.Do(
		func() error {
			var err error
			resp, err = http.Get(path)
			if err != nil {
				log.Printf("[%s] Get request Error: %v", shared.LogLevelError, err)
				return err
			}

			if resp.StatusCode != http.StatusOK {
				log.Printf("[%s] Received non-200 status code: %d", shared.LogLevelError, resp.StatusCode)
				return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
			}

			err = json.NewDecoder(resp.Body).Decode(data)
			if err != nil {
				log.Printf("[%s] JSON Parse Error: %v", shared.LogLevelError, err)
				return fmt.Errorf("JSON Parse Error: %w", err)
			}

			return nil
		},
		retry.Attempts(RetryAttempts),
		retry.DelayType(retry.BackOffDelay),
		retry.Delay(RetryDelay),
		retry.OnRetry(func(n uint, err error) {
			if n == uint(RetryAttempts-1) {
				log.Printf("[%s] Warning: Reached max retry attempts - 1.", shared.LogLevelWarning)
			}
		}),
	)

	if err != nil {
		log.Printf("[%s] Failed to get image from random quote API after retries: %v", shared.LogLevelError, err)
		return "", fmt.Errorf("failed to get image from random quote API after retries: %s", err)
	}

	return data.QuoteText, nil
}

// buildPath constructs the URL path for fetching a quote based on the provided configuration.
func (api *quoteAPI) buildPath(txtcnfg quoteConfig) string {
	baseURL, _ := url.Parse(api.baseURL)
	query := url.Values{
		"method": []string{api.method},
		"lang":   []string{api.language},
		"format": []string{api.format},
	}

	if txtcnfg.Key > 0 {
		query.Set("key", strconv.Itoa(txtcnfg.Key))
	}

	baseURL.RawQuery = query.Encode()
	return baseURL.String()
}
