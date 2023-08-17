package web

import (
	"errors"
	"image"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ramyad/tucows/internal/api/facade"
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

func TestHandleRandomImageQuote_Success(t *testing.T) {
	mockAPI := new(MockAPIFacade)
	mockAPI.On("GetRandomQuoteWithImage", mock.Anything, mock.Anything).
		Return("Random Quote", image.NewRGBA(image.Rect(0, 0, 1, 1)), nil)
	app := &WebApp{
		API: mockAPI,
	}
	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()
	app.HandleRandomImageQuote(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected status OK (200)")
}

func TestHandleRandomImageQuote_FetchQuoteAndImageError(t *testing.T) {
	mockAPI := new(MockAPIFacade)
	mockAPI.On("GetRandomQuoteWithImage", mock.Anything, mock.Anything).
		Return("", image.NewRGBA(image.Rect(0, 0, 1, 1)), errors.New("Failed to fetch random quote image"))
	app := &WebApp{
		API: mockAPI,
	}
	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()
	app.HandleRandomImageQuote(recorder, req)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Expected error 500")
	assert.Contains(t, recorder.Body.String(), "Failed to fetch data", "Error message should be in the response body")
}

func TestParseRequest(t *testing.T) {
	assert := assert.New(t)
	api := facade.NewAPIFacade()
	app := NewWebApp(api, 8080)

	// Simulate query parameters
	queryParams := url.Values{}
	queryParams["key"] = []string{"123"}
	queryParams["width"] = []string{"600"}
	queryParams["height"] = []string{"400"}
	queryParams["filters"] = []string{"grayscale,blur"}
	request := httptest.NewRequest("GET", "/random-image-quote?key=123&width=600&height=400&filters=grayscale,blur", nil)
	request.URL.RawQuery = queryParams.Encode()

	app.(*WebApp).IncomingRequest = request
	err := app.ParseRequest()
	assert.Nil(err, "Expected no error")
	assert.Equal(123, app.(*WebApp).AppOptions.QuoteCategory, "QuoteCategory should be parsed correctly")
	assert.Equal(600, app.(*WebApp).AppOptions.ImageWidth, "ImageWidth should be parsed correctly")
	assert.Equal(400, app.(*WebApp).AppOptions.ImageHeight, "ImageHeight should be parsed correctly")
	assert.Equal([]string{"grayscale", "blur"}, app.(*WebApp).AppOptions.Filters, "Filters should be parsed correctly")
}
