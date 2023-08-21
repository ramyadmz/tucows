// Package imageapi provides functions for interacting with image APIs and handling image configurations.
package imageapi

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strings"
	"time"

	retry "github.com/avast/retry-go"
	"github.com/ramyad/tucows/internal/shared"
)

type ImageFilters []string

const (
	DefaultImageWidth  = 200
	DefaultImageHeight = 300
	MaxImageWidth      = 1920
	MaxImageHeight     = 1080
	defaultBaseUrl     = "https://picsum.photos"
	RetryAttempts      = 3
	RetryDelay         = time.Second
)

const (
	// ImageFilterGrayscale and ImageFilterBlur are constants representing available image filter options.
	ImageFilterGrayscale = "grayscale"
	ImageFilterBlur      = "blur"
)

// ImageProvider is an interface that defines the contract for fetching random images.
type ImageProvider interface {
	GetRandomImage(imgCnfg *ImageConfigBuilder) (image.Image, error)
}

// ImageAPIBuilder provides methods for building an imageAPI instance.
type ImageAPIBuilder struct {
	api *imageAPI
}

// imageAPI represents an image API with a base URL.
type imageAPI struct {
	baseURL string
}

// ImageConfigBuilder provides methods for building an imageConfig instance.
type ImageConfigBuilder struct {
	config imageConfig
}

// imageConfig represents configuration options for fetching images.
type imageConfig struct {
	Filters ImageFilters
	Width   int
	Height  int
}

// NewImageAPIBuilder creates a new ImageAPIBuilder instance with the default base URL.
func NewImageAPIBuilder() *ImageAPIBuilder {
	return &ImageAPIBuilder{
		api: &imageAPI{
			baseURL: defaultBaseUrl,
		},
	}
}

// WithBaseURL sets the base URL for the image API and returns the builder instance.
func (iab *ImageAPIBuilder) WithBaseURL(baseURL string) *ImageAPIBuilder {
	iab.api.baseURL = baseURL
	return iab
}

// Build constructs and returns an ImageProvider interface.
func (iab *ImageAPIBuilder) Build() ImageProvider {
	return iab.api
}

// NewImageConfigBuilder creates a new ImageConfigBuilder instance with default dimensions.
func NewImageConfigBuilder() *ImageConfigBuilder {
	return &ImageConfigBuilder{
		config: imageConfig{
			Width:  DefaultImageWidth,
			Height: DefaultImageHeight,
		},
	}
}

// WithWidth sets the image width in the configuration and returns the builder instance.
func (icb *ImageConfigBuilder) WithWidth(w int) *ImageConfigBuilder {
	if w > MaxImageWidth {
		log.Printf("[%s] Requested image width exceeds maximum allowed width.", shared.LogLevelWarning)
		w = MaxImageWidth
	}
	icb.config.Width = w
	return icb
}

// WithHeight sets the image height in the configuration and returns the builder instance.
func (icb *ImageConfigBuilder) WithHeight(h int) *ImageConfigBuilder {
	if h > MaxImageHeight {
		log.Printf("[%s] Requested image height exceeds maximum allowed height.", shared.LogLevelWarning)
		h = MaxImageHeight
	}
	icb.config.Height = h
	return icb
}

// WithFilters adds image filters to the configuration and returns the builder instance.
func (icb *ImageConfigBuilder) WithFilters(filters ImageFilters) *ImageConfigBuilder {
	icb.config.Filters = append(icb.config.Filters, filters...)
	return icb
}

// Build constructs and returns an imageConfig instance.
func (icb *ImageConfigBuilder) Build() imageConfig {
	return icb.config
}

// GetRandomImage fetches a random image using the provided configuration from the image API.
func (api *imageAPI) GetRandomImage(imgCnfg *ImageConfigBuilder) (image.Image, error) {
	path := api.buildPath(imgCnfg.Build())
	var resp *http.Response
	var image image.Image

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
			image, err = jpeg.Decode(resp.Body)
			if err != nil {
				log.Printf("[%s] JSON Parse Error: %v", shared.LogLevelError, err)
				return err
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
		log.Printf("[%s] Failed to get image from random image API after retries: %v", shared.LogLevelError, err)
		return nil, fmt.Errorf("failed to get image from random image API after retries: %s", err)
	}

	return image, nil
}

// buildPath constructs the URL path for fetching an image based on the provided configuration.
func (api *imageAPI) buildPath(imgCnfg imageConfig) string {
	sizeOptions := fmt.Sprintf("%d/%d", imgCnfg.Width, imgCnfg.Height)
	var filterBuilder strings.Builder
	for i, filter := range imgCnfg.Filters {
		if i != 0 {
			filterBuilder.WriteString("&")
		}
		switch filter {
		case ImageFilterGrayscale:
			filterBuilder.WriteString(ImageFilterGrayscale)
		case ImageFilterBlur:
			filterBuilder.WriteString(ImageFilterBlur)
		}
	}
	filterOptions := filterBuilder.String()
	var pathBuilder strings.Builder
	pathBuilder.WriteString(api.baseURL)
	pathBuilder.WriteString("/")
	pathBuilder.WriteString(sizeOptions)
	pathBuilder.WriteString(".jpg")
	if filterOptions != "" {
		pathBuilder.WriteString("?")
		pathBuilder.WriteString(filterOptions)
	}
	return pathBuilder.String()
}
