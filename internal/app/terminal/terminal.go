package terminal

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/ramyad/tucows/internal/api/facade"
	"github.com/ramyad/tucows/internal/api/imageapi"
	"github.com/ramyad/tucows/internal/api/quoteapi"
	"github.com/ramyad/tucows/internal/app"
	"github.com/ramyad/tucows/internal/shared"
)

const (
	DefaultImageWidth  = 40
	DefaultImageHeight = 30
)

// TerminalApp implements the AppInterface for the terminal application.
type TerminalApp struct {
	api     facade.APIInterface
	options app.Options
}

// Ensure that *TerminalApp implements app.APP interface
var _ app.App = (*TerminalApp)(nil)

// NewTerminalApp creates a new instance of the TerminalApp.
func NewTerminalApp(apiFacade facade.APIInterface) app.App {
	return &TerminalApp{
		api: apiFacade,
	}
}

// Run executes the main logic for the terminal application,
// including parsing input, fetching a random quote and image,
// and displaying the content in the terminal.
func (t *TerminalApp) Run() error {
	logFile, err := os.Create("app.log")
	if err != nil {
		return fmt.Errorf("failed to create log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Printf("[%s] Starting terminal application...", shared.LogLevelInfo)

	if err := t.ParseRequest(); err != nil {
		log.Printf("[%s] Failed to parse input: %v", shared.LogLevelError, err)
		return err
	}

	if t.options.ImageWidth > DefaultImageWidth || t.options.ImageHeight > DefaultImageHeight {
		log.Printf("[%s] Requested image size exceeds default terminal size.", shared.LogLevelWarning)
	}

	log.Printf("[%s] Fetching random quote and image...", shared.LogLevelInfo)
	randomQuote, randomImage, err := t.FetchQuoteAndImage()
	if err != nil {
		log.Printf("[%s] Failed to fetch random quote and image: %v", shared.LogLevelError, err)
		return err
	}

	log.Printf("[%s] Displaying quote and image content in terminal...", shared.LogLevelInfo)
	if err := t.DisplayContent(randomQuote, randomImage); err != nil {
		log.Printf("[%s] Failed to display terminal image: %v", shared.LogLevelError, err)
		return err
	}

	log.Printf("[%s] Terminal application completed successfully.", shared.LogLevelInfo)
	return nil
}

// ParseRequest parses the command-line flags and returns the Options.
func (t *TerminalApp) ParseRequest() error {
	var filters string

	flag.IntVar(&t.options.QuoteCategory, "category", 0, "Specify the quote category")
	flag.IntVar(&t.options.ImageWidth, "width", DefaultImageWidth, "Specify the image width")
	flag.IntVar(&t.options.ImageHeight, "height", DefaultImageHeight, "Specify the image height")
	flag.StringVar(&filters, "filters", "", "Specify image filters as a comma-separated list: grayscale, blur")
	flag.Parse()

	t.options.Filters = imageapi.ImageFilters(strings.Split(filters, ","))
	return nil
}

// FetchQuoteAndImage fetches a random quote and image for the terminal application.
func (t *TerminalApp) FetchQuoteAndImage() (string, image.Image, error) {
	quoteConfigBuilder := quoteapi.NewQuoteConfigBuilder().WithKey(t.options.QuoteCategory)
	imageConfigBuilder := imageapi.NewImageConfigBuilder().WithWidth(t.options.ImageWidth).WithHeight(t.options.ImageHeight).WithFilters(t.options.Filters)

	randomQuote, randomImage, err := t.api.GetRandomQuoteWithImage(quoteConfigBuilder, imageConfigBuilder)
	if err != nil {
		return "", nil, err
	}

	return randomQuote, randomImage, nil
}

// DisplayContent displays the quote and image content for the terminal application.
func (t *TerminalApp) DisplayContent(quote string, img image.Image) error {
	displayRandomQuote(quote)
	return displayImageInTerminal(img, t.options.ImageWidth, t.options.ImageHeight)
}

// displayRandomQuote displays the random quote in the terminal.
func displayRandomQuote(quote string) {
	fmt.Println(quote)
}

// displayImageInTerminal displays the image in the terminal using ASCII art.
func displayImageInTerminal(img image.Image, width, height int) error {
	dc := gg.NewContext(width, height)
	dc.DrawImage(img, 0, 0)

	for y := 0; y < height; y += 1 {
		for x := 0; x < width; x += 1 {
			r, g, b, _ := dc.Image().At(x, y).RGBA()
			fmt.Printf("\x1b[48;2;%d;%d;%dm  \x1b[0m", r>>8, g>>8, b>>8)
		}
		fmt.Println()
	}
	return nil
}
