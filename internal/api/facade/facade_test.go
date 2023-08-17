package facade_test

import (
	"errors"
	"image"
	"testing"

	"github.com/ramyad/tucows/internal/api/imageapi"
	"github.com/ramyad/tucows/internal/api/quoteapi"
)

type MockAPIFacade struct {
	ShouldSimulateError bool
}

func (m *MockAPIFacade) GetRandomQuoteWithImage(qtcnfbldr *quoteapi.QuoteConfigBuilder, imgCnfgBldr *imageapi.ImageConfigBuilder) (string, image.Image, error) {
	if m.ShouldSimulateError {
		return "", nil, errors.New("simulated error")
	}
	return "Mocked quote", image.NewRGBA(image.Rect(0, 0, 1, 1)), nil
}

func TestGetRandomQuoteWithImage_Success(t *testing.T) {
	mockQuoteConfigBuilder := quoteapi.NewQuoteConfigBuilder().WithKey(123)
	mockImageConfigBuilder := imageapi.NewImageConfigBuilder().
		WithWidth(300).
		WithHeight(200).
		WithFilters([]string{imageapi.ImageFilterGrayscale})

	mockAPIFacade := &MockAPIFacade{}

	quote, img, err := mockAPIFacade.GetRandomQuoteWithImage(mockQuoteConfigBuilder, mockImageConfigBuilder)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if quote == "" {
		t.Errorf("Expected a non-empty quote, but got an empty string")
	}

	if img == nil {
		t.Errorf("Expected a non-nil image, but got nil")
	}
}

func TestGetRandomQuoteWithImage_Error(t *testing.T) {
	mockQuoteConfigBuilder := quoteapi.NewQuoteConfigBuilder().WithKey(123)
	mockImageConfigBuilder := imageapi.NewImageConfigBuilder().
		WithWidth(300).
		WithHeight(200).
		WithFilters([]string{imageapi.ImageFilterGrayscale})

	mockAPIFacade := &MockAPIFacade{ShouldSimulateError: true}

	_, _, err := mockAPIFacade.GetRandomQuoteWithImage(mockQuoteConfigBuilder, mockImageConfigBuilder)
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
}
