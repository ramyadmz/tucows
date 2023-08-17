package terminal

import (
	"errors"
	"image"
	"testing"

	"github.com/ramyad/tucows/internal/api/imageapi"
	"github.com/ramyad/tucows/internal/api/quoteapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAPIFacade struct {
	mock.Mock
}

func (m *MockAPIFacade) GetRandomQuoteWithImage(qtcnfbldr *quoteapi.QuoteConfigBuilder, imgCnfgBldr *imageapi.ImageConfigBuilder) (string, image.Image, error) {
	args := m.Called(qtcnfbldr, imgCnfgBldr)
	return args.String(0), args.Get(1).(image.Image), args.Error(2)
}

func TestRun_success(t *testing.T) {
	mockAPI := new(MockAPIFacade)
	app := NewTerminalApp(mockAPI)
	mockAPI.On("GetRandomQuoteWithImage", mock.Anything, mock.Anything).Return("Random Quote", image.NewRGBA(image.Rect(0, 0, 1, 1)), nil)
	err := app.Run()
	assert.Nil(t, err, "Expected no error")

}

func TestRun_FetchQuoteAndImageError(t *testing.T) {
	mockAPI := new(MockAPIFacade)
	app := NewTerminalApp(mockAPI)
	mockAPI.On("GetRandomQuoteWithImage", mock.Anything, mock.Anything).Return("", image.NewRGBA(image.Rect(0, 0, 0, 0)), errors.New("Failed to fetch random quote image"))
	err := app.Run()
	assert.Error(t, err, "Expected error as GetRandomQuoteWithImage returned error")
	assert.EqualError(t, err, "Failed to fetch random quote image")
}
