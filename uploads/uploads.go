package uploads

import (
	"context"
	"errors"
	"fmt"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/gcerrors"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strconv"
)

var (
	supportedTypes = []string{"application/jpeg", "application/tiff", "application/png"}
)

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

	ResourceRequest struct {
		File FileHeader `formFile:"image"`
		Name string     `form:"name"`
	}
)

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

func (bs *BlobStore) GetImage(name string, size string) (b []byte, err error) {
	var buck *blob.Bucket
	if buck, err = bs.getBucket(name); err == nil {
		defer func() { _ = buck.Close() }()
		b, err = buck.ReadAll(bs.ctx, name)
		if gcerrors.Code(err) == gcerrors.NotFound {
			err = errors.New("not found")
		}
	}
	//TODO: if size not there yet, resize it and store
	//ResizeImage()
	err = bs.StoreImage(name, size, b)
	return
}

func (bs *BlobStore) StoreRawImage(name string, size string, req FileHeader) error {
	data, _, err := ReadImageAndValidate(req)
	if err != nil {
		return err
	}

	return bs.StoreImage(name, size, data)
}

func (bs *BlobStore) StoreImage(name string, size string, data []byte) error {

	buck, err := bs.getBucket(name)
	if err != nil {
		return err
	}

	defer func() { _ = buck.Close() }()
	return buck.WriteAll(bs.ctx, fmt.Sprintf("%s/%s", name, size), data, nil)
}

func (bs *BlobStore) keyString(key uint64) string {
	return strconv.FormatUint(key, 10)
}

func (bs *BlobStore) getBucket(name string) (*blob.Bucket, error) {
	return bs.bucket(name)
}

func (bs *BlobStore) bucket(seg ...string) (*blob.Bucket, error) {
	path := filepath.Join(seg...)
	path = filepath.Join(bs.base, path)
	return bs.store.OpenBucketURL(bs.ctx, &url.URL{Path: path})
}
