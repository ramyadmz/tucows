package facade

import (
	"fmt"
	"image"
	"testing"

	"github.com/ramyad/tucows/internal/api/imageapi"
	"github.com/ramyad/tucows/internal/api/quoteapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockQuoteProvider struct {
	mock.Mock
}

func (m *MockQuoteProvider) GetRandomQuote(qc *quoteapi.QuoteConfigBuilder) (string, error) {
	args := m.Called(qc)
	return args.String(0), args.Error(1)
}

type MockImageProvider struct {
	mock.Mock
}

func (m *MockImageProvider) GetRandomImage(ic *imageapi.ImageConfigBuilder) (image.Image, error) {
	args := m.Called(ic)
	return args.Get(0).(image.Image), args.Error(1)
}

func TestGetRandomQuoteWithImage_success(t *testing.T) {

	mockQuoteProvider := new(MockQuoteProvider)
	mockImageProvider := new(MockImageProvider)
	mockImage := image.NewRGBA(image.Rect(0, 0, 100, 100))

	mockQuoteProvider.On("GetRandomQuote", mock.Anything).Return("Random Quote", nil)
	mockImageProvider.On("GetRandomImage", mock.Anything).Return(mockImage, nil)

	apiFacade := APIFacade{
		quoteProvider: mockQuoteProvider,
		imageProvider: mockImageProvider,
	}

	quote, image, err := apiFacade.GetRandomQuoteWithImage(quoteapi.NewQuoteConfigBuilder(), imageapi.NewImageConfigBuilder())
	assert.NoError(t, err)
	assert.Equal(t, "Random Quote", quote)
	assert.Equal(t, mockImage, image)
}

func TestGetRandomQuoteWithImage_quoteProviderReturnError(t *testing.T) {

	mockQuoteProvider := new(MockQuoteProvider)
	mockImageProvider := new(MockImageProvider)
	mockImage := image.NewRGBA(image.Rect(0, 0, 100, 100))

	mockQuoteProvider.On("GetRandomQuote", mock.Anything).Return("", fmt.Errorf("fetch quote failed"))
	mockImageProvider.On("GetRandomImage", mock.Anything).Return(mockImage, nil)

	apiFacade := APIFacade{
		quoteProvider: mockQuoteProvider,
		imageProvider: mockImageProvider,
	}

	_, _, err := apiFacade.GetRandomQuoteWithImage(quoteapi.NewQuoteConfigBuilder(), imageapi.NewImageConfigBuilder())
	assert.Error(t, err, fmt.Errorf("fetch quote failed"))
}

func TestGetRandomQuoteWithImage_imageProviderReturnError(t *testing.T) {

	mockQuoteProvider := new(MockQuoteProvider)
	mockImageProvider := new(MockImageProvider)
	mockImage := image.NewRGBA(image.Rect(0, 0, 100, 100))

	mockQuoteProvider.On("GetRandomQuote", mock.Anything).Return("Random Quote", nil)
	mockImageProvider.On("GetRandomImage", mock.Anything).Return(mockImage, fmt.Errorf("fetch image failed"))

	apiFacade := APIFacade{
		quoteProvider: mockQuoteProvider,
		imageProvider: mockImageProvider,
	}

	_, _, err := apiFacade.GetRandomQuoteWithImage(quoteapi.NewQuoteConfigBuilder(), imageapi.NewImageConfigBuilder())
	assert.Error(t, err, fmt.Errorf("fetch image failed"))
}
