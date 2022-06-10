package uploads

import (
	"bytes"
	"errors"
	"github.com/disintegration/imaging"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	MaxSize = 2*10 ^ 7 // 20 mb
)

// ResizeImage extracts the image type from the file upload and then resizes it to specific type
func ResizeImage(mimeType string, size int, res []byte) ([]byte, error) {
	if size < 30 || size > 3000 {
		return nil, errors.New("specify size between {30,3000} px ")
	}
	if format, err := imaging.FormatFromExtension(strings.TrimPrefix(mimeType, "image/")); err != nil {
		return nil, errors.New("unreadable")
	} else if img, err2 := imaging.Decode(bytes.NewReader(res)); err2 != nil {
		return nil, errors.New("unreadable")
	} else {
		buffer := new(bytes.Buffer)
		croppedImg := imaging.Fill(img, size, size, imaging.Center, imaging.NearestNeighbor)
		err = imaging.Encode(buffer, croppedImg, format)
		return buffer.Bytes(), err
	}
}

// ReadImageAndValidate extracts content and file type from the raw image
func ReadImageAndValidate(req FileHeader) ([]byte, string, error) {
	if req == nil {
		return nil, "", errors.New("invalid payload")
	}

	file, err := req.Open()
	if err != nil {
		return nil, "", errors.New("unreadable")
	}
	defer func() { _ = file.Close() }()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, "", errors.New("unreadable")
	}

	if l := len(data); l == 0 || (l > MaxSize) {
		return nil, "", errors.New("file too big")
	}

	switch mimeType := http.DetectContentType(data); mimeType {
	case "image/tiff",
		"image/jpeg",
		"image/png":
		return data, mimeType, nil
	default:
		return nil, "", errors.New("invalid file type")
	}
}
