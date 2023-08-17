package facade

import (
	"image"
	"log"

	"github.com/ramyad/tucows/internal/api/imageapi"
	"github.com/ramyad/tucows/internal/api/quoteapi"
	"github.com/ramyad/tucows/internal/shared"
)

// APIInterface represents the interface for your API methods.
type APIInterface interface {
	GetRandomQuoteWithImage(qtcnfbldr *quoteapi.QuoteConfigBuilder, imgCnfgBldr *imageapi.ImageConfigBuilder) (string, image.Image, error)
}

// APIFacade encapsulates interactions with quote and image APIs.
type APIFacade struct {
	quoteAPIBuilder *quoteapi.QuoteApiBuilder
	imageAPIBuilder *imageapi.ImageAPIBuilder
}

// NewAPIFacade creates a new instance of APIFacade.
func NewAPIFacade() APIInterface {
	return &APIFacade{
		quoteAPIBuilder: quoteapi.NewQuoteApiBuilder(),
		imageAPIBuilder: imageapi.NewImageAPIBuilder(),
	}
}

// GetRandomQuoteWithImage fetches a random quote and image concurrently using the provided configurations.
// It returns the fetched quote, image, and any error encountered during the fetching process.
func (facade *APIFacade) GetRandomQuoteWithImage(qtcnfbldr *quoteapi.QuoteConfigBuilder, imgCnfgBldr *imageapi.ImageConfigBuilder) (string, image.Image, error) {
	quoteAPI := facade.quoteAPIBuilder.Build()
	imageAPI := facade.imageAPIBuilder.Build()

	quoteConfig := qtcnfbldr.Build()
	imageConfig := imgCnfgBldr.Build()

	// Create channels for results and errors
	quoteChan := make(chan string)
	imageChan := make(chan image.Image)
	errorChan := make(chan error, 2) // Two possible errors: one for quote and one for image

	// Fetch quote concurrently
	go func() {
		quote, err := quoteAPI.GetRandomQuote(quoteConfig)
		if err != nil {
			errorChan <- err
			return
		}
		quoteChan <- quote
	}()

	// Fetch image concurrently
	go func() {
		image, err := imageAPI.GetRandomImage(imageConfig)
		if err != nil {
			errorChan <- err
			return
		}
		imageChan <- image
	}()

	// Receive results
	var quote string
	var image image.Image
	var errors []error

	for i := 0; i < 2; i++ {
		select {
		case q := <-quoteChan:
			quote = q
		case img := <-imageChan:
			image = img
		case err := <-errorChan:
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		log.Printf("[%s] Errors during concurrent fetching: %v", shared.LogLevelError, errors)
		return "", nil, errors[0] // Return the first error encountered
	}

	return quote, image, nil
}
