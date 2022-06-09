package uploads

import (
	"bytes"
	"context"
	"errors"
	"github.com/disintegration/imaging"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/gcerrors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

// add caching, repository

func ResizeImage(mimeType string, size int, res []byte) ([]byte, error) {
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

type (
	// BlobStore is the representation of the file persistence
	BlobStore struct {
		base  string
		store blob.BucketURLOpener
		ctx   context.Context
	}

	FileHeader interface {
		Open() (multipart.File, error)
	}

	Store struct {
		Save     func() error
		Close    func()
		MimeType string
		Length   uint64
	}
)

//TODO: if in map / bucket, retrieve, otherwise resize and add to bucket

func NewStorage(ctx context.Context) (*BlobStore, error) {
	//adding  a  cloud store is possible
	opener :=
		&fileblob.URLOpener{
			Options: fileblob.Options{
				CreateDir: true,
				Metadata:  fileblob.MetadataDontWrite,
			},
		}
	path := filepath.Join("data")
	return &BlobStore{ctx: ctx, store: opener, base: filepath.Clean(path)}, nil
}

func (bs *BlobStore) GetImage(name uint64, size uint) (b []byte, err error) {
	var buck *blob.Bucket
	if buck, err = bs.getBucket(name); err == nil {
		defer func() { _ = buck.Close() }()
		b, err = buck.ReadAll(bs.ctx, name)
		if gcerrors.Code(err) == gcerrors.NotFound {
			err = errors.New("not found")
		}
	}
	return
}
func (bs *BlobStore) StoreSample(appName uint64, req FileHeader) error {
	data, mimeType, err := ReadImage(req)
	if err != nil {
		return err
	}

	buck, err := bs.getBucket(appName)
	if err != nil {
		return err
	}

	defer func() { _ = buck.Close() }()
	return buck.WriteAll(bs.ctx, "0", data, nil)
}

func (bs *BlobStore) keyString(key uint64) string {
	return strconv.FormatUint(key, 10)
}

func (bs *BlobStore) getBucket(name string) (*blob.Bucket, error) {
	return bs.bucket(keyString(name)), nil
}

func ReadImage(req FileHeader) ([]byte, string, error) {
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

	if l := len(data); l == 0 || (l > 2*10^7) {
		return nil, "", errors.New("invalid payload")
	}

	switch mimeType := http.DetectContentType(data); mimeType {
	case "image/tiff",
		"image/jpeg",
		"image/png":
		return data, mimeType, nil
	default:
		return nil, "", errors.New("invalid payload")
	}
}

func (bs *BlobStore) bucket(seg ...string) (*blob.Bucket, error) {
	path := filepath.Join(seg...)
	path = filepath.Join(bs.base, path)
	return bs.store.OpenBucketURL(bs.ctx, &url.URL{Path: path})
}
