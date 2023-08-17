package imageapi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRandomImage_success(t *testing.T) {
	api := NewImageAPIBuilder().Build()
	imageConfig := NewImageConfigBuilder().WithWidth(800).WithHeight(1200).WithFilters(ImageFilters{"grayscale"}).Build()

	_, err := api.GetRandomImage(imageConfig)
	assert.NoError(t, err, "Expected no error for this image config")
}

func TestGetRandomImage_Error(t *testing.T) {
	api := NewImageAPIBuilder().WithBaseURL("http://unavailable.unavailable").Build()
	imageConfig := NewImageConfigBuilder().Build()
	expectedError := fmt.Errorf("failed to get image from random image API after retries")
	_, err := api.GetRandomImage(imageConfig)
	assert.Error(t, err, "Expected error due to unavailable API")
	assert.Contains(t, err.Error(), expectedError.Error(), "Expected error message mismatch")
}

func TestBuildPathWithImageSizeAndFilters(t *testing.T) {
	api := NewImageAPIBuilder().Build()
	imageConfig := NewImageConfigBuilder().WithWidth(400).WithHeight(600).WithFilters(ImageFilters{"blur", "grayscale"}).Build()
	expectedPath := "https://picsum.photos/400/600.jpg?blur&grayscale"
	resultPath := api.buildPath(imageConfig)
	assert.Equal(t, expectedPath, resultPath, "Path built with incorrect format")
}

func TestBuildPathWithoutImageConfig(t *testing.T) {
	api := NewImageAPIBuilder().Build()
	imageConfig := NewImageConfigBuilder().Build()
	expectedPath := "https://picsum.photos/200/300.jpg"
	resultPath := api.buildPath(imageConfig)
	assert.Equal(t, expectedPath, resultPath, "Path built with incorrect format")
}

func TestImageAPIBuilder(t *testing.T) {
	expected := imageAPI{
		baseURL: "test",
	}
	result := NewImageAPIBuilder().WithBaseURL(expected.baseURL).Build()
	assert.Equal(t, expected, result, "API Builder does not create the expected instance")
}

func TestImageConfigBuilder(t *testing.T) {
	expected := imageConfig{
		Width:   100,
		Height:  200,
		Filters: ImageFilters{"grayscale"},
	}
	result := NewImageConfigBuilder().WithWidth(expected.Width).WithHeight(expected.Height).WithFilters(expected.Filters).Build()
	assert.Equal(t, expected, result, "Config Builder does not create the expected instance")
}

func TestImageConfigBuilder_WidthReplacesHigherValueWithMax(t *testing.T) {
	builder := NewImageConfigBuilder()
	newWidth := MaxImageWidth + 100 // A value higher than the maximum

	config := builder.WithWidth(newWidth).Build()

	assert.Equal(t, MaxImageWidth, config.Width, "Width should be replaced with MaxImageWidth")
}

func TestImageConfigBuilder_HeightReplacesHigherValueWithMax(t *testing.T) {
	builder := NewImageConfigBuilder()
	newHeight := MaxImageHeight + 100 // A value higher than the maximum

	config := builder.WithHeight(newHeight).Build()

	assert.Equal(t, MaxImageHeight, config.Height, "Height should be replaced with MaxImageHeight")
}
