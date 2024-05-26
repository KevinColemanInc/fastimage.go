package fastimage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const maxBytes = 12

// ImageType represents the type of image
type ImageType int

//go:generate moq -out mock_http_client.go . HttpClient
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Enumerate the different image types
const (
	JPEG ImageType = iota
	JPEG2000
	PNG
	GIF
	BMP
	WEBP
	JFIF
	EXIF
	HEIC
)

var (
	ErrUnsupportedImageType = errors.New("unsupported image type")
	ErrFilesizeTooSmall     = errors.New("file size too small")
	ErrFilesizeTooLarge     = errors.New("file size too large")
	ErrHTTPCodeInvalid      = errors.New("invalid HTTP status code")
)

type fastimageImpl struct {
	httpClient HttpClient
}

func New(httpClient HttpClient) *fastimageImpl {
	return &fastimageImpl{
		httpClient: httpClient,
	}
}

func (i *fastimageImpl) ValidateURL(url string, policy *Policy) (bool, []byte, *http.Response, error) { // todo return response
	// Fetch the image from the URL
	resp, err := i.fetchFirstBytes(url)
	respBytesFirst, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil && (resp != nil && resp.StatusCode != http.StatusPartialContent) {
		return false, respBytesFirst, resp, err
	}
	err = nil
	// if resp is not 2xx
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return false, respBytesFirst, resp, ErrHTTPCodeInvalid
	}

	isAllowedImageType := IsImageHeader(respBytesFirst, policy.imageTypes)

	if !isAllowedImageType {
		return false, respBytesFirst, resp, ErrUnsupportedImageType
	}

	// Check image file size
	// Some CDNs (like google) use content-range instead of content length
	contentLength := int64(0)
	if resp.StatusCode == http.StatusPartialContent {
		contentRange := resp.Header.Get("Content-Range")
		byteRanges := strings.Split(contentRange, "/")
		if len(contentRange) > 2 {
			rangeLength, err := strconv.Atoi(byteRanges[1])
			if err != nil {
				contentLength = resp.ContentLength
			}
			contentLength = int64(rangeLength)
		}
	} else {
		contentLength = resp.ContentLength
	}
	if contentLength < int64(policy.minImageFileBytesSize) {
		return false, respBytesFirst, resp, ErrFilesizeTooSmall
	}

	// Check image file size
	if policy.maxImageFileBytesSize > 0 && contentLength > int64(policy.maxImageFileBytesSize) {
		return false, respBytesFirst, resp, ErrFilesizeTooLarge
	}

	if !policy.fetchFullImage {
		return true, respBytesFirst, resp, nil
	}
	resp, err = i.fetchRemainingBytes(url)
	if err != nil {
		return false, respBytesFirst, resp, err
	}
	defer resp.Body.Close()

	// Read the remaining bytes
	respBytesRemaining, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, respBytesFirst, resp, err
	}
	// Combine the bytes
	combinedBytes := append(respBytesFirst, respBytesRemaining...)

	// Return true if all validations passed
	return true, combinedBytes, resp, nil
}

// IsImageHeader checks if the provided bytes represent an image header for the specified image types
func IsImageHeader(data []byte, imageTypes []ImageType) bool {
	for _, imageType := range imageTypes {
		if ok := testImageHeader(data, imageType); ok {
			return true
		}
	}
	return false
}

// testImageHeader returns the header bytes for the specified image type
func testImageHeader(data []byte, imageType ImageType) bool {
	switch imageType {
	case JPEG:
		return bytes.HasPrefix(data, []byte{0xFF, 0xD8})
	case PNG:
		return bytes.HasPrefix(data, []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1A, '\n'})
	case GIF:
		return bytes.HasPrefix(data, []byte("GIF89a")) || bytes.HasPrefix(data, []byte("GIF87a"))
	case BMP:
		return bytes.HasPrefix(data, []byte("BM"))
	case WEBP:
		return bytes.HasPrefix(data, []byte("RIFF")) && bytes.HasPrefix(data[8:len(data)-1], []byte("WEBP"))
	case JFIF:
		return bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF})
	case JPEG2000:
		return bytes.HasPrefix(data, []byte{0x00, 0x00, 0x00, 0x0C, 0x6A, 0x50, 0x20, 0x20, 0x0D, 0x0A, 0x87, 0x0A})
	case EXIF:
		return bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF})
	case HEIC:
		return bytes.HasPrefix(data, []byte("ftypheic"))
	default:
		return false
	}
}

// FetchFirstBytes fetches the first n bytes of the file from the given URL using HTTP Range header
func (i *fastimageImpl) fetchFirstBytes(url string) (*http.Response, error) {
	return i.fetchBytes(url, 0, maxBytes)
}

func (i *fastimageImpl) fetchRemainingBytes(url string) (*http.Response, error) {
	return i.fetchBytes(url, maxBytes, -1)
}

// fetchBytes fetches the bytes from the given URL using HTTP Range header
func (i *fastimageImpl) fetchBytes(url string, start, end int) (*http.Response, error) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set Range header
	if end == -1 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", start))
	} else {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	}

	// Send the request
	resp, err := i.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
