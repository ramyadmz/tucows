package facade

import (
	"fmt"
	"image"
	"sync"

	"github.com/ramyad/tucows/internal/api/imageapi"
	"github.com/ramyad/tucows/internal/api/quoteapi"
)

// API represents an interface for interacting with various APIs to fetch random quotes and images.
type API interface {
	GetRandomQuoteWithImage(qtcnfbldr *quoteapi.QuoteConfigBuilder, imgCnfgBldr *imageapi.ImageConfigBuilder) (string, image.Image, error)
}

// APIFacade encapsulates interactions with quote and image APIs.
type APIFacade struct {
	quoteProvider quoteapi.QuoteProvider
	imageProvider imageapi.ImageProvider
}

// NewAPIFacade creates a new instance of APIFacade.
func NewAPIFacade() API {
	return &APIFacade{
		quoteProvider: quoteapi.NewQuoteApiBuilder().Build(),
		imageProvider: imageapi.NewImageAPIBuilder().Build(),
	}
}

// GetRandomQuoteWithImage fetches a random quote and image concurrently using the provided configurations.
// It returns the fetched quote, image, and any error encountered during the fetching process.
func (facade *APIFacade) GetRandomQuoteWithImage(qtcnfbldr *quoteapi.QuoteConfigBuilder, imgCnfgBldr *imageapi.ImageConfigBuilder) (string, image.Image, error) {
	var wg sync.WaitGroup
	var quote string
	var image image.Image
	var quoteErr, imageErr error

	quoteConfig := qtcnfbldr.Build()
	imageConfig := imgCnfgBldr.Build()

	wg.Add(2)
	go func() {
		defer wg.Done()
		quote, quoteErr = facade.quoteProvider.GetRandomQuote(quoteConfig)
	}()

	go func() {
		defer wg.Done()
		image, imageErr = facade.imageProvider.GetRandomImage(imageConfig)
	}()

	wg.Wait()

	if quoteErr != nil {
		return "", nil, fmt.Errorf("error calling quote api: %w", quoteErr)
	}

	if imageErr != nil {
		return "", nil, fmt.Errorf("error calling quote api: %w", imageErr)
	}

	return quote, image, nil
}
