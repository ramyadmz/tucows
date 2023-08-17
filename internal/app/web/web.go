package web

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/ramyad/tucows/internal/api/facade"
	"github.com/ramyad/tucows/internal/api/imageapi"
	"github.com/ramyad/tucows/internal/api/quoteapi"
	"github.com/ramyad/tucows/internal/app"
	"github.com/ramyad/tucows/internal/shared"
	"github.com/ramyad/tucows/internal/static"
)

const (
	DefaultWebImageWidth  = 600
	DefaultWebImageHeight = 400
	DefaultTextCategory   = 0
)

// Data represents the data to be passed to the template.
type Data struct {
	Text  string
	Image string
}

// WebApp implements the AppInterface for the web application.
type WebApp struct {
	IncomingRequest *http.Request
	ResponseWriter  http.ResponseWriter
	API             facade.APIInterface
	AppOptions      app.Options
	RenderedContent Data
	Port            int
}

// Ensure that *WebApp implements app.APP interface
var _ app.App = (*WebApp)(nil)

// NewTerminalApp creates a new instance of the TerminalApp.
func NewWebApp(apiFacade facade.APIInterface, port int) app.App {
	return &WebApp{
		API:  apiFacade,
		Port: port,
	}
}

// Run starts the web application by handling the HTTP requests
// and serving the random image and quote content.
func (w *WebApp) Run() error {
	addr := fmt.Sprintf(":%d", w.Port)

	log.Printf("[%s] [%s] Starting web application on %s ...", time.Now(), shared.LogLevelInfo, addr)

	http.HandleFunc("/", w.HandleRandomImageQuote)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Printf("[%s] Failed to start web application: %v", shared.LogLevelError, err)
		return err
	}
	return nil
}

// HandleRandomImageQuote handles the HTTP request for random image and quote.
func (w *WebApp) HandleRandomImageQuote(responseWriter http.ResponseWriter, request *http.Request) {
	log.Printf("[%s] [%s] Handling random image and quote request...", time.Now(), shared.LogLevelInfo)

	w.IncomingRequest = request
	w.ResponseWriter = responseWriter

	err := w.ParseRequest()
	if err != nil {
		log.Printf("[%s] Invalid request parameteres: %v\n", shared.LogLevelError, err)
		http.Error(w.ResponseWriter, "Invalid request parameters", http.StatusBadRequest)
		return
	}

	quote, image, err := w.FetchQuoteAndImage()
	if err != nil {
		log.Printf("[%s] Failed to fetch data %v\n", shared.LogLevelError, err)
		http.Error(w.ResponseWriter, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	err = w.DisplayContent(quote, image)
	if err != nil {
		log.Printf("[%s] Failed to display data %v\n", shared.LogLevelError, err)
		http.Error(w.ResponseWriter, "Failed to display data", http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] [%s] Request handled successfully.", shared.LogLevelInfo, time.Now())
}

// ParseRequest parses the web request and returns the Options.
func (w *WebApp) ParseRequest() error {
	queryParams := w.IncomingRequest.URL.Query()
	w.AppOptions = *app.NewOptions(DefaultTextCategory, DefaultWebImageWidth, DefaultWebImageHeight, nil)

	keyParam := queryParams.Get("key")
	if len(keyParam) > 0 {
		var err error
		w.AppOptions.QuoteCategory, err = strconv.Atoi(keyParam)
		if err != nil {
			log.Printf("[%s] Invalid value for key parameter: %v\n", shared.LogLevelError, err)
			return fmt.Errorf("invalid value for key parameter: %s", err)
		}
	}

	log.Printf("[%s] Quote configuration: key=%d\n", shared.LogLevelInfo, w.AppOptions.QuoteCategory)

	widthParam := queryParams.Get("width")
	if len(widthParam) > 0 {
		var err error
		w.AppOptions.ImageWidth, err = strconv.Atoi(widthParam)
		if err != nil {
			log.Printf("[%s] Invalid value for width parameter: %v\n", shared.LogLevelError, err)
			return fmt.Errorf("invalid value for width parameter: %s", err)
		}
	}

	heightParam := queryParams.Get("height")
	if len(heightParam) > 0 {
		var err error
		w.AppOptions.ImageHeight, err = strconv.Atoi(heightParam)
		if err != nil {
			log.Printf("[%s] Invalid value for height parameter: %v\n", shared.LogLevelError, err)
			return fmt.Errorf("invalid value for height parameter: %s", err)
		}
	}

	filtersParam := queryParams.Get("filters")
	w.AppOptions.Filters = imageapi.ImageFilters(strings.Split(filtersParam, ","))

	log.Printf("[%s] Image configuration: width=%d, height=%d, filters=%v\n", shared.LogLevelInfo, w.AppOptions.ImageWidth, w.AppOptions.ImageHeight, w.AppOptions.Filters)

	return nil
}

// FetchQuoteAndImage fetches a random quote and image for the terminal application.
func (w *WebApp) FetchQuoteAndImage() (string, image.Image, error) {
	quoteConfigBuilder := quoteapi.NewQuoteConfigBuilder().WithKey(w.AppOptions.QuoteCategory)
	imageConfigBuilder := imageapi.NewImageConfigBuilder().WithWidth(w.AppOptions.ImageWidth).WithHeight(w.AppOptions.ImageHeight).WithFilters(w.AppOptions.Filters)

	randomQuote, randomImage, err := w.API.GetRandomQuoteWithImage(quoteConfigBuilder, imageConfigBuilder)
	if err != nil {
		log.Printf("[%s] Failed to GetRandomQuoteWithImage: %v\n", shared.LogLevelError, err)
		return "", nil, fmt.Errorf("failed to get random quote with image: %s", err)
	}

	return randomQuote, randomImage, nil
}

// DisplayContent displays the quote and image content for the web application.
func (w *WebApp) DisplayContent(quote string, img image.Image) error {
	image, err := encodeImageToBase64(img)
	if err != nil {
		log.Printf("[%s] Failed to encode image to base64: %v\n", shared.LogLevelError, err)
		return fmt.Errorf("failed to encode image to base64: %v", err)
	}

	w.RenderedContent.Image = image
	w.RenderedContent.Text = quote

	err = executeTemplate(w.ResponseWriter, w.RenderedContent)
	if err != nil {
		log.Printf("[%s] Failed to execute template: %v\n", shared.LogLevelError, err)
		return fmt.Errorf("failed to execute template: %v", err)
	}
	return nil
}

func encodeImageToBase64(img image.Image) (string, error) {
	imgBuffer := new(bytes.Buffer)
	if err := jpeg.Encode(imgBuffer, img, nil); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(imgBuffer.Bytes()), nil
}

func executeTemplate(w http.ResponseWriter, data Data) error {
	template, err := template.New("myTemplate").Parse(static.Template)
	if err != nil {
		return err
	}
	return template.Execute(w, data)
}
