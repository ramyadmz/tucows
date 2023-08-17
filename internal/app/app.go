package app

import (
	"image"
)

type Options struct {
	QuoteCategory int
	ImageWidth    int
	ImageHeight   int
	Filters       []string
}

func NewOptions(quoteCategory, imageWidth, imageHeight int, filters []string) *Options {
	return &Options{quoteCategory, imageWidth, imageHeight, filters}
}

type App interface {
	ParseRequest() error
	FetchQuoteAndImage() (string, image.Image, error)
	DisplayContent(quote string, img image.Image) error
	Run() error
}
