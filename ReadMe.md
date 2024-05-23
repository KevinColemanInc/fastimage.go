# FastImage.go

This library helps developers validate untrusted image urls before passing the urls or bytes to upstream services.

## Quickstart

```go

policy := fastimage.NewPolicy().
  WithImageTypes(fastImage.JPG, fastImage.WEBP, fastImage.PNG).
  WithMaxImageFileSizeInBytes(10000).
  WithMinImageFileSizeInBytes(100).
  WithFetchFullImage(true) // default false. Fetches the rest of the image if the validations pass
  
httpClient := http.Client{}
ok, imgBytes, imgResponse, errs := fastImage.New(httpClient).ValidateURL(url, policy)

if err != nil {
    fmt.Println("error occured", string(imgBytes), err.Error())
}

if !ok {
    for err := range errs {
        fmt.Println("Image is not valid", string(imgBytes), err.Error())
    }
}

```

## Validations

### Verify Image type

FastImage.go downloads the first 12 bytes of the image file to verify the url is actually an image (and not text with an image file extension)

### Verify Image size

FastImage.go uses the `Content-Length` http header to prevent downloading images that are too big and could cause OOM crashes for unexpectantly large images

### Verify Image exists

FastImage.go performs a GET RANGE for 12 bytes request to verify the image is fetchable (not not going to 4xx)

## Roadmap

- [ ] Support b64 and byte array images
- [ ] Improve test coverage
