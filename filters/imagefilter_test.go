package filters

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
	"github.com/zalando/skipper/filters/filtertest"
	"gopkg.in/h2non/bimg.v1"
)

const (
	widthTarget  = 400
	heightTarget = 200
)

var optionsTarget = bimg.Options{
	Width:  widthTarget,
	Height: heightTarget,
}

func createSampleImageReader(t *testing.T) io.ReadCloser {
	buffer, err := bimg.Read(imagefiltertest.LandscapeImageFile)
	assert.Nil(t, err, "Failed to read sample image")
	return ioutil.NopCloser(bytes.NewReader(buffer))
}

func readResultImage(r io.Reader, t *testing.T) *bimg.Image {
	result, err := ioutil.ReadAll(r)
	assert.Nil(t, err, "Error reading the result image")
	return bimg.NewImage(result)
}

func assertCorrectImageSize(r io.Reader, t *testing.T) {
	resultImage := readResultImage(r, t)
	size, _ := resultImage.Size()

	assert.Equal(t, widthTarget, size.Width, "The width is not correct")
	assert.Equal(t, heightTarget, size.Height, "The height is not correct")

}

func TestHandleResponse(t *testing.T) {
	imageReader := createSampleImageReader(t)
	response := &http.Response{Body: imageReader}
	response.Header = make(http.Header)
	response.Header.Add("Content-Length", "100")
	fc := &filtertest.Context{FResponse: response}
	imageFilter := imagefiltertest.FakeImageFilter(optionsTarget)

	HandleImageResponse(fc, &imageFilter)

	assertCorrectImageSize(fc.Response().Body, t)
}
