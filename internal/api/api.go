package api

import (
	"image"

	"github.com/ramyad/tucows/internal/api/imageapi"
	"github.com/ramyad/tucows/internal/api/quoteapi"
)

// API represents an interface for interacting with various APIs to fetch random quotes and images.
type API interface {
	GetRandomQuoteWithImage(qtcnfbldr *quoteapi.QuoteConfigBuilder, imgCnfgBldr *imageapi.ImageConfigBuilder) (string, image.Image, error)
}
