package quoteapi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRandomQuote_success(t *testing.T) {
	quoteAPI := NewQuoteApiBuilder().Build()
	quoteConfig := NewQuoteConfigBuilder().WithKey(100).Build()

	result, err := quoteAPI.GetRandomQuote(quoteConfig)
	assert.NoError(t, err, "Expected no error from GetRandomQuote")
	assert.NotEmpty(t, result, "Expected a non-empty quote result")
}

func TestGetRandomQuote_Error(t *testing.T) {
	quoteAPI := NewQuoteApiBuilder().WithBaseURL("http://unavailable.api").Build()
	quoteConfig := NewQuoteConfigBuilder().Build()
	expectedError := fmt.Errorf("failed to get image from random quote API after retries")

	_, err := quoteAPI.GetRandomQuote(quoteConfig)
	assert.Error(t, err, "Expected error due to unavailable API")
	assert.Contains(t, err.Error(), expectedError.Error(), "Expected error message mismatch")
}

func TestBuildPath(t *testing.T) {
	api := NewQuoteApiBuilder().Build()
	quoteConfig := NewQuoteConfigBuilder().WithKey(100).Build()
	expectedPath := "http://api.forismatic.com/api/1.0/?format=json&key=100&lang=en&method=getQuote"

	result := api.(*quoteAPI).buildPath(quoteConfig)
	assert.Equal(t, expectedPath, result, "Path built with incorrect format")
}

func TestQuoteAPIBuilder(t *testing.T) {
	expected := &quoteAPI{
		baseURL:  "test",
		method:   "get",
		format:   "xml",
		language: "ru",
	}
	result := NewQuoteApiBuilder().WithBaseURL(expected.baseURL).WithMethod(expected.method).WithFormat(expected.format).WithLanguage(expected.language).Build()
	assert.Equal(t, expected, result, "API Builder does not create the expected instance")
}

func TestQuoteConfigBuilder_WidthReplacesKeyValueWithMax(t *testing.T) {
	assert := assert.New(t)
	builder := NewQuoteConfigBuilder()
	maxKey := MaxKeyValue + 100 // A key value greater than the maximum allowed value
	config := builder.WithKey(maxKey).Build()
	assert.Equal(MaxKeyValue, config.Key, "Key should be set to the maximum allowed value")
}
